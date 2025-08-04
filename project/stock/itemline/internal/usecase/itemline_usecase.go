package itemline_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/logger"
	itemline_repo "erp/project/stock/itemline/internal/repository"
)

type ItemLineUseCase interface {
	EditItemLine(req *common.RequestContext, i *dto.EditLineItemRequest) error
	GetItemLines(req *common.RequestContext, d *dto.RequestLineItems) (
		res []dto.LineItemDto, err error)
	DeleteLineItem(req *common.RequestContext, d *dto.DeleteLineItemRequest) (err error)
	AddLineItem(req *common.RequestContext,d *dto.AddLineItemRequest)(err error)
	UpsertProductList(req *common.RequestContext,dto dto.ProductListData) (err error)
}

type itemLineUseCase struct {
	emitLog      logger.EmitLog
	permission   repository.PermissionService
	itemLineRepo itemline_repo.ItemLineRepository
	bus          bus.Bus
	c            di.Container
}

func NewItemLineUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	itemLineRepo itemline_repo.ItemLineRepository,
	bus bus.Bus,
	c di.Container,
) ItemLineUseCase {
	return &itemLineUseCase{
		emitLog:      logger.EmitLog("itemline-usecase"),
		permission:   permission,
		itemLineRepo: itemLineRepo,
		bus:          bus,
		c:            c,
	}
}

func(u *itemLineUseCase) UpsertProductList(req *common.RequestContext,d dto.ProductListData) (err error) {
	defer func ()  {
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("UpsertProductList"))
		}	
	}()
	err =  u.itemLineRepo.UpsertProductList(req,d)
	return 
}

func(u *itemLineUseCase)AddLineItem(req *common.RequestContext,d *dto.AddLineItemRequest)(err error){
	defer func(){
		if err != nil {
			u.emitLog.Err(err,logger.OptionsLog.WithMethod("AddLineItem"))
		}
	}()
	entity, err := u.GetItemLinePartyEntity(d.Body.DocPartyType)
	if err != nil {
		return
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.itemLineRepo.AddLineItem(req, d)
	if err != nil {
		return
	}
	return 
}

func (u *itemLineUseCase) GetItemLines(req *common.RequestContext, d *dto.RequestLineItems) (
	res []dto.LineItemDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetItemLines"))
		}
	}()
	res, err = u.itemLineRepo.GetItemLines(req, d)
	return
}

func (u *itemLineUseCase) DeleteLineItem(req *common.RequestContext, d *dto.DeleteLineItemRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("DeleteLineItem"))
		}
	}()
	err = u.itemLineRepo.DeleteLineItem(req, d)
	return
}

func (u *itemLineUseCase) EditItemLine(req *common.RequestContext, i *dto.EditLineItemRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditItemLine"))
		}
	}()
	entity, err := u.GetItemLinePartyEntity(i.Body.DocPartyType)
	if err != nil {
		return
	}
	if allow := u.permission.CheckPermission(req.Ctx, req, entity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.itemLineRepo.EditItemLine(req, i)
	if err != nil {
		return
	}

	// err = u.bus.Emit(req.Ctx,)
	return
}

func (s *itemLineUseCase) closeTx(tx *query.QueryTx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}

// Returns an entity based on the party type, which can be an order, invoice, etc.
func (u *itemLineUseCase) GetItemLinePartyEntity(partyType string) (domain.EntityTemplate, error) {
	switch partyType {
	case proto.PartyType_saleOrder.String():
		return domain.SALE_ORDER, nil
	case proto.PartyType_saleInvoice.String():
		return domain.SALE_INVOICE, nil
	case proto.PartyType_purchaseOrder.String():
		return domain.PURCHASE_ORDER, nil
	case proto.PartyType_purchaseInvoice.String():
		return domain.PURCHASE_INVOICE, nil
	case proto.PartyType_purchaseReceipt.String():
		return domain.PURCHASE_RECEIPT, nil
	case proto.PartyType_deliveryNote.String():
		return domain.DELIVERY_NOTE, nil
	case proto.PartyType_salesQuotation.String():
		return domain.QUOTATION,nil
	case proto.PartyType_supplierQuotation.String():
		return domain.SUPPLIER_QUOTATION,nil
	}
	return domain.EntityTemplate{}, domain.PARTY_TYPE_NOT_FOUND
}
