package chat_message_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"

	"github.com/samber/lo"
)

type ChatMessageRepo interface {
	GetMessages(req *common.RequestContext, d dto.ChatMessagesRequest) (res []dto.ChatMessageDto, err error)
	CreateMessage(req *common.RequestContext, d dto.ChatMessageData) (res dto.ChatMessageDto,err error)
	GetMembersIds(chatID int64)(res []int64,err error)
}

type chatMessageRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewChatMessageRepo(
	conn db.Connection,
	helpers *helpers.Helpers,
) ChatMessageRepo {
	return &chatMessageRepo{
		convertor: helpers.Convertor,
		Q: conn.GetQ(),
	}
}
func(r *chatMessageRepo)GetMembersIds(chatID int64)(res []int64,err error) {
	members,err := r.Q.ChatMember.Select(r.Q.ChatMember.ProfileID).Where(
		r.Q.ChatMember.ChatID.Eq(chatID),
	).Find()
	if err != nil {
		return
	}
	res = lo.Map(members,func(m *model.ChatMember,i int) int64 {
		return m.ProfileID
	})
	return
}

func (r *chatMessageRepo) GetMessages(req *common.RequestContext, d dto.ChatMessagesRequest) (res []dto.ChatMessageDto, err error) {
	chatID := r.convertor.StrtoInt(d.ID)
	mQ := r.Q.ChatMessage
	pQ := r.Q.Profile
	err = mQ.WithContext(req.Ctx).Select(
		mQ.ID, mQ.CreatedAt, mQ.Content, mQ.Type,
		mQ.ProfileID, pQ.FamilyName.As("profile_fn"), pQ.GivenName.As("profile_gn"),
	).Join(
		pQ,pQ.ID.EqCol(mQ.ProfileID),
	).Where(
		mQ.ChatID.Eq(chatID),
	).Limit(domain.DEFAULT_LIMIT).Offset(d.Page * domain.DEFAULT_LIMIT).Order(
	 mQ.CreatedAt.Desc(),
	).Scan(&res)
	return
}

func (r *chatMessageRepo) CreateMessage(req *common.RequestContext, d dto.ChatMessageData) (res dto.ChatMessageDto,err error) {
	fields := d.Fields
	message := model.ChatMessage{}
	if err = r.convertor.CopyStructData(fields, &message); err != nil {
		return
	}
	err = r.Q.WithContext(req.Ctx).ChatMessage.Save(&message)
	if err != nil {
		return
	}
	res = dto.ChatMessageToDto(message)
	res.ProfileFn = d.ProfileFn
	res.ProfileGn = d.ProfileGn
	return
}
