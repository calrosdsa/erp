package item_inventory_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	item_inventory_repo "erp/project/stock/item_inventory/internal/repository"
)

type ItemInventoryUcase interface {
	GetItemInventory(req *common.RequestContext, d *dto.RequestEntity) (res dto.ItemInventoryDto, err error)
	EditItemInventory(req *common.RequestContext, d dto.ItemInventoryFields) (err error)
}

type itemInventoryUcase struct {
	emitLog           logger.EmitLog
	permission        repository.PermissionService
	itemInventoryRepo item_inventory_repo.ItemInventoryRepo
}

func NewItemInventoryUcase(
	logger logger.Logger,
	permission repository.PermissionService,
	itemInventoryRepo item_inventory_repo.ItemInventoryRepo,
) ItemInventoryUcase {
	return &itemInventoryUcase{
		permission:        permission,
		emitLog:           logger.EmitLog("item-inventory-ucase"),
		itemInventoryRepo: itemInventoryRepo,
	}
}

func (u *itemInventoryUcase) EditItemInventory(req *common.RequestContext, d dto.ItemInventoryFields) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditItemInventory"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM, domain.EDIT); err != nil {
		return
	}
	err = u.itemInventoryRepo.EditItemInventory(req, d)
	return
}

func (u *itemInventoryUcase) GetItemInventory(req *common.RequestContext, d *dto.RequestEntity) (res dto.ItemInventoryDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemInventory"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.ITEM, domain.VIEW); err != nil {
		return
	}
	res, err = u.itemInventoryRepo.GetItemInventory(req, d)
	return
}
