package address_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	address_repo "erp/project/core/address/repository"
)

type AddressUseCase interface {
	CreateAddress(req *common.RequestContext, d dto.AddressData) (dto.AddressDto, error)
	EditAddress(req *common.RequestContext, d dto.AddressData) error

	GetAddresses(req *common.RequestContext, d dto.RequestAddresses) (dto.ResponseDataList[[]dto.AddressDto], error)
	GetAddress(req *common.RequestContext, d *dto.RequestEntity) (dto.ResultEntity[dto.AddressDto], error)
	UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (
		err error)
	GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto]
}

type addressUseCase struct {
	emitLog    logger.EmitLog
	repo       address_repo.AddressRepository
	permission repository.PermissionService
	core       repository.CoreService
	fsm        fsm.FsmState
}

func NewAddressUseCase(
	logger logger.Logger,
	repo address_repo.AddressRepository,
	permission repository.PermissionService,
	core repository.CoreService,
	fsm fsm.FsmState,
) AddressUseCase {
	return &addressUseCase{
		repo:       repo,
		emitLog:    logger.EmitLog("address-service"),
		permission: permission,
		core:       core,
		fsm:        fsm,
	}
}

func (u *addressUseCase) UpdateStatus(req *common.RequestContext, d dto.UpdateStatusWithEvent) (
	err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}		
	}()
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.ADDRESS, domain.EDIT); err != nil {
		return err
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.repo.UpdateStatus(req, d, nextState)
	
	return
}

func (s *addressUseCase) GetAddresses(req *common.RequestContext, d dto.RequestAddresses) (
	res dto.ResponseDataList[[]dto.AddressDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAddresses"))
		}
	}()
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.ADDRESS, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = s.repo.GetAddresses(req, d)
	return res, err
}
func (s *addressUseCase) GetAddress(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.AddressDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetAddress"))
		}
	}()
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.ADDRESS, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = s.repo.GetAddress(req, d)
	if err != nil {
		return
	}
	res.Activities = s.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return res, err
}

func (s *addressUseCase) EditAddress(req *common.RequestContext, d dto.AddressData) (err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("EditAddress"))
		}
	}()
	if err := s.permission.CheckPermissionEntity(req.Ctx, req, domain.ADDRESS, domain.EDIT); err != nil {
		return err
	}
	err = s.repo.EditAddress(req, d)
	if err != nil {
		return
	}
	return
}

func (s *addressUseCase) CreateAddress(req *common.RequestContext, d dto.AddressData) (res dto.AddressDto, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateAddress"))
		}
	}()
	if err := s.permission.CheckPermissionEntity(req.Ctx, req, domain.ADDRESS, domain.CREATE); err != nil {
		return res, err
	}
	res, err = s.repo.CreateAddress(req, d)
	if err != nil {
		return
	}
	return
}

func (s *addressUseCase) GetAllowedPartiesForAddress(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto] {
	res := s.repo.GetAllowedPartiesForAddress(req)
	return res
}
