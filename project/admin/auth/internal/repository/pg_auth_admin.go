package auth_admin_repo

import (
	"context"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/pkg/config"
	"erp/pkg/db"
	"fmt"
)

type AuthAdminRepository interface {
	SignIn(ctx context.Context,d *dto.SignInRequest)(res dto.SignInResponse, err error)
}

type authAdminRepository struct {
	Q *query.Query
	appConfing *config.AppConfig
}

func NewAdminAuthRepository(
	db db.Connection,
	appConfig *config.AppConfig,
)AuthAdminRepository {
	return &authAdminRepository{
		Q: db.GetQ(),
		appConfing: appConfig,
	}
}


func (r *authAdminRepository)SignIn(ctx context.Context,d *dto.SignInRequest)(res dto.SignInResponse, err error){
   cryptoPass := r.appConfing.PG.CryptoPass
	var user model.User
	tx := r.Q.User.UnderlyingDB().WithContext(ctx).First(&user, 
		"identifier = $1 and pgp_sym_decrypt(password_hash::bytea, $2) = $3",
		d.Body.Email, cryptoPass, d.Body.Password)
	if tx.Error != nil {
		fmt.Println("ERROR", tx.Error)
		err = tx.Error
		return
	}
	res.Body.U = user
	return
}
