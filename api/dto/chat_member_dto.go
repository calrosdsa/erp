package dto

import "time"

type (

	
	ChatMemberData struct {
		ID     int64            `json:"id"`
		Fields ChatMemberFields `json:"fields"`
	}
	ChatMemberFields struct {
		ProfileID int64 `json:"profile_id"`
		ChatID    int64 `json:"chat_id"`
	}
	ChatMemberDto struct {
		ID         int64     `json:"id"`
		ProfileID  int64     `json:"profile_id"`
		ProfileGn  string    `json:"profile_gn"`
		ProfileFn  string    `json:"profile_fn"`
		ChatID     int64     `json:"chat_id"`
		Chat       string    `json:"chat"`
		LastReadAt time.Time `json:"last_read_at"`
	}
)
