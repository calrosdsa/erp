package dto

import "time"

type (
	ChatsRequest struct {
		Name string `query:"name"`
	}
	ChatData struct {
		ID     int64      `json:"id"`
		Fields ChatFields `json:"fields"`
	}

	ChatFields struct {
		Name      string `json:"name"`
		PartyID   int64  `json:"party_id"`
		Type      string `json:"type"`
		PartyType string `json:"party_type"`
	}

	ChatDetailDto struct {
		ID         int64           `json:"id"`
		CreatedAt  time.Time       `json:"created_at"`
		Name       string          `json:"name"`
		Type       string          `json:"type"`
		Members    []ChatMemberDto `gorm:"-" json:"members"`
		PartyID    int64           `json:"party_id"`
		EntityHref string          `json:"entity_href"`
	}

	ChatDto struct {
		ID                   int64      `json:"id"`
		Name                 *string    `json:"name"`
		Type                 string     `json:"type"`
		UnreadCount          int        `json:"unread_count"`
		ProfileID            *int64     `json:"profile_id"`
		ProfileGn            *string    `json:"profile_gn"`
		ProfileFn            *string    `json:"profile_fn"`
		LastMessageContent   *string    `json:"last_message_content"`
		LastMessageType      *string    `json:"last_message_type"`
		LastMessageCreatedAt *time.Time `json:"last_message_created_at"`
	}
)
