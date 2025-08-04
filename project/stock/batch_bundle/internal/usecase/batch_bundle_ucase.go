package batch_bundle_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	batch_bundle_repo "erp/project/stock/batch_bundle/internal/repository"
)

type BatchBundleUseCase interface {
	GetBatchBundle(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.BatchBundleDto], err error)

	GetBatchBundles(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.BatchBundleDto], err error)
}

type batchBundleUcase struct {
	emitLog     logger.EmitLog
	batchBundle batch_bundle_repo.BatchBundleRepository
	permission  repository.PermissionService
	core        repository.CoreService
}

func NewBatchBundleUcase(
	logger logger.Logger,
	batchBundle batch_bundle_repo.BatchBundleRepository,
	permission repository.PermissionService,
	core repository.CoreService,
) BatchBundleUseCase {
	return &batchBundleUcase{
		emitLog: logger.EmitLog("batch-bundle-usecase"),
		batchBundle: batchBundle,
		permission: permission,
		core: core,
	}
}

func (u *batchBundleUcase) GetBatchBundle(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.BatchBundleDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetBatchBundle"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.BATCH_BUNDLE, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.batchBundle.GetBatchBundle(req, d)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}

func (u *batchBundleUcase) GetBatchBundles(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.BatchBundleDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetBatchBundles"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.BATCH_BUNDLE, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.batchBundle.GetBatchBundles(req, d)
	if err != nil {
		return
	}
	return
}
