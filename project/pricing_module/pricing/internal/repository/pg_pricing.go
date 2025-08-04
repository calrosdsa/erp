package pricing_repo

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

type PricingRepository interface {
	GetPricing(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.PricingDetailDto], err error)
	CreatePricing(req *common.RequestContext, d *dto.CreatePricingRequest) (
		res model.Pricing, err error)
	GetPricings(req *common.RequestContext, d *dto.RequestPricings) (
		res dto.PaginationResult[[]dto.PricingDto], err error)
	EditPricing(req *common.RequestContext, d *dto.EditPricingRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
}

type pricingRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	query     helpers.QueryHelper
	generator helpers.Generator
	currency  helpers.CurrencyHelper
}

const PRICING_CODE_TEMPLATE = "PR-#######"

func NewPricingRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) PricingRepository {
	return &pricingRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
		generator: helpers.Generator,
		currency:  helpers.Currency,
		query:     helpers.Query,
	}
}

func (r *pricingRepository) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	pricingQ := r.Q.Pricing
	_, err = r.Q.Pricing.WithContext(req.Ctx).Where(
		pricingQ.CompanyID.Eq(req.ActiveCompany.ID),
		pricingQ.Status.Eq(d.Body.CurrentState),
		pricingQ.Code.Eq(d.Body.PartyID),
	).UpdateSimple(
		pricingQ.Status.Value(nextState),
	)
	return
}

func (r *pricingRepository) GetPricing(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.PricingDetailDto], err error) {
	pricingQ := r.Q.Pricing
	customerQ := r.Q.Customer
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	err = pricingQ.WithContext(req.Ctx).Select(
		pricingQ.ID, pricingQ.Code, pricingQ.Status,
		pricingQ.CreatedAt,
		pricingQ.CustomerID, customerQ.Name.As("customer"), customerQ.UUID.As("customer_uuid"),
		projectQ.Name.As("project"), projectQ.ID.As("project_id"), projectQ.UUID.As("project_uuid"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"), costCenterQ.UUID.As("cost_center_uuid"),
	).
		Where(pricingQ.CompanyID.Eq(req.ActiveCompany.ID), pricingQ.Code.Eq(d.ID)).
		LeftJoin(customerQ, customerQ.ID.EqCol(pricingQ.CustomerID)).
		LeftJoin(projectQ, projectQ.ID.EqCol(pricingQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(pricingQ.CostCenterID)).
		Scan(&res.Entity.PricingDto)
	if err != nil {
		return
	}
	prLineQ := r.Q.PricingLineItem
	supplierQ := r.Q.Supplier
	err = prLineQ.WithContext(req.Ctx).Select(
		prLineQ.ID, prLineQ.Description, prLineQ.PartNumber, prLineQ.PlUnit, prLineQ.Quantity,
		prLineQ.SupplierID, supplierQ.Name.As("supplier"), prLineQ.FobUnitFn,
		prLineQ.RetentionFn, prLineQ.CostZfFn, prLineQ.CostAlmFn,
		prLineQ.TvaFn, prLineQ.CantidadFn, prLineQ.TvaFn, prLineQ.PrecioUnitarioFn,
		prLineQ.PrecioTotalFn, prLineQ.PrecioUnitarioTcFn, prLineQ.TvaFn, prLineQ.PrecioTotalTcFn,
		prLineQ.FobTotalFn, prLineQ.GplTotalFn, prLineQ.TvaFn, prLineQ.TvaTotalFn,
		prLineQ.IsTitle, prLineQ.Color,
	).
		LeftJoin(supplierQ, supplierQ.ID.EqCol(prLineQ.SupplierID)).
		Where(prLineQ.PricingID.Eq(res.Entity.PricingDto.ID)).
		Order(prLineQ.ID.Asc()).
		Scan(&res.Entity.PricingLineItems)
	if err != nil {
		return
	}
	chargeLineQ := r.Q.PricingCharge
	err = chargeLineQ.WithContext(req.Ctx).Select(
		chargeLineQ.ID, chargeLineQ.Name, chargeLineQ.Rate,
	).
		Where(chargeLineQ.PricingID.Eq(res.Entity.PricingDto.ID)).
		Scan(&res.Entity.PricingCharges)
	if err != nil {
		return
	}
	return
}

func (r *pricingRepository) CreatePricing(req *common.RequestContext, d *dto.CreatePricingRequest) (
	res model.Pricing, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	pricingQ := tx.Pricing
	count, err := pricingQ.WithContext(req.Ctx).Where(
		pricingQ.CompanyID.Eq(req.ActiveCompany.ID),
	).Count()
	if err != nil {
		return
	}
	res.Code, err = r.generator.GenerateCodeAutoIncrement(PRICING_CODE_TEMPLATE, count)
	if err != nil {
		return
	}
	res.CustomerID = d.Body.PricingFields.CustomerID
	pricingID, err := pricingQ.InsertParty(proto.PartyType_pricing.String())
	res.CompanyID = req.ActiveCompany.ID
	res.Status = proto.State_DRAFT.String()
	res.ProjectID = d.Body.PricingFields.ProjectID
	res.CostCenterID = d.Body.PricingFields.CostCenterID
	res.ID = pricingID
	err = tx.Pricing.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return
	}
	err = r.createPricingLineItems(tx, req.Ctx, d.Body.PricingLineItems, pricingID)
	if err != nil {
		return
	}
	err = r.createPricingCharges(tx, req.Ctx, d.Body.PricingCharges, pricingID)
	if err != nil {
		return
	}
	err = tx.SalesRecord.InsertActivity(res.ID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *pricingRepository) createPricingCharges(tx *query.QueryTx, ctx context.Context, lines []dto.PricingChargeData,
	pricingID int64) (err error) {
	pricingCharges := make([]*model.PricingCharge, len(lines))
	for i, line := range lines {
		pricingCharge := &model.PricingCharge{}
		pricingCharge.Name = line.Name
		pricingCharge.Rate = int32(r.currency.FloatToInt(line.Rate))
		pricingCharge.PricingID = pricingID
		pricingCharges[i] = pricingCharge
	}
	err = tx.WithContext(ctx).PricingCharge.CreateInBatches(pricingCharges, len(pricingCharges))
	return
}

func (r *pricingRepository) createPricingLineItems(tx *query.QueryTx, ctx context.Context, lines []dto.PricingLineItemData,
	pricingID int64) (err error) {
	pricingLines := make([]*model.PricingLineItem, len(lines))
	fmt.Println("PRICING LINES", len(lines))
	for i, line := range lines {
		pricingLine := &model.PricingLineItem{}
		fmt.Println("LINE", pricingLine, line)
		pricingLine.PricingID = pricingID
		pricingLine.Description = line.Description
		pricingLine.PartNumber = line.PartNumber
		pricingLine.Quantity = line.Quantity
		pricingLine.SupplierID = line.SupplierID
		if line.PlUnit != nil {
			plUnit := int32(r.currency.FloatToInt(*line.PlUnit))
			pricingLine.PlUnit = &plUnit
		}

		pricingLine.FobUnitFn = line.FobUnitFn
		pricingLine.RetentionFn = line.RetentionFn
		pricingLine.CostZfFn = line.CostZfFn
		pricingLine.CostAlmFn = line.CostAlmFn
		pricingLine.TvaFn = line.TvaFn
		pricingLine.CantidadFn = line.CantidadFn
		pricingLine.PrecioUnitarioFn = line.PrecioUnitarioFn
		pricingLine.PrecioTotalFn = line.PrecioTotalFn
		pricingLine.PrecioUnitarioTcFn = line.PrecioUnitarioTcFn
		pricingLine.PrecioTotalTcFn = line.PrecioTotalTcFn
		pricingLine.FobTotalFn = line.FobTotalFn
		pricingLine.GplTotalFn = line.GplTotalFn
		pricingLine.TvaTotalFn = line.TvaTotalFn
		pricingLine.IsTitle = line.IsTitle
		pricingLine.Color = line.Color

		fmt.Println("LINE", pricingLine)

		pricingLines[i] = pricingLine
	}
	err = tx.WithContext(ctx).PricingLineItem.CreateInBatches(pricingLines, len(pricingLines))
	return
}

func (r *pricingRepository) EditPricing(req *common.RequestContext, d *dto.EditPricingRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Body.PricingFields)
	if err != nil {
		return
	}
	fmt.Println("PRICING DATA", data)
	err = tx.Pricing.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Pricing{ID: d.Body.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	_, err = tx.PricingLineItem.Unscoped().WithContext(req.Ctx).Where(
		tx.PricingLineItem.PricingID.Eq(d.Body.ID),
	).Delete()
	if err != nil {
		return
	}
	_, err = tx.PricingCharge.Unscoped().WithContext(req.Ctx).Where(
		tx.PricingCharge.PricingID.Eq(d.Body.ID),
	).Delete()
	if err != nil {
		return
	}
	err = r.createPricingLineItems(tx, req.Ctx, d.Body.PricingLineItems, d.Body.ID)
	if err != nil {
		return
	}
	err = r.createPricingCharges(tx, req.Ctx, d.Body.PricingCharges, d.Body.ID)
	if err != nil {
		return
	}
	err = tx.SalesRecord.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *pricingRepository) GetPricings(req *common.RequestContext, d *dto.RequestPricings) (
	res dto.PaginationResult[[]dto.PricingDto], err error) {
	var (
		generateSQL strings.Builder
	)
	queryData := r.convertor.GenerateQueryMap(d)
	builder := r.Q.Pricing.WithContext(req.Ctx)
	params := r.pricingQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res.Results).Error
	return
}

func (r *pricingRepository) pricingQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		t.id,t.created_at,t.code,t.status
		from pricings as t
	`)
	whereSQL.WriteString(` t.deleted_at is null and t.company_id = ?`)
	params = append(params, req.ActiveCompany.ID)
	if code, ok := d["code"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(code, "t.code", &params))
	}
	if status, ok := d["status"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(status, "t.status", &params))
	}
	if id, ok := d["id"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(id, "t.id", &params))
	}
	if status, ok := d["created_at"]; ok {
		whereSQL.WriteString(r.convertor.GetConditionFromQuery(status, "t.created_at", &params))
	}

	// if invoiceID, ok := d["invoice_id"]; ok {
	// 	whereSQL.WriteString(r.convertor.GetConditionFromQuery(invoiceID, "invoice_id", &params))
	// }
	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}
