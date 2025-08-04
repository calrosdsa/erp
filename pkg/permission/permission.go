package permission

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/domain"

	"gorm.io/gorm"
)

type PermissionService interface {
	GetActions(ctx context.Context, entityID int64) []dto.ActionDto
	GetParentCompany(req *common.RequestContext, tx *gorm.DB) (model.Company, error)
	DeleteTemplateAuthData(req *common.RequestContext, roleID int64,
		entityName string, actionName string) error
	WriteTemplateAuthData(req *common.RequestContext, roleID int64,
		entityName string, actionName string) error
	CheckPermission(ctx context.Context, req *common.RequestContext,
		entity domain.EntityTemplate, actionName domain.ActionType) bool
	CheckCustomPermission(ctx context.Context, req *common.RequestContext, permissionOpts ...domain.PermissionOpt) bool
	CheckIfUserIsCompanyAdmin(ctx context.Context, req *common.RequestContext,
		companyID int64, permission domain.Permission) bool
}
