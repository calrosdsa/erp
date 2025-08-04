package account_service

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/connection"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"time"
)

type SessionService struct {
	conn    *connection.Connection
	timeout time.Duration
	emitLog helpers.EmitLog
}

func NewSessionService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
) *SessionService {
	return &SessionService{
		conn:    conn,
		timeout: timeout,
		emitLog: helpers.Logger.EmitLog("user-relation-service"),
	}
}



func (s *SessionService) GetUserRelations(req *common.RequestContext) ([]dto.UserRelationDto, error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	res, err := s.conn.Q.UserRelation.WithContext(ctx).Where(
		s.conn.Q.UserRelation.UserID.Eq(req.User.ID),
	).Preload(s.conn.Q.UserRelation.Company).
	Preload(s.conn.Q.UserRelation.Profile).
	Preload(s.conn.Q.UserRelation.Role).
		Limit(50).
		Find()
	if err != nil {
		s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
		return []dto.UserRelationDto{}, err
	}
	userRelationDtos := make([]dto.UserRelationDto,len(res))
	for i,userRelation := range res {
		userRelationDtos[i] = dto.UserRelationDtoFromModel(userRelation)
	}
	return userRelationDtos, err
}

func (s *SessionService) GetUserRelation(ctx context.Context, uuid string) (model.UserRelation, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	var (
		res *model.UserRelation
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
		}
	}()
	res, err = s.conn.Q.UserRelation.Where(
		s.conn.Q.UserRelation.UUID.Eq(uuid),
	).Preload(s.conn.Q.UserRelation.Company).
		Preload(s.conn.Q.UserRelation.User).
		Preload(s.conn.Q.UserRelation.Profile).
		Preload(s.conn.Q.UserRelation.Role).First()
	if err != nil {
		return model.UserRelation{},err
	}
	// err = s.conn.Db.WithContext(ctx).Where(&entity.UserRelation{Uuid: uuid}).Preload(clause.Associations).First(&res).Error
	// if err != nil {
	// 	s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
	// 	return res, err
	// }
	return *res, err
}
