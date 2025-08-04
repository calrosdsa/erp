package account_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/config"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gorm"
)

type AccountRepository interface {
	ResetPassword(ctx context.Context, i *dto.ResetPasswordRequest) (model.User,error)
	GetUserAccount(ctx context.Context, claims *common.Claims) (model.User, error)
	SignIn(ctx context.Context, d *dto.SignInRequest, ipAddress, userAgent string) (
		res dto.SignInResponse, err error)
	UpdatePassword(req *common.RequestContext, i *dto.UpdatePasswordRequest) (
		res *string, err error)
	ChangePassword(ctx context.Context, i *dto.ChangePasswordRequest,c *common.Claims) (
		err error)
}

type accountRepository struct {
	Q  *query.Query
	DB *gorm.DB
	// conn          *connection.Connection
	appConfing *config.AppConfig
	locale   helpers.Locale
	sessionRepo SessionRepository
}

func NewAccountRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
	appConfing *config.AppConfig,
) AccountRepository {
	return &accountRepository{
		// conn:          conn,
		Q:        conn.GetQ(),
		DB:       conn.GetDB(),
		locale:   helpers.Locale,
		appConfing: appConfing,
		sessionRepo: NewSessionRepository(conn,helpers),
	}
}

func (s *accountRepository) ResetPassword(ctx context.Context, i *dto.ResetPasswordRequest) (
	res model.User,err error) {
	user,err := s.Q.User.WithContext(ctx).Where(
		s.Q.User.Identifier.Eq(i.Body.Email),
	).First()
	if err != nil {
		return 
	}
	return *user,err
}


func (s *accountRepository) GetUserAccount(ctx context.Context, claims *common.Claims) (model.User, error) {
	var (
		err error
	)
	user, err := s.Q.User.WithContext(ctx).Where(
		s.Q.User.UUID.Eq(claims.Uuid),
	).First()
	return *user, err
}

func (s *accountRepository) SignIn(ctx context.Context, d *dto.SignInRequest, ipAddress, userAgent string) (
	res dto.SignInResponse, err error,
) {
	cryptoPass := s.appConfing.PG.CryptoPass
	var user model.User
	tx := s.DB.WithContext(ctx).First(&user, "identifier = $1 and pgp_sym_decrypt(password_hash::bytea, $2) = $3",
		d.Body.Email, cryptoPass, d.Body.Password)
	if tx.Error != nil {
		fmt.Println("ERROR", tx.Error)
		err = tx.Error
		return
	}
	res.Body.U = user
	res.Body.User = dto.UserDTOFromModel(&user)
	userRelation, err := s.Q.UserRelation.Where(
		s.Q.UserRelation.UserID.Eq(user.ID),
	).First()
	if err != nil {
		return
	}
	userRelationDto := dto.UserRelationDtoFromModel(userRelation)
	res.Body.UserRelation = userRelationDto

	// Create session for the user
	activeCompanyID := userRelation.CompanyID
	session, err := s.sessionRepo.CreateSession(ctx, user.ID,user.UUID, ipAddress, userAgent, &activeCompanyID)
	if err != nil {
		fmt.Println("ERROR creating session", err)
		// Don't fail the login if session creation fails, just log the error
		err = nil
	} else {
		// Add session token to response (could be used for session management)
		fmt.Printf("Session created: ID=%d, Token=%s\n", session.ID, session.Token)
	}
	res.Body.AccessToken = session.Token

	return
}

func (s *accountRepository) ChangePassword(ctx context.Context, i *dto.ChangePasswordRequest,
	c *common.Claims) (err error){
	cryptoPass := s.appConfing.PG.CryptoPass
	err = s.DB.WithContext(ctx).Exec("update users set password_hash = pgp_sym_encrypt($1,$3) where id = $2",
	i.Body.Password, c.ID, cryptoPass).Error
	return err
}

func (s *accountRepository) UpdatePassword(req *common.RequestContext, i *dto.UpdatePasswordRequest) (
	*string, error) {
	var user model.User
	var responseMsg string
	cryptoPass := s.appConfing.PG.CryptoPass
	err := s.DB.WithContext(req.Ctx).First(&user, "identifier = $1 and pgp_sym_decrypt(password_hash::bytea, $2) = $3",
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
	err = s.DB.Exec("update users set password_hash = pgp_sym_encrypt($1,$3) where id = $2",
		i.Body.NewPassword, user.ID, cryptoPass).Error
	return nil, err
}

// func (s *accountRepository) getUserKeyCache(id int64) string {
// 	return fmt.Sprintf("user-%d", id)
// }
// func (s *accountRepository) getUserFromCache(ctx context.Context, user *model.User) (err error) {
// 	key := s.getUserKeyCache(user.ID)
// 	err = s.cache.Get(ctx, key, &user)
// 	if err != nil {
// 		fmt.Println("USER NOT PRESENT ON CACHE")
// 	}
// 	return
// }

// func (s *accountRepository) GetAccountCompanies(ctx context.Context,)
