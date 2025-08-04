package party_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/pkg/logger"
	party_repo "erp/project/core/party/internal/repository"
	"fmt"
)

type PartyUseCase interface {
	GetUserPartyTypes(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto]
	GetPartiesByReference(req *common.RequestContext, i *dto.RequestPartyReference) (
		dto.ResultEntity[[]dto.PartyDto], error,
	)
	AddPartyReference(req *common.RequestContext, i *dto.RequestAddPartyReference) error
	GetPartyReferences(req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]dto.PartyReferenceDto], error)

	GetAllowedPartiesForReferences(req *common.RequestContext) []dto.PartyTypeDto
	GetPartyConnections(req *common.RequestContext, i *dto.RequestEntityWithParty) (
		res []dto.PartyConnections, err error,
	)
}

type partyUseCase struct {
	emitLog         logger.EmitLog
	partyRepository party_repo.PartyRepository
}

func NewPartyUseCase(
	logger logger.Logger,
	partyRepository party_repo.PartyRepository,
) PartyUseCase {
	return &partyUseCase{
		emitLog:         logger.EmitLog("party-usecase"),
		partyRepository: partyRepository,
	}
}

func (s *partyUseCase) GetPartyConnections(req *common.RequestContext, i *dto.RequestEntityWithParty) (
	res []dto.PartyConnections, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPartyConnections"))
		}
	}()
	res,err = s.partyRepository.GetPartyConnections(req,i)
	fmt.Println("CONNECTIONS",res)
	return 
}

func (s *partyUseCase) AddPartyReference(req *common.RequestContext, i *dto.RequestAddPartyReference) (err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("AddPartyReference"))
		}
	}()
	err = s.partyRepository.AddPartyReference(req.Ctx, req, i)
	return err
}

func (s *partyUseCase) GetPartyReferences(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.PartyReferenceDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPartyReferences"))
		}
	}()
	res, err = s.partyRepository.GetPartyReferences(req.Ctx, req, i)
	return res, err
}

func (s *partyUseCase) GetPartiesByReference(req *common.RequestContext, i *dto.RequestPartyReference) (
	res dto.ResultEntity[[]dto.PartyDto], err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetPartiesByReference"))
		}
	}()
	res, err = s.partyRepository.GetPartiesByReference(req.Ctx, req, i)
	if err != nil {
		return
	}
	return
}

func (s *partyUseCase) GetAllowedPartiesForReferences(req *common.RequestContext) []dto.PartyTypeDto {
	res := s.partyRepository.GetAllowedPartiesForReferences(req)
	return res
}

// func (s *partyUseCase) CreateClientParty(ctx context.Context, db *gorm.DB, client entity.Client) (model.Party, error) {
// 	ctx, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()
// 	var party model.Party
// 	party.PartyTypeCode = string(entity.PARTY_CLIENT_CODE)
// 	party.ID = int64(client.ID)
// 	err := db.WithContext(ctx).Save(&party).Error
// 	return party, err
// }

func (s *partyUseCase) GetUserPartyTypes(req *common.RequestContext) dto.ResultEntity[[]dto.PartyTypeDto] {
	res := s.partyRepository.GetUserPartyTypes(req)
	return res
}
