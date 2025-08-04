package auth_admin_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	auth_admin_repo "erp/project/admin/auth/internal/repository"
)

type RoleTemplateUseCase interface {
	GetRoleTemplates(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.RoleTemplateDto], error)
	GetRoleTemplate(req *common.AdminRequestContext, d *dto.RequestEntity) (
		dto.ResultEntity[dto.RoleTemplateDto], error)
	CreateRoleTemplate(req *common.AdminRequestContext, d *dto.CreateRoleTemplateRequest) error
}

type roleTemplateUseCase struct {
	emitLog          logger.EmitLog
	roleTemplateRepo auth_admin_repo.RoleTemplateRepository
}

func NewRoleTemplateUcase(
	logger logger.Logger,
	roleTemplateRepo auth_admin_repo.RoleTemplateRepository,
) RoleTemplateUseCase {
	return &roleTemplateUseCase{
		emitLog: logger.EmitLog("role-template-usecase"),
		roleTemplateRepo: roleTemplateRepo,
	}
}

func (u *roleTemplateUseCase) GetRoleTemplates(req *common.AdminRequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.RoleTemplateDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoleTemplates"))
		}
	}()
	res, err = u.roleTemplateRepo.GetRoleTemplates(req, d)
	return
}
func (u *roleTemplateUseCase) GetRoleTemplate(req *common.AdminRequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.RoleTemplateDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetRoleTemplate"))
		}
	}()
	res, err = u.roleTemplateRepo.GetRoleTemplate(req, d)
	return
}
func (u *roleTemplateUseCase) CreateRoleTemplate(req *common.AdminRequestContext,
	d *dto.CreateRoleTemplateRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateRoleTemplate"))
		}
	}()
	err = u.roleTemplateRepo.CreateRoleTemplate(req, d)
	return
}
