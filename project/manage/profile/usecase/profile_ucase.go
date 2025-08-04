package profile_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/internal/domain"
	"erp/internal/domain/repository"
	"erp/pkg/bus"
	"erp/pkg/di"
	"erp/pkg/logger"
	profile_repo "erp/project/manage/profile/repository"
)

type ProfileUseCase interface {
	UpdateProfileSession(req *common.RequestContext, i *dto.UpdateProfileRequest) error
	GetUserProfileDetail(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ResultEntity[dto.ProfileDto], error)
	GetProfileSession(req *common.RequestContext) (
		res dto.ResultEntity[dto.ProfileDto], err error)
	GetProfiles(req *common.RequestContext, d dto.ProfilesRequest) (res []dto.ProfileDto, err error)
	GetCompanyUserProfiles(req *common.RequestContext, d *dto.RequestPaginationData) (
		res []dto.ProfileDto, err error)
}

type profileUcase struct {
	permission repository.PermissionService
	core       repository.CoreService
	repo       profile_repo.ProfileRepository
	emitLog    logger.EmitLog
	bus        bus.Bus
	c          di.Container
}

func NewProfileUseCase(
	permission repository.PermissionService,
	core repository.CoreService,
	repo profile_repo.ProfileRepository,
	logger logger.Logger,
	bus bus.Bus,
	c di.Container,
) ProfileUseCase {
	return &profileUcase{
		core:    core,
		permission: permission,
		repo:    repo,
		emitLog: logger.EmitLog("profile-usecase"),
		bus:     bus,
		c:       c,
	}
}

func (u *profileUcase) UpdateProfileSession(req *common.RequestContext, i *dto.UpdateProfileRequest) (err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("UpdateProfileSession"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.USER, domain.EDIT); err != nil {
		return
	}
	err = u.repo.UpdateProfileSession(req, i)
	return
}

func (u *profileUcase) GetUserProfileDetail(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[dto.ProfileDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetUserProfileDetail"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.USER, domain.VIEW); err != nil {
		return
	}
	res.Entity, err = u.repo.GetUserProfileDetail(req, i)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}
func (u *profileUcase) GetProfileSession(req *common.RequestContext) (
	res dto.ResultEntity[dto.ProfileDto], err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetProfileSession"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.USER, domain.VIEW); err != nil {
		return
	}
	res.Entity, err = u.repo.GetProfileSession(req)
	return
}
func (u *profileUcase) GetProfiles(req *common.RequestContext, d dto.ProfilesRequest) (res []dto.ProfileDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetProfiles"))
		}
	}()
	if err = u.permission.CheckPermissionEntity(req.Ctx, req, domain.USER, domain.VIEW); err != nil {
		return
	}
	res, err = u.repo.GetProfiles(req, d)
	return
}
func (u *profileUcase) GetCompanyUserProfiles(req *common.RequestContext, d *dto.RequestPaginationData) (
	res []dto.ProfileDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetCompanyUserProfiles"))
		}
	}()
	res,err = u.repo.GetCompanyUserProfiles(req,d)
	return
}

func (u *profileUcase) closeTx(tx *query.QueryTx, err error) error {
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
