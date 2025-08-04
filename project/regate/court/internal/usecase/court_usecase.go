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

type CourtUseCase interface {
	CreateCourt(req *common.RequestContext, i dto.CreateCourtBody) (dto.CourtDto, error)
	GetCourt(req *common.RequestContext, i *dto.RequestEntity) (dto.ResultEntity[dto.CourtDto], error)
	GetCourts(req *common.RequestContext, i dto.CourtsRequest) (
		res dto.ResponseDataList[[]dto.CourtDto], err error)
	EditCourt(req *common.RequestContext,d dto.EditCourtBody)(err error)
}

type courtUseCase struct {
	emitLog    logger.EmitLog
	permission repository.PermissionService
	courtRepo  court_repo.CourtRepository
	core repository.CoreService
}

func NewCourtUseCase(
	logger logger.Logger,
	permission repository.PermissionService,
	courtRepo court_repo.CourtRepository,
	core repository.CoreService,
) CourtUseCase {
	return &courtUseCase{
		emitLog:    logger.EmitLog("court-usecase"),
		permission: permission,
		courtRepo:  courtRepo,
		core: core,
	}
}

func (r *courtUseCase) EditCourt(req *common.RequestContext,d dto.EditCourtBody)(err error){
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("EditCourt"))
		}
	}()	
	err = r.permission.CheckPermissionEntity(req.Ctx,req,regate_domain.COURT,domain.EDIT)
	if err != nil {
		return
	}
	err = r.courtRepo.EditCourt(req,d)
	return
}

func (r *courtUseCase) CreateCourt(req *common.RequestContext, i dto.CreateCourtBody) (res dto.CourtDto, err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateCourt"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.COURT, domain.CREATE); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	court, err := r.courtRepo.CreateCourt(req, i)
	res = dto.CourtDtoFromModel(&court)
	return
}
func (r *courtUseCase) GetCourt(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.CourtDto], err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCourts"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.COURT, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res, err = r.courtRepo.GetCourt(req, i)
	if err != nil {
		return
	}
	res.Activities = r.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}
func (r *courtUseCase) GetCourts(req *common.RequestContext, d dto.CourtsRequest) (
	res dto.ResponseDataList[[]dto.CourtDto], err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCourts"))
		}
	}()
	if allow := r.permission.CheckPermission(req.Ctx, req, regate_domain.COURT, domain.VIEW); !allow {
		return res, domain.ACTION_NOT_ALLOWED
	}
	res.Body.Result, err = r.courtRepo.GetCourts(req, d)
	if err != nil {
		return
	}
	res.Body.FilterOptions = r.courtRepo.GetFilterOptions(string(req.LanguageCode))
	return
}
