package session

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	"erp/pkg/config"
	"erp/pkg/db"
	"erp/pkg/logger"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type sessionService struct {
	conn    db.Connection
	timeout time.Duration
	emitLog logger.EmitLog
	Q       *query.Query
	pgConfig config.PGConfig
	generator helpers.Generator
}

func New(
	conn db.Connection,
	logger logger.Logger,
	timeout time.Duration,
	helpers *helpers.Helpers,
	pgConfig config.PGConfig,
) repository.SessionService {
	return &sessionService{
		conn:    conn,
		emitLog: logger.EmitLog("session-service"),
		timeout: timeout,
		Q:       conn.GetQ(),
		pgConfig: pgConfig,
		generator:helpers.Generator,
	}
}

func (s *sessionService) GetCompanyDefaults(ctx context.Context, companyID int64) (
	res model.CompanyDefault, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompanyDefaults"))
		}
	}()
	companyDefaults, err := s.Q.CompanyDefault.WithContext(ctx).Where(
		s.Q.CompanyDefault.CompanyID.Eq(companyID),
	).First()
	if err != nil {
		return
	}
	res = *companyDefaults
	return
}

func (s *sessionService) GetUser(ctx context.Context, id int64) (res model.User, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUser"))
		}
	}()
	user, err := s.Q.User.WithContext(ctx).Where(
		s.Q.User.ID.Eq(id),
	).First()
	if err != nil {
		return
	}
	return *user, err
}



func (s *sessionService) GetUserRelationByUserID(ctx context.Context, userID int64) (model.UserRelation, error) {
	var (
		res *model.UserRelation
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
		}
	}()
	res, err = s.Q.UserRelation.Where(
		s.Q.UserRelation.UserID.Eq(userID),
	).Preload(s.Q.UserRelation.Company).
		Preload(s.Q.UserRelation.User).
		Preload(s.Q.UserRelation.Profile).
		Preload(s.Q.UserRelation.Role).First()
	if err != nil {
		return model.UserRelation{}, err
	}
	// err = s.conn.Db.WithContext(ctx).Where(&entity.UserRelation{Uuid: uuid}).Preload(clause.Associations).First(&res).Error
	// if err != nil {
	// 	s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
	// 	return res, err
	// }
	return *res, err
}

func (s *sessionService) GetRoleActionsByRole(req *common.RequestContext, roleID int64) (
	[]dto.RoleActionDto, error) {
	roleActions, err := s.Q.RoleAction.WithContext(req.Ctx).Where(
		s.Q.RoleAction.RoleID.Eq(roleID),
	).Preload(
		s.Q.RoleAction.Action,
	).Find()
	roleActionsDto := make([]dto.RoleActionDto, len(roleActions))
	for i, roleAction := range roleActions {
		roleActionsDto[i] = dto.RoleActionDTOFromModel(roleAction)
	}
	return roleActionsDto, err
}

func (s *sessionService) GetUserRelations(req *common.RequestContext) ([]dto.UserRelationDto, error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	Q := s.conn.GetQ()
	res, err := Q.UserRelation.WithContext(ctx).Where(
		Q.UserRelation.UserID.Eq(req.User.ID),
	).Preload(Q.UserRelation.Company).
		Preload(Q.UserRelation.Profile).
		Preload(Q.UserRelation.Role).
		Limit(50).
		Find()
	if err != nil {
		s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
		return []dto.UserRelationDto{}, err
	}
	userRelationDtos := make([]dto.UserRelationDto, len(res))
	for i, userRelation := range res {
		userRelationDtos[i] = dto.UserRelationDtoFromModel(userRelation)
	}
	return userRelationDtos, err
}

func (s *sessionService) GetUserRelation(ctx context.Context, uuid string) (model.UserRelation, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()
	Q := s.conn.GetQ()
	var (
		res *model.UserRelation
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
		}
	}()
	res, err = Q.UserRelation.Where(
		Q.UserRelation.UUID.Eq(uuid),
	).Preload(Q.UserRelation.Company).
		Preload(Q.UserRelation.User).
		Preload(Q.UserRelation.Profile).
		Preload(Q.UserRelation.Role).First()
	if err != nil {
		return model.UserRelation{}, err
	}
	// err = s.conn.Db.WithContext(ctx).Where(&entity.UserRelation{Uuid: uuid}).Preload(clause.Associations).First(&res).Error
	// if err != nil {
	// 	s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserRelation"))
	// 	return res, err
	// }
	return *res, err
}


func (s *sessionService) InsertUser(ctx context.Context, tx *query.QueryTx, identifier string) (model.User, error) {
	userPassword := s.generator.GeneratePassword()
	pass := s.pgConfig.CryptoPass
	// Check if the user already exists
	var (
		res model.User
		err error
	)
	user, err := tx.User.WithContext(ctx).Where(
		s.Q.User.Identifier.Eq(identifier),
	).First()
	// .(`SELECT id FROM users WHERE identifier = ?`, user.Identifier).Scan(&existingUserID).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("FAIL TO CHECK USER EXISTENCE", err)
		return res, err
	}
	if user != nil {
		return *user, err
	}
	userID,err := tx.User.InsertParty(proto.PartyType_user.String())
	if err != nil {
		return res,err
	}

	err = tx.User.UnderlyingDB().Raw(`
			INSERT INTO users (id,identifier, password_hash)
			VALUES (?,?, pgp_sym_encrypt(?, ?))
			RETURNING id
		`, userID, identifier, userPassword, pass).Scan(&res.ID).Error
	if err != nil {
		fmt.Println("FAIL TO CREATE USER", err)
		return res, err
	}

	res.Identifier = identifier
	return res, err
}
