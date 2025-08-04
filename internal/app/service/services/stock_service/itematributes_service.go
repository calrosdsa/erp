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

type ItemAttributeService struct {
	conn    *connection.Connection
	timeout time.Duration
	emitLog logger.EmitLog
	permissionService permission.PermissionService
}

func NewItemAttributeService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	permissionService permission.PermissionService,
	logger logger.Logger,
) *ItemAttributeService {
	return &ItemAttributeService{
		conn:    conn,
		timeout: timeout,
		emitLog: logger.EmitLog("item-attributes-service"),
		permissionService: permissionService,
	}
}

func (s *ItemAttributeService) UpsertItemAttributeValue(req *common.RequestContext, d *dto.UpsertItemAttributeValueRequest) (err error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("UpsertItemAttributeValue"))
		}
	}()
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM_ATTRIBUTE, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	// itemAttributeValue := s.conn.Q.ItemAttributeValue.WithContext(ctx).Where(
	// 	s.conn.Q.ItemAttribute.UUID.Eq(d.Body.ID),
	// )
	var itemAttributeValue model.ItemAttributeValue
	itemAttributeValue.ID = d.Body.ID
	itemAttributeValue.ItemAttributeID = d.Body.ItemAttributeID
	itemAttributeValue.Abbreviation = d.Body.Abbreviation
	itemAttributeValue.Value = d.Body.Value
	itemAttributeValue.Ordinal = d.Body.Ordinal
	err = s.conn.Q.ItemAttributeValue.WithContext(ctx).Save(&itemAttributeValue)

	return
}

func (s *ItemAttributeService) GetItemAttributeDetail(req *common.RequestContext, d *dto.RequestEntity) (dto.ResultEntity[dto.ItemAttributeDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		err error
		res dto.ResultEntity[dto.ItemAttributeDto]
		itemAttribute model.ItemAttribute
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemAttributeDetail"))
		}
	}()

	err = s.conn.Db.WithContext(ctx).Where(&model.ItemAttribute{
		UUID: d.ID,CompanyID: req.ActiveCompany.ID,
		}).First(&itemAttribute).Error
	if err != nil {
		return res, nil
	}
	res.Entity = dto.ItemAttributeDtoFromModel(&itemAttribute)
	itemAttributeValues,err := s.conn.Q.ItemAttributeValue.Where(
		s.conn.Q.ItemAttributeValue.ItemAttributeID.Eq(itemAttribute.ID),
	).Find()
	itemAttributeValueDtos := make([]dto.ItemAttributeValueDto,len(itemAttributeValues))
	for i,itemAttributeValue := range itemAttributeValues {
		itemAttributeValueDtos[i] = dto.ItemAttributeValueDtoFromModel(itemAttributeValue)
	}
	res.Entity.ItemAttributeValues = itemAttributeValueDtos
	return res, err
}

func (s *ItemAttributeService) GetItemAttributes(req *common.RequestContext, d *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.ItemAttributeDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.PaginationResult[[]dto.ItemAttributeDto]
		itemAttributes []model.ItemAttribute
		err error
	)
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM_ATTRIBUTE, domain.VIEW); !allow {
		return res,domain.ACTION_NOT_ALLOWED
	}
	defer func ()  {
		if err != nil {
			s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetItemAttributeDtos"))
		}
	}()
	queryBuilder := s.conn.Db.WithContext(ctx).Model(&model.ItemAttribute{}).
		Where(&model.ItemAttribute{CompanyID: req.ActiveCompany.ID})

	err = queryBuilder.
		Count(&res.Total).Error

	if d.Query != "" {
		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
	}

	err = queryBuilder.
		Scopes(s.conn.Paginate(req.Params)).
		Find(&itemAttributes).Error
	if err != nil {
		return res, err
	}
	itemAttributesDto := make([]dto.ItemAttributeDto,len(itemAttributes))
	for i,itemAttribute := range  itemAttributes {
		itemAttributesDto[i] = dto.ItemAttributeDtoFromModel(&itemAttribute)
	}
	res.Results = itemAttributesDto
	return res, err
}

func (s *ItemAttributeService) CreateAttributeItem(req *common.RequestContext, d *dto.CreateItemAttributeRequest) (
	dto.ItemAttributeDto,error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	tx := s.conn.Q.Begin()
	var (
		err error
		res dto.ItemAttributeDto
	)
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM_ATTRIBUTE, domain.CREATE); !allow {
		return res,domain.ACTION_NOT_ALLOWED
	}
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateAttributeItem"))
			tx.Rollback()
		}
	}()
	if tx.Error != nil {
		return res,err
	}
	var itemAttribute model.ItemAttribute
	partyId,err := s.conn.Q.ItemAttribute.InsertParty(domain.PARTY_ITEM_ATTRIBUTE)
	itemAttribute.Name = d.Body.Name
	itemAttribute.ID = partyId
	itemAttribute.CompanyID = req.ActiveCompany.ID
	err = tx.ItemAttribute.WithContext(ctx).Save(&itemAttribute)
	if err != nil {
		return res,err
	}
	itemAttributeValues := make([]*model.ItemAttributeValue, len(d.Body.Values))
	for i, value := range d.Body.Values {
		itemAttributeValue := &model.ItemAttributeValue{}
		itemAttributeValue.Value = value.Value
		itemAttributeValue.Abbreviation = value.Abbreviation
		itemAttributeValue.Ordinal = value.Ordinal
		itemAttributeValue.ItemAttributeID = itemAttribute.ID
		itemAttributeValues[i] = itemAttributeValue
	}
	err = tx.ItemAttributeValue.WithContext(ctx).CreateInBatches(itemAttributeValues, len(itemAttributeValues))
	if err != nil{
		return res,err
	}
	err = tx.Commit()
	return res,err
}
