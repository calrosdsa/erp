package connection_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	connection_repo "erp/project/core/connection/repository"
)

type ConnectionUcase interface {
	GetConnectionsEntity(req *common.RequestContext, d dto.RequestEntity) (res []dto.ConnectionDto, err error)
}

type connectionUcase struct {
	emitLog logger.EmitLog
	repo    connection_repo.ConnectionRepository
}

func NewConnectionUcase(
	logger logger.Logger,
	repo connection_repo.ConnectionRepository,
) ConnectionUcase {
	return &connectionUcase{
		emitLog: logger.EmitLog("connection-ucase"),
		repo:    repo,
	}
}

func (u *connectionUcase) GetConnectionsEntity(req *common.RequestContext, d dto.RequestEntity) (
	res []dto.ConnectionDto, err error) {
	defer func() {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GetConnectionsEntity"))
		}
	}()
	res, err = u.repo.GetConnectionsEntity(req, d)
	if err != nil {
		return
	}
	return
}
