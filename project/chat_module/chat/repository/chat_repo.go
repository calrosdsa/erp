package chat_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"time"
)

type ChatRepository interface {
	GetChat(req *common.RequestContext, d dto.RequestEntity) (res dto.ChatDetailDto, err error)
	GetChats(req *common.RequestContext, d dto.ChatsRequest) (res []dto.ChatDto, err error)
	UpdateMemberLastRead(req *common.RequestContext,d dto.RequestEntity)(err error)
}

type chatRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewChatRepo(
	conn db.Connection,
	helpers *helpers.Helpers,
) ChatRepository {
	return &chatRepo{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}

// func (r *chatRepo) CreateChat(req *common.RequestContext,)

func (r *chatRepo) GetChat(req *common.RequestContext, d dto.RequestEntity) (res dto.ChatDetailDto, err error) {
	chatQ := r.Q.Chat
	chatID := r.convertor.StrtoInt(d.ID)
	query := `	SELECT
    cm.chat_id as id,
    coalesce(c.name,concat(p.given_name,' ',p.family_name),'') as name,
	c.type,
    cm.last_read_at,
	c.party_id,
	e.href as entity_href
	FROM chats c
	JOIN chat_members cm ON c.id = cm.chat_id
	JOIN entities e ON e.id = c.entity_id
	LEFT JOIN chat_members cm2 ON c.type = $1 AND cm2.profile_id != cm.profile_id 
	AND cm2.chat_id = cm.chat_id
	LEFT JOIN profiles p ON p.id = cm2.profile_id 
	WHERE cm.chat_id = $2 and cm.profile_id = $3`
	err = chatQ.WithContext(req.Ctx).UnderlyingDB().Raw(query, proto.ChatType_Conversation.String(),
		chatID, req.Profile.ID).Scan(&res).Error
	if err != nil {
		return
	}
	err = r.updateMemberLastRead(req.Ctx, req.Profile.ID, chatID)
	if err != nil {
		return
	}

	var members []dto.ChatMemberDto
	cmQ := r.Q.ChatMember
	profileQ := r.Q.Profile
	err = cmQ.WithContext(req.Ctx).Select(
		cmQ.ChatID, cmQ.ProfileID,
		profileQ.FamilyName.As("profile_fn"), profileQ.GivenName.As("profile_gn"),
	).Join(
		profileQ, profileQ.ID.EqCol(cmQ.ProfileID),
	).Where(
		cmQ.ChatID.Eq(chatID),
	).Scan(&members)
	if err != nil {
		return
	}
	res.Members = members
	return
}

func(r *chatRepo) UpdateMemberLastRead(req *common.RequestContext,d dto.RequestEntity)(err error){
	chatID := r.convertor.StrtoInt(d.ID)
	err = r.updateMemberLastRead(req.Ctx,req.Profile.ID,chatID)
	return
}

func (r *chatRepo) updateMemberLastRead(ctx context.Context, profileID int64, chatID int64) (err error) {
	mQ := r.Q.ChatMember
	_, err = mQ.WithContext(ctx).Where(
		mQ.ChatID.Eq(chatID),
		mQ.ProfileID.Eq(profileID),
	).UpdateSimple(mQ.LastReadAt.Value(time.Now()))
	return
}

func (r *chatRepo) GetChats(req *common.RequestContext, d dto.ChatsRequest) (res []dto.ChatDto, err error) {
	query := `
		SELECT
    cm.chat_id as id,
    coalesce(c.name,concat(p.given_name,' ',p.family_name),'') as name,
	c.type,
    cm.last_read_at,
	cm.profile_id,p.given_name as profile_gn,p.family_name as profile_fn,
    last_message.content AS last_message_content,
    last_message.created_at AS last_message_created_at,
	last_message.type AS last_message_type,
    (SELECT COUNT(*) 
     FROM chat_messages 
     WHERE chat_id = cm.chat_id
       AND created_at > cm.last_read_at
       AND deleted_at IS NULL) AS unread_count
FROM chat_members cm
JOIN chats c ON c.id = cm.chat_id
LEFT JOIN chat_members cm2 ON c.type = $1 AND cm2.profile_id != cm.profile_id 
AND cm2.chat_id = cm.chat_id
LEFT JOIN profiles p ON p.id = cm2.profile_id 
LEFT JOIN LATERAL (
    SELECT content, created_at,type
    FROM chat_messages
    WHERE chat_id = cm.chat_id
      AND deleted_at IS NULL
    ORDER BY created_at DESC
    LIMIT 1
) last_message ON true
 WHERE cm.profile_id = $2 and coalesce(c.name,concat(p.given_name,' ',p.family_name),'') LIKE $3
ORDER BY last_message_created_at DESC NULLS LAST LIMIT 100;
	`
	err = r.Q.Chat.UnderlyingDB().Raw(query, proto.ChatType_Conversation.String(), req.Profile.ID,
		"%"+d.Name+"%").Scan(&res).Error
	return
}
