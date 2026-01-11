package sender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
)

type TelegramConfig struct {
	Token string
}

type TelegramSender struct {
	token  string
	log    log.Log
	client *http.Client
}

func NewTelegramSender(token string, log log.Log) domain.Sender {
	return &TelegramSender{
		token: token,
		log:   log,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type tgMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func (s *TelegramSender) Send(ctx context.Context, n *domain.Notify) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", s.token)

	msg := tgMessage{
		ChatID: n.Target,
		Text:   string(n.Payload),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal tg message: %w", err)
	}

	reqCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram api request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned non-200 status: %d", resp.StatusCode)
	}

	s.log.Info().Str("target", n.Target).Msg("[TELEGRAM] Message sent successfully")
	return nil
}
