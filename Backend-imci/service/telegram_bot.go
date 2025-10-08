package service

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/Afomiat/Digital-IMCI/domain"
	"github.com/Afomiat/Digital-IMCI/internal/userutil"
)

type telegramBotService struct {
	bot          *tgbotapi.BotAPI
	telegramRepo domain.TelegramRepository
	otpRepo      domain.OtpRepository
	botUsername  string
	token        string
	isRunning    bool
	mu           sync.RWMutex
}

func NewTelegramBotService(token string, telegramRepo domain.TelegramRepository, otpRepo domain.OtpRepository) (domain.TelegramService, error) {
	// Don't initialize the bot immediately, just store the token
	service := &telegramBotService{
		token:        token,
		telegramRepo: telegramRepo,
		otpRepo:      otpRepo,
		isRunning:    false,
	}

	// Initialize bot API but don't start polling yet
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	bot.Debug = true
	service.bot = bot
	service.botUsername = bot.Self.UserName

	log.Printf("Telegram bot service created for account %s (polling not started)", service.botUsername)

	return service, nil
}

func (t *telegramBotService) StartPolling() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isRunning {
		return nil // Already running
	}

	if t.bot == nil {
		return fmt.Errorf("bot not initialized")
	}

	go t.startPolling()
	t.isRunning = true
	log.Println("Started Telegram bot polling...")
	return nil
}

func (t *telegramBotService) StopPolling() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.bot != nil {
		t.bot.StopReceivingUpdates()
		t.isRunning = false
		log.Println("Stopped Telegram bot polling")
	}
}

func (t *telegramBotService) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.isRunning
}

func (t *telegramBotService) startPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		t.handleUpdate(update)
	}
}

// Rest of your existing methods remain the same...
func (t *telegramBotService) GetStartLink() string {
	return fmt.Sprintf("https://t.me/%s", t.botUsername)
}

func (t *telegramBotService) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	// Handle /start command
	if update.Message.IsCommand() && update.Message.Command() == "start" {
		t.handleStartCommand(update)
		return
	}

	// Handle contact sharing
	if update.Message.Contact != nil {
		t.handleContact(update)
		return
	}

	// Handle help command
	if update.Message.IsCommand() && update.Message.Command() == "help" {
		t.sendHelpMessage(update.Message.Chat.ID)
		return
	}
}


func (t *telegramBotService) handleStartCommand(update tgbotapi.Update) {


	args := strings.TrimSpace(update.Message.CommandArguments())
	chatID := update.Message.Chat.ID

	if args == "" {
		// Ask for phone number
		btn := tgbotapi.NewKeyboardButtonContact("📱 Share phone number")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(btn),
		)

		msg := tgbotapi.NewMessage(chatID, "Please share your phone number to continue.")
		msg.ReplyMarkup = keyboard
		t.bot.Send(msg)
		return
	}

	// Deep link path → use provided phone directly
	t.linkAccountAndSendOTP(update.Message.From.UserName, chatID, args)
}

func (t *telegramBotService) handleContact(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	contact := update.Message.Contact

	if contact == nil {
		return
	}

	phone := contact.PhoneNumber
	username := update.Message.From.UserName

	// Remove the keyboard after sharing
	msg := tgbotapi.NewMessage(chatID, "✅ Phone received. Linking your account...")
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	t.bot.Send(msg)

	t.linkAccountAndSendOTP(username, chatID, phone)
}

func (t *telegramBotService) linkAccountAndSendOTP(username string, chatID int64, phone string) {
	ctx := context.Background()

	// Save Telegram mapping
	if err := t.telegramRepo.SaveChatID(ctx, username, chatID, phone); err != nil {
		log.Printf("Failed to save chat ID: %v", err)
		t.reply(chatID, "❌ Failed to link your account. Please try again.")
		return
	}

	log.Printf("Successfully saved Telegram mapping for @%s", username)

	// Generate OTP
	otpCode := userutil.GenerateOTP()
	otp := &domain.OTP{
		Phone:     phone,
		Code:      otpCode,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}
	if err := t.otpRepo.SaveOTP(ctx, otp); err != nil {
		log.Printf("Failed to save OTP: %v", err)
		t.reply(chatID, "❌ Error generating OTP. Please try again.")
		return
	}

	// Send OTP (HTML mode)
	otpMessage := fmt.Sprintf(
		"🔐 <b>Digital IMCI Verification</b>\n\n"+
			"Your OTP code is: <code>%s</code>\n\n"+
			"⏰ <b>Expires in:</b> 5 minutes\n\n"+
			"Enter this code in the app to complete your registration!",
		otpCode,
	)

	msg := tgbotapi.NewMessage(chatID, otpMessage)
	msg.ParseMode = "HTML"
	if _, err := t.bot.Send(msg); err != nil {
		log.Printf("Failed to send OTP message: %v", err)
		t.reply(chatID, "❌ Failed to send OTP. Please try again.")
		return
	}

	// Success summary
	// successMsg := fmt.Sprintf(
	// 	"✅ <b>Account Linked Successfully!</b>\n\n"+
	// 		"• Telegram: @%s\n"+
	// 		"• Phone: %s\n\n"+
	// 		"Return to the Digital IMCI app and enter the OTP code to complete your registration.",
	// 	username, phone,
	// )
	// msg2 := tgbotapi.NewMessage(chatID, successMsg)
	// msg2.ParseMode = "HTML"
	// t.bot.Send(msg2)

	// log.Printf("Successfully completed Telegram linking for @%s", username)
}

func (t *telegramBotService) reply(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = "HTML"
	_, err := t.bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
}

func (t *telegramBotService) sendHelpMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID,
		"🤖 <b>Digital IMCI Telegram Bot</b>\n\n"+
			"To link your account:\n"+
			"1. Visit our app and start signup\n"+
			"2. Click 'Link Telegram' button\n"+
			"3. Use the provided link\n\n"+
			"Your OTP will be sent automatically after linking!",
	)
	msg.ParseMode = "HTML"
	t.bot.Send(msg)
}

func (t *telegramBotService) SendOTP(ctx context.Context, telegramUsername, code string) error {
	if len(telegramUsername) > 0 && telegramUsername[0] == '@' {
		telegramUsername = telegramUsername[1:]
	}

	chatID, err := t.telegramRepo.GetChatIDByUsername(ctx, telegramUsername)
	if err != nil {
		return fmt.Errorf("user @%s has not linked their Telegram account", telegramUsername)
	}

	messageText := fmt.Sprintf(
		"🔐 <b>Digital IMCI Verification</b>\n\nYour code is: <code>%s</code>\n⏰ Expires in 5 minutes",
		code,
	)
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "HTML"

	_, err = t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message to @%s: %w", telegramUsername, err)
	}

	log.Printf("OTP %s sent successfully to @%s (Chat ID: %d)", code, telegramUsername, chatID)
	return nil
}

func (t *telegramBotService) SendPasswordResetOTP(ctx context.Context, telegramUsername, code string) error {
	if len(telegramUsername) > 0 && telegramUsername[0] == '@' {
		telegramUsername = telegramUsername[1:]
	}

	chatID, err := t.telegramRepo.GetChatIDByUsername(ctx, telegramUsername)
	if err != nil {
		return fmt.Errorf("user @%s has not linked their Telegram account", telegramUsername)
	}

	messageText := fmt.Sprintf(
		"🔐 <b>Password Reset Request</b>\n\n"+
			"Your password reset code is: <code>%s</code>\n\n"+
			"⏰ <b>Expires in:</b> 5 minutes\n\n"+
			"Enter this code in the app to reset your password.",
		code,
	)
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "HTML"

	_, err = t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message to @%s: %w", telegramUsername, err)
	}

	log.Printf("Password reset OTP %s sent successfully to @%s (Chat ID: %d)", code, telegramUsername, chatID)
	return nil
}