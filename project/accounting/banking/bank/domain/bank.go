package bank_domain

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
)

type UseCase interface {
	GetList(req *common.RequestContext, d dto.BanksRequest) (
		res dto.ResponseDataList[[]dto.BankDto], err error)
	Get(req *common.RequestContext, d dto.RequestEntity) (
		res dto.ResultEntity[dto.BankDto], err error)
	Create(req *common.RequestContext, d dto.BankData) (res dto.BankDto, err error)
	Edit(req *common.RequestContext, d dto.BankData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent) (err error)
}

type Repository interface {
	GetFilterOptions(lang string) (res []dto.FilterOptionDto)
	GetList(req *common.RequestContext, d dto.BanksRequest) (
		res []dto.BankDto, err error)
	Get(req *common.RequestContext, d dto.RequestEntity) (
		res dto.BankDto, err error)
	Create(req *common.RequestContext, d dto.BankData) (res model.Bank, err error)
	Edit(req *common.RequestContext, d dto.BankData) (err error)
	UpdateStatus(req *common.RequestContext, d *dto.UpdateStatusWithEvent,nextState string) (err error)
	
}

type Handler interface {
	GetList(ctx context.Context, d *dto.BanksRequest) (
		*dto.ResponseDataList[[]dto.BankDto],error) 
	Get(ctx context.Context, d *dto.RequestEntity) (
		*dto.EntityResponse[dto.ResultEntity[dto.BankDto]],error) 
	Create(ctx context.Context, d *dto.BankDataRequest) (
		*dto.ResponseData[dto.BankDto],error)
	Edit(ctx context.Context, d *dto.BankDataRequest) (
		*dto.ResponseMessage,error)
	UpdateStatus(ctx context.Context, d *dto.UpdateStatusWithEvent) (
		*dto.ResponseMessage, error)
}
