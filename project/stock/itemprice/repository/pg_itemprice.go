package itemprice_repo

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
	"strings"

	"github.com/samber/lo"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"
	"gorm.io/gorm"
)

type ItemPriceRepository interface {
	GetItemPrice(req *common.RequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.ItemPriceDto], error)
	GetItemPrices(req *common.RequestContext, i dto.ItemPricesRequest) (
		res []dto.ItemPriceDto, err error)
	GetItemPricesByItemCode(req *common.RequestContext, d *dto.RequestItemPriceByCode) (
		res dto.PaginationResult[[]dto.ItemPriceDto], err error)
	GetItemPricesForOrders(req *common.RequestContext, i *dto.RequestItemPricesForOrder) (
		dto.ResultEntity[[]dto.ItemPriceDto], error)
	UpsertItemPrice(req *common.RequestContext, d *dto.UpsertItemPriceRequest) error
	CreateItemPrice(req *common.RequestContext, d dto.ItemPriceData) (res dto.ItemPriceDto, err error)
	CreateItemPriceTx(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (
		res dto.ItemPriceDto, err error)
	EditItemPrice(req *common.RequestContext, d dto.ItemPriceData) (err error)
	EditItemPriceTx(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (err error)
	DeleteItemPriceTx(tx *query.QueryTx, req *common.RequestContext, id int64) (err error)
}

type itemPriceRepository struct {
	Q              *query.Query
	DB             *gorm.DB
	convertor      helpers.ConvertorHelper
	currencyHelper helpers.CurrencyHelper
	query          helpers.QueryHelper
}

func NewItemPriceService(
	conn db.Connection,
	helpers *helpers.Helpers,
) ItemPriceRepository {
	return &itemPriceRepository{
		currencyHelper: helpers.Currency,
		Q:              conn.GetQ(),
		DB:             conn.GetDB(),
		convertor:      helpers.Convertor,
		query:          helpers.Query,
	}
}

func (r *itemPriceRepository) DeleteItemPriceTx(tx *query.QueryTx, req *common.RequestContext,
	id int64) (err error) {
	e := tx.ItemPrice
	e.WithContext(req.Ctx).Where(
		e.ID.Eq(id),
		e.CompanyID.Eq(req.ActiveCompany.ID),
	).Delete()
	return
}

func (r *itemPriceRepository) EditItemPrice(req *common.RequestContext, d dto.ItemPriceData) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = r.editItemPrice(tx, req, d)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *itemPriceRepository) EditItemPriceTx(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (err error) {
	err = r.editItemPrice(tx, req, d)
	return
}

func (s *itemPriceRepository) CreateItemPriceTx(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (
	res dto.ItemPriceDto, err error) {
	res, err = s.createItemPrice(tx, req, d)
	return res, err
}

func (s *itemPriceRepository) CreateItemPrice(req *common.RequestContext, d dto.ItemPriceData) (
	res dto.ItemPriceDto, err error) {
	tx := s.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	res, err = s.createItemPrice(tx, req, d)
	if err != nil {
		return
	}
	err = tx.Commit()
	if err != nil {
		return
	}
	return res, err
}

func (r *itemPriceRepository) editItemPrice(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	fmt.Println("ITEM PRICE UNIT OF MEASURE ID:",d.Fields.UnitOfMeasureID)

	err = tx.ItemPrice.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.ItemPrice{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.ItemPrice.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

func (r *itemPriceRepository) createItemPrice(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (
	res dto.ItemPriceDto, err error) {
	var itemPrice model.ItemPrice
	id, err := tx.ItemPrice.InsertParty(proto.PartyType_itemPrice.String())
	if err != nil {
		return
	}
	fields := d.Fields
	itemPrice.ID = id
	itemPrice.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &itemPrice); err != nil {
		return
	}
	err = tx.WithContext(req.Ctx).ItemPrice.Save(&itemPrice)
	if err != nil {
		return
	}
	if err = tx.ItemPrice.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.ItemPriceDtoFromModel(&itemPrice)
	return
}

func (s *itemPriceRepository) GetItemPrice(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ItemPriceDto], err error) {
	id := s.convertor.StrtoInt(d.ID)
	itemPriceQ := s.Q.ItemPrice
	itemQ := s.Q.Item
	priceListQ := s.Q.PriceList
	// uomQ := s.Q.UnitOfMeasure
	uomTranslateQ := s.Q.UnitOfMeasureTranslation

	err = s.Q.WithContext(req.Ctx).ItemPrice.
		Select(itemPriceQ.ID, itemPriceQ.UUID, itemPriceQ.ItemQuantity, itemPriceQ.Rate,
			itemQ.Name.As("item_name"), itemQ.UUID.As("item_uuid"), itemQ.Code.As("item_code"),
			itemQ.ID.As("item_id"),
			itemPriceQ.UnitOfMeasureID.As("uom_id"), uomTranslateQ.Name.As("uom"),
			priceListQ.Name.As("price_list_name"), priceListQ.UUID.As("price_list_uuid"),
			priceListQ.Currency.As("price_list_currency"), priceListQ.ID.As("price_list_id"),
		).
		Join(itemQ, itemQ.ID.EqCol(itemPriceQ.ItemID)).
		Join(uomTranslateQ, uomTranslateQ.BaseID.EqCol(itemPriceQ.UnitOfMeasureID)).
		LeftJoin(priceListQ, priceListQ.ID.EqCol(itemPriceQ.PriceListID)).
		Where(
			itemPriceQ.ID.Eq(id),
			itemPriceQ.CompanyID.Eq(req.ActiveCompany.ID),
		).Limit(1).Scan(&res.Entity)
	return res, err
}

func (r *itemPriceRepository) GetItemPrices(req *common.RequestContext, d dto.ItemPricesRequest) (
	res []dto.ItemPriceDto, err error,
) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Deal
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.itemPricesQuery(req, queryData, &generateSQL)

	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *itemPriceRepository) itemPricesQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
					e.id,e.item_id,e.created_at,e.item_quantity,e.rate,
					e.item_id,it.name as item_name,
					e.price_list_id,pl.name as price_list_name,pl.currency as price_list_currency,
					uom.base_id as item_price_uom_id,uom.name as item_price_uom
					from item_prices as e 
					join price_lists as pl on pl.id = e.price_list_id
					join items as it on it.id = e.item_id
					left join unit_of_measure_translations as uom on 
					uom.base_id = e.unit_of_measure_id and uom.language_code = ?
		`)
	//Add first param for filter laguage code on uoms
	params = append(params, req.LanguageCode)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{
		"item_id",
		"price_list_id",
	}

	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}

// func (r *itemPriceRepository) GetItemPrices(req *common.RequestContext, i *dto.RequestPaginationData) (
// 	res dto.PaginationResult[[]dto.ItemPriceDto], err error) {
// 	var (
// 		conds []gen.Condition
// 		order field.Expr
// 	)
// 	itemPriceQ := r.Q.ItemPrice
// 	builder := r.Q.WithContext(req.Ctx).ItemPrice

// 	//ADDING CONDITIONS
// 	conds = append(conds, itemPriceQ.CompanyID.Eq(req.ActiveCompany.ID))
// 	// if i.Query  != "" {
// 	// 	conds = append(conds, itemPriceQ..Like("%"+i.Query+"%"))
// 	// }

// 	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
// 	orderCol, ok := r.Q.ItemPrice.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
// 	if ok {
// 		if i.Order == "ASC" {
// 			order = orderCol.Asc()
// 		} else {
// 			order = orderCol.Desc()
// 		}
// 		builder = builder.Order(order)
// 		// User doesn't contains orderColStr
// 	}
// 	itemQ := r.Q.Item
// 	priceListQ := r.Q.PriceList
// 	builder = builder.Select(
// 		itemPriceQ.ID, itemPriceQ.UUID, itemPriceQ.Rate, itemPriceQ.ItemQuantity, itemPriceQ.CreatedAt,
// 		itemPriceQ.UnitOfMeasureID.As("uom_id"),
// 		itemQ.Name.As("item_name"), itemQ.UUID.As("item_uuid"), itemQ.Pn.As("item_code"),
// 		priceListQ.Name.As("price_list_name"), priceListQ.UUID.As("price_list_uuid"),
// 		priceListQ.Currency.As("price_list_currency"),
// 	).
// 		Join(itemQ, itemPriceQ.ItemID.EqCol(itemQ.ID)).
// 		LeftJoin(priceListQ, itemPriceQ.PriceListID.EqCol(priceListQ.ID)).
// 		Where(conds...)
// 	total, err := builder.Count()
// 	if err != nil {
// 		return res, err
// 	}
// 	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
// 	res.Total = total
// 	return
// 	// var result dto.PaginationResult[[]dto.ItemPriceDto]
// 	// var itemPrices []model.ItemPrice
// 	// queryBuilder := s.conn.Db.WithContext(req.Ctx).Model(&model.ItemPrice{}).
// 	// 	Where(&model.ItemPrice{CompanyID: req.ActiveCompany.ID})

// 	// err := queryBuilder.
// 	// 	Count(&result.Total).Error

// 	// if d.Query != "" {
// 	// 	queryBuilder = queryBuilder.Where("code ILIKE ?", "%"+d.Query+"%")
// 	// }

// 	// err = queryBuilder.
// 	// 	Scopes(s.conn.Paginate(req.Params)).
// 	// 	Preload("Item", s.conn.PreloadColumns([]string{"ID", "Code", "Name"})).
// 	// 	Preload("PriceList", s.conn.PreloadColumns([]string{"ID", "Currency"})).
// 	// 	Find(&itemPrices).Error
// 	// if err != nil {
// 	// 	return result, err
// 	// }
// 	// itemPriceDtos := make([]dto.ItemPriceDto, len(itemPrices))
// 	// for i, itemPrice := range itemPrices {
// 	// 	itemPriceDtos[i] = dto.ItemPriceDtoFromModel(&itemPrice)
// 	// }
// 	// result.Results = itemPriceDtos
// 	// return result, err
// }

func (s *itemPriceRepository) GetItemPricesByItemCode(req *common.RequestContext, d *dto.RequestItemPriceByCode) (
	res dto.PaginationResult[[]dto.ItemPriceDto], err error) {

	// var itemPrices []model.ItemPrice
	// itemID := s.convertor.StrtoInt(d.ItemID)
	// queryBuilder := s.DB.WithContext(req.Ctx).Model(&model.ItemPrice{}).
	// 	Where(&model.ItemPrice{CompanyID: req.ActiveCompany.ID, ItemID: itemID})

	// err = queryBuilder.Count(&res.Total).Error

	// if d.Query != "" {
	// 	queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
	// }

	// err = queryBuilder.
	// 	Scopes(s.conn.Paginate(req.Params)).
	// 	Preload("PriceList").
	// 	Find(&itemPrices).Error
	// if err != nil {
	// 	return res, err
	// }
	// itemPricesDtos := make([]dto.ItemPriceDto, len(itemPrices))

	// for i, itemPrice := range itemPrices {
	// 	itemPricesDtos[i] = dto.ItemPriceDtoFromModel(&itemPrice)
	// }
	// res.Results = itemPricesDtos
	return res, err
}

func (s *itemPriceRepository) GetItemPricesForOrders(req *common.RequestContext, reqDTO *dto.RequestItemPricesForOrder) (
	result dto.ResultEntity[[]dto.ItemPriceDto], err error) {

	var columns []field.Expr

	itemPrice := s.Q.ItemPrice
	item := s.Q.Item
	uom := s.Q.UnitOfMeasureTranslation
	itemPriceUom := s.Q.UnitOfMeasureTranslation.As("item_price_uom")
	priceListQ := s.Q.PriceList

	// Base conditions for filtering items
	itemConds := []gen.Condition{
		item.CompanyID.Eq(req.ActiveCompany.ID),
	}

	// Apply search query condition if provided
	if reqDTO.Query != "" {
		itemConds = append(itemConds, item.Name.Like("%"+reqDTO.Query+"%"))
	}

	// Common item-related columns
	columns = append(columns,
		item.Name.As("item_name"),
		item.ID.As("item_id"),
		item.UUID.As("item_uuid"),
		item.Code.As("item_code"),
		uom.Name.As("uom"),
		item.UnitOfMeasureID.As("uom_id"),
	)

	// Determine applicable price list IDs
	priceListID := s.convertor.StrtoInt(reqDTO.PriceListID)
	var priceListIDs []int64

	if priceListID != 0 {
		priceListIDs = []int64{priceListID}
	} else {
		var priceLists []*model.PriceList
		priceLists, err = priceListQ.WithContext(req.Ctx).Select(priceListQ.ID).Where(
			priceListQ.CompanyID.Eq(req.ActiveCompany.ID),
			priceListQ.Status.Eq(proto.State_ENABLED.String()),
			priceListQ.Currency.Eq(reqDTO.Currency),
		).Find()
		if err != nil {
			return result, err
		}
		priceListIDs = lo.Map(priceLists, func(pl *model.PriceList, _ int) int64 {
			return pl.ID
		})
	}

	// Base query with item and UoM translation join
	builder := item.WithContext(req.Ctx).
		Join(uom, uom.BaseID.EqCol(item.UnitOfMeasureID), uom.LanguageCode.Eq(string(req.LanguageCode)))

	// Add price list and item price related joins/columns if priceListIDs exist
	if len(priceListIDs) > 0 {
		columns = append(columns,
			itemPrice.ID,
			itemPrice.UUID,
			itemPrice.ItemQuantity,
			itemPrice.Rate,
			priceListQ.Name.As("price_list_name"),
			priceListQ.Currency.As("price_list_currency"),
			itemPrice.UnitOfMeasureID.As("item_price_uom_id"),
			itemPriceUom.Name.As("item_price_uom"),
		)

		builder = builder.
			LeftJoin(itemPrice, itemPrice.ItemID.EqCol(item.ID), itemPrice.PriceListID.In(priceListIDs...)).
			LeftJoin(priceListQ, priceListQ.ID.EqCol(itemPrice.PriceListID)).
			LeftJoin(itemPriceUom, itemPriceUom.BaseID.EqCol(itemPrice.UnitOfMeasureID), itemPriceUom.LanguageCode.Eq(string(req.LanguageCode))).
			Order(itemPrice.ID.Asc())
	}

	// Finalize and execute the query
	err = builder.Select(columns...).
		Where(itemConds...).
		Limit(domain.DEFAULT_LIMIT).
		Scan(&result.Entity)

	return result, err
}

func (s *itemPriceRepository) UpsertItemPrice(req *common.RequestContext, d *dto.UpsertItemPriceRequest) error {
	d.Body.ItemPrice.CompanyID = req.ActiveCompany.ID
	err := s.DB.WithContext(req.Ctx).Save(&d.Body.ItemPrice).Error
	return err
}
