package dto

import (
	"erp/gen/db/model"
	"time"
)

type (
	ChatMessagesRequest struct {
		DefaultListParams
		ID   string `query:"id"`
		Page int    `query:"page"`
	}
	ChatMessageDataRequest struct {
		Body ChatMessageData
	}
	ChatMessageData struct {
		ID     int64             `json:"id" required:"false"`
		ProfileGn string `json:"profile_gn" required:"false"`
		ProfileFn string `json:"profile_fn" required:"false"`
		Fields ChatMessageFields `json:"fields"`
	}

	ChatMessageFields struct {
		ProfileID int64  `json:"profile_id"`
		Content   string `json:"content"`
		ChatID    int64  `json:"chat_id"`
		Type      string `json:"type"`
	}
	ChatMessageDto struct {
		ID        int64     `json:"id"`
		Content   string    `json:"content"`
		Type      string    `json:"type"`
		ProfileID int64     `json:"profile_id"`
		ProfileGn string    `json:"profile_gn"`
		ProfileFn string    `json:"profile_fn"`
		CreatedAt time.Time `json:"created_at"`
		ChatID    int64     `json:"chat_id"`
	}
)

func ChatMessageToDto(m model.ChatMessage) ChatMessageDto {
	return ChatMessageDto{
		Content:   m.Content,
		Type:      m.Type,
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		ChatID: m.ChatID,
		ProfileID: m.ProfileID,
	}
}
