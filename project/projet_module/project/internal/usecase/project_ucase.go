package project_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/logger"
	project_repo "erp/project/projet_module/project/internal/repository"
)

type ProjectUseCase interface {
	GetProject(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.ProjectDto], err error)
	CreateProject(req *common.RequestContext, d *dto.CreateProjectRequest) (
		res dto.ProjectDto, err error)
	GetProjects(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.ProjectDto], err error)
}

type projectUcase struct {
	emitLog     logger.EmitLog
	projectRepo project_repo.ProjectRepository
	permission  repository.PermissionService
	core        repository.CoreService
}

func NewProjectUcase(
	logger logger.Logger,
	projectRepo project_repo.ProjectRepository,
	permission repository.PermissionService,
	core repository.CoreService,
) ProjectUseCase {
	return &projectUcase{
		emitLog: logger.EmitLog("project-usecase"),
		projectRepo: projectRepo,
		permission: permission,
		core: core,
	}
}

func (u *projectUcase) GetProject(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.ProjectDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetProject"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PROJECT, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.projectRepo.GetProject(req, d)
	if err != nil {
		return 
	}
	res.Activities = u.core.GerActivitiesByPartyID(req,res.Entity.ID)
	return
}
func (u *projectUcase) CreateProject(req *common.RequestContext, d *dto.CreateProjectRequest) (
	res dto.ProjectDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateProject"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PROJECT, domain.CREATE)
	if err != nil {
		return
	}
	res, err = u.projectRepo.CreateProject(req, d)
	if err != nil {
		return
	}
	return
}
func (u *projectUcase) GetProjects(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.ProjectDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetProjects"))
		}
	}()
	err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.PROJECT, domain.VIEW)
	if err != nil {
		return
	}
	res, err = u.projectRepo.GetProjects(req, d)
	if err != nil {
		return
	}
	return
}
