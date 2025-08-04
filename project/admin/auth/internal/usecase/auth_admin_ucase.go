package auth_admin_ucase

import (
	"context"
	"erp/api/dto"
	"erp/pkg/logger"
	auth_admin_repo "erp/project/admin/auth/internal/repository"
)

type AuthAdminUseCase interface {
	SignIn(ctx context.Context, d *dto.SignInRequest) (res dto.SignInResponse, err error)
}

type authAdminUseCase struct {
	emitLog logger.EmitLog
	authAdminRepo auth_admin_repo.AuthAdminRepository
}

func NewAdminAuthUseCase(
	logger logger.Logger,
	authAdminRepo auth_admin_repo.AuthAdminRepository,
) AuthAdminUseCase {
	return &authAdminUseCase{
		emitLog: logger.EmitLog("admin-auth-usecase"),
		authAdminRepo: authAdminRepo,
	}
}

func (u *authAdminUseCase) SignIn(ctx context.Context, d *dto.SignInRequest) (res dto.SignInResponse, err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("SignIn"))
		}
	}()
	res,err = u.authAdminRepo.SignIn(ctx,d)
	return
}