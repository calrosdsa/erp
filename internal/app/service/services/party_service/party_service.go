package partyservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
	"erp/pkg/logger"
	"time"
)

type partyService struct {
	conn            *connection.Connection
	timeout         time.Duration
	emitLog         helpers.EmitLog
	partyRepository repository.PartyRepository
}

func NewPartyService(
	conn *connection.Connection,
	timeout time.Duration,
	helpers *helpers.Helpers,
	repositories *repository.Repositories,
) repository.PartyService {
	return &partyService{
		conn:            conn,
		timeout:         timeout,
		emitLog:         helpers.Logger.EmitLog("party-service"),
		partyRepository: repositories.PartyRepositories.PartyRepository,
	}
}

func (s *partyService) AddPartyReference(req *common.RequestContext, i *dto.RequestAddPartyReference) (err error) {
	ctx,cancel := context.WithTimeout(req.Ctx,s.timeout)
	defer func ()  {
		cancel()
		if err != nil {
			s.emitLog.Err(err,logger.OptionsLog.WithMethod("AddPartyReference"))
		}
	}()
	err = s.partyRepository.AddPartyReference(ctx,req,i)
	return err
}

func (s *partyService) GetPartyReferences(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.PartyReferenceDto],err error) {
	ctx,cancel := context.WithTimeout(req.Ctx,s.timeout)
	defer func ()  {
		cancel()
		if err != nil {
			s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetPartyReferences"))
		}
	}()
	res,err = s.partyRepository.GetPartyReferences(ctx,req,i)
	return res,err
}

func (s *partyService) GetPartiesByReference(req *common.RequestContext, i *dto.RequestPartyReference) (
	res dto.ResultEntity[[]dto.PartyDto], err error) {
	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
	defer func() {
		cancel()
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPartiesByReference"))
		}
	}()
	res, err = s.partyRepository.GetPartiesByReference(ctx, req, i)
	if err != nil {
		return
	}
	return
}

func (s *partyService) GetAllowedPartiesForReferences(req *common.RequestContext) []dto.PartyTypeDto {
	res := s.partyRepository.GetAllowedPartiesForReferences(req)
	return res
}


// func (s *partyService) CreateClientParty(ctx context.Context, db *gorm.DB, client entity.Client) (model.Party, error) {
// 	ctx, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()
// 	var party model.Party
// 	party.PartyTypeCode = string(entity.PARTY_CLIENT_CODE)
// 	party.ID = int64(client.ID)
// 	err := db.WithContext(ctx).Save(&party).Error
// 	return party, err
// }

func (s *partyService) GetUserPartyTypes(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto] {
	res := s.partyRepository.GetUserPartyTypes(req)
	return res
}
