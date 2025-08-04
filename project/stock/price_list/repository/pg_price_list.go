package price_list_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"erp/pkg/logger"
	"fmt"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type PriceListRepository interface {
	CreatePriceList(req *common.RequestContext, i *dto.CreatePriceListRequest) (res model.PriceList, err error)
	CreatePriceListTx(tx *query.QueryTx, req *common.RequestContext, i *dto.CreatePriceListRequest) (
		res model.PriceList, err error)
	GetListPriceDetail(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.PriceListDto], err error)
	GetPriceLists(req *common.RequestContext, d *dto.RequestPriceLists) (
		res dto.PaginationResult[[]dto.PriceListDto], err error)
	EditPriceList(req *common.RequestContext, d *dto.EditPriceListRequest) (err error)
}

type priceListRepository struct {
	Q         *query.Query
	dbHelper  db.DbHelper
	convertor helpers.ConvertorHelper
}

func NewPriceListServer(
	db db.Connection,
	helpers *helpers.Helpers,
	logger logger.Logger,
) *priceListRepository {
	return &priceListRepository{
		dbHelper:  db.GetDbHelper(),
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
	}
}

func (r *priceListRepository) EditPriceList(req *common.RequestContext, d *dto.EditPriceListRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	priceListQ := tx.PriceList
	var columns []field.AssignExpr

	columns = append(columns, priceListQ.Name.Value(d.Body.Name),
		priceListQ.Currency.Value(d.Body.Currency), priceListQ.IsBuying.Value(d.Body.IsBuying),
		priceListQ.IsSelling.Value(d.Body.IsSelling),
	)
	_, err = tx.PriceList.WithContext(req.Ctx).Where(
		priceListQ.ID.Eq(d.Body.ID), priceListQ.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}
	err = tx.PriceList.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (s *priceListRepository) CreatePriceList(req *common.RequestContext, i *dto.CreatePriceListRequest) (
	model.PriceList, error) {
	tx := s.Q.Begin()
	var err error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	priceList, err := s.createPriceList(tx, req, i)
	err = tx.Commit()
	return priceList, err
}

func (s *priceListRepository) CreatePriceListTx(tx *query.QueryTx, req *common.RequestContext, i *dto.CreatePriceListRequest) (
	model.PriceList, error) {
	priceList, err := s.createPriceList(tx, req, i)
	return priceList, err
}

func (s *priceListRepository) createPriceList(tx *query.QueryTx, req *common.RequestContext, i *dto.CreatePriceListRequest) (
	model.PriceList, error) {
	var priceList model.PriceList
	priceList.Name = i.Body.Name
	priceList.CompanyID = req.ActiveCompany.ID
	priceListID, err := tx.PriceList.InsertParty(proto.PartyType_priceList.String())
	if err != nil {
		return priceList, err
	}
	priceList.ID = priceListID
	priceList.Currency = i.Body.Currency
	priceList.IsBuying = i.Body.IsBuying
	priceList.IsSelling = i.Body.IsSelling
	priceList.Status = proto.State_ENABLED.String()
	err = tx.PriceList.WithContext(req.Ctx).Save(&priceList)
	if err != nil {
		return priceList, err
	}
	err = tx.PriceList.InsertActivity(priceListID, req.Profile.ID, proto.ActivityType_CREATE.String(), nil)
	if err != nil {
		return priceList, err
	}
	return priceList, err
}

func (r *priceListRepository) GetListPriceDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.PriceListDto], err error) {
	id := r.convertor.StrtoInt(i.ID)
	priceListQ := r.Q.PriceList
	err = priceListQ.WithContext(req.Ctx).Select(
		priceListQ.ID, priceListQ.Name, priceListQ.Status, priceListQ.Currency,
		priceListQ.IsBuying, priceListQ.IsSelling,
	).
		Where(
			priceListQ.CompanyID.Eq(req.ActiveCompany.ID),
			priceListQ.ID.Eq(id),
		).
		Scan(&res.Entity)
	return
}

func (r *priceListRepository) GetPriceLists(req *common.RequestContext, d *dto.RequestPriceLists) (
	res dto.PaginationResult[[]dto.PriceListDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	priceListQ := r.Q.PriceList
	builder := r.Q.WithContext(req.Ctx).PriceList

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.PriceList.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	fmt.Println("QUERY D", d)
	conds = append(conds, priceListQ.CompanyID.Eq(req.ActiveCompany.ID))
	if d.Query != "" {
		conds = append(conds, priceListQ.Name.Like("%"+d.Query+"%"))
	}
	if d.IsBuying != "" {
		conds = append(conds, priceListQ.IsBuying.Is(r.convertor.StrToBool(d.IsBuying)))
	}
	if d.IsSelling != "" {
		conds = append(conds, priceListQ.IsSelling.Is(r.convertor.StrToBool(d.IsSelling)))
	}

	builder = builder.Select(
		priceListQ.ID, priceListQ.Name, priceListQ.Status, priceListQ.Currency,
		priceListQ.IsBuying, priceListQ.IsSelling,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}

func (s *priceListRepository) UpsertPriceList(req *common.RequestContext, d *dto.UpsertPriceListRequest) error {

	d.Body.ItemPriceList.CompanyID = req.ActiveCompany.ID
	err := s.Q.PriceList.WithContext(req.Ctx).Save(&d.Body.ItemPriceList)
	return err
}
