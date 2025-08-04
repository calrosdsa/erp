package documentinfo_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	documentinfo_repo "erp/project/document/document_info/repository"
	"fmt"
)

type DocumentInfoUseCase interface {
	EditAddressAndContact(req *common.RequestContext, d dto.AddressAndContactData) (err error)
	GetAddressAndContact(req *common.RequestContext,d dto.RequestEntityWithParty)(res dto.AddressAndContactDto,err error)
	EditDocTerms(req *common.RequestContext, d dto.DocTermsData) (err error)
	GetDocTerms(req *common.RequestContext, d dto.RequestEntityWithParty) (res dto.DocTermsDto, err error)
	EditDocAccounts(req *common.RequestContext, d dto.DocAccountingData) (err error)
	GetDocAccounts(req *common.RequestContext, d dto.RequestEntityWithParty) (res dto.DocAccountingDto, err error)
}

type documentInfoUseCase struct {
	emitLog    logger.EmitLog
	repo       documentinfo_repo.DocumentInfoRepository
	permission repository.PermissionService
}

func NewDocumentInfoUseCase(
	logger logger.Logger,
	repo documentinfo_repo.DocumentInfoRepository,
	permission repository.PermissionService,
) DocumentInfoUseCase {
	return &documentInfoUseCase{
		emitLog:    logger.EmitLog("docuemnt-info-usecase"),
		repo:       repo,
		permission: permission,
	}
}


func(u *documentInfoUseCase) GetDocAccounts(req *common.RequestContext,d dto.RequestEntityWithParty)(res dto.DocAccountingDto,err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetDocAccounts"))
		}
	}()
	entity, err := u.permission.GetDocumentEntity(d.PartyType)
	if err != nil {
		return
	}
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, entity, domain.VIEW); err != nil {
		return
	}
	res,err = u.repo.GetDocAccounts(req, d)

	return
}

func (u *documentInfoUseCase) EditDocAccounts(req *common.RequestContext, d dto.DocAccountingData) (
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditDocAccounts"))
		}
	}()
	entity, err := u.permission.GetDocumentEntity(d.DocPartyType)
	if err != nil {
		return
	}
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, entity, domain.EDIT); err != nil {
		return
	}
	err = u.repo.EditDocAccounts(req, d)
	return
}

func(u *documentInfoUseCase) GetDocTerms(req *common.RequestContext,d dto.RequestEntityWithParty)(res dto.DocTermsDto,err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetDocTerms"))
		}
	}()
	entity, err := u.permission.GetDocumentEntity(d.PartyType)
	fmt.Println("DOC TERMS PERM",err)
	if err != nil {
		return
	}
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, entity, domain.VIEW); err != nil {
		return
	}
	res,err = u.repo.GetDocTerms(req, d)

	return
}

func (u *documentInfoUseCase) EditDocTerms(req *common.RequestContext, d dto.DocTermsData) (
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditDocTerms"))
		}
	}()
	entity, err := u.permission.GetDocumentEntity(d.DocPartyType)
	if err != nil {
		return
	}
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, entity, domain.EDIT); err != nil {
		return
	}
	err = u.repo.EditDocTerms(req, d)
	return
}


func(u *documentInfoUseCase) GetAddressAndContact(req *common.RequestContext,d dto.RequestEntityWithParty)(res dto.AddressAndContactDto,err error) {
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetAddressAndContact"))
		}
	}()
	entity, err := u.permission.GetDocumentEntity(d.PartyType)
	if err != nil {
		return
	}
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, entity, domain.VIEW); err != nil {
		return
	}
	res,err = u.repo.GetAddressAndContact(req, d)

	return
}

func (u *documentInfoUseCase) EditAddressAndContact(req *common.RequestContext, d dto.AddressAndContactData) (
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditAddressAndContact"))
		}
	}()
	entity, err := u.permission.GetDocumentEntity(d.PartyType)
	if err != nil {
		return
	}
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, entity, domain.EDIT); err != nil {
		return
	}
	err = u.repo.EditAddressAndContact(req, d)
	return
}
