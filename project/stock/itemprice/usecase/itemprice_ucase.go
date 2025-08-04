package itemprice_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	itemprice_repo "erp/project/stock/itemprice/repository"
)

type ItemPriceUseCase interface {
	CreateItemPrice(req *common.RequestContext, d dto.ItemPriceData) (res dto.ItemPriceDto, err error)
	CreateItemPriceTx(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (res dto.ItemPriceDto, err error)
	GetItemPrice(req *common.RequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.ItemPriceDto], error)
	GetItemPrices(req *common.RequestContext, i dto.ItemPricesRequest) (
		res dto.ResponseDataList[[]dto.ItemPriceDto], err error)
	GetItemPricesByItemCode(req *common.RequestContext, d *dto.RequestItemPriceByCode) (
		res dto.PaginationResult[[]dto.ItemPriceDto], err error)
	GetItemPricesForOrders(req *common.RequestContext, i *dto.RequestItemPricesForOrder) (
		dto.ResultEntity[[]dto.ItemPriceDto], error)
	UpsertItemPrice(req *common.RequestContext, d *dto.UpsertItemPriceRequest) error
	EditItemPrice(req *common.RequestContext, d dto.ItemPriceData) (err error)
}

type itemPriceUseCase struct {
	emitLog    logger.EmitLog
	permission repository.PermissionService
	core       repository.CoreService
	itemPrice  itemprice_repo.ItemPriceRepository
}

func NewItemPriceUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	core repository.CoreService,
	itemPrice itemprice_repo.ItemPriceRepository,
) ItemPriceUseCase {
	return &itemPriceUseCase{
		emitLog:    logger.EmitLog("itemprice-usecase"),
		permission: permission,
		core:       core,
		itemPrice:  itemPrice,
	}
}

func (u *itemPriceUseCase) EditItemPrice(req *common.RequestContext, d dto.ItemPriceData) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditItemPrice"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM_PRICE, domain.EDIT)
	if err != nil {
		return
	}
	err = u.itemPrice.EditItemPrice(req, d)
	return
}

func (u *itemPriceUseCase) CreateItemPrice(req *common.RequestContext, d dto.ItemPriceData) (res dto.ItemPriceDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateItemPrice"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM_PRICE, domain.CREATE); err != nil {
		return
	}
	res, err = u.itemPrice.CreateItemPrice(req, d)
	if err != nil {
		return
	}

	return
}
func (u *itemPriceUseCase) CreateItemPriceTx(tx *query.QueryTx, req *common.RequestContext, d dto.ItemPriceData) (res dto.ItemPriceDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateItemPriceTx"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM_PRICE, domain.CREATE); err != nil {
		return
	}
	res, err = u.itemPrice.CreateItemPriceTx(tx, req, d)
	if err != nil {
		return
	}

	return
}
func (u *itemPriceUseCase) GetItemPrice(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ItemPriceDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemPrice"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.ITEM_PRICE, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.itemPrice.GetItemPrice(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *itemPriceUseCase) GetItemPrices(req *common.RequestContext, i dto.ItemPricesRequest) (
	res dto.ResponseDataList[[]dto.ItemPriceDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemPrices"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.ITEM_PRICE, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res.Body.Result, err = u.itemPrice.GetItemPrices(req, i)
	return
}
func (u *itemPriceUseCase) GetItemPricesByItemCode(req *common.RequestContext, i *dto.RequestItemPriceByCode) (
	res dto.PaginationResult[[]dto.ItemPriceDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemPricesByItemCode"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.ITEM_PRICE, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = u.itemPrice.GetItemPricesByItemCode(req, i)
	return
}
func (u *itemPriceUseCase) GetItemPricesForOrders(req *common.RequestContext, i *dto.RequestItemPricesForOrder) (
	res dto.ResultEntity[[]dto.ItemPriceDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemPricesForOrders"))
		}
	}()
	res, err = u.itemPrice.GetItemPricesForOrders(req, i)
	return
}
func (u *itemPriceUseCase) UpsertItemPrice(req *common.RequestContext, d *dto.UpsertItemPriceRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpsertItemPrice"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.ITEM_PRICE, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.itemPrice.UpsertItemPrice(req, d)
	return
}
