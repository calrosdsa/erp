package user_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/config"
	"erp/pkg/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(tx *query.QueryTx,req *common.RequestContext, i *dto.CreateUserRequest) (model.UserRelation, error)
}

type userRepository struct {
	conn      db.Connection
	Q         *query.Query
	generator helpers.Generator
	pgConfig  config.PGConfig
}

func NewUserRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
	pgConfig config.PGConfig,
) UserRepository {
	return &userRepository{
		conn:      conn,
		Q:         conn.GetQ(),
		generator: helpers.Generator,
		pgConfig:  pgConfig,
	}
}

func (r *userRepository) CreateUser(tx *query.QueryTx,req *common.RequestContext, i *dto.CreateUserRequest) (
	userRelation model.UserRelation, err error) {
	var (
		user model.User
	)
	user, err = r.InsertUser(req.Ctx, tx, i.Body.Email)
	if err != nil {
		return
	}
	role, err := r.Q.Role.WithContext(req.Ctx).Where(
		r.Q.Role.UUID.Eq(i.Body.RoleUUID),
	).First()
	if err != nil {
		return
	}
	userRelation.RoleID = role.ID
	userRelation.UserID = int64(user.ID)
	userRelation.User = user
	profile, err := r.createProfile(tx, i,req.ActiveCompany.ID)
	if err != nil {
		return
	}
	fmt.Println("PROFILE",profile.GivenName)
	userRelation.ProfileID = int64(profile.ID)
	userRelation.Profile = profile
	// if req.Profile.Party.PartyTypeCode == domain.PARTY_ADMIN {
	// for _, companyId := range i.Body.CompanyIds {
	// 	fmt.Println("COMPANY ID",companyId)
	// 	// if allow := r.permissionService.CheckIfUserIsCompanyAdmin(ctx, req, int64(companyId), domain.CREATE_USER); !allow {
	// 	// 	fmt.Println("NO ALLOW CUSTOM PERMISSION")
	// 	// 	return domain.ACTION_NOT_ALLOWED
	// 	// }
	// 	userRelation.CompanyID = int64(companyId)
	// 	err = r.createUserRelation(req.Ctx,tx, &userRelation)
	// 	if err != nil && err != domain.USER_ALREDY_EXIST {
	// 		return
	// 	}
	// }
	// } else {
	userRelation.CompanyID = req.ActiveCompany.ID
	userRelation.Company = req.ActiveCompany
	// if allow := r.permissionService.CheckIfUserIsCompanyAdmin(req.Ctx, req, req.ActiveCompany.ID, domain.CREATE_USER); !allow {
	// 	return domain.ACTION_NOT_ALLOWED
	// }
	err = r.createUserRelation(req.Ctx, tx, &userRelation)
	if err != nil && err != domain.USER_ALREDY_EXIST {
		return
	}
	// }
	// r.eventHelper.Publish(event.NOTIFICATION_EVENT, event.NotificationData{
	// 	NotificationEventType: event.NOTIFY_NEW_CLIENT,
	// 	Data: event.NotificationPayload{
	// 		Payload:        userRelation,
	// 		RequestContext: *req,
	// 	},
	// },
	// )
	return
}

func (s *userRepository) createUserRelation(ctx context.Context, tx *query.QueryTx, userRelation *model.UserRelation) error {
	_, err := tx.UserRelation.WithContext(ctx).Where(
		s.Q.UserRelation.RoleID.Eq(userRelation.RoleID),
		s.Q.UserRelation.CompanyID.Eq(userRelation.CompanyID),
		s.Q.UserRelation.ProfileID.Eq(userRelation.ProfileID),
		s.Q.UserRelation.UserID.Eq(userRelation.UserID),
	).First()
	if err == nil {
		return domain.USER_ALREDY_EXIST
	}

	err = tx.UserRelation.WithContext(ctx).Create(&model.UserRelation{
		UserID:    userRelation.UserID,
		ProfileID: userRelation.ProfileID,
		CompanyID: userRelation.CompanyID,
		RoleID:    userRelation.RoleID,
	})
	if err != nil {
		return err
	}

	return err
}

func (s *userRepository) createProfile(tx *query.QueryTx, i *dto.CreateUserRequest,companyID int64) (model.Profile, error) {
	var (
		profile model.Profile
	)

	partyID, err := tx.Profile.InsertParty(i.Body.PartyCode)
	if err != nil {
		return profile, err
	}
	profile.GivenName = i.Body.GivenName
	profile.FamilyName = i.Body.FamilyName
	profile.EmailAddress = i.Body.Email
	profile.ID = partyID
	profile.CompanyID = companyID
	profile.PhoneNumber = &i.Body.PhoneNumber
	err = tx.Profile.Save(&profile)
	if err != nil {
		return profile, err
	}
	keyValues := make([]*model.KeyValue, len(i.Body.KeyValueData))
	if len(keyValues) == 0 {
		return profile, nil
	}
	for i, keyValue := range i.Body.KeyValueData {
		kV := &model.KeyValue{}
		kV.Key = keyValue.Key
		kV.Value = keyValue.Value
		kV.PartyID = profile.ID
		keyValues[i] = kV
	}
	err = tx.KeyValue.CreateInBatches(keyValues, len(keyValues))
	if err != nil {
		return profile, err
	}
	return profile, nil
}

func (s *userRepository) InsertUser(ctx context.Context, tx *query.QueryTx, identifier string) (model.User, error) {
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
