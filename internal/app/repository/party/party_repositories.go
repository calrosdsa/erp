package party

import (
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
)

func NewPartyRepositories(
	conn *connection.Connection,
	helpers *helpers.Helpers,
) repository.PartyRepositories {
	partyAddressRepository := NewPartyAddressRepository(conn, helpers)
	partyRepository := NewPartyRepository(conn, helpers)
	return repository.PartyRepositories{
		PartyAddress:    partyAddressRepository,
		PartyRepository: partyRepository,
	}
}
