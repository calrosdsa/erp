package buying

import (
	"erp/internal/app/connection"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
)

func NewBuyingRepository(
	conn *connection.Connection,
	helpers *helpers.Helpers,
) repository.BuyingRepository {
	purchaseRepository := NewPurchaseRepository(conn, helpers)
	return repository.BuyingRepository{
		PurchaseRepository: purchaseRepository,
	}
}
