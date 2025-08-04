package accountingservice

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

type TaxService struct {
	conn    *connection.Connection
	timeout time.Duration
	emitLog logger.EmitLog
	permissionService permission.PermissionService
}

func NewTaxService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	permissionService permission.PermissionService,
	logger logger.Logger,
) *TaxService {
	return &TaxService{
		conn:    conn,
		timeout: timeout,
		emitLog: logger.EmitLog("tax-service"),
		permissionService: permissionService,
	}
}

func (s *TaxService) GetTaxDetail(req *common.RequestContext,i *dto.RequestEntity)(
	dto.ResultEntity[dto.TaxDto],error){
	ctx,cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		result dto.ResultEntity[dto.TaxDto]
		err error
	)
	if allow := s.permissionService.CheckPermission(ctx,req,domain.TAX,domain.VIEW); !allow {
		return result,domain.ACTION_NOT_ALLOWED
	}
	var tax model.Tax
	err = s.conn.Db.WithContext(ctx).Where(&model.Tax{UUID: i.ID}).First(&tax).Error
	if err != nil {
		s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetTaxDetail"))
	}
	result.Entity = dto.TaxDtoFromModel(&tax)
	return result,err
}

func (s *TaxService) CreateTax(req *common.RequestContext,i *dto.CreateTaxRequest)(error){
	ctx,cancel := context.WithTimeout(req.Ctx,s.timeout)
	defer cancel()
	if allow := s.permissionService.CheckPermission(ctx,req,domain.TAX,domain.CREATE); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	var (
		tax model.Tax
		err error
	)
	defer func(){
		if err != nil {
			s.emitLog.Err(err,logger.OptionsLog.WithMethod("CreateTax"))
		}
	}()
	partyId,err := s.conn.Q.Tax.InsertParty(domain.PARTY_TAX)
	tax.Name = i.Body.Name
	tax.Value = i.Body.Value
	tax.Enabled = i.Body.Enabled
	tax.ID = partyId
	tax.CompanyID = req.ActiveCompany.ID
	err = s.conn.Q.Tax.WithContext(ctx).Save(&tax)
	if err != nil {
		s.emitLog.Err(err,logger.OptionsLog.WithMethod("CreateTax"))
	}
	return err
}

func (s *TaxService) GetTaxes(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.TaxDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.PaginationResult[[]dto.TaxDto]
		err    error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetTaxes"))
		}
	}()
	if allow := s.permissionService.CheckPermission(ctx,req,domain.TAX,domain.VIEW); !allow {
		return res,domain.ACTION_NOT_ALLOWED
	}
	var taxes []model.Tax
	queryBuilder := s.conn.Db.WithContext(ctx).Model(&model.Tax{}).
		Where(&model.Tax{CompanyID: req.ActiveCompany.ID})

	err = queryBuilder.
		Count(&res.Total).Error

	if i.Query != "" {
		queryBuilder = queryBuilder.Where("code ILIKE ?", "%"+i.Query+"%")
	}

	err = queryBuilder.
		Scopes(s.conn.Paginate(req.Params)).
		Order(s.conn.Order(req.Params)).
		Find(&taxes).Error
	if err != nil {
		return res, err
	}
	taxDtos := make([]dto.TaxDto,len(taxes))
	for i,tax := range taxes {
		taxDtos[i] = dto.TaxDtoFromModel(&tax)
	}
	res.Results = taxDtos
	return res, err

}
