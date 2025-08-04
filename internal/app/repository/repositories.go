package repository

import (
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/repository/buying"
	"erp/internal/app/repository/party"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
)

func NewRepositories(
	conn *connection.Connection,
	db db.Connection,
	helpers *helpers.Helpers,
) *repository.Repositories {
	partyRepositories := party.NewPartyRepositories(conn, helpers)
	buyingRepository := buying.NewBuyingRepository(conn, helpers)

	
	return &repository.Repositories{
		PartyRepositories:  partyRepositories,
		BuyingRepository:   buyingRepository,
		
	}
}
