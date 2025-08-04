package warehouse_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	warehouse_repo "erp/project/stock/warehouse/internal/repository"
)

type WarehouseUseCase interface {
	GetWareHouseDetail(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.WareHouseDto], err error)
	GetWareHouses(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.WareHouseDto], err error)
	CreateWareHouse(req *common.RequestContext, i *dto.CreateWareHouseRequest) (err error)
	GetWarehouseTreeView(req *common.RequestContext) (
		res []dto.TreeEntryDto, err error)
	EditWarehouse(req *common.RequestContext, d *dto.EditWarehouseRequest) (err error)
}

type warehouseUcase struct {
	emitLog       logger.EmitLog
	permission    repository.PermissionService
	core          repository.CoreService
	warehouseRepo warehouse_repo.WarehouseRepository
}

func NewWarehouseUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	core repository.CoreService,
	warehouseRepo warehouse_repo.WarehouseRepository,
) WarehouseUseCase {
	return &warehouseUcase{
		emitLog:       logger.EmitLog("warehouse-usecase"),
		permission:    permission,
		core:          core,
		warehouseRepo: warehouseRepo,
	}
}
func(u *warehouseUcase) EditWarehouse(req *common.RequestContext, d *dto.EditWarehouseRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditWarehouse"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx,req,domain.WAREHOUSE,domain.EDIT)
	if err != nil {
		return
	}
	err = u.warehouseRepo.EditWarehouse(req,d)
	return 
}
func (u *warehouseUcase) GetWarehouseTreeView(req *common.RequestContext) (
	res []dto.TreeEntryDto, err error){
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("GetWarehouseTreeView"))
		}
	}()
	res,err = u.warehouseRepo.GetWarehouseTreeView(req)
	return 
}

func (u *warehouseUcase) GetWareHouseDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.WareHouseDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetWarehouseDetial"))
		}
	}()
	res, err = u.warehouseRepo.GetWareHouseDetail(req, i)
	if err != nil {
		return
	}
	res.Addresses = u.core.GetPartyAddresses(req, res.Entity.ID)
	res.Contacts = u.core.GetPartyContacts(req, res.Entity.ID)
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *warehouseUcase) GetWareHouses(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.WareHouseDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetWarehouses"))
		}
	}()
	res, err = u.warehouseRepo.GetWareHouses(req, d)
	if err != nil {
		return
	}
	return
}
func (u *warehouseUcase) CreateWareHouse(req *common.RequestContext, i *dto.CreateWareHouseRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateWareHouse"))
		}
	}()
	err = u.warehouseRepo.CreateWareHouse(req, i)
	if err != nil {
		return
	}
	return
}
