package account_service

import (
	"context"
	"database/sql"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/connection"
	"erp/internal/app/service/helpers"
	"erp/internal/app/service/services/company_service"
	"erp/internal/domain"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type RoleService struct {
	conn              *connection.Connection
	timeout           time.Duration
	emitLog           logger.EmitLog
	companyService    *company_service.CompanyService
	convertor         helpers.ConvertorHelper
	permissionService permission.PermissionService
}

func NewRoleService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	companyService *company_service.CompanyService,
	permissionService permission.PermissionService,
	logger logger.Logger,
) *RoleService {
	return &RoleService{
		conn:              conn,
		timeout:           timeout,
		emitLog:           logger.EmitLog("role-service"),
		companyService:    companyService,
		convertor:         helpers.Convertor,
		permissionService: permissionService,
	}
}

func (s *RoleService) CreateRole(req *common.RequestContext, i *dto.RoleRequestData) (err error) {
	// ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	// defer cancel()
	// var (
	// 	err error
	// )
	// defer func() {
	// 	if err != nil {
	// 		s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateRole"))
	// 	}
	// }()
	// var role model.Role
	// if allow := s.permissionService.CheckPermission(ctx, req, domain.ROLE, domain.CREATE); !allow {
	// 	return domain.ACTION_NOT_ALLOWED
	// }
	// parentCompany, err := s.permissionService.GetParentCompany(req, s.conn.Db)
	// if err != nil {
	// 	return err
	// }
	// role.CompanyID = parentCompany.ID
	// role.Code = i.Body.Name
	// role.Description = &i.Body.Description
	// err = s.conn.Db.WithContext(ctx).Save(&role).Error
	return err
}

func (s *RoleService) EditRolePermissionActions(req *common.RequestContext, i *dto.EditRolePermissionActions) error {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		err error
	)
	tx := s.conn.Db.WithContext(ctx).Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("EditRilePermission"))
		}
	}()
	role, err := s.getRoleByUuid(ctx, i.Body.RoleUUID)
	if err != nil {
		return err
	}
	for _, actionSelected := range i.Body.ActionSelecteds {
		// err = tx.WithContext(ctx).Where(&model.RoleAction{ActionID: actionSelected.ActionID})
		if actionSelected.Selected {
			s.addRoleAction(req, tx, role.ID, actionSelected.ActionID, actionSelected.ActionName, i.Body.EntityName)
		} else {
			s.removeRoleAction(req, tx, role.ID, actionSelected.ActionID, actionSelected.ActionName, i.Body.EntityName)
		}
		fmt.Println(actionSelected)
	}
	return tx.Commit().Error
}

func (s *RoleService) addRoleAction(req *common.RequestContext, tx *gorm.DB, roleID, actionID int64,
	actionName, entityName string) {
	var roleAction model.RoleAction
	roleAction.ActionID = actionID
	roleAction.RoleID = roleID
	err := tx.Where(&model.RoleAction{RoleID: roleID, ActionID: actionID}).First(&roleAction).Error
	if err == gorm.ErrRecordNotFound {
		roleAction.DeletedAt = gorm.DeletedAt(sql.NullTime{
			Valid: false,
		})
		fmt.Println(roleAction)
		err := tx.Save(&roleAction).Error
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("addRoleAction"))
			// fmt.Println("Fail to add role", err)
		}
		err = s.permissionService.WriteTemplateAuthData(req, roleID, entityName, actionName)
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("addRoleAction"))
		}
	}
}

func (s *RoleService) removeRoleAction(req *common.RequestContext, tx *gorm.DB, roleID, actionID int64, actionName, entityName string) {
	var roleAction model.RoleAction
	fmt.Println(roleID, actionID)
	if err := tx.Where(&model.RoleAction{RoleID: roleID, ActionID: actionID}).First(&roleAction).Error; err == nil {
		fmt.Println("DELETING ROLE ACTIION ")
		roleAction.DeletedAt = gorm.DeletedAt(sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		})
		err := tx.Save(&roleAction).Error
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("removeRoleAction"))
		}
		err = s.permissionService.DeleteTemplateAuthData(req, roleID, entityName, actionName)
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("removeRoleAction"))
		}
	}
}

func (s *RoleService) GetRole(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.RoleDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res  dto.ResultEntity[dto.RoleDto]
		role model.Role
		err  error
	)

	err = s.conn.Db.WithContext(ctx).Where(&model.Role{
		UUID: i.ID,
	}).First(&role).Error
	if err != nil {
		s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRole"))
		return res, err
	}
	res.Entity = dto.RoleDTOFromModel(&role)

	// res.model.RoleDefinitions,err = s.getRoleDefinitions(ctx,uint(roleID))
	// if err != nil {
	// 	s.emitLog.Err(err,logger.OptionsLog.WithMethod("getRoleDefinitions"))
	// }
	return res, err
}

func (s *RoleService) GetEntityActions(req *common.RequestContext) (dto.ResultEntity[[]dto.EntityActionsDto], error) {
	var (
		res      dto.ResultEntity[[]dto.EntityActionsDto]
		entities []dto.EntityDto
		err      error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetEntityActions"))
		}
	}()
	cEntities := s.conn.Q.CompanyEntity
	entityQ := s.conn.Q.Entity
	err = s.conn.Q.CompanyEntity.WithContext(req.Ctx).Select(
		entityQ.ID,entityQ.Name,
	).Where(
		cEntities.CompanyID.Eq(req.ActiveCompany.ID),
	).Order(entityQ.Name.Asc()).Scan(&entities)
	if err != nil {
		return res,err
	}
	entityActions := make([]dto.EntityActionsDto, len(entities))
	for i, entity := range entities {
		actions, err := s.getActions(req.Ctx, entity.ID)
		if err != nil {
			continue
		}
		entityActions[i].Entity = entity
		entityActions[i].Actions = actions
	}
	res.Entity = entityActions
	return res, err
}

func (s *RoleService) getActions(ctx context.Context, entityID int64) ([]dto.ActionDto, error) {
	var (
		actions []dto.ActionDto
	)
	actionQ := s.conn.Q.Action
	err := s.conn.Q.Action.WithContext(ctx).Select(
		actionQ.ID,actionQ.Name,
	).Where(actionQ.EntityID.Eq(entityID)).
		Scan(&actions)
	if err != nil {
		return actions, err
	}
	return actions, err
}

func (s *RoleService) GetRoleActions(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.RoleActionDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.PaginationResult[[]dto.RoleActionDto]
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoleActions"))
		}
	}()
	role, err := s.getRoleByUuid(ctx, i.FilterID)
	if err != nil {
		return res, err
	}
	roleActions, err := s.conn.Q.RoleAction.WithContext(ctx).Where(
		s.conn.Q.RoleAction.RoleID.Eq(role.ID),
	).Preload(
		s.conn.Q.RoleAction.Action,
	).Find()
	roleActionsDto := make([]dto.RoleActionDto, len(roleActions))
	for i, roleAction := range roleActions {
		roleActionsDto[i] = dto.RoleActionDTOFromModel(roleAction)
	}
	res.Results = roleActionsDto
	return res, err
}

func (s *RoleService) getRoleByUuid(ctx context.Context, uuid string) (*model.Role, error) {
	role, err := s.conn.Q.Role.WithContext(ctx).Where(
		s.conn.Q.Role.UUID.Eq(uuid),
	).First()
	return role, err
}

func (s *RoleService) GetRoleActionsByRole(req *common.RequestContext, roleID int64) (
	[]dto.RoleActionDto, error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	roleActions, err := s.conn.Q.RoleAction.WithContext(ctx).Where(
		s.conn.Q.RoleAction.RoleID.Eq(roleID),
	).Preload(
		s.conn.Q.RoleAction.Action,
	).Find()
	roleActionsDto := make([]dto.RoleActionDto, len(roleActions))
	for i, roleAction := range roleActions {
		roleActionsDto[i] = dto.RoleActionDTOFromModel(roleAction)
	}
	return roleActionsDto, err
}

func (s *RoleService) GetRoles(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.RoleDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		result dto.PaginationResult[[]dto.RoleDto]
		err    error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoles"))
		}
	}()
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ROLE, domain.VIEW); !allow {
		return result, domain.ACTION_NOT_ALLOWED
	}
	limit, offset := s.convertor.ToPaginationParams(i.Page, i.Size)
	result.Total, err = s.conn.Q.Role.WithContext(ctx).CountPaginate("company_id", int64(req.ActiveCompany.ID),
		i.Query, i.Enabled)
	if err != nil {
		return result, err
	}
	roles, err := s.conn.Q.Role.WithContext(ctx).Paginate("company_id", int64(req.ActiveCompany.ID), i.Query,
		limit, offset, i.OrderColumn, i.Order, i.Enabled)
	if err != nil {
		return result, err
	}
	roleDtos := make([]dto.RoleDto, len(roles))
	for i, role := range roles {
		roleDtos[i] = dto.RoleDTOFromModel(&role)
	}
	result.Results = roleDtos
	// parentCompany, err := s.permissionService.GetParentCompany(req, s.conn.Db)
	// queryBuilder := s.conn.Db.WithContext(ctx).Model(&model.Role{}).
	// 	Where(&model.Role{CompanyID: parentCompany.ID})

	// err = queryBuilder.
	// 	Count(&result.Total).Error

	// if i.Query != "" {
	// 	queryBuilder = queryBuilder.Where("code ILIKE ?", "%"+i.Query+"%")
	// }

	// err = queryBuilder.
	// 	Scopes(s.conn.Paginate(req.Params)).
	// 	Find(&result.Results).Error
	// if err != nil {
	// 	return result, err
	// }

	return result, err
}

func (s *RoleService) GetActions(req *common.RequestContext, entityT domain.EntityTemplate) []dto.ActionDto {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var actions []model.Action
	err := s.conn.Db.WithContext(ctx).Where(&model.Action{EntityID: entityT.ID}).Find(&actions).Error
	if err != nil {
		s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetActions"))
	}
	actionDtos := make([]dto.ActionDto, len(actions))
	for i, action := range actions {
		actionDtos[i] = dto.ActionDtoFromModel(&action)
	}

	return actionDtos
}

// func (s *RoleService) GetItems(req *common.RequestContext, d *dto.RequestPaginationData) (dto.PaginationResult[[]model.Role], error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
// 	defer cancel()
// 	var result dto.PaginationResult[[]model.Role]
// 	queryBuilder := s.conn.Db.WithContext(ctx).Model(&model.Role{}).
// 		Where(&model.Role{CompanyID: req.ActiveCompany.ID})

// 	err := queryBuilder.
// 		Count(&result.Total).Error

// 	if d.Query != "" {
// 		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
// 	}

// 	err = queryBuilder.
// 		Scopes(s.conn.Paginate(req.Params)).
// 		Find(&result.Results).Error
// 	if err != nil {
// 		return result, err
// 	}

// 	return result, err
// }
