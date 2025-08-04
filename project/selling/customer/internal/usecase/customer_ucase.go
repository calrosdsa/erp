package customer_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/internal/domain/event"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/fsm"
	"erp/pkg/logger"
	customer_repo "erp/project/selling/customer/internal/repository"
	"fmt"
)

type CustomerUseCase interface {
	GetCustomerTypes(req *common.RequestContext) []dto.CustomerType
	CreateCustomer(req *common.RequestContext, i dto.CustomerData) (res dto.CustomerDto,err error)
	GetCustomerDetail(req *common.RequestContext, i *dto.RequestEntity) (
		res dto.ResultEntity[dto.CustomerDto], err error)
	GetCustomers(req *common.RequestContext, i *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.CustomerDto], err error)
	EditCustomer(req *common.RequestContext, d dto.CustomerData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error)
}

type customerUseCase struct {
	customerRepository customer_repo.CustomerRepository
	permission         repository.PermissionService
	emitLog            logger.EmitLog
	locale             helpers.Locale
	bus                bus.Bus
	c                  di.Container
	core               repository.CoreService
	fsm                fsm.FsmState
}

func NewCustomerUseCase(
	helpers *helpers.Helpers,
	logger logger.Logger,
	customerRepo customer_repo.CustomerRepository,
	permission repository.PermissionService,
	bus bus.Bus,
	c di.Container,
	core repository.CoreService,
	fsm fsm.FsmState,
) CustomerUseCase {
	return &customerUseCase{
		emitLog:            logger.EmitLog("customer-service"),
		customerRepository: customerRepo,
		permission:         permission,
		locale:             helpers.Locale,
		bus:                bus,
		c:                  c,
		core:               core,
		fsm:                fsm,
	}
}
func (u *customerUseCase) UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateStatus"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.SUPPLIER, domain.EDIT); err != nil {
		return
	}
	nextState, err := u.fsm.NextState(d.Body.CurrentState, d.Body.Events)
	if err != nil {
		return err
	}
	err = u.customerRepository.UpdateStatus(req, d, nextState)
	return
}

func (u *customerUseCase) EditCustomer(req *common.RequestContext, d dto.CustomerData) (err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("EditCustomer"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.CUSTOMER, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.customerRepository.EditCustomer(tx, req, d)
	if err != nil {
		return
	}

	err = u.bus.Emit(ctx, domain.EventCustomerEdited, event.CustomerEventData{
		Req:  *req,
		Tx:   tx,
		Data: d,
	})
	return err
}

func (s *customerUseCase) GetCustomerTypes(req *common.RequestContext) []dto.CustomerType {
	customerCodes := s.customerRepository.GetCustomerTypes()
	customerTypes := make([]dto.CustomerType, len(customerCodes))
	for i, code := range customerCodes {
		customerType := dto.CustomerType{}
		customerType.Code = code
		customerType.Name = s.locale.MustLocalize(
			helpers.OptionsLocale.WithID(fmt.Sprintf("Options.%s", code)),
			helpers.OptionsLocale.WithLang(string(req.LanguageCode)),
		)
		customerTypes[i] = customerType
	}
	return customerTypes
}

func (u *customerUseCase) CreateCustomer(req *common.RequestContext, i dto.CustomerData) (res dto.CustomerDto,err error) {
	ctx := u.c.Scoped(req.Ctx)
	tx := di.Get(ctx, domain.DatabaseTransactionKey).(*query.QueryTx)
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateCustomer"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	if allow := u.permission.CheckPermission(req.Ctx, req, domain.CUSTOMER, domain.CREATE); !allow {
		return res,domain.ACTION_NOT_ALLOWED
	}
	customer, err := u.customerRepository.CreateCustomer(req, tx, i)
	if err != nil {
		return
	}

	err = u.bus.Emit(ctx, domain.EventCustomerCreated, event.CustomerEventData{
		Req:      *req,
		Tx:       tx,
		Data:     i,
		Customer: customer,
	})
	res = dto.CustomerDtoFromModel(&customer)
	return 
}

func (s *customerUseCase) GetCustomerDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.CustomerDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCustomerDetail"))
		}
	}()
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.CUSTOMER, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = s.customerRepository.GetCustomerDetail(req, i)

	res.Activities = s.core.GerActivitiesByPartyID(req, res.Entity.ID)
	res.Addresses = s.core.GetPartyAddresses(req, res.Entity.ID)
	res.Contacts = s.core.GetPartyContacts(req, res.Entity.ID)
	return res, err
}

func (s *customerUseCase) GetCustomers(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.CustomerDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCustomerDetail"))
		}
	}()
	if allow := s.permission.CheckPermission(req.Ctx, req, domain.CUSTOMER, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = s.customerRepository.GetCustomers(req, i)
	return res, err
}

func (s *customerUseCase) closeTx(tx *query.QueryTx, err error) error {
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
