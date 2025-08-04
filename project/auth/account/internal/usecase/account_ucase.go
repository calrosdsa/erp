package account_ucase

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/logger"
	account_repo "erp/project/auth/account/internal/repository"
)

type AccountUseCase interface {
	ResetPassword(ctx context.Context, i *dto.ResetPasswordRequest) (error)
	ChangePassword(ctx context.Context, i *dto.ChangePasswordRequest,c *common.Claims) (error)
	GetUserAccount(ctx context.Context, claims *common.Claims) (model.User, error)
	SignIn(ctx context.Context, d *dto.SignInRequest, ipAddress, userAgent string) (
		res dto.SignInResponse, err error)
	UpdatePassword(req *common.RequestContext, i *dto.UpdatePasswordRequest) (
		res *string, err error)
}

type accountUseCase struct {
	emiLog      logger.EmitLog
	accountRepo account_repo.AccountRepository
	bus bus.Bus
	session repository.SessionService
}

func NewAccountUseCase(
	logger logger.Logger,
	accountRepo account_repo.AccountRepository,
	bus bus.Bus,
	session repository.SessionService,
) AccountUseCase {
	return &accountUseCase{
		emiLog:      logger.EmitLog("account-ucase"),
		accountRepo: accountRepo,
		session: session,
		bus: bus,
	}
}

func (u *accountUseCase) ChangePassword(ctx context.Context, i *dto.ChangePasswordRequest,c *common.Claims) (err error) {
	defer func(){
		if err != nil {
			u.emiLog.Err(err,logger.OptionsLog.WithMethod("ChangePassword"))
		}
	}()
	err = u.accountRepo.ChangePassword(ctx,i,c)
	return
}

func (u *accountUseCase) ResetPassword(ctx context.Context, i *dto.ResetPasswordRequest) (err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("ResetPassword"))
		}
	}()
	user,err := u.accountRepo.ResetPassword(ctx, i)
	if err != nil {
		return
	}
	userRelation,err := u.session.GetUserRelationByUserID(ctx,user.ID)
	if err != nil {
		return
	}
	err = u.bus.Emit(ctx,domain.PasswordResetEvent,event.PasswordResetEventData{
		LanguageCode: i.AcceptLanguage,
		User:user,
		Profile: userRelation.Profile,
		Company: userRelation.Company,
	})
	return
}

func (u *accountUseCase) GetUserAccount(ctx context.Context, claims *common.Claims) (
	res model.User, err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("ResetPassword"))
		}
	}()
	res, err = u.accountRepo.GetUserAccount(ctx, claims)
	return
}

func (u *accountUseCase) SignIn(ctx context.Context, d *dto.SignInRequest, ipAddress, userAgent string) (
	res dto.SignInResponse, err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("SignIn"))
		}
	}()
	res, err = u.accountRepo.SignIn(ctx, d, ipAddress, userAgent)
	return
}

func (u *accountUseCase) UpdatePassword(req *common.RequestContext, i *dto.UpdatePasswordRequest) (
	res *string, err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("UpdatePassword"))
		}
	}()
	res, err = u.accountRepo.UpdatePassword(req, i)
	return
}
