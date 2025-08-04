package price_list_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	price_list_repo "erp/project/stock/price_list/repository"
)

type PriceListUseCase interface {
	CreatePriceList(req *common.RequestContext, i *dto.CreatePriceListRequest) (res dto.PriceListDto, err error)
	CreatePriceListTx(tx *query.QueryTx,req *common.RequestContext, i *dto.CreatePriceListRequest) (
		res dto.PriceListDto, err error)
	GetListPriceDetail(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.PriceListDto], err error)
	GetPriceLists(req *common.RequestContext, d *dto.RequestPriceLists) (
		res dto.PaginationResult[[]dto.PriceListDto], err error)
	
	EditPriceList(req *common.RequestContext, d *dto.EditPriceListRequest) (err error)
}

type priceListUseCase struct {
	emitLog       logger.EmitLog
	priceListRepo price_list_repo.PriceListRepository
	permission    repository.PermissionService
	core          repository.CoreService
}

func NewPriceListUseCase(
	logger logger.Logger,
	priceListRepo price_list_repo.PriceListRepository,
	permission repository.PermissionService,
	core repository.CoreService,
) PriceListUseCase {
	return &priceListUseCase{
		emitLog:       logger.EmitLog("price-list"),
		priceListRepo: priceListRepo,
		core:          core,
		permission:    permission,
	}
}
func (u *priceListUseCase) EditPriceList(req *common.RequestContext, d *dto.EditPriceListRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditPriceList"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICE_LIST, domain.EDIT)
	if err != nil {
		return
	}
	err = u.priceListRepo.EditPriceList(req, d)
	return
}

func (u *priceListUseCase) CreatePriceList(req *common.RequestContext, i *dto.CreatePriceListRequest) (
	res dto.PriceListDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePriceList"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICE_LIST, domain.CREATE)
	if err != nil {
		return 
	}
	priceList,err := u.priceListRepo.CreatePriceList(req, i)
	res = dto.PriceListDtoFromModel(&priceList)
	return 
}

func (u *priceListUseCase) CreatePriceListTx(tx *query.QueryTx,req *common.RequestContext, i *dto.CreatePriceListRequest) (
	res dto.PriceListDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePriceListTx"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICE_LIST, domain.CREATE)
	if err != nil {
		return 
	}
	priceList,err := u.priceListRepo.CreatePriceListTx(tx,req, i)
	res = dto.PriceListDtoFromModel(&priceList)
	return 
}
func (u *priceListUseCase) GetListPriceDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.PriceListDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetListPriceDetail"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICE_LIST, domain.VIEW)
	if err != nil {
		return res, err
	}
	res, err = u.priceListRepo.GetListPriceDetail(req, i)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *priceListUseCase) GetPriceLists(req *common.RequestContext, d *dto.RequestPriceLists) (
	res dto.PaginationResult[[]dto.PriceListDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPriceLists"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PRICE_LIST, domain.VIEW)
	if err != nil {
		return res, err
	}
	res, err = u.priceListRepo.GetPriceLists(req, d)
	return
}
