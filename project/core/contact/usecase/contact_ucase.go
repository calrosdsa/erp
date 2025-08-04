package contact_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	contact_repo "erp/project/core/contact/repository"
)

type ContactUseCase interface {
	CreateContact(req *common.RequestContext, i dto.ContactData) (dto.ContactDto, error)
	GetContacts(req *common.RequestContext, i dto.ContactsRequest) (
		[]dto.ContactDto, error)
	GetContact(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.ContactDto], error)
	EditContact(req *common.RequestContext, d dto.ContactData) (err error)
	ContactBulk(req *common.RequestContext,d dto.ContactBulkData) (err error)
}

type partyContactUseCase struct {
	contactRepo contact_repo.ContactRepository
	emitLog     helpers.EmitLog
	permission repository.PermissionService
	core repository.CoreService
}

func NewContactUseCase(
	helpers *helpers.Helpers,
	contactRepo contact_repo.ContactRepository,
	permission repository.PermissionService,
	core repository.CoreService,
) ContactUseCase {
	return &partyContactUseCase{
		emitLog:     helpers.Logger.EmitLog("contact-usecase"),
		contactRepo: contactRepo,
		permission: permission,
		core: core,
	}
}

func(u *partyContactUseCase) ContactBulk(req *common.RequestContext,d dto.ContactBulkData) (err error) {
	defer func (){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("ContactBulk"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.CONTACT,domain.CREATE);err != nil {
		return
	}
	if err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.CONTACT,domain.EDIT);err != nil {
		return
	}
	
	err = u.contactRepo.ContactBulk(req,d)
	return
}

func(s *partyContactUseCase)EditContact(req *common.RequestContext, d dto.ContactData) (err error){
	defer func (){
		if err != nil {
			s.emitLog.Err(err,logger.OptionsLog.WithMethod("EditContact"))
		}
	}()
	err = s.permission.CheckPermissionEntity(req.Ctx,req,domain.CONTACT,domain.EDIT)
	if err != nil {
		return
	}
	err = s.contactRepo.EditContact(req,d)
	return
}

func (s *partyContactUseCase) CreateContact(req *common.RequestContext, i dto.ContactData) (
	res dto.ContactDto, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateContact"))
		}
	}()
	err = s.permission.CheckPermissionEntity(req.Ctx,req,domain.CONTACT,domain.CREATE)
	if err != nil {
		return
	}
	contact, err := s.contactRepo.CreateContact(req, i)
	if err != nil {
		return
	}
	res = dto.ContactDtoFromModel(&contact)
	return res, err
}
func (s *partyContactUseCase) GetContacts(req *common.RequestContext, i dto.ContactsRequest) (
	res []dto.ContactDto,err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetContacts"))
		}
	}()
	err = s.permission.CheckPermissionEntity(req.Ctx,req,domain.CONTACT,domain.VIEW)
	if err != nil {
		return
	}
	res, err = s.contactRepo.GetContacts(req, i)
	return res, err
}
func (s *partyContactUseCase) GetContact(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.ContactDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetContact"))
		}
	}()
	err = s.permission.CheckPermissionEntity(req.Ctx,req,domain.CONTACT,domain.VIEW)
	if err != nil {
		return
	}
	res, err = s.contactRepo.GetContact(req, i)
	if err != nil {
		return
	}
	res.Activities = s.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return res, err
}
