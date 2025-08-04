package account_service

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type AccountService struct {
	timeout       *time.Duration
	conn          *connection.Connection
	configService *config.ConfigService
	locale        helpers.Locale
	cache         *helpers.CacheHelper
	emitLog       helpers.EmitLog
}

func NewAccountService(timeout *time.Duration, conn *connection.Connection, configService *config.ConfigService,
	helpers *helpers.Helpers,
) *AccountService {
	return &AccountService{
		timeout:       timeout,
		conn:          conn,
		configService: configService,
		locale:        helpers.Locale,
		cache:         helpers.Cache,
		emitLog:       helpers.Logger.EmitLog("account-service"),
	}
}

func (s *AccountService) ResetPassword(ctx context.Context, i *dto.ResetPasswordRequest) error {
	return nil
}

func (s *AccountService) GetUserAccount(ctx context.Context, claims *common.Claims) (model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, *s.timeout)
	defer cancel()
	var (
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserAccount"))

		}
	}()
	// err = s.getUserFromCache(ctx, &user)
	// if err != nil {
	// userM,err := s.conn.Q.User.Where(
	// 	s.conn.Q.User.UUID.Eq(claims.Uuid),
	// ).First()

	// userRelation, err := s.conn.Q.UserRelation.Where(
	// 	s.conn.Q.UserRelation.UUID.Eq(claims.Uuid),
	// ).Preload(s.conn.Q.UserRelation.Company).
	// 	Preload(s.conn.Q.UserRelation.User).
	// 	Preload(s.conn.Q.UserRelation.Profile).
	// 	Preload(s.conn.Q.UserRelation.Role).First()

	// fmt.Println("USERRELATION", userRelation)
	user,err := s.conn.Q.User.WithContext(ctx).Where(
		s.conn.Q.User.UUID.Eq(claims.Uuid),
	).First()

	// err = s.conn.Db.WithContext(ctx).Preload("UserRelation.Company").
	// 	Preload("UserRelation.Role").Preload("UserRelation.Role.RoleActions").
	// 	Preload("UserRelation.Role.RoleActions.Action").
	// 	Preload("UserRelation.Profile").First(&user, "uuid = $1", claims.Uuid).Error
	// if err != nil {
	// 	return
	// }
	// var party entity.Party
	// err = s.conn.Db.WithContext(ctx).Preload("PartyType").Where(&entity.Party{Base: entity.Base{
	// 	ID: user.UserRelation.Profile.ID,
	// }}).First(&party).Error
	// if err != nil {
	// 	return
	// }
	// user.UserRelation.Profile.Party = party

	// err = s.cache.Set(ctx, s.getUserKeyCache(uint(claims.ID)), user)
	// if err != nil {
	// 	return
	// }
	// } 
	return *user,err
}

func (s *AccountService) SignIn(ctx context.Context, d *dto.SignInRequest) (
	res dto.SignInResponse, err error,
) {
	ctx, cancel := context.WithTimeout(ctx, *s.timeout)
	defer cancel()
	cryptoPass := s.configService.GetDbConfig().CryptoPass
	var user model.User
	tx := s.conn.Db.WithContext(ctx).First(&user, "identifier = $1 and pgp_sym_decrypt(password_hash::bytea, $2) = $3",
		d.Body.Email, cryptoPass, d.Body.Password)
	if tx.Error != nil {
		fmt.Println("ERROR", tx.Error)
		err = tx.Error
		return
	}
	res.Body.U = user
	res.Body.User = dto.UserDTOFromModel(&user)
	userRelation,err := s.conn.Q.UserRelation.Where(
		s.conn.Q.UserRelation.UserID.Eq(user.ID),
	).First()
	if err != nil {
		return 
	}
	userRelationDto := dto.UserRelationDtoFromModel(userRelation)
	res.Body.UserRelation = userRelationDto
	return
}

func (s *AccountService) UpdatePassword(req *common.RequestContext, i *dto.UpdatePasswordRequest) (
	*string, error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var user model.User
	var responseMsg string
	cryptoPass := s.configService.GetDbConfig().CryptoPass
	err := s.conn.Db.WithContext(ctx).First(&user, "identifier = $1 and pgp_sym_decrypt(password_hash::bytea, $2) = $3",
		req.User.Identifier, cryptoPass, i.Body.Password).Error
	if err == gorm.ErrRecordNotFound {
		responseMsg = s.locale.MustLocalize(
			helpers.OptionsLocale.WithID("Error.IncorrectPassword"),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)))
		return &responseMsg, err
	}
	if err != nil {
		return nil, err
	}
	err = s.conn.Db.Exec("update users set password_hash = pgp_sym_encrypt($1,$3) where id = $2",
		i.Body.NewPassword, user.ID, cryptoPass).Error
	return nil, err
}

func (s *AccountService) getUserKeyCache(id int64) string {
	return fmt.Sprintf("user-%d", id)
}
func (s *AccountService) getUserFromCache(ctx context.Context, user *model.User) (err error) {
	key := s.getUserKeyCache(user.ID)
	err = s.cache.Get(ctx, key, &user)
	if err != nil {
		fmt.Println("USER NOT PRESENT ON CACHE")
	}
	return
}

// func (s *AccountService) GetAccountCompanies(ctx context.Context,)
