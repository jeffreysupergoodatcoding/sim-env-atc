package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/yegors/co-atc/internal/models"
	"github.com/yegors/co-atc/pkg/logger"
)

// We don't need local ChatMessage struct anymore

// ChatStorage handles storage of chat messages
type ChatStorage struct {
	db     *sql.DB
	logger *logger.Logger
}

// NewChatStorage creates a new SQLite chat storage
func NewChatStorage(db *sql.DB, log *logger.Logger) *ChatStorage {
	storage := &ChatStorage{
		db:     db,
		logger: log.Named("sqlite-chat"),
	}

	// Initialize database
	if err := storage.initDB(); err != nil {
		storage.logger.Error("Failed to initialize chat storage", logger.Error(err))
	}

	return storage
}

// initDB initializes the database tables
func (s *ChatStorage) initDB() error {
	// Create chat_messages table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS chat_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create chat_messages table: %w", err)
	}

	// Create indexes
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_chat_session_id ON chat_messages(session_id)`)
	if err != nil {
		return fmt.Errorf("failed to create session_id index: %w", err)
	}

	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_chat_timestamp ON chat_messages(timestamp)`)
	if err != nil {
		return fmt.Errorf("failed to create timestamp index: %w", err)
	}

	return nil
}

// StoreMessage stores a chat message
func (s *ChatStorage) StoreMessage(ctx context.Context, sessionID, role, content string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO chat_messages (session_id, role, content, timestamp) VALUES (?, ?, ?, ?)`,
		sessionID, role, content, time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return fmt.Errorf("failed to insert chat message: %w", err)
	}
	return nil
}

// GetMessagesBySession returns messages for a specific session
func (s *ChatStorage) GetMessagesBySession(sessionID string) ([]*models.ChatMessage, error) {
	rows, err := s.db.Query(
		`SELECT id, session_id, role, content, timestamp FROM chat_messages WHERE session_id = ? ORDER BY timestamp ASC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query chat messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		var id int64
		var timestamp string
		if err := rows.Scan(&id, &m.SessionID, &m.Type, &m.Content, &timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan chat message: %w", err)
		}

		m.ID = fmt.Sprintf("%d", id)

		t, err := time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		m.Timestamp = t
		messages = append(messages, &m)
	}

	return messages, nil
}
