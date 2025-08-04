package role_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"erp/pkg/logger"
	"fmt"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type RoleRepository interface {
	CreateRole(tx *query.QueryTx, req *common.RequestContext, i dto.RoleData) (dto.RoleDto, error)
	EditRole(tx *query.QueryTx, req *common.RequestContext, i dto.RoleData) error

	EditRolePermissionActions(req *common.RequestContext, i *dto.EditRolePermissionActions) error
	GetRole(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.RoleDto], error)
	GetEntityActions(req *common.RequestContext) (dto.ResultEntity[[]dto.EntityActionsDto], error)
	GetRoleActions(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.RoleActionDto], error)
	// GetRoleActionsByRole(req *common.RequestContext, roleID int64) ([]dto.RoleActionDto, error)
	GetRoles(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.RoleDto], error)
}

type roleRepository struct {
	conn       db.Connection
	Q          *query.Query
	DB         *gorm.DB
	emitLog    logger.EmitLog
	permission repository.PermissionService
	convertor  helpers.ConvertorHelper
	locale     helpers.Locale
}

func NewRoleRepository(
	conn db.Connection,
	logger logger.Logger,
	helpers *helpers.Helpers,
	permission repository.PermissionService,
) RoleRepository {
	return &roleRepository{
		conn:       conn,
		Q:          conn.GetQ(),
		DB:         conn.GetDB(),
		emitLog:    logger.EmitLog("role-repository"),
		locale:     helpers.Locale,
		permission: permission,
		convertor:  helpers.Convertor,
	}
}

func (r *roleRepository) EditRole(tx *query.QueryTx, req *common.RequestContext, d dto.RoleData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Role.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Role{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Role.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

func (r *roleRepository) CreateRole(tx *query.QueryTx, req *common.RequestContext, d dto.RoleData) (res dto.RoleDto, err error) {
	var role model.Role
	id, err := tx.Deal.InsertParty(proto.PartyType_role.String())
	if err != nil {
		return
	}
	fields := d.Fields
	role.ID = id
	role.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &role); err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Role.Save(&role)
	if err != nil {
		return
	}
	if err = tx.Deal.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.RoleDTOFromModel(&role)
	return
}

func (r *roleRepository) EditRolePermissionActions(req *common.RequestContext, i *dto.EditRolePermissionActions) error {

	var (
		err error
	)
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	role, err := r.getRoleByUuid(req.Ctx, i.Body.RoleUUID)
	if err != nil {
		return err
	}
	roleAction := tx.RoleAction
	for _, action := range i.Body.EntityActions.Actions {
		_, err = roleAction.WithContext(req.Ctx).Unscoped().Where(
			roleAction.RoleID.Eq(role.ID),
			roleAction.ActionID.Eq(action.ID),
		).Delete()
		if err != nil {
			return err
		}
	}
	for _, selected := range i.Body.ActionSelecteds {
		roleAction := &model.RoleAction{}
		roleAction.ActionID = selected.ActionID
		roleAction.RoleID = role.ID
		err = tx.WithContext(req.Ctx).RoleAction.Save(roleAction)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// func (r *roleRepository) addRoleAction(req *common.RequestContext, tx *query.QueryTx, roleID, actionID int64,
// 	actionName, entityName string) {
// 	var roleAction model.RoleAction
// 	roleAction.ActionID = actionID
// 	roleAction.RoleID = roleID
// 	tx.RoleAction.Save(&roleAction)
// 	// err := tx.Where(&model.RoleAction{RoleID: roleID, ActionID: actionID}).First(&roleAction).Error
// 	// if err == gorm.ErrRecordNotFound {
// 	// 	roleAction.DeletedAt = gorm.DeletedAt(sql.NullTime{
// 	// 		Valid: false,
// 	// 	})
// 	// 	fmt.Println(roleAction)
// 	// 	err := tx.WithContext(req.Ctx).Save(&roleAction).Error
// 	// 	if err != nil {
// 	// 		// fmt.Println("Fail to add role", err)
// 	// 	}
// 	// 	err = r.permission.WriteTemplateAuthData(req, roleID, entityName, actionName)
// 	// 	if err != nil {
// 	// 		r.emitLog.Err(err, logger.OptionsLog.WithMethod("addRoleAction"))
// 	// 	}
// 	// }
// }

// func (r *roleRepository) removeRoleAction(req *common.RequestContext, tx *query.QueryTx, roleID, actionID int64, actionName, entityName string) {
// 	var roleAction model.RoleAction
// 	fmt.Println(roleID, actionID)
// 	if err := tx.Where(&model.RoleAction{RoleID: roleID, ActionID: actionID}).First(&roleAction).Error; err == nil {
// 		fmt.Println("DELETING ROLE ACTIION ")
// 		roleAction.DeletedAt = gorm.DeletedAt(sql.NullTime{
// 			Time:  time.Now(),
// 			Valid: true,
// 		})
// 		err := tx.Save(&roleAction).Error
// 		if err != nil {
// 			r.emitLog.Err(err, logger.OptionsLog.WithMethod("removeRoleAction"))
// 		}
// 		err = r.permission.DeleteTemplateAuthData(req, roleID, entityName, actionName)
// 		if err != nil {
// 			r.emitLog.Err(err, logger.OptionsLog.WithMethod("removeRoleAction"))
// 		}
// 	}
// }

func (r *roleRepository) GetRole(req *common.RequestContext, d *dto.RequestEntity) (dto.ResultEntity[dto.RoleDto], error) {
	var (
		res dto.ResultEntity[dto.RoleDto]
		err error
	)
	id := r.convertor.StrtoInt(d.ID)
	e := r.Q.Role
	w := r.Q.Workspace
	err = r.Q.WithContext(req.Ctx).Role.Select(
		e.ID, e.UUID, e.Code, e.CreatedAt, e.Description,
		w.Name.As("workspace"), e.WorkspaceID,
	).LeftJoin(
		w, w.ID.EqCol(e.WorkspaceID),
	).
		Where(
			e.CompanyID.Eq(req.ActiveCompany.ID),
			e.ID.Eq(id),
		).Scan(&res.Entity)
	if err != nil {
		r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRole"))
		return res, err
	}
	return res, err
}

func (r *roleRepository) GetEntityActions(req *common.RequestContext) (dto.ResultEntity[[]dto.EntityActionsDto], error) {
	var (
		res      dto.ResultEntity[[]dto.EntityActionsDto]
		entities []dto.EntityDto
		err      error
	)
	cEntities := r.Q.CompanyEntity
	entityQ := r.Q.Entity
	err = r.Q.CompanyEntity.WithContext(req.Ctx).Select(
		entityQ.ID, entityQ.Name,
	).
		Join(entityQ, cEntities.EntityID.EqCol(entityQ.ID)).
		Where(
			cEntities.CompanyID.Eq(req.ActiveCompany.ID),
			cEntities.Enabled.Is(true),
		).Order(entityQ.Name.Asc()).Scan(&entities)
	if err != nil {
		return res, err
	}
	ids := lo.Map(entities, func(item dto.EntityDto, index int) int64 {
		return item.ID
	})
	entityActions := make([]dto.EntityActionsDto, len(entities))
	actions, err := r.getActions(req.Ctx, ids)
	if err != nil {
		return res, err
	}
	groups := lo.GroupBy(actions, func(item dto.ActionDto) int64 {
		return item.EntityID
	})
	fmt.Println("GROUPS", groups)
	for i, entity := range entities {
		// entity.Name = r.locale.MustLocalize(
		// 	helpers.OptionsLocale.WithID(fmt.Sprintf("Entity.%s",entity.Name)),
		// 	helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		// )
		if val, ok := groups[entity.ID]; ok {
			fmt.Println("ENTITY", val, entity)
			entityActions[i].Entity = entity
			entityActions[i].Actions = val
		}

	}
	res.Entity = entityActions
	return res, err
}

func (s *roleRepository) getActions(ctx context.Context, entityIDs []int64) ([]dto.ActionDto, error) {
	var (
		actions []dto.ActionDto
	)
	actionQ := s.Q.Action
	err := s.Q.Action.WithContext(ctx).Select(
		actionQ.EntityID, actionQ.ID, actionQ.Name,
	).Where(
		actionQ.EntityID.In(entityIDs...),
	).
		Scan(&actions)
	if err != nil {
		return actions, err
	}
	return actions, err
}

func (r *roleRepository) GetRoleActions(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.RoleActionDto], error) {
	var (
		res dto.PaginationResult[[]dto.RoleActionDto]
		err error
	)
	id := r.convertor.StrtoInt(i.FilterID)

	roleActions, err := r.Q.RoleAction.WithContext(req.Ctx).Where(
		r.Q.RoleAction.RoleID.Eq(id),
	).Preload(
		r.Q.RoleAction.Action,
	).Find()
	roleActionsDto := make([]dto.RoleActionDto, len(roleActions))
	for i, roleAction := range roleActions {
		roleActionsDto[i] = dto.RoleActionDTOFromModel(roleAction)
	}
	res.Results = roleActionsDto
	return res, err
}

func (r *roleRepository) getRoleByUuid(ctx context.Context, uuid string) (*model.Role, error) {
	role, err := r.Q.Role.WithContext(ctx).Where(
		r.Q.Role.UUID.Eq(uuid),
	).First()
	return role, err
}

// func (r *roleRepository) GetRoleActionsByRole(req *common.RequestContext, roleID int64) (
// 	[]dto.RoleActionDto, error) {
// 	roleActions, err := r.Q.RoleAction.WithContext(req.Ctx).Where(
// 		r.Q.RoleAction.RoleID.Eq(roleID),
// 	).Preload(
// 		r.Q.RoleAction.Action,
// 	).Find()
// 	roleActionsDto := make([]dto.RoleActionDto, len(roleActions))
// 	for i, roleAction := range roleActions {
// 		roleActionsDto[i] = dto.RoleActionDTOFromModel(roleAction)
// 	}
// 	return roleActionsDto, err
// }

func (r *roleRepository) GetRoles(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.RoleDto], error) {
	var (
		result dto.PaginationResult[[]dto.RoleDto]
		err    error
	)
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoles"))
		}
	}()
	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
	result.Total, err = r.Q.Role.WithContext(req.Ctx).CountPaginate("company_id", int64(req.ActiveCompany.ID),
		i.Query, i.Enabled)
	if err != nil {
		return result, err
	}
	roles, err := r.Q.Role.WithContext(req.Ctx).Paginate("company_id", int64(req.ActiveCompany.ID), i.Query,
		limit, offset, i.OrderColumn, i.Order, i.Enabled)
	if err != nil {
		return result, err
	}
	roleDtos := make([]dto.RoleDto, len(roles))
	for i, role := range roles {
		roleDtos[i] = dto.RoleDTOFromModel(&role)
	}
	result.Results = roleDtos
	// parentCompany, err := r.permission.GetParentCompany(req, r.conn.Db)
	// queryBuilder := r.conn.Db.WithContext(ctx).Model(&model.Role{}).
	// 	Where(&model.Role{CompanyID: parentCompany.ID})

	// err = queryBuilder.
	// 	Count(&result.Total).Error

	// if i.Query != "" {
	// 	queryBuilder = queryBuilder.Where("code ILIKE ?", "%"+i.Query+"%")
	// }

	// err = queryBuilder.
	// 	Scopes(r.conn.Paginate(req.Params)).
	// 	Find(&result.Results).Error
	// if err != nil {
	// 	return result, err
	// }

	return result, err
}

func (r *roleRepository) GetActions(req *common.RequestContext, entityT domain.EntityTemplate) []dto.ActionDto {
	var actions []model.Action
	err := r.DB.WithContext(req.Ctx).Where(&model.Action{EntityID: entityT.ID}).Find(&actions).Error
	if err != nil {
		r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetActions"))
	}
	actionDtos := make([]dto.ActionDto, len(actions))
	for i, action := range actions {
		actionDtos[i] = dto.ActionDtoFromModel(&action)
	}

	return actionDtos
}
