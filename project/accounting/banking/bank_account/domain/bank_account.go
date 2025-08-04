package bank_account_domain
import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
)

type UseCase interface {
	GetList(req *common.RequestContext, d dto.BankAccountsRequest) (
		res dto.ResponseDataList[[]dto.BankAccountDto], err error)
	Get(req *common.RequestContext, d dto.RequestEntity) (
		res dto.ResultEntity[dto.BankAccountDto], err error)
	Create(req *common.RequestContext, d dto.BankAccountData) (res dto.BankAccountDto, err error)
	Edit(req *common.RequestContext, d dto.BankAccountData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error)
}

type Repository interface {
	GetFilterOptions(lang string) (res []dto.FilterOptionDto)
	GetList(req *common.RequestContext, d dto.BankAccountsRequest) (
		res []dto.BankAccountDto, err error)
	Get(req *common.RequestContext, d dto.RequestEntity) (
		res dto.BankAccountDto, err error)
	Create(req *common.RequestContext, d dto.BankAccountData) (res model.BankAccount, err error)
	Edit(req *common.RequestContext, d dto.BankAccountData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,nextState string) (err error)
}

type Handler interface {
	GetList(ctx context.Context, d *dto.BankAccountsRequest) (
		*dto.ResponseDataList[[]dto.BankAccountDto],error) 
	Get(ctx context.Context, d *dto.RequestEntity) (
		*dto.EntityResponse[dto.ResultEntity[dto.BankAccountDto]],error) 
	Create(ctx context.Context, d *dto.BankAccountDataRequest) (
		*dto.ResponseData[dto.BankAccountDto],error)
	Edit(ctx context.Context, d *dto.BankAccountDataRequest) (
		*dto.ResponseMessage,error)
	UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
		*dto.ResponseMessage, error)
}
