package stockservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/connection"
	"erp/internal/domain"
	"erp/internal/app/service/helpers"
	"erp/pkg/cache"
	"erp/pkg/logger"
	"erp/pkg/permission"
	"time"
)

type ItemVariantService struct {
	conn              *connection.Connection
	timeout           time.Duration
	convertor         helpers.ConvertorHelper
	permissionService permission.PermissionService
	emitLog           logger.EmitLog
	cache             *cache.Cache
}

func NewItemVaraintService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	permissionService permission.PermissionService,
	cache *cache.Cache,
	logger logger.Logger,
) *ItemVariantService {
	return &ItemVariantService{
		timeout:           timeout,
		conn:              conn,
		convertor:         helpers.Convertor,
		permissionService: permissionService,
		emitLog:           logger.EmitLog("item-variant-service"),
		cache:             cache,
	}
}

func (s *ItemVariantService) GetVariantsFromItem(req *common.RequestContext, i *dto.RequestPaginationData) (
	dto.PaginationResult[[]dto.ItemVariantDto], error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	var (
		res dto.PaginationResult[[]dto.ItemVariantDto]
		err error
	)
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetVariantsFromItem"))
		}
	}()
	var item model.Item
	item.UUID = i.FilterID
	item.CompanyID = req.ActiveCompany.ID
	if err := s.cache.GetEntity(
		cache.Options.WithData(&item),
		cache.Options.WithKeyEntity(domain.CACHE_ITEM_KEY),
		cache.Options.WithID(i.FilterID),
		cache.Options.WithTypeAssertion(func(d interface{}) {
			val, ok := d.(*model.Item)
			if ok {
				item = *val
			}
		}),
	); err != nil {
		return res, err
	}
	variantQ := s.conn.Q.Item
	itemVariantQ := s.conn.Q.ItemVariant
	attributeValueQ := s.conn.Q.ItemAttributeValue
	limit, offset := s.convertor.ToPaginationParams(i.Page, i.Size)
	s.conn.Q.ItemVariant.WithContext(ctx).Select(
		variantQ.Name, variantQ.UUID, variantQ.Code, variantQ.CreatedAt,
		attributeValueQ.Value, attributeValueQ.Abbreviation,
	).Join(attributeValueQ, attributeValueQ.ID.EqCol(itemVariantQ.ItemAttributeValueID)).
		Join(variantQ, variantQ.ID.EqCol(itemVariantQ.VariantID)).Limit(limit).Offset(offset).Where(
		itemVariantQ.ItemID.Eq(item.ID),
	).Scan(&res.Results)

	total, err := s.conn.Q.ItemVariant.Where(
		itemVariantQ.ItemID.Eq(item.ID),
	).Count()
	res.Total = total
	return res, err
}

func (s *ItemVariantService) CreateItemVariant(req *common.RequestContext, i *dto.CreateItemVariantRequest) (err error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer cancel()
	tx := s.conn.Db.Begin()
	if allow := s.permissionService.CheckPermission(ctx, req, domain.ITEM, domain.CREATE); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var item model.Item
	item.UUID = i.Body.ItemUUID
	item.CompanyID = req.ActiveCompany.ID
	err = s.conn.GetEntity(ctx, tx, &item)
	if err != nil {
		return
	}
	var newItem model.Item

	newItem.Name = i.Body.Name
	code := s.conn.GenerateCode(ctx, tx, &model.Item{}, req.ActiveCompany.ID)
	newItem.Code = &code
	newItem.CompanyID = req.ActiveCompany.ID
	newItem.UnitOfMeasureID = item.UnitOfMeasureID
	newItem.ItemType = domain.ITEM_VARIANT_TYPE
	newItem.GroupID = item.GroupID
	err = tx.WithContext(ctx).Save(&newItem).Error
	if err != nil {
		return
	}

	var itemVariant model.ItemVariant
	itemVariant.ItemAttributeValueID = i.Body.AttributeValueValueID
	itemVariant.ItemID = item.ID
	itemVariant.VariantID = newItem.ID
	err = tx.WithContext(ctx).Save(&itemVariant).Error
	if err != nil {
		return
	}
	return tx.Commit().Error
}
