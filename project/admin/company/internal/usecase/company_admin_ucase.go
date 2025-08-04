package a_company_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/logger"
	a_company_repo "erp/project/admin/company/internal/repository"
	"fmt"
)

type AdminCompanyUCase interface {
	GetParentCompanies(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.CompanyDto], error)
	GetCompany(req *common.AdminRequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.CompanyDto], error)
	GetCompanyModules(req *common.AdminRequestContext, d *dto.RequestData) (
		res []dto.CompanyEntityDto, err error)
	AddCompanyModules(req *common.AdminRequestContext, d *dto.AddCompanyModules) error
	GetCompanyUsers(req *common.AdminRequestContext, d *dto.RequestData) ([]dto.UserDto, error)
	AddCompanyUser(req *common.AdminRequestContext, d *dto.CreateUserAdminRequest) (error)
	CreateCompany(req *common.AdminRequestContext, d *dto.CreateCompanyAdminRequest) (err error)
}

type adminCompanyUCase struct {
	companyAdminRepo a_company_repo.AdminCompanyRepository
	emitLog          logger.EmitLog
	bus bus.Bus
	c di.Container
}

func NewAdminCompanyUCase(
	logger logger.Logger,
	companyAdminRepo a_company_repo.AdminCompanyRepository,
	bus bus.Bus,
	c di.Container,
) AdminCompanyUCase {
	return &adminCompanyUCase{
		emitLog:          logger.EmitLog("company-admin-ucase"),
		companyAdminRepo: companyAdminRepo,
		bus: bus,
		c: c,
	}
}


func (u *adminCompanyUCase) AddCompanyUser(req *common.AdminRequestContext, d *dto.CreateUserAdminRequest) (err error){
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("AddCompanyUser"))
		}
	}()
	userRelation,err := u.companyAdminRepo.AddCompanyUser(req,d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx,domain.UserCreatedEvent,event.UserCreatedEventData{
		UseRelation: userRelation,
		LanguageCode: req.LanguageCode,
	})
	return
}

func (u *adminCompanyUCase) GetCompanyUsers(req *common.AdminRequestContext, d *dto.RequestData) (
	res []dto.UserDto,err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetCompanyUsers"))
		}
	}()
	res,err = u.companyAdminRepo.GetCompanyUsers(req,d)
	return 
}

func (u *adminCompanyUCase) AddCompanyModules(req *common.AdminRequestContext, d *dto.AddCompanyModules) (err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("AddCompanyModules"))
		}
	}()
	err = u.companyAdminRepo.AddCompanyModules(req,d)
	return 
}

func (u *adminCompanyUCase) GetCompanyModules(req *common.AdminRequestContext, d *dto.RequestData) (
	res []dto.CompanyEntityDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompanyModules"))
		}
	}()
	res, err = u.companyAdminRepo.GetCompanyModules(req, d)
	return
}

func (u *adminCompanyUCase) GetCompany(req *common.AdminRequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.CompanyDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompany"))
		}
	}()
	res, err = u.companyAdminRepo.GetCompany(req, d)
	return
}

func (u *adminCompanyUCase) GetParentCompanies(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CompanyDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetParentCompanies"))
		}
	}()
	res, err = u.companyAdminRepo.GetParentCompanies(req, d)
	return
}


func (u *adminCompanyUCase)CreateCompany(req *common.AdminRequestContext, d *dto.CreateCompanyAdminRequest) (err error){
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateCompany"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}()
	company,companyDefaults,err := u.companyAdminRepo.CreateCompany(tx,req, d)
	if err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx,domain.EventCompanyCreated,event.CreatedCompanyEventData{
		Tx: tx,
		LanguageCode: string(req.LanguageCode),
		Company:company,
		CompanyDefaults: companyDefaults,
		Body: d.Body,
		IsRoot: true,
	})
	return
}

func (u *adminCompanyUCase) closeTx(tx *query.QueryTx, err error) error {
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