package order_repo

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

	"gorm.io/gen/field"
	"gorm.io/gen/helper"
	"gorm.io/gorm/schema"
)

type OrderRepository interface {
	CreateOrderTx(req *common.RequestContext, tx *query.QueryTx, i dto.OrderBody) (
		model.Order, error,
	)
	GetOrder(req *common.RequestContext, i *dto.RequestEntityWithParty) (dto.ResultEntity[dto.OrderDetailDto], error)
	GetOrders(req *common.RequestContext, d *dto.RequestOrders) (
		dto.PaginationResult[[]dto.OrderDto], error)
	EditOrder(tx *query.QueryTx, req *common.RequestContext, d dto.OrderBody) (err error)
	UpdateOrderStatus(req *common.RequestContext, tx *query.QueryTx,
		id, prevState, nextState string) (model.Order, error)
	GetFilterOptions(partyType string) []dto.FilterOptionDto
}

const PO_CODE_TEMPLATE = "OC-#######"	
const SO_CODE_TEMPLATE = "OV-#######"	


type orderRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	currency  helpers.CurrencyHelper
	query     helpers.QueryHelper
	generator helpers.Generator
}

func NewOrderRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) OrderRepository {
	return &orderRepository{
		query:     helpers.Query,
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  conn.GetDbHelper(),
		currency:  helpers.Currency,
		generator: helpers.Generator,
	}
}

func (r *orderRepository) UpdateOrderStatus(req *common.RequestContext, tx *query.QueryTx,
	id, prevState, nextState string) (res model.Order, err error) {
	orderQ := tx.Order
	_, err = tx.Order.WithContext(req.Ctx).Where(
		orderQ.CompanyID.Eq(req.ActiveCompany.ID),
		orderQ.Status.Eq(prevState),
		orderQ.Code.Eq(id),
	).UpdateSimple(orderQ.Status.Value(nextState))
	if err != nil {
		return
	}
	order, err := tx.Order.WithContext(req.Ctx).Where(
		orderQ.CompanyID.Eq(req.ActiveCompany.ID),
		orderQ.Code.Eq(id),
	).First()
	if err != nil {
		return
	}

	return *order, err
}

func (s *orderRepository) CreateOrderTx(req *common.RequestContext, tx *query.QueryTx, i dto.OrderBody) (
	res model.Order, err error,
) {
	partyID, err := tx.WithContext(req.Ctx).Order.InsertParty(i.Order.OrderPartyType)
	if err != nil {
		return
	}
	res.CompanyID = req.ActiveCompany.ID
	res.ID = partyID
	fields := i.Order.Fields
	if err = s.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	count, err := tx.Order.WithContext(req.Ctx).Where(
		tx.Order.CompanyID.Eq(req.ActiveCompany.ID),
		tx.Party.PartyTypeCode.Eq(i.Order.OrderPartyType),
	).Join(
		tx.Party, tx.Party.ID.EqCol(tx.Order.ID),
	).Count()
	if err != nil {
		return
	}
	var orderTemplate  string
	if i.Order.OrderPartyType == proto.PartyType_purchaseOrder.String() {
		orderTemplate = PO_CODE_TEMPLATE
	}
	if i.Order.OrderPartyType == proto.PartyType_saleOrder.String() {
		orderTemplate = SO_CODE_TEMPLATE
	}
	res.Code, err = s.generator.GenerateCodeAutoIncrement(orderTemplate, count)
	if err != nil {
		return
	}

	err = tx.Order.WithContext(req.Ctx).Save(&res)
	if err != nil {
		return res, err
	}
	totalQuantity, totalAmount := s.calculateTotal(i.CreateItemLines,i.CreateTaxAndCharges)


	err = s.createProgressOrder(req.Ctx, tx, res.ID, totalQuantity, totalAmount)
	if err != nil {
		return res, err
	}
	//Create Reference to party
	references := i.Order.References
	references = append(references, &res.PartyID,
		res.ProjectID, res.CostCenterID,
	)
	err = s.dbHelper.InsertReferences(req.Ctx, tx, res.ID, references)
	if err != nil {
		return
	}

	return res, nil
}
func (r *orderRepository) createProgressOrder(ctx context.Context, tx *query.QueryTx,
	orderID int64, totalQuantity int, totalAmount float64) error {
	progressOrder := model.ProgressOrder{}
	progressOrder.OrderID = orderID
	progressOrder.TotalItems = int32(totalQuantity)
	progressOrder.TotalAmount = r.currency.FloatToInt(totalAmount)
	err := tx.ProgressOrder.WithContext(ctx).Save(&progressOrder)
	return err
	// progressOrder.TotalItems = d.Body
}

func (s *orderRepository) calculateTotal(
	d dto.CreateItemLines,c dto.CreateTaxAndChanges) (int, float64) {
	var totalQuantity int
	var totalAmount float64
	var totalAmountCharges float64
	for _, line := range d.Lines {
		totalQuantity += int(line.Quantity)
		totalAmount += line.Rate * float64(line.Quantity)
	}
	for _, line := range c.TaxAndCharges {
		if line.IsDeducted {
			totalAmount -= line.Amount
		}else {
			totalAmount += line.Amount

		}
	}
	return totalQuantity, totalAmount+totalAmountCharges
}

func (s *orderRepository) GetOrder(req *common.RequestContext, i *dto.RequestEntityWithParty) (
	dto.ResultEntity[dto.OrderDetailDto], error) {
	var (
		err            error
		res            dto.ResultEntity[dto.OrderDetailDto]
		partyTable     schema.Tabler
		partyJoinExprs []field.Expr
	)
	orderQ := s.Q.Order
	projectQ := s.Q.Project
	costCenterQ := s.Q.CostCenter
	priceListQ := s.Q.PriceList
	var columns []field.Expr
	columns = append(columns, orderQ.ID, orderQ.CreatedAt, orderQ.PostingDate, orderQ.PostingTime,
		orderQ.Tz, orderQ.DeliveryDate,
		orderQ.Code, orderQ.Currency, orderQ.Status, orderQ.PartyID,
		projectQ.Name.As("project"), projectQ.ID.As("project_id"), projectQ.UUID.As("project_uuid"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"), costCenterQ.UUID.As("cost_center_uuid"),
		priceListQ.Name.As("price_list"), priceListQ.ID.As("price_list_id"), priceListQ.UUID.As("price_list_uuid"),
	)

	switch i.PartyType {
	case proto.PartyType_purchaseOrder.String():
		supplierQ := s.Q.Supplier
		partyTable = supplierQ
		partyJoinExprs = append(partyJoinExprs, supplierQ.ID.EqCol(orderQ.PartyID))
		columns = append(columns, supplierQ.Name.As("party_name"),
			supplierQ.UUID.As("party_uuid"))
	case proto.PartyType_saleOrder.String():
		customer := s.Q.Customer
		partyTable = customer
		partyJoinExprs = append(partyJoinExprs, customer.ID.EqCol(orderQ.PartyID))
		columns = append(columns, customer.Name.As("party_name"),
			customer.UUID.As("party_uuid"))
	}
	s.Q.Order.Select(columns...).WithContext(req.Ctx).Where(
		s.Q.Order.Code.Eq(i.ID),
	).
		Join(partyTable, partyJoinExprs...).
		LeftJoin(projectQ, projectQ.ID.EqCol(orderQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(orderQ.CostCenterID)).
		LeftJoin(priceListQ, priceListQ.ID.EqCol(orderQ.PriceListID)).
		Scan(&res.Entity.Order)
	return res, err
}

func (r *orderRepository) EditOrder(tx *query.QueryTx, req *common.RequestContext, d dto.OrderBody) (err error) {
	data, err := r.convertor.DataMap(d.Order.Fields)
	if err != nil {
		return
	}
	err = tx.Order.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Order{ID: d.Order.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	totalQuantity, totalAmount := r.calculateTotal(d.CreateItemLines,d.CreateTaxAndCharges)
	_,err = tx.ProgressOrder.WithContext(req.Ctx).Where(
		tx.ProgressOrder.OrderID.Eq(d.Order.ID),
	).UpdateSimple(
		tx.ProgressOrder.TotalAmount.Value(r.currency.FloatToInt(totalAmount)),
		tx.ProgressOrder.TotalItems.Value(int32(totalQuantity)),
	)
	if err != nil {
		return
	}


	err = tx.Order.InsertActivity(d.Order.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	var references []*int64
	references = append(references,
		&d.Order.Fields.PartyID,
		d.Order.Fields.CostCenterID,
		d.Order.Fields.ProjectID,
	)
	if err = r.dbHelper.InsertReferences(req.Ctx, tx, d.Order.ID, references,true); err != nil {
		return
	}
	return
}

func (r *orderRepository) GetOrders(req *common.RequestContext, d *dto.RequestOrders) (
	res dto.PaginationResult[[]dto.OrderDto], err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Invoice
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.ordersQuery(req, queryData, &generateSQL, d.PartyType)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res.Results).Error
	return
}

func (r *orderRepository) ordersQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
	docPartyType string) (
	params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.code,e.created_at,e.status,e.posting_date,e.posting_time,e.tz,e.delivery_date,
		e.currency,
		p.party_type_code as party_type,party.name as party_name,party.uuid as party_uuid,
		po.total_items,po.received_items,po.total_amount,po.billed_amount
		from orders as e 
		join parties as invoice_party on invoice_party.id = e.id
		join parties as p on p.id = e.party_id
		join progress_orders as po on po.order_id = e.id
	`)
	if docPartyType == proto.PartyType_purchaseOrder.String() {
		generateSQL.WriteString(`join suppliers as party on party.id = e.party_id `)
	}
	if docPartyType == proto.PartyType_saleOrder.String() {
		generateSQL.WriteString(`join customers as party on party.id = e.party_id `)
	}
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? and invoice_party.party_type_code = ?`)
	params = append(params, req.ActiveCompany.ID, docPartyType)

	columnFilters := []string{
		"id",
		"project_id",
		"cost_center_id",
		"status",
		"code",
		"party_id",
		"posting_date",
		"delivery_date",
	}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	r.query.ReferenceFilterBuilder(generateSQL, &whereSQL, &params, d, "pricing_id")

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

// get entity that the group is related
func (h *orderRepository) GetEntityOrder(partyCode string) (domain.EntityTemplate, error) {
	switch partyCode {
	case domain.PARTY_PURCHASE_ORDER:
		return domain.PURCHASE_ORDER, nil
	default:
		return domain.EntityTemplate{}, domain.PARTY_TYPE_NOT_FOUND
	}
}

func (r *orderRepository) GetFilterOptions(partyType string) []dto.FilterOptionDto {

	filterOptions := []dto.FilterOptionDto{}
	statusList := []string{proto.State_DRAFT.String(), proto.State_COMPLETED.String(),
		proto.State_CANCELLED.String(), proto.State_TO_BILL.String(),
	}
	if partyType == proto.PartyType_purchaseOrder.String() {
		statusList = append(statusList, proto.State_TO_RECEIVE_AND_BILL.String(), proto.State_TO_RECEIVE.String())
	}
	if partyType == proto.PartyType_saleOrder.String() {
		statusList = append(statusList, proto.State_TO_DELIVER_AND_BILL.String(), proto.State_TO_DELIVER.String())
	}
	status := dto.FilterOptionDto{
		Param:     "status",
		Name:      "Estado",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		Options:   statusList,
	}
	dueDate := dto.FilterOptionDto{
		Name:      "Fecha de Vencimiento",
		Param:     "due_date",
		Type:      dto.FILTER_TYPE_DATE,
		Operators: dto.DateOperators,
	}

	postingDate := dto.FilterOptionDto{
		Name:      "Fecha de Publicacion",
		Param:     "posting_date",
		Type:      dto.FILTER_TYPE_DATE,
		Operators: dto.DateOperators,
	}

	code := dto.FilterOptionDto{
		Name:      "ID",
		Param:     "code",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
	}

	pricing := dto.FilterOptionDto{
		Name:      "Pricing",
		Param:     "pricing_id",
		Type:      dto.FILTER_TYPE_STRING,
		Operators: dto.StringOperators,
		PartyType: proto.PartyType_pricing.String(),
	}

	filterOptions = append(filterOptions, status, dueDate, postingDate, code, pricing)
	return filterOptions
}
