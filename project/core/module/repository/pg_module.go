package module_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/helper"
)

type ModuleRepository interface {
	GetModule(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ModuleDetailDto], err error)
	GetModuleDetail(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ModuleDetailDto], err error)
	GetEntitiesSearch(req *common.RequestContext, d *dto.ModuleSearchRequest) (
		res []dto.EntityDto, err error,
	)
	CreateModule(req *common.RequestContext, d dto.ModuleData) (
		res model.Module, err error)
	GetModules(req *common.RequestContext, d dto.ModulesRequest) (
		res []dto.ModuleDto, err error)
	EditModule(req *common.RequestContext, d dto.ModuleData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
}

type moduleRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	currency  helpers.CurrencyHelper
	query     helpers.QueryHelper
}

func NewModuleRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) ModuleRepository {
	return &moduleRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
		currency:  helpers.Currency,
		query:     helpers.Query,
	}
}
func (r *moduleRepository) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	moduleQ := r.Q.Module
	_, err = r.Q.Module.WithContext(req.Ctx).Where(
		moduleQ.CompanyID.Eq(req.ActiveCompany.ID),
		moduleQ.Status.Eq(d.Body.CurrentState),
		moduleQ.ID.Eq(r.convertor.StrtoInt(d.Body.PartyID)),
	).UpdateSimple(
		moduleQ.Status.Value(nextState),
	)
	return
}

func (r *moduleRepository) GetModule(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ModuleDetailDto], err error) {
	moduleQ := r.Q.Module
	err = moduleQ.WithContext(req.Ctx).Select(
		moduleQ.ID, moduleQ.Href, moduleQ.Label, moduleQ.IconCode,
		moduleQ.IconName, moduleQ.Status, moduleQ.UUID,
	).
		Where(moduleQ.CompanyID.Eq(req.ActiveCompany.ID),
			moduleQ.Href.Eq(d.ID)).
		Scan(&res.Entity.Module)
	if err != nil {
		return
	}
	res.Entity.Sections, err = r.getModuleSections(req.Ctx, res.Entity.Module.ID)

	return
}

func (r *moduleRepository) GetModuleDetail(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ModuleDetailDto], err error) {
	moduleQ := r.Q.Module
	err = moduleQ.WithContext(req.Ctx).Select(
		moduleQ.ID, moduleQ.Href, moduleQ.Label, moduleQ.IconCode,
		moduleQ.IconName, moduleQ.Status, moduleQ.UUID, moduleQ.HasDirectAccess,
		moduleQ.Priority,
	).
		Where(moduleQ.CompanyID.Eq(req.ActiveCompany.ID), moduleQ.UUID.Eq(d.ID)).
		Scan(&res.Entity.Module)
	if err != nil {
		return
	}
	res.Entity.Sections, err = r.getModuleSections(req.Ctx, res.Entity.Module.ID)
	return
}

func (r *moduleRepository) getModuleSections(ctx context.Context, moduleID int64) (res []dto.ModuleSectionDto, err error) {
	moduleSectionQ := r.Q.ModuleSection
	entityQ := r.Q.Entity
	moduleSectionQ.WithContext(ctx).Select(
		moduleSectionQ.Name.As("section_name"), entityQ.ID,
		entityQ.Name, entityQ.Href,
	).Join(entityQ, entityQ.ID.EqCol(moduleSectionQ.EntityID)).
		Where(
			moduleSectionQ.ModuleID.Eq(moduleID),
		).
		Scan(&res)
	return
}

func (r *moduleRepository) CreateModule(req *common.RequestContext, d dto.ModuleData) (
	res model.Module, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	moduleID, err := tx.Module.InsertParty(proto.PartyType_module.String())
	if err != nil {
		return
	}
	res.ID = moduleID
	res.Status = proto.State_ENABLED.String()
	res.CompanyID = req.ActiveCompany.ID
	fields := d.Fields
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}

	err = tx.Module.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return
	}
	err = tx.Module.InsertActivity(res.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return
	}
	if err = r.createModuleSections(tx, req.Ctx, res.ID, d.Sections); err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *moduleRepository) EditModule(req *common.RequestContext, d dto.ModuleData) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}

	err = tx.Module.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Module{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = tx.Module.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	if err = r.deleteModuleSections(tx, req.Ctx, d.ID); err != nil {
		return
	}
	if err = r.createModuleSections(tx, req.Ctx, d.ID, d.Sections); err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *moduleRepository) createModuleSections(tx *query.QueryTx, ctx context.Context,
	moduleID int64, d []dto.ModuleSectionData) (err error) {
	moduleSections := make([]*model.ModuleSection, len(d))
	for i, section := range d {
		moduleSection := &model.ModuleSection{}
		moduleSection.ModuleID = moduleID
		moduleSection.Name = section.Name
		moduleSection.EntityID = section.EntityID
		moduleSections[i] = moduleSection
	}
	err = tx.ModuleSection.WithContext(ctx).CreateInBatches(moduleSections, len(moduleSections))
	return
}

func (r *moduleRepository) deleteModuleSections(tx *query.QueryTx, ctx context.Context,
	moduleID int64) (err error) {
	_, err = tx.ModuleSection.WithContext(ctx).Unscoped().Where(
		tx.ModuleSection.ModuleID.Eq(moduleID),
	).Delete()
	return
}

func (r *moduleRepository) GetModules(req *common.RequestContext, d dto.ModulesRequest) (
	res []dto.ModuleDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Deal
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.moduleQuery(req, queryData, &generateSQL)
	fmt.Println("MODULES QUERY", generateSQL.String())
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *moduleRepository) moduleQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
				e.id,e.href,e.icon_code,e.icon_name,e.label,e.uuid,
				e.status,e.priority
				from modules as e 
			`)
	if v, ok := d["workspace_id"]; ok {
		generateSQL.WriteString(`
				inner join workspace_modules as wm on wm.module_id = e.id and wm.workspace_id = ?
				`)
		params = append(params, v)
	}
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"label",
		"status",
		"created_at",
		"updated_at",
	}

	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

func (r *moduleRepository) GetEntitiesSearch(req *common.RequestContext, d *dto.ModuleSearchRequest) (
	res []dto.EntityDto, err error) {
	//d.Query is convert to lowercase in the client side
	var (
		generateSQL strings.Builder
		params      []interface{}
	)

	if d.LoadEntities {
		generateSQL.WriteString(`select e.id, e.name,e.href 
		from entities as e
		inner join actions as a on a.entity_id = e.id and a.name = 'view'
		inner join role_actions as ra on ra.action_id = a.id
		where lower(e.name) like ? and ra.role_id = ? and e.href != '' `)
		params = append(params, "%"+d.Query+"%", req.Role.ID)
	}
	if d.LoadEntities && d.LoadModules {
		generateSQL.WriteString(`union all `)
	}
	if d.LoadModules {
		generateSQL.WriteString(`select m.id,m.label as name,m.href 
		from modules as m
		where lower(m.label) like ? and m.company_id = ? and m.status = ? `)
		params = append(params, "%"+d.Query+"%", req.ActiveCompany.ID, proto.State_ENABLED.String())
	}
	generateSQL.WriteString("limit ?")
	params = append(params, d.Size)

	err = r.Q.Entity.WithContext(req.Ctx).UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}
