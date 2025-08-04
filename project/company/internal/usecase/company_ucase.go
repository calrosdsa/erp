package company_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/logger"
	company_repo "erp/project/company/internal/repository"
	"fmt"
)

type CompanyUseCase interface {
	CreateCompany(req *common.RequestContext, d *dto.CreateCompanyRequest) (err error)
	GetAllUserCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.CompanyDto], err error)
	// GetCompanyByUuid(ctx context.Context, uuid string) (res model.Company, err error)
	GetCompanyDetail(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.CompanyDto], err error)
	GetAccountSetting(req *common.RequestContext, i *dto.RequestData) (
		res dto.AccountSettingsDto, err error)
	EditAccountSetting(req *common.RequestContext,d dto.AccountSettingData) (err error)
}

type companyUseCase struct {
	emitLog     logger.EmitLog
	companyRepo company_repo.CompanyRepository
	permission  repository.PermissionService
	core        repository.CoreService
	bus bus.Bus
	c           di.Container
}

func NewCompanyUseCase(
	logger logger.Logger,
	companyRepo company_repo.CompanyRepository,
	permission repository.PermissionService,
	core repository.CoreService,
	bus bus.Bus,
	c di.Container,
) CompanyUseCase {
	return &companyUseCase{
		emitLog:     logger.EmitLog("company-usecase"),
		companyRepo: companyRepo,
		permission:  permission,
		core:        core,
		c:           c,
		bus: bus,
	}
}
func (u *companyUseCase)EditAccountSetting(req *common.RequestContext,d dto.AccountSettingData) (err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("EditAccountSetting"))
		}
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.COMPANY, domain.EDIT); err != nil {
		return  err
	}
	err = u.companyRepo.EditAccountSetting(req, d)
	return
}

func (u *companyUseCase) GetAccountSetting(req *common.RequestContext, i *dto.RequestData) (
	res dto.AccountSettingsDto, err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetAccountSetting"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.COMPANY, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.companyRepo.GetAccountSetting(req, i)
	return
}

func (u *companyUseCase) GetAllUserCompanies(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CompanyDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAllUserCompanies"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.COMPANY, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.companyRepo.GetAllUserCompanies(req, d)
	return
}
// func (u *companyUseCase) GetCompanyByUuid(ctx context.Context, uuid string) (res model.Company, err error) {
// 	defer func() {
// 		if err != nil {
// 			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompanyByUuid"))
// 		}
// 	}()
// 	res, err = u.companyRepo.GetCompanyByUuid(ctx, uuid)
// 	return
// }
func (u *companyUseCase) GetCompanyDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.CompanyDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompanyDetail"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.COMPANY, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.companyRepo.GetCompanyDetail(req,i)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}

func (u *companyUseCase) CreateCompany(req *common.RequestContext, d *dto.CreateCompanyRequest) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateInvoiceState"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	company,err := u.companyRepo.CreateCompany(tx,req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx,domain.EventCompanyCreated,event.CreatedCompanyEventData{
		Tx: tx,
		LanguageCode: string(req.LanguageCode),
		Company:company,
	})
	return
}

func (u *companyUseCase) closeTx(tx *query.QueryTx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}
