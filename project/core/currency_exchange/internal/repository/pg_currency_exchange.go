package currency_exchange_repo

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

type CurrencyExchangeRepo interface {
	GetCurrencyExchange(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.CurrencyExchangeDto], err error)
	CreateCurrencyExchange(req *common.RequestContext, d *dto.CreateCurrencyExchangeRequest) (
		res model.CurrencyExchange, err error)
	GetCurrencyExchanges(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.CurrencyExchangeDto], err error)
	EditCurrencyExchange(req *common.RequestContext, d *dto.EditCurrencyExchangeRequest) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent, nextState string) (err error)
}

type currencyExchangeRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	currency  helpers.CurrencyHelper
}

func NewCurrencyExchangeRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) CurrencyExchangeRepo {
	return &currencyExchangeRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
		currency:  helpers.Currency,
	}
}
func (r *currencyExchangeRepo) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,
	nextState string) (err error) {
	cExchangeQ := r.Q.CurrencyExchange
	_, err = r.Q.CurrencyExchange.WithContext(req.Ctx).Where(
		cExchangeQ.CompanyID.Eq(req.ActiveCompany.ID),
		cExchangeQ.Status.Eq(d.Body.CurrentState),
		cExchangeQ.UUID.Eq(d.Body.PartyID),
	).UpdateSimple(
		cExchangeQ.Status.Value(nextState),
	)
	return
}

func (r *currencyExchangeRepo) GetCurrencyExchange(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.CurrencyExchangeDto], err error) {
	id := r.convertor.StrtoInt(d.ID)
	cExchangeQ := r.Q.CurrencyExchange
	err = cExchangeQ.WithContext(req.Ctx).Select(
		cExchangeQ.ID, cExchangeQ.UUID, cExchangeQ.Status,
		cExchangeQ.Name, cExchangeQ.FromCurrency, cExchangeQ.ToCurrency,
		cExchangeQ.ForBuying, cExchangeQ.ForSelling, cExchangeQ.ExchangeRate,
	).
		Where(cExchangeQ.CompanyID.Eq(req.ActiveCompany.ID), cExchangeQ.ID.Eq(id)).
		Scan(&res.Entity)
	return
}

func (r *currencyExchangeRepo) CreateCurrencyExchange(req *common.RequestContext, d *dto.CreateCurrencyExchangeRequest) (
	res model.CurrencyExchange, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	res.Name = d.Body.Name
	res.CompanyID = req.ActiveCompany.ID
	err = r.dbHelper.ValidateName(tx.CurrencyExchange.UnderlyingDB(), &res)
	if err != nil {
		return res, domain.ERROR_NAME_TAKEN
	}
	currencyExchangeID, err := tx.CurrencyExchange.InsertParty(proto.PartyType_currencyExchange.String())
	if err != nil {
		return
	}
	res.ID = currencyExchangeID
	res.Status = proto.State_ENABLED.String()
	res.FromCurrency = d.Body.FromCurrency
	res.ToCurrency = d.Body.ToCurrency
	res.ForBuying = d.Body.ForBuying
	res.ForSelling = d.Body.ForSelling
	res.ExchangeRate = int32(r.currency.FloatToInt(d.Body.ExchangeRate))

	err = tx.CurrencyExchange.WithContext(req.Ctx).Save(&res)
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

func (r *currencyExchangeRepo) EditCurrencyExchange(req *common.RequestContext, d *dto.EditCurrencyExchangeRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	cExchangeQ := tx.CurrencyExchange
	var columns []field.AssignExpr

	columns = append(columns, cExchangeQ.Name.Value(d.Body.Name), cExchangeQ.FromCurrency.Value(d.Body.FromCurrency),
		cExchangeQ.ToCurrency.Value(d.Body.ToCurrency), cExchangeQ.ForBuying.Value(d.Body.ForBuying),
		cExchangeQ.ForSelling.Value(d.Body.ForSelling),
		cExchangeQ.ExchangeRate.Value(int32(r.currency.FloatToInt(d.Body.ExchangeRate))),
	)
	_, err = tx.CurrencyExchange.WithContext(req.Ctx).Where(
		cExchangeQ.ID.Eq(d.Body.ID), cExchangeQ.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}
	err = tx.CurrencyExchange.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *currencyExchangeRepo) GetCurrencyExchanges(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CurrencyExchangeDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	cExchangeQ := r.Q.CurrencyExchange
	builder := r.Q.WithContext(req.Ctx).CurrencyExchange

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.CurrencyExchange.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	conds = append(conds, cExchangeQ.CompanyID.Eq(req.ActiveCompany.ID))
	fmt.Println("QUERY", d.Query)
	if d.Query != "" {
		conds = append(conds, cExchangeQ.Name.Like("%"+d.Query+"%"))
	}

	builder = builder.Select(
		cExchangeQ.ID, cExchangeQ.UUID, cExchangeQ.Name, cExchangeQ.Status, cExchangeQ.CreatedAt,
		cExchangeQ.FromCurrency, cExchangeQ.ToCurrency,cExchangeQ.ExchangeRate,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
