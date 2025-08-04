package role_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/di"
	"erp/pkg/logger"
	role_repo "erp/project/manage/role/internal/repository"
)

type RoleUseCase interface {
	CreateRole(req *common.RequestContext, d dto.RoleData) (dto.RoleDto,error)
	EditRole(req *common.RequestContext, d dto.RoleData) error

	EditRolePermissionActions(req *common.RequestContext, i *dto.EditRolePermissionActions) error
	GetRole(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.RoleDto], error)
	GetEntityActions(req *common.RequestContext) (dto.ResultEntity[[]dto.EntityActionsDto], error)
	GetRoleActions(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.RoleActionDto], error)
	// GetRoleActionsByRole(req *common.RequestContext, roleID int64) ([]dto.RoleActionDto, error)
	GetRoles(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.RoleDto], error)
}

type roleUseCase struct {
	emitLog    logger.EmitLog
	permission repository.PermissionService
	roleRepo   role_repo.RoleRepository
    c di.Container
}

func NewRoleUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	roleRepo role_repo.RoleRepository,
	c di.Container,
) RoleUseCase {
	return &roleUseCase{
		emitLog:    logger.EmitLog("role-usecase"),
		permission: permission,
		roleRepo:   roleRepo,
		c: c,
	}
}

func (u *roleUseCase) EditRole(req *common.RequestContext, d dto.RoleData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()

	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditRole"))
		}
		err = domain.CloseTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.ROLE, domain.EDIT); err != nil {
		return err
	}
	err = u.roleRepo.EditRole(tx, req, d)
	return
}

func (u *roleUseCase) CreateRole(req *common.RequestContext, d dto.RoleData) (res dto.RoleDto,err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()

	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateRole"))
		}
		err = domain.CloseTx(tx, err)
	}(tx)
	if err := u.permission.CheckPermissionEntity(req.Ctx, req, domain.ROLE, domain.CREATE); err != nil {
		return res, err
	}
	res, err = u.roleRepo.CreateRole(tx, req, d)
	if err != nil {
		return
	}
	return 
}

func (s *roleUseCase) EditRolePermissionActions(req *common.RequestContext, i *dto.EditRolePermissionActions) (err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("EditRilePermission"))
		}
	}()
	err = s.roleRepo.EditRolePermissionActions(req, i)
	return
}

func (s *roleUseCase) GetRole(req *common.RequestContext, i *dto.RequestEntity) (res dto.ResultEntity[dto.RoleDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRole"))
		}
	}()
	res, err = s.roleRepo.GetRole(req, i)
	// res.model.RoleDefinitions,err = s.getRoleDefinitions(ctx,uint(roleID))
	// if err != nil {
	// 	s.emitLog.Err(err,logger.OptionsLog.WithMethod("getRoleDefinitions"))
	// }
	return res, err
}

func (s *roleUseCase) GetEntityActions(req *common.RequestContext) (res dto.ResultEntity[[]dto.EntityActionsDto], err error) {

	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetEntityActions"))
		}
	}()
	res, err = s.roleRepo.GetEntityActions(req)
	return
}

func (s *roleUseCase) GetRoleActions(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.RoleActionDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoleActions"))
		}
	}()
	res, err = s.roleRepo.GetRoleActions(req, i)
	return
}

// func (s *roleUseCase) GetRoleActionsByRole(req *common.RequestContext, roleID int64) (
// 	res []dto.RoleActionDto, err error) {
// 	defer func() {
// 		if err != nil {
// 			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoleActionsByRole"))
// 		}
// 	}()
// 	res, err = s.roleRepo.GetRoleActionsByRole(req, roleID)
// 	return res, err
// }

func (s *roleUseCase) GetRoles(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.RoleDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoles"))
		}
	}()
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.ROLE, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = s.roleRepo.GetRoles(req, i)
	return
}

// func (s *roleUseCase) GetActions(req *common.RequestContext, entityT domain.EntityTemplate) []dto.ActionDto {
// 	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
// 	defer cancel()
// 	var actions []model.Action
// 	err := s.conn.Db.WithContext(ctx).Where(&model.Action{EntityID: entityT.ID}).Find(&actions).Error
// 	if err != nil {
// 		s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetActions"))
// 	}
// 	actionDtos := make([]dto.ActionDto, len(actions))
// 	for i, action := range actions {
// 		actionDtos[i] = dto.ActionDtoFromModel(&action)
// 	}

// 	return actionDtos
// }
