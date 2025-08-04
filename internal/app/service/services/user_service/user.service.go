package userservice

import (
	"context"
	"erp/api/common"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/service/helpers"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserService struct {
	conn          *connection.Connection
	timeout       time.Duration
	configService *config.ConfigService
	generator     helpers.Generator
}

func NewUserService(conn *connection.Connection,
	timeout time.Duration,
	configService *config.ConfigService, helpers *helpers.Helpers) *UserService {
	return &UserService{
		conn:          conn,
		configService: configService,
		generator:     helpers.Generator,
		timeout:       timeout,
	}
}

// func (s *UserService) GetUserCompanies(req *common.RequestContext,i *dto.RequestPaginationData)(
// 	dto.PaginationResult[[]entity.UserRelation],error){
// 	ctx,cancel := context.WithTimeout(req.Ctx,s.timeout)
// 	defer cancel()
// 	var result dto.PaginationResult[[]entity.UserRelation]
// 	queryBuilder := s.conn.Db.WithContext(ctx).Model(&entity.Item{}).
// 		Where(&entity.Item{CompanyID: req.ActiveCompany.ID})

// 	err := queryBuilder.
// 			Count(&result.Total).Error

// 	if i.Query != "" {
// 		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+i.Query+"%")
// 	}

// 	err = queryBuilder.
// 		Scopes(s.conn.Paginate(req.Params)).
// 		Find(&result.Results).Error
// 	if err != nil {
// 		return result, err
// 	}

// 	return result, err
// }

func (s *UserService) GetUserPassword(id int64) (string, error) {
	dbConfig := s.configService.GetDbConfig()
	var password string
	err := s.conn.Db.Raw("select pgp_sym_decrypt(password_hash::bytea, ?) as password_hash from users where id = ?", dbConfig.CryptoPass, id).
		Scan(&password).Error
	return password, err
}

// func (s *UserService) InsertUser(ctx context.Context, tx *gorm.DB, user *entity.User) (err error) {
// 	userPassword := s.generator.GeneratePassword()
// 	pass := s.configService.GetDbConfig().CryptoPass
// 	// Check if the user already exists
// 	var existingUserID uint
// 	err = tx.WithContext(ctx).Raw(`SELECT id FROM users WHERE identifier = ?`, user.Identifier).Scan(&existingUserID).Error

// 	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
// 		fmt.Println("FAIL TO CHECK USER EXISTENCE", err)
// 		return
// 	}

// 	// If the user does not exist, insert the new user
// 	if existingUserID == 0 {
// 		err = tx.WithContext(ctx).Raw(`
// 			INSERT INTO users (identifier, password_hash)
// 			VALUES (?, pgp_sym_encrypt(?, ?))
// 			RETURNING id
// 		`, user.Identifier, userPassword, pass).Scan(&user.ID).Error

// 		if err != nil {
// 			fmt.Println("FAIL TO CREATE USER", err)
// 			return
// 		}
// 	} else {
// 		// If the user exists, populate the user ID
// 		user.ID = existingUserID
// 	}

// 	return nil
// }

func (s *UserService) InsertUser(ctx context.Context, tx *query.QueryTx, identifier string) (model.User, error) {
	userPassword := s.generator.GeneratePassword()
	pass := s.configService.GetDbConfig().CryptoPass
	// Check if the user already exists
	var (
		res model.User
		err error
	)
	user, err := tx.User.WithContext(ctx).Where(
		s.conn.Q.User.Identifier.Eq(identifier),
	).First()
	// .(`SELECT id FROM users WHERE identifier = ?`, user.Identifier).Scan(&existingUserID).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Println("FAIL TO CHECK USER EXISTENCE", err)
		return res, err
	}
	if user != nil {
		return *user, err
	}
	err = s.conn.Db.Raw(`
			INSERT INTO users (identifier, password_hash)
			VALUES (?, pgp_sym_encrypt(?, ?))
			RETURNING id
		`, res.Identifier, userPassword, pass).Scan(&res.ID).Error
	if err != nil {
		fmt.Println("FAIL TO CREATE USER", err)
		return res,err
	}
	
	return res, err
}

func (s *UserService) InsertUserCompany(req *common.RequestContext,ctx context.Context, tx *gorm.DB, companyId int64) (err error) {
	err = tx.WithContext(ctx).Exec(`insert into user_relations(company_id,user_id,role_id,profile_id) values($1,$2,$3,$4)
	ON CONFLICT (company_id, user_id,role_id,profile_id) DO NOTHING`, companyId, req.User.ID,req.Role.ID,req.Profile.ID).Error
	return
}

func (s *UserService) InsertUserRole(ctx context.Context, tx *gorm.DB, userId int64, rol int64) (err error) {
	err = tx.WithContext(ctx).Exec(`insert into user_roles(role_id,user_id) values($1,$2)
	ON CONFLICT (role_id, user_id) DO NOTHING`, rol, userId).Error
	return
}

func (s *UserService) GetUser() string {
	return "sdas"
}
