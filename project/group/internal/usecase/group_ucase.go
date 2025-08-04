package group_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	group_repo "erp/project/group/internal/repository"
	"fmt"
)

type GroupUseCase interface {
	GetGroupDescendents(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		res dto.ResultEntity[[]dto.GroupHierarchyDto], err error)
	GetEntityGroup(partyCode string) (domain.EntityTemplate, error)
	GetGroup(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		res dto.ResultEntity[dto.GroupDto], err error)
	GetGroups(req *common.RequestContext, d *dto.RequestPaginationPartyData) (
		res dto.PaginationResult[[]dto.GroupDto], err error)
	CreateGroup(req *common.RequestContext, i *dto.CreateGroupRequest) (err error)
	EditGroup(req *common.RequestContext, d *dto.EditGroupRequest) (err error)
}

type groupUseCase struct {
	emitLog    logger.EmitLog
	groupRepo  group_repo.GroupRepository
	permission repository.PermissionService
	core       repository.CoreService
}

func NewGroupUcaseCase(
	helpers *helpers.Helpers,
	logger logger.Logger,
	groupRepo group_repo.GroupRepository,
	permission repository.PermissionService,
	core repository.CoreService,
) GroupUseCase {
	return &groupUseCase{
		emitLog:    logger.EmitLog("group-usecase"),
		permission: permission,
		groupRepo:  groupRepo,
		core:       core,
	}
}
func (u *groupUseCase) GetTreeView(req *common.RequestContext, d *dto.RequestDataWithPartyType) (
	res []dto.TreeEntryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetTreeView"))
		}
	}()
	entityT, err := u.GetEntityGroup(d.PartyType)
	if err != nil {
		return res, err
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entityT, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.groupRepo.GetTreeView(req, d)
	if err != nil {
		return
	}
	return res, err
}

func (u *groupUseCase) EditGroup(req *common.RequestContext, d *dto.EditGroupRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditGroup"))
		}
	}()
	entityT, err := u.GetEntityGroup(d.Body.PartyTypeCode)
	if err != nil {
		return err
	}
	err = u.permission.CheckPermissionEntity(req.Ctx, req, entityT, domain.EDIT)
	if err != nil {
		return
	}
	err = u.groupRepo.EditGroup(req, d)
	return
}

func (u *groupUseCase) GetGroupDescendents(req *common.RequestContext, d *dto.RequestEntityWithParty) (
	res dto.ResultEntity[[]dto.GroupHierarchyDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetGroupDescendents"))
		}
	}()
	entityT, err := u.GetEntityGroup(d.PartyType)
	if err != nil {
		return res, err
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entityT, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.groupRepo.GetGroupDescendents(req, d)
	if err != nil {
		return
	}
	return res, err
}

func (u *groupUseCase) GetGroup(req *common.RequestContext, d *dto.RequestEntityWithParty) (
	res dto.ResultEntity[dto.GroupDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetGroup"))
		}
	}()
	entityT, err := u.GetEntityGroup(d.PartyType)
	if err != nil {
		return res, err
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entityT, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.groupRepo.GetGroup(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return res, err
}

func (u *groupUseCase) GetGroups(req *common.RequestContext, d *dto.RequestPaginationPartyData) (
	res dto.PaginationResult[[]dto.GroupDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetGroups"))
		}
	}()
	entityT, err := u.GetEntityGroup(d.PartyType)
	if err != nil {
		return res, err
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entityT, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.groupRepo.GetGroups(req, d)
	return res, err
}

func (u *groupUseCase) CreateGroup(req *common.RequestContext, i *dto.CreateGroupRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetGroups"))
		}
	}()
	fmt.Println("GROUP PARTY TYPE", i.Body.PartyTypeCode)
	entityT, err := u.GetEntityGroup(i.Body.PartyTypeCode)
	if err != nil {
		return err
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entityT, domain.CREATE); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.groupRepo.CreateGroup(req, i)
	return
}

func (h *groupUseCase) GetEntityGroup(partyCode string) (domain.EntityTemplate, error) {
	switch partyCode {
	case domain.PARTY_CUSTOMER_GROUP:
		return domain.CUSTOMER, nil
	case domain.PARTY_SUPPLIER_GROUP:
		return domain.SUPPLIER, nil
	case domain.PARTY_ITEM_GROUP:
		return domain.ITEM, nil
	default:
		return domain.EntityTemplate{}, domain.PARTY_TYPE_NOT_FOUND
	}
}
