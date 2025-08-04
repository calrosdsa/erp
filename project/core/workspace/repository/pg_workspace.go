package workspace_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"
	"strings"

	"gorm.io/gen/helper"
)

type WorkSpaceRepository interface {
	EditWorkSpace(tx *query.QueryTx, req *common.RequestContext, d dto.WorkSpaceData) (err error)
	CreateWorkSpace(tx *query.QueryTx, req *common.RequestContext, d dto.WorkSpaceData) (res dto.WorkSpaceDto, err error)

	GetWorkSpaces(req *common.RequestContext, d dto.WorkSpaceRequest) (res []dto.WorkSpaceDto, err error)
	GetWorkSpace(req *common.RequestContext, d dto.RequestEntity) (res dto.WorkSpaceDto, err error)
	UpdateStatus(tx *query.QueryTx, req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (
		res *model.Workspace, err error)
}

type workspaceRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	query     helpers.QueryHelper
}

func NewWorkSpaceRepo(
	conn db.Connection,
	helpers *helpers.Helpers,
) WorkSpaceRepository {
	return &workspaceRepo{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
		query:     helpers.Query,
	}
}

func (r *workspaceRepo) UpdateStatus(tx *query.QueryTx, req *common.RequestContext, d dto.UpdateStatusWithEvent, nextState string) (
	res *model.Workspace, err error) {
	e := r.Q.Workspace
	id:= r.convertor.StrtoInt(d.Body.PartyID)
	_, err = tx.Workspace.WithContext(req.Ctx).Where(
		e.CompanyID.Eq(req.ActiveCompany.ID),
		e.Status.Eq(d.Body.CurrentState),
		e.ID.Eq(id),
	).UpdateSimple(
		e.Status.Value(nextState),
	)
	if err != nil {
		return
	}
	res, err = tx.Workspace.WithContext(req.Ctx).Where(
		e.ID.Eq(r.convertor.StrtoInt(d.Body.PartyID)),
		e.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}

	return
}

func (r *workspaceRepo) GetWorkSpace(req *common.RequestContext, d dto.RequestEntity) (
	res dto.WorkSpaceDto, err error) {
	id := r.convertor.StrtoInt(d.ID)
	e := r.Q.Workspace
	err = r.Q.Workspace.WithContext(req.Ctx).Select(
		e.ID, e.Name, e.Status, e.CreatedAt,
	).Where(
		e.CompanyID.Eq(req.ActiveCompany.ID),
		e.ID.Eq(id),
	).Scan(&res)
	return
}

func (r *workspaceRepo) EditWorkSpace(tx *query.QueryTx, req *common.RequestContext, d dto.WorkSpaceData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Workspace.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Workspace{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = r.processModulesWorkspace(tx, req.Ctx, d, d.ID)
	if err != nil {
		return
	}

	err = tx.Workspace.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}

	return
}

func (r *workspaceRepo) CreateWorkSpace(tx *query.QueryTx, req *common.RequestContext, d dto.WorkSpaceData) (
	res dto.WorkSpaceDto, err error) {
	var workSpace model.Workspace
	id, err := tx.Address.InsertParty(proto.PartyType_workspace.String())
	if err != nil {
		return
	}
	fields := d.Fields
	workSpace.ID = id
	workSpace.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &workSpace); err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Workspace.Save(&workSpace)
	if err != nil {
		return
	}
	err = r.processModulesWorkspace(tx, req.Ctx, d, workSpace.ID)
	if err != nil {
		return
	}
	if err = tx.Address.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.WorkspaceFromModel(workSpace)
	return
}

func (r *workspaceRepo) GetWorkSpaces(req *common.RequestContext, d dto.WorkSpaceRequest) (
	res []dto.WorkSpaceDto, err error,
) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Deal
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.workSpaceQuery(req, queryData, &generateSQL)

	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *workspaceRepo) workSpaceQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
					e.id,e.name,e.created_at,e.status
					from workspaces as e 
		`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"name",
		"status",
		"created_at",
	}

	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	r.query.ReferenceFilterBuilder(generateSQL, &whereSQL, &params, d, "workspace_id")

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *workspaceRepo) processModulesWorkspace(tx *query.QueryTx, ctx context.Context, d dto.WorkSpaceData,
	workspaceID int64) (err error) {
	e := tx.WorkspaceModule
	_,err = e.WithContext(ctx).Unscoped().Where(
		e.WorkspaceID.Eq(workspaceID),
	).Delete()
	if err != nil {
		return
	}
	workspaces := make([]*model.WorkspaceModule, len(d.Modules))
	for i, m := range d.Modules {
		workspaceModule := &model.WorkspaceModule{
			ModuleID:    m,
			WorkspaceID: workspaceID,
		}
		workspaces[i] = workspaceModule
	}
	err = e.WithContext(ctx).CreateInBatches(workspaces,domain.DEFAULT_BATCH_SIZE)
	return
}
