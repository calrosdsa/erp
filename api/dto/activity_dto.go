package dto

import "time"

// import "erp/gen/db/model"

type (
	ActivityDataRequest struct {
		Body ActivityData
	}
	ActivityData struct {
		PartyID          int64                `json:"party_id" required:"true"`
		PartyName        string               `json:"party_name"`
		EntityID         int                  `json:"entity_id"`
		Type             string               `json:"type" required:"true"`
		IsPinned         *bool                `json:"is_pinned" required:"false"`
		ActivityDeadLine ActivityDeadlineData `json:"activity_deadline" required:"false"`
		ActivityComment  ActivityCommentData  `json:"activity_comment" required:"false"`
		//Activitye
	}
	ActivityDeadlineData struct {
		Fields ActivityDeadlineFields `json:"fields"`
	}

	ActivityDeadlineFields struct {
		ActivityID  int32     `json:"activity_id"`
		Link        *string   `json:"link" required:"false"`
		PartyID     *int64    `json:"party_id" required:"false"`
		Deadline    time.Time `json:"deadline" required:"true"`
		Address     *string   `json:"address" required:"false"`
		Title       *string   `json:"title" required:"false"`
		Content     *string   `json:"content" required:"false"`
		Color       string    `json:"color"`
		IsCompleted bool      `json:"is_completed"`
		ProfileID   *int64     `json:"profile_id"`
	}

	ActivityCommentData struct {
		Mentions []ActivityMentionData `json:"mentions"`
		Fields   ActivityCommentFields `json:"fields"`
	}

	ActivityCommentFields struct {
		ActivityID int32  `json:"activity_id"`
		Comment    string `json:"comment" required:"true" minLength:"1"`
	}

	ActivityMentionData struct {
		Action string                `json:"action"`
		ID     int32                 `json:"id"`
		Fields ActivityMentionFields `json:"fields"`
	}

	ActivityMentionFields struct {
		ProfileID  int64 `json:"profile_id"`
		ActivityID int32 `json:"activity_id"`
		StartIndex int32 `json:"start_index" required:"true"`
		EndIndex   int32 `json:"end_index" required:"true"`
	}

	ActivityMentionDto struct {
		ID         int32 `json:"id"`
		ActivityID int32 `json:"activity_id"`
		StartIndex int32 `json:"start_index"`
		EndIndex   int32 `json:"end_index"`

		ProfileUUID string `json:"profile_uuid"`
		GivenName   string `json:"given_name"`
		FamilyName  string `json:"family_name"`
	}

	ActivityDto struct {
		ID        int32     `json:"id"`
		Type      string    `json:"type"`
		// Action    string    `json:"action"`
		CreatedAt time.Time `json:"created_at"`
		IsPinned  *bool     `json:"is_pinned"`
		Data      *string   `json:"data"`
		ProfileID int64	 `json:"profile_id"`
		// Message   string    `json:"message"`

		//Comment
		Comment *string `json:"comment"`
		//Mentions
		Mentions []ActivityMentionDto `json:"mentions" gorm:"-"`

		//Acvity
		Link        *string    `json:"link"`
		PartyID     *int64     `json:"party_id"`
		Deadline    *time.Time `json:"deadline"`
		Address     *string    `json:"address"`
		Title       *string    `json:"title"`
		Content     *string    `json:"content"`
		Color       *string    `json:"color"`
		IsCompleted *bool      `json:"is_completed"`

		//Profile
		ProfileGivenName  string  `json:"profile_given_name"`
		ProfileFamilyName string  `json:"profile_family_name"`
		ProfileUUID       string  `json:"profile_uuid"`
		ProfileAvatar     *string `json:"profile_avatar"`
	}
)

// func ActivityDtoFromModel(m *model.Activity) ActivityDto {
// 	return ActivityDto{
// 		ID: m.ID,
// 		Type: ,
// 	}
// }
