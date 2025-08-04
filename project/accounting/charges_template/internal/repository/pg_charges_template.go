package charges_template_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type ChargesTemplateRepository interface {
	GetChargesTemplate(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ChargesTemplateDto], err error)
	CreateChargesTemplate(tx *query.QueryTx, req *common.RequestContext, d *dto.CreateChargesTemplateRequest) (
		res model.ChargesTemplate, err error)
	GetChargesTemplates(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.ChargesTemplateDto], err error)
	EditChargesTemplate(req *common.RequestContext, d *dto.EditChargesTemplateRequest) (err error)
}

type chargesTemplateRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
}

func NewChargesTemplateRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) ChargesTemplateRepository {
	return &chargesTemplateRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
	}
}

func (r *chargesTemplateRepo) GetChargesTemplate(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ChargesTemplateDto], err error) {
	id := r.convertor.StrtoInt(d.ID)
	chargesTemplateQ := r.Q.ChargesTemplate
	err = chargesTemplateQ.WithContext(req.Ctx).Select(
		chargesTemplateQ.ID, chargesTemplateQ.UUID, chargesTemplateQ.Status,
		chargesTemplateQ.Name,
	).
		Where(
			chargesTemplateQ.CompanyID.Eq(req.ActiveCompany.ID),
			chargesTemplateQ.ID.Eq(id),
		).
		Scan(&res.Entity)
	return
}

func (r *chargesTemplateRepo) CreateChargesTemplate(tx *query.QueryTx, req *common.RequestContext, d *dto.CreateChargesTemplateRequest) (
	res model.ChargesTemplate, err error) {
	res.Name = d.Body.ChargesTemplate.Name
	res.CompanyID = req.ActiveCompany.ID
	err = r.dbHelper.ValidateName(tx.ChargesTemplate.UnderlyingDB(), &res)
	if err != nil {
		return res, domain.ERROR_NAME_TAKEN
	}
	chargesTemplateID, err := tx.ChargesTemplate.InsertParty(proto.PartyType_chargesTemplate.String())
	if err != nil {
		return
	}
	res.ID = chargesTemplateID
	res.Status = proto.State_ENABLED.String()
	err = tx.ChargesTemplate.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return
	}
	return
}

func (r *chargesTemplateRepo) EditChargesTemplate(req *common.RequestContext, d *dto.EditChargesTemplateRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	chargesTemplateQ := tx.ChargesTemplate
	var columns []field.AssignExpr

	columns = append(columns, chargesTemplateQ.Name.Value(d.Body.Name))
	_, err = tx.ChargesTemplate.WithContext(req.Ctx).Where(
		chargesTemplateQ.ID.Eq(d.Body.ID), chargesTemplateQ.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}
	err = tx.ChargesTemplate.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *chargesTemplateRepo) GetChargesTemplates(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.ChargesTemplateDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	chargesTemplateQ := r.Q.ChargesTemplate
	builder := r.Q.WithContext(req.Ctx).ChargesTemplate

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.ChargesTemplate.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	conds = append(conds, chargesTemplateQ.CompanyID.Eq(req.ActiveCompany.ID))
	fmt.Println("QUERY", d.Query)
	if d.Query != "" {
		conds = append(conds, chargesTemplateQ.Name.Like("%"+d.Query+"%"))
	}

	builder = builder.Select(
		chargesTemplateQ.ID, chargesTemplateQ.UUID, chargesTemplateQ.Name, chargesTemplateQ.Status, chargesTemplateQ.CreatedAt,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
