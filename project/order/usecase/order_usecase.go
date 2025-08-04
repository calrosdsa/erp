package order_usecase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"fmt"

	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"

	// "erp/pkg/exporter"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	order_repo "erp/project/order/repository"
)

type OrderUseCase interface {
	CreateOrder(req *common.RequestContext, i dto.OrderBody) (
		dto.ResultEntity[dto.OrderDto], error,
	)
	CreateOrderTx(tx *query.QueryTx, req *common.RequestContext, i dto.OrderBody) (
		dto.OrderDto, error,
	)
	GetOrder(req *common.RequestContext, i *dto.RequestEntityWithParty) (dto.ResultEntity[dto.OrderDetailDto], error)
	GetOrders(req *common.RequestContext, d *dto.RequestOrders) (
		dto.PaginationResult[[]dto.OrderDto], error)
	GetEntityOrder(partyCode string) (domain.EntityTemplate, error)
	EditOrder(req *common.RequestContext, d dto.OrderBody) (err error)
	UpdateOrderStatus(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error)

	ExportOrder(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error)
}

type orderUseCase struct {
	permissionService repository.PermissionService
	emitLog           logger.EmitLog
	convertor         helpers.ConvertorHelper
	orderRepository   order_repo.OrderRepository
	orderFsm          fsm.FsmState
	core              repository.CoreService
	c                 di.Container
	bus               bus.Bus
	document repository.DocumentService
}

func NewOrderUseCase(
	helpers *helpers.Helpers,
	permissionService repository.PermissionService,
	logger logger.Logger,
	orderRepository order_repo.OrderRepository,
	orderFsm fsm.FsmState,
	core repository.CoreService,
	c di.Container,
	bus bus.Bus,
	document repository.DocumentService,
) OrderUseCase {
	return &orderUseCase{
		permissionService: permissionService,
		convertor:         helpers.Convertor,
		emitLog:           logger.EmitLog("order-service"),
		orderRepository:   orderRepository,
		orderFsm:          orderFsm,
		core:              core,
		c:                 c,
		bus:               bus,
		document: document,
	}
}

func (u *orderUseCase) ExportOrder(req *common.RequestContext, i dto.ExportDocumentData) (res []byte, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("ExportOrder"))
		}
	}()
	order,err := u.GetOrder(req,&dto.RequestEntityWithParty{
		ID: i.ID,
		PartyType: i.PartyType,
	})
	if err != nil {
		return
	}
	res,err  =u.document.GenerateOrderDocumentPdf(req,order.Entity.Order,i.PartyType)
	return
}

func (u *orderUseCase) UpdateOrderStatus(req *common.RequestContext, i *dto.UpdateStatusWithEvent) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpstateOrderState"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	orderEntity, err := u.GetEntityOrder(i.Body.PartyType)
	if err != nil {
		return domain.PARTY_TYPE_NOT_FOUND
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, orderEntity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	nextState, err := u.orderFsm.NextState(i.Body.CurrentState, i.Body.Events, i.Body.PartyType)
	if err != nil {
		return err
	}
	fmt.Println("NEXT STATE", nextState)
	_, err = u.orderRepository.UpdateOrderStatus(req, tx, i.Body.PartyID, i.Body.CurrentState, nextState)
	if err != nil {
		return
	}
	// switch nextState {
	// case proto.State_.String():
	// 	s.bus.Emit(req.Ctx, domain.InvoiceSubmittedEvent, event.OnSubmitInvoiceEventData{
	// 		Invoice: *invoice,
	// 		Tx:      tx,
	// 	})
	// }

	return
}

func (u *orderUseCase) EditOrder(req *common.RequestContext, d dto.OrderBody) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditOrder"))
		}
		err = u.closeTx(tx, err)
	}(tx)
	entity, err := u.GetEntityOrder(d.Order.OrderPartyType)
	if err != nil {
		return err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entity, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	if err = u.orderRepository.EditOrder(tx, req, d); err != nil {
		return
	}
	err = u.bus.Emit(req.Ctx, domain.OrderEditEvent, event.OrderEventData{
		Tx:   tx,
		Body: d,
	})
	return
}

func (u *orderUseCase) CreateOrder(req *common.RequestContext, i dto.OrderBody) (
	res dto.ResultEntity[dto.OrderDto], err error,
) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateOrder"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	res.Entity, err = u.createOrder(tx, req, i)
	return res, err
}

func (u *orderUseCase) CreateOrderTx(tx *query.QueryTx, req *common.RequestContext, i dto.OrderBody) (
	res dto.OrderDto, err error,
) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateOrder"))
		}
		fmt.Println("TRANSACTION ERROR", err)
	}()
	res, err = u.createOrder(tx, req, i)
	return res, err
}

func (u *orderUseCase) createOrder(tx *query.QueryTx, req *common.RequestContext, i dto.OrderBody) (
	res dto.OrderDto, err error) {
	entityT, err := u.GetEntityOrder(i.Order.OrderPartyType)
	if err != nil {
		return res, err
	}
	if allow := u.permissionService.CheckPermission(req.Ctx, req, entityT, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	order, err := u.orderRepository.CreateOrderTx(req, tx, i)
	if err != nil {
		return
	}
	res = dto.OrderDtoFromModel(&order)

	u.bus.Emit(req.Ctx, domain.OrderCreatedEvent, event.OrderEventData{
		Tx:    tx,
		Body:  i,
		Order: order,
	})
	return res, err
}

func (s *orderUseCase) GetOrder(req *common.RequestContext, i *dto.RequestEntityWithParty) (
	dto.ResultEntity[dto.OrderDetailDto], error) {
	var (
		err error
		res dto.ResultEntity[dto.OrderDetailDto]
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetOrder"))
		}
	}()
	entityT, err := s.GetEntityOrder(i.PartyType)
	if err != nil {
		return res, err
	}
	if allow := s.permissionService.CheckPermission(req.Ctx, req, entityT, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}

	res, err = s.orderRepository.GetOrder(req, i)
	if err != nil {
		return res, err
	}
	res.Activities = s.core.GerActivitiesByPartyID(req, res.Entity.Order.ID)
	return res, err
}

func (s *orderUseCase) GetOrders(req *common.RequestContext, d *dto.RequestOrders) (
	dto.PaginationResult[[]dto.OrderDto], error) {
	var (
		result dto.PaginationResult[[]dto.OrderDto]
		err    error
	)
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetOrders"))
		}
	}()
	entityT, err := s.GetEntityOrder(d.PartyType)
	if err != nil {
		return result, err
	}
	if allow := s.permissionService.CheckPermission(req.Ctx, req, entityT, domain.VIEW); !allow {
		return result, domain.ACTION_NOT_ALLOWED
	}
	result, err = s.orderRepository.GetOrders(req, d)
	if err != nil {
		return result, err
	}
	result.FilterOptions = s.orderRepository.GetFilterOptions(d.PartyType)
	return result, err
}

// get entity that the group is related
func (h *orderUseCase) GetEntityOrder(partyCode string) (domain.EntityTemplate, error) {
	switch partyCode {
	case proto.PartyType_saleOrder.String():
		return domain.SALE_ORDER, nil
	case proto.PartyType_purchaseOrder.String():
		return domain.PURCHASE_ORDER, nil
	default:
		return domain.EntityTemplate{}, domain.PARTY_TYPE_NOT_FOUND
	}
}

func (s *orderUseCase) closeTx(tx *query.QueryTx, err error) error {
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
