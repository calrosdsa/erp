package customer_repo

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

type CustomerRepository interface {
	CreateCustomer(req *common.RequestContext, tx *query.QueryTx,
		i dto.CustomerData) (customer model.Customer, err error)
	GetCustomers(req *common.RequestContext, i *dto.RequestPaginationData) (dto.PaginationResult[[]dto.CustomerDto], error)
	GetCustomerDetail(req *common.RequestContext, id *dto.RequestEntity) (dto.ResultEntity[dto.CustomerDto], error)
	GetCustomerTypes() []string
	EditCustomer(tx *query.QueryTx,req *common.RequestContext, d dto.CustomerData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
}

type customerRepository struct {
	customerTypes []string
	convertor     helpers.ConvertorHelper
	dbHelper      db.DbHelper
	Q             *query.Query
}

func NewCustomerRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) CustomerRepository {
	return &customerRepository{
		customerTypes: []string{"company", "individual"},
		convertor:     helpers.Convertor,
		Q:             conn.GetQ(),
		dbHelper:      conn.GetDbHelper(),
	}
}
func (r *customerRepository) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	customerQ := r.Q.Customer
	id := r.convertor.StrtoInt(d.Body.PartyID)
	_, err = r.Q.Customer.WithContext(req.Ctx).Where(
		customerQ.CompanyID.Eq(req.ActiveCompany.ID),
		customerQ.Status.Eq(d.Body.CurrentState),
		customerQ.ID.Eq(id),
	).UpdateSimple(
		customerQ.Status.Value(nextState),
	)
	return
}

func (r *customerRepository) CreateCustomer(req *common.RequestContext, tx *query.QueryTx, i dto.CustomerData) (
	res model.Customer, err error) {
		id, err := tx.Customer.InsertParty(proto.PartyType_customer.String())
		if err != nil {
			return
		}
		fields := i.Fields
		res.ID = id
		res.CompanyID = req.ActiveCompany.ID
		if err = r.convertor.CopyStructData(fields, &res); err != nil {
			return
		}
	
	
		err = tx.WithContext(req.Ctx).Customer.Save(&res)
		if err != nil {
			return
		}
		if err = tx.Customer.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
			return
		}

	return
}

func (r *customerRepository) EditCustomer(tx *query.QueryTx,req *common.RequestContext, d dto.CustomerData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Customer.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Customer{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Customer.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

// func (r *customerRepository) createCustomer(req *common.RequestContext,
// 	tx *query.QueryTx,i dto.ContactData)(err error) {
// 	contact := model.Contact{}
// 	return
// }

func (r *customerRepository) GetCustomerDetail(req *common.RequestContext, i *dto.RequestEntity) (
	dto.ResultEntity[dto.CustomerDto], error) {
	var (
		err error
		res dto.ResultEntity[dto.CustomerDto]
	)
	customerQ := r.Q.Customer
	groupQ := r.Q.Group
	customerID := r.convertor.StrtoInt(i.ID)
	err = r.Q.Customer.WithContext(req.Ctx).Select(
		customerQ.ID, customerQ.UUID, customerQ.Name, customerQ.CreatedAt, customerQ.CustomerType,
		customerQ.Status,
		groupQ.ID.As("group_id"), groupQ.UUID.As("group_uuid"), groupQ.Name.As("group_name"),
	).LeftJoin(groupQ, groupQ.ID.EqCol(customerQ.GroupID)).Where(
		customerQ.ID.Eq(customerID),
		customerQ.CompanyID.Eq(req.ActiveCompany.ID),
	).Limit(1).Scan(&res.Entity)
	return res, err
}

func (r *customerRepository) GetCustomers(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.CustomerDto], error) {
	var (
		res   dto.PaginationResult[[]dto.CustomerDto]
		conds []gen.Condition
		order field.Expr
	)
	builder := r.Q.WithContext(req.Ctx).Customer
	customerQ := r.Q.Customer
	conds = append(conds, customerQ.CompanyID.Eq(req.ActiveCompany.ID))

	builder = builder.Where(conds...)

	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
	orderCol, ok := r.Q.Customer.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if i.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		fmt.Println("ORDER COLUMN EXIST")
		// User doesn't contains orderColStr
	} else {
		fmt.Println("Order column dosent exist")
	}

	customers, err := builder.Limit(limit).Offset(offset).Find()
	customerDtos := make([]dto.CustomerDto, len(customers))
	for i, customer := range customers {
		customerDtos[i] = dto.CustomerDtoFromModel(customer)
	}
	res.Total = total
	res.Results = customerDtos
	return res, err
}

func (r *customerRepository) GetCustomerTypes() []string {
	return r.customerTypes
}

func (r *customerRepository) validateCustomerType(t string) error {
	for _, customerType := range r.customerTypes {
		if customerType == t {
			return nil
		}
	}
	return domain.TYPE_NOT_FOUND
}
