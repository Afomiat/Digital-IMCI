package service

import (
	"context"
	"fmt"
	"log"

	"github.com/Afomiat/Digital-IMCI/repository"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/Afomiat/Digital-IMCI/domain"
)

type telegramBotService struct {
	bot              *tgbotapi.BotAPI
	telegramRepo     repository.TelegramRepository
	botUsername      string // Store bot username for deep links
}

func NewTelegramBotService(token string, telegramRepo repository.TelegramRepository) (domain.TelegramService, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	bot.Debug = true
	log.Printf("Authorized on Telegram account %s", bot.Self.UserName)

	service := &telegramBotService{
		bot:          bot,
		telegramRepo: telegramRepo,
		botUsername:  bot.Self.UserName, // Store the bot username
	}

	// Start the polling loop
	go service.startPolling()

	return service, nil
}

// Add this one simple function to generate start links
func (t *telegramBotService) GetStartLink() string {
	return fmt.Sprintf("https://t.me/%s?start=signup", t.botUsername)
}

// The rest of your existing code remains exactly the same...
func (t *telegramBotService) startPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	log.Println("Started Telegram bot polling for /start commands...")
	for update := range updates {
		t.handleUpdate(update)
	}
}

func (t *telegramBotService) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	user := update.Message.From
	if user == nil {
		return
	}

	// Handle /start command or "Start" button click
	if update.Message.IsCommand() && update.Message.Command() == "start" || update.Message.Text == "🚀 Start" {
		// Save to database
		err := t.telegramRepo.SaveChatID(context.Background(), user.UserName, update.Message.Chat.ID)
		if err != nil {
			log.Printf("Failed to save chat ID for @%s: %v", user.UserName, err)
			return
		}

		log.Printf("User @%s started the bot. Chat ID: %d (saved to DB)", user.UserName, update.Message.Chat.ID)

		// Create MAIN keyboard with persistent buttons
		mainKeyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("📝 Sign Up"),
				tgbotapi.NewKeyboardButton("ℹ️ Help"),
			),
		)
		mainKeyboard.OneTimeKeyboard = false // Keyboard stays visible

		welcomeMsg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"👋 *Welcome to Digital IMCI!*\n\n"+
				"✅ Your Telegram is connected: @@"+user.UserName+"\n"+
				"🆔 Your Chat ID: `"+fmt.Sprintf("%d", update.Message.Chat.ID)+"`\n\n"+
				"*Ready for verification!* Use the buttons below:",
		)
		welcomeMsg.ParseMode = "Markdown"
		welcomeMsg.ReplyMarkup = mainKeyboard

		t.bot.Send(welcomeMsg)
	}

	// Handle button clicks
	if update.Message.Text == "📝 Sign Up" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"🎯 *Let's get you registered!*\n\n"+
				"1. Visit our website\n"+
				"2. Use your Telegram username: *@@"+user.UserName+"*\n"+
				"3. You'll receive your OTP here!",
		)
		msg.ParseMode = "Markdown"
		t.bot.Send(msg)
	}

	if update.Message.Text == "ℹ️ Help" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID,
			"❓ *Need Help?*\n\n"+
				"• *Sign Up:* Click the 'Sign Up' button \n"+
				"• *OTP Issues:* Make sure you started the bot  by sending '/start'\n"+
     			"*Your Telegram ID:* @@"+user.UserName,
		)
		msg.ParseMode = "Markdown"
		t.bot.Send(msg)
	}
}
func (t *telegramBotService) SendOTP(ctx context.Context, telegramUsername, code string) error {
	// Remove @ prefix if present
	if len(telegramUsername) > 0 && telegramUsername[0] == '@' {
		telegramUsername = telegramUsername[1:]
	}

	// Get chat ID from DATABASE instead of memory
	chatID, err := t.telegramRepo.GetChatIDByUsername(ctx, telegramUsername)
	if err != nil {
		return fmt.Errorf("user @%s has not started the bot. Please ask them to send /start to @%s first", 
			telegramUsername, t.bot.Self.UserName)
	}

	// ✅ FIX: Escape special characters for MarkdownV2
	// In MarkdownV2, these characters must be escaped: _ * [ ] ( ) ~ ` > # + - = | { } . !
	escapedCode := tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, code)
	
	// Send the OTP via Telegram
	messageText := fmt.Sprintf(
		"🔐 *Digital IMCI Verification*\n\nYour verification code is: `%s`\n\nThis code will expire in 5 minutes\\.",
		escapedCode,
	)
	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "MarkdownV2"

	_, err = t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send Telegram message to @%s: %w", telegramUsername, err)
	}

	log.Printf("OTP %s sent successfully to @%s (Chat ID: %d)", code, telegramUsername, chatID)
	return nil
}