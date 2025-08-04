package role_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/proto"
	"erp/internal/domain/event"
)

type RoleEventRepository interface {
	OnCompanyCreated(ctx context.Context, payload event.CreatedCompanyEventData) (err error)
}

type roleEventRepository struct {
}

func NewRoleEventRepository() RoleEventRepository {
	return &roleEventRepository{}
}

func (r *roleEventRepository) OnCompanyCreated(ctx context.Context, payload event.CreatedCompanyEventData) (err error) {
	//Create role root admin
	tx := payload.Tx
	roleQ := tx.Role
	role := model.Role{}
	roleID, err := roleQ.WithContext(ctx).InsertParty(proto.PartyType_role.String())
	if err != nil {
		return
	}
	role.ID = roleID
	role.Code = proto.PartyType_admin.String()
	role.CompanyID = payload.Company.ID
	err = roleQ.WithContext(ctx).Save(&role)
	if err != nil {
		return
	}
	return
}
