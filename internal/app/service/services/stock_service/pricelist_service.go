package stockservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/connection"
	"erp/internal/domain"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"time"
)

type PriceListService struct {
	conn              *connection.Connection
	timeout           time.Duration
	emitLog           logger.EmitLog
	permissionService permission.PermissionService

}

func NewPriceListServer(conn *connection.Connection, timeout time.Duration,
	helpers *helpers.Helpers,
	permissionService permission.PermissionService,
	logger logger.Logger,
) *PriceListService {
	return &PriceListService{
		conn:              conn,
		timeout:           timeout,
		emitLog:           logger.EmitLog("price-list"),
		permissionService: permissionService,
	}
}

func (s *PriceListService) CreatePriceList(req *common.RequestContext, i *dto.CreatePriceListRequest) error {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	tx := s.conn.Db.Begin()
	var err error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	err = tx.Error
	if err != nil {
		return err
	}
	if allow := s.permissionService.CheckPermission(ctx, req, domain.PRICE_LIST, domain.CREATE); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	var priceList model.PriceList
	// fmt.Println("PRICE LIST BODY",i.Body)
	priceList.Name = i.Body.Name
	priceList.Currency = i.Body.Currency
	priceList.CompanyID = req.ActiveCompany.ID
	priceList.IsBuying = i.Body.IsBuying
	priceList.IsSelling = i.Body.IsSelling
	err = tx.WithContext(ctx).Save(&priceList).Error
	if err != nil {
		s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePriceList"))
	}
	return tx.Commit().Error
}

func (s *PriceListService) GetListPriceDetail(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.PriceListDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.ResultEntity[dto.PriceListDto]
		err error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetListPriceDetail"))
		}
	}()
	priceList, err := s.conn.Q.PriceList.WithContext(ctx).Where(
		s.conn.Q.PriceList.UUID.Eq(i.ID),
	).First()
	res.Entity = dto.PriceListDtoFromModel(priceList)
	return res, err
}

func (s *PriceListService) GetPriceLists(req *common.RequestContext, d *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.PriceListDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.PaginationResult[[]dto.PriceListDto]
		err error
		priceLists []model.PriceList
	)
	if allow := s.permissionService.CheckPermission(ctx, req, domain.PRICE_LIST, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	queryBuilder := s.conn.Db.WithContext(ctx).Model(&model.PriceList{}).
		Where("company_id = ?", req.ActiveCompany.ID)
	err = queryBuilder.Count(&res.Total).Error
	if err != nil {
		return res, err
	}

	// Add search query if provided
	if d.Query != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
	}

	// Apply pagination and get the results
	err = queryBuilder.Scopes(s.conn.Paginate(req.Params)).Find(&priceLists).Error
	if err != nil {
		return res, err
	}

	priceListDtos := make([]dto.PriceListDto,len(priceLists))
	for i,priceList := range priceLists {
		priceListDtos[i] = dto.PriceListDtoFromModel(&priceList)
	}
	res.Results = priceListDtos

	return res, nil
}

func (s *PriceListService) UpsertPriceList(req *common.RequestContext, d *dto.UpsertPriceListRequest) error {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	if allow := s.permissionService.CheckPermission(ctx, req, domain.PRICE_LIST, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	d.Body.ItemPriceList.CompanyID = req.ActiveCompany.ID
	err := s.conn.Db.WithContext(ctx).Save(&d.Body.ItemPriceList).Error
	return err
}
