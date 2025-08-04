package stage_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"
	"strings"

	"gorm.io/gen/helper"
)

type StageRepository interface {
	CreateStage(req *common.RequestContext, d dto.StageData) (res model.Stage, err error)
	EditStage(req *common.RequestContext, d dto.StageData) (err error)
	GetStages(req *common.RequestContext, d dto.StagesRequest) (res []dto.StageDto, err error)
	StageTransition(req *common.RequestContext,d dto.StageTransitionData)(err error)
	DeleteStage(req *common.RequestContext,d *dto.DeleteRequest)(err error)
}

type stageRepo struct {
	query     helpers.QueryHelper
	convertor helpers.ConvertorHelper
	Q         *query.Query
}

func NewStageRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) StageRepository {
	return &stageRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		query:     helpers.Query,
	}
}

func(r *stageRepo)DeleteStage(req *common.RequestContext,d *dto.DeleteRequest)(err error){
	id := r.convertor.StrtoInt(d.ID)
	_,err = r.Q.Stage.Where(
		r.Q.Stage.ID.Eq(int32(id)),
		r.Q.Stage.CompanyID.Eq(req.ActiveCompany.ID),
	).Delete()
	return
}

func (r *stageRepo)StageTransition(req *common.RequestContext,d dto.StageTransitionData)(err error) {
	tx := r.Q.Begin()
	defer func(){
		if err != nil {
			tx.Rollback()
		}
	}()
	//Change indexes base on source and destination
	fmt.Println("TRANSITION DATA", d)
	_,err = tx.Stage.Where(
		tx.Stage.ID.Eq(d.SourceID),
		tx.Stage.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(tx.Stage.Index.Value(d.DestinationIndex))
	if err != nil {
		return
	}
	_,err = tx.Stage.Where(
		tx.Stage.ID.Eq(d.DestinationID),
		tx.Stage.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(tx.Stage.Index.Value(d.SourceIndex))
	if err != nil {
		return
	}
	
	return tx.Commit()
}

func (r *stageRepo) CreateStage(req *common.RequestContext, d dto.StageData) (res model.Stage, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	fields := d.Fields
	res.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Stage.Save(&res)
	if err != nil {
		return
	}
	// if err = tx.Stage.InsertActivity(int64(res.ID), req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
	// 	return
	// }
	err = r.updateIndexes(tx, req.ActiveCompany.ID, d.Fields.EntityID, d.Fields.Index)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

// func(r *stageRepo) UpdateIndexStages(req *common.RequestContext,)

func (r *stageRepo) updateIndexes(tx *query.QueryTx, companyID int64, entityID int32, index int32) (err error) {
	_, err = tx.Stage.Where(
		tx.Stage.CompanyID.Eq(companyID),
		tx.Stage.EntityID.Eq(entityID),
		tx.Stage.Index.Gt(index),
	).UpdateSimple(tx.Stage.Index.Add(1))
	return
}

func (r *stageRepo) EditStage(req *common.RequestContext, d dto.StageData) (err error) {
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
	err = tx.Stage.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Stage{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	// err = tx.Stage.InsertActivity(int64(d.ID), req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	// if err != nil {
	// 	return
	// }
	return tx.Commit()
}

func (r *stageRepo) GetStages(req *common.RequestContext, d dto.StagesRequest) (res []dto.StageDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Stage
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.stageQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *stageRepo) stageQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
			e.id,e.name,e.entity_id,e.color,e.index,e.entity_id
			from stages as e 
		`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}
