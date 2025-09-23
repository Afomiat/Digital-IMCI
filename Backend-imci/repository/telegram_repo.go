package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Afomiat/Digital-IMCI/domain"
)


type telegramRepository struct {
	db *pgxpool.Pool
}

func NewTelegramRepository(db *pgxpool.Pool) domain.TelegramRepository {
	return &telegramRepository{db: db}
}

func (t *telegramRepository) SaveChatID(ctx context.Context, username string, chatID int64, phone string) error {
	query := `
		INSERT INTO telegram_chat_ids (telegram_username, chat_id, phone_number)
		VALUES ($1, $2, $3)
		ON CONFLICT (telegram_username) 
		DO UPDATE SET chat_id = $2, phone_number = $3, updated_at = NOW()
	`
	
	_, err := t.db.Exec(ctx, query, username, chatID, phone)
	if err != nil {
		return fmt.Errorf("failed to save telegram chat ID: %w", err)
	}
	return nil
}

func (t *telegramRepository) GetChatIDByUsername(ctx context.Context, username string) (int64, error) {
	var chatID int64
	query := `SELECT chat_id FROM telegram_chat_ids WHERE telegram_username = $1`
	
	err := t.db.QueryRow(ctx, query, username).Scan(&chatID)
	if err != nil {
		return 0, fmt.Errorf("failed to get chat ID for username %s: %w", username, err)
	}
	
	return chatID, nil
}

func (t *telegramRepository) GetUsernameByPhone(ctx context.Context, phone string) (string, error) {
	var username string
	query := `SELECT telegram_username FROM telegram_chat_ids WHERE phone_number = $1`
	
	err := t.db.QueryRow(ctx, query, phone).Scan(&username)
	if err != nil {
		return "", fmt.Errorf("failed to get username for phone %s: %w", phone, err)
	}
	
	return username, nil
}

func (t *telegramRepository) DeleteChatID(ctx context.Context, username string) error {
	query := `DELETE FROM telegram_chat_ids WHERE telegram_username = $1`
	_, err := t.db.Exec(ctx, query, username)
	if err != nil {
		return fmt.Errorf("failed to delete chat ID for username %s: %w", username, err)
	}
	return nil
}