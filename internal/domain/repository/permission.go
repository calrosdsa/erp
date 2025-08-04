package repository

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
	GetEntitiesActions(ctx context.Context, entityIDS []int64) map[int][]dto.ActionDto
	
	GetParentCompany(req *common.RequestContext, tx *gorm.DB) (model.Company, error)
	DeleteTemplateAuthData(req *common.RequestContext, roleID int64,
		entityName string, actionName string) error
	WriteTemplateAuthData(req *common.RequestContext, roleID int64,
		entityName string, actionName string) error
	CheckPermission(ctx context.Context, req *common.RequestContext,
		entity domain.EntityTemplate, actionName domain.ActionType) bool
	CheckPermissionEntity(ctx context.Context, req *common.RequestContext,
		entity domain.EntityTemplate, actionName domain.ActionType) error
	CheckCustomPermission(ctx context.Context, req *common.RequestContext, permissionOpts ...PermissionOpt) bool
	CheckIfUserIsCompanyAdmin(ctx context.Context, req *common.RequestContext,
		companyID int64, permission domain.Permission) bool
	GetDocumentEntity(partyType string) (domain.EntityTemplate,error)
}



type permissionOpts struct {
	EntityType      domain.PermifyEntity
	EntityID        string
	Permission      domain.Permission
	SubjectType     domain.PermifyEntity
	SubjectID       string
	SubjectRelation string
}

var PermissionOpts permissionOpts

type PermissionOpt func(opts *permissionOpts)

func (*permissionOpts) WithEntityType(entityType domain.PermifyEntity) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.EntityType = entityType
	}
}
func (*permissionOpts) WithEntityID(entityId string) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.EntityID = entityId
	}
}
func (*permissionOpts) WithPermission(permission domain.Permission) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.Permission = permission
	}
}
func (*permissionOpts) WithSubjectType(subjectType domain.PermifyEntity) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.SubjectType = subjectType
	}
}
func (*permissionOpts) WithSubjectId(subjectId string) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.SubjectID = subjectId
	}
}
func (*permissionOpts) WithSubjectRelation(subjectRelation string) PermissionOpt {
	return func(opts *permissionOpts) {
		opts.SubjectRelation = subjectRelation
	}
}

func (*permissionOpts) Apply(opts ...PermissionOpt) permissionOpts {
	options := permissionOpts{}
	for _, opt := range opts {
		opt(&options)
	}
	return options
}
