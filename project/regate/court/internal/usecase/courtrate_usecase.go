package court_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	court_repo "erp/project/regate/court/internal/repository"
	regate_domain "erp/project/regate/internal/domain"
)

type CourtRateUseCase interface {
	UpdateCourtRatesSchedule(req *common.RequestContext, i dto.UpdateCourtRatesBody) (err error)
	GetCourtRates(req *common.RequestContext, i *dto.RequestEntity) (
		res []dto.CourtRateDto, err error)
}

type courtRateUseCase struct {
	permission    repository.PermissionService
	emitLog       logger.EmitLog
	courtRateRepo court_repo.CourtRateRepository
}

func NewCourtRateUseCase(
	permission repository.PermissionService,
	logger logger.Logger,
	courtRateRepo court_repo.CourtRateRepository,
) CourtRateUseCase {
	return &courtRateUseCase{
		permission:    permission,
		emitLog:       logger.EmitLog("court-rate-usecase"),
		courtRateRepo: courtRateRepo,
	}
}

func (u *courtRateUseCase) GetCourtRates(req *common.RequestContext, i *dto.RequestEntity) (
	res []dto.CourtRateDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCourtRates"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.COURT, domain.VIEW); !allow {
		return res,domain.ACTION_NOT_ALLOWED
	}
	res,err = u.courtRateRepo.GetCourtRates(req, i)
	return
}

func (u *courtRateUseCase) UpdateCourtRatesSchedule(req *common.RequestContext, i dto.UpdateCourtRatesBody) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateCourtRatesSchedule"))
		}
	}()
	if allow := u.permission.CheckPermission(req.Ctx, req, regate_domain.COURT, domain.EDIT); !allow {
		return domain.ACTION_NOT_ALLOWED
	}
	err = u.courtRateRepo.UpdateCourtRatesSchedule(req, i)
	return
}
