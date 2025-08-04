package stockservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/connection"
	"erp/internal/domain"
	"erp/internal/app/service/helpers"
	"erp/pkg/cache"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"time"
)

type ItemStockService struct {
	conn              *connection.Connection
	timeout           time.Duration
	emitLog           logger.EmitLog
	cache             *cache.Cache
	convertor         helpers.ConvertorHelper
	permissionService permission.PermissionService
}

func NewItemStockService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	cache *cache.Cache,
	permissionService permission.PermissionService,
	logger logger.Logger,
) *ItemStockService {
	return &ItemStockService{
		conn:              conn,
		timeout:           timeout,
		emitLog:           logger.EmitLog("item-stock-service"),
		cache:             cache,
		convertor:         helpers.Convertor,
		permissionService: permissionService,
	}
}
func (s *ItemStockService) GetWarehouseItemStockLevels(req *common.RequestContext, d *dto.RequestPaginationData) (dto.PaginationResult[[]dto.StockLevelDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		err error
		res dto.PaginationResult[[]dto.StockLevelDto]
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetWarehouseItemStockLevels"))
		}
	}()
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM_STOCK, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	stocklevel := s.conn.Q.StockLevel
	item := s.conn.Q.Item
	var warehouse model.WareHouse
	warehouse.UUID = d.FilterID
	warehouse.CompanyID = req.ActiveCompany.ID
	if err := s.cache.GetEntity(
		cache.Options.WithData(&warehouse),
		cache.Options.WithKeyEntity(domain.CACHE_WAREHOUSE_KEY),
		cache.Options.WithID(d.FilterID),
		cache.Options.WithTypeAssertion(func(d interface{}) {
			val, ok := d.(*model.WareHouse)
			if ok {
				warehouse = *val
			}
		}),
	); err != nil {
		return res, err
	}

	limit, offset := s.convertor.ToPaginationParams(d.Page, d.Size)

	// item := s.conn.Q.Item
	s.conn.Q.StockLevel.WithContext(ctx).Select(
		stocklevel.UUID, stocklevel.CreatedAt, stocklevel.Enabled, stocklevel.Stock, stocklevel.OutOfStockThreshold,
		stocklevel.UUID, item.Name.As("item_name"), item.UUID.As("item_uuid"),
	).Join(item, item.ID.EqCol(stocklevel.ItemID)).Where(
		stocklevel.WareHouseID.Eq(warehouse.ID),
	).Limit(limit).Offset(offset).
		Scan(&res.Results)

	count, err := stocklevel.WithContext(ctx).Where(
		stocklevel.WareHouseID.Eq(warehouse.ID),
	).Count()
	if err != nil {
		return res, err
	}
	res.Total = count

	return res, err
}

func (s *ItemStockService) GetItemStockLevels(req *common.RequestContext, d *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.StockLevelDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		err error
		res dto.PaginationResult[[]dto.StockLevelDto]
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemStockLevels"))
		}
	}()
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM_STOCK, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	stocklevel := s.conn.Q.StockLevel
	warehouse := s.conn.Q.WareHouse
	var item model.Item
	item.UUID = d.FilterID
	item.CompanyID = req.ActiveCompany.ID
	if err := s.cache.GetEntity(
		cache.Options.WithData(&item),
		cache.Options.WithKeyEntity(domain.CACHE_ITEM_KEY),
		cache.Options.WithID(d.FilterID),
		cache.Options.WithTypeAssertion(func(d interface{}) {
			val, ok := d.(*model.Item)
			if ok {
				item = *val
			}
		}),
	); err != nil {
		return res, err
	}

	limit, offset := s.convertor.ToPaginationParams(d.Page, d.Size)

	// item := s.conn.Q.Item
	s.conn.Q.StockLevel.WithContext(ctx).Select(
		stocklevel.UUID, stocklevel.CreatedAt, stocklevel.Enabled, stocklevel.Stock, stocklevel.OutOfStockThreshold,
		stocklevel.UUID, warehouse.Name.As("warehouse_name"), warehouse.UUID.As("warehouse_uuid"),
	).Join(warehouse, warehouse.ID.EqCol(stocklevel.WareHouseID)).Where(
		stocklevel.ItemID.Eq(item.ID),
	).Limit(limit).Offset(offset).
		Scan(&res.Results)

	count, err := stocklevel.WithContext(ctx).Where(
		stocklevel.ItemID.Eq(item.ID),
	).Count()
	if err != nil {
		return res, err
	}
	res.Total = count
	return res, err
}

func (s *ItemStockService) AddStockLevel(req *common.RequestContext, i *dto.AddStockLevelRequest) error {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	tx := s.conn.Q.Begin()
	var (
		err        error
		stockLevel model.StockLevel
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("AddStockLevel"))
			tx.Rollback()
		}
	}()
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM_STOCK, domain.CREATE); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	var item model.Item
	item.UUID = i.Body.ItemUUID
	if err := s.cache.GetEntity(
		cache.Options.WithData(&item),
		cache.Options.WithKeyEntity(domain.CACHE_ITEM_KEY),
		cache.Options.WithID(i.Body.ItemUUID),
		cache.Options.WithTypeAssertion(func(d interface{}) {
			val, ok := d.(*model.Item)
			if ok {
				item = *val
			}
		}),
	); err != nil {
		return err
	}
	var warehouse model.WareHouse
	warehouse.UUID = i.Body.WareHouseUUID
	if err := s.cache.GetEntity(
		cache.Options.WithData(&warehouse),
		cache.Options.WithKeyEntity(domain.CACHE_WAREHOUSE_KEY),
		cache.Options.WithID(i.Body.WareHouseUUID),
		cache.Options.WithTypeAssertion(func(d interface{}) {
			val, ok := d.(*model.WareHouse)
			if ok {
				warehouse = *val
			}
		}),
	); err != nil {
		return err
	}

	partyId, err := tx.StockLevel.InsertParty(domain.PARTY_STOCK_LEVEL)
	stockLevel.ItemID = item.ID
	stockLevel.ID = partyId
	stockLevel.WareHouseID = warehouse.ID
	stockLevel.Enabled = i.Body.Enabled
	stockLevel.OutOfStockThreshold = i.Body.OutOfStockThreshold
	// stockLevel.Stock = i.Body.Stock
	err = tx.StockLevel.Save(&stockLevel)
	if err != nil {
		return err
	}
	var stockMovement model.StockMovement
	stockMovement.ItemID = stockLevel.ItemID
	stockMovement.WareHouseID = stockLevel.WareHouseID
	stockMovement.Quantity = i.Body.Stock
	stockMovement.StockMovementType = string(domain.ADJUSTMENT)
	err = s.StockMovement(tx, stockMovement)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *ItemStockService) StockMovement(tx *query.QueryTx, stockMovement model.StockMovement) error {
	var (
		stockLevel *model.StockLevel
		err        error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("StockMovement"))
		}
	}()
	stockLevel, err = tx.StockLevel.Where(
		s.conn.Q.StockLevel.ItemID.Eq(stockMovement.ItemID),
		s.conn.Q.StockLevel.WareHouseID.Eq(stockMovement.WareHouseID),
	).First()
	if err != nil {
		return err
	}
	diff := stockLevel.Stock + stockMovement.Quantity

	if diff < stockLevel.OutOfStockThreshold {
		return domain.ERROR_OUT_OF_STOCK
	}
	stockLevel.Stock = diff
	err = tx.StockMovement.Save(&stockMovement)
	if err != nil {
		return err
	}
	err = tx.StockLevel.Save(stockLevel)

	return err
}
