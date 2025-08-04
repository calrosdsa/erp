package repository

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
)

type SessionService interface {
	GetUserRelations(req *common.RequestContext) ([]dto.UserRelationDto, error)
	GetUserRelation(ctx context.Context, uuid string) (model.UserRelation, error)
	GetCompanyDefaults(ctx context.Context,companyID int64) (model.CompanyDefault, error)
	GetUser(ctx context.Context,id int64) (model.User, error)
	GetUserRelationByUserID(ctx context.Context,id int64) (model.UserRelation, error)
	GetRoleActionsByRole(req *common.RequestContext, roleID int64) ([]dto.RoleActionDto, error)
	InsertUser(ctx context.Context, tx *query.QueryTx, identifier string) (model.User, error)
}