package chat_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/event"

	"github.com/samber/lo"
)

type ChatEventRepo interface {
	OnUserAdded(ctx context.Context, payload event.UserCreatedEventData) (err error)
}

type chatEventRepo struct {
}

func NewChatEventRepo() ChatEventRepo {
	return &chatEventRepo{}
}

func (r *chatEventRepo) OnUserAdded(ctx context.Context, payload event.UserCreatedEventData) (err error) {
	tx := payload.Tx
	profiles, err := tx.Profile.Select(tx.Profile.ID).Where(
		tx.Profile.CompanyID.Eq(payload.UseRelation.CompanyID),
	).Find()
	if err != nil {
		return
	}
	profiles = lo.Filter(profiles, func(p *model.Profile, i int) bool {
		return p.ID != payload.UseRelation.ProfileID
	})
	chats := make([]*model.Chat, len(profiles))
	var chatMembers []*model.ChatMember
	for i, profile := range profiles {
		chatID, err1 := tx.Chat.InsertParty(proto.PartyType_chat.String())
		if err1 != nil {
			return err1
		}
		chat := &model.Chat{
			ID:        chatID,
			CompanyID: payload.UseRelation.CompanyID,
			Type:      proto.ChatType_Conversation.String(),
			PartyID:   &profile.ID,
			EntityID:  int32(domain.USER.ID),
		}
		chats[i] = chat

		member1 := &model.ChatMember{
			ChatID:    chatID,
			ProfileID: profile.ID,
		}
		member2 := &model.ChatMember{
			ChatID:    chatID,
			ProfileID: payload.UseRelation.ProfileID,
		}
		chatMembers = append(chatMembers, member1, member2)
	}
	if err = tx.Chat.CreateInBatches(chats, domain.DEFAULT_BATCH_SIZE); err != nil {
		return
	}
	if err = tx.ChatMember.CreateInBatches(chatMembers, domain.DEFAULT_BATCH_SIZE); err != nil {
		return
	}

	return
}
