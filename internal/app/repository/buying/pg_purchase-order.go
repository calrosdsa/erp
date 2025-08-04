package buying

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/connection"
	"erp/internal/app/domain"
	"erp/internal/app/domain/repository"
	"erp/internal/app/service/helpers"
)

type purchaseRepository struct {
	conn *connection.Connection
}

func NewPurchaseRepository(
	conn *connection.Connection,
	helpers *helpers.Helpers,
) repository.PurchaseRepository {
	return &purchaseRepository{
		conn: conn,
	}
}

func (s *purchaseRepository) CreatePurchaseOrder(ctx context.Context, req *common.RequestContext, i *dto.CreatePurchaseOrderRequest) (
	dto.ResultEntity[dto.OrderDto], error) {
	var (
		err error
		res dto.ResultEntity[dto.OrderDto]
	)
	tx := s.conn.Q.Begin()
	var order model.Order
	partyID, err := s.conn.Q.Order.InsertParty(domain.PARTY_PURCHASE_ORDER)
	if err != nil {
		return res, err
	}
	order.Code = s.conn.GenerateCode(ctx, s.conn.Db, &model.Order{}, req.ActiveCompany.ID)
	order.DeliveryDate = i.Body.DeliveryDate
	order.CompanyID = req.ActiveCompany.ID
	// order.Name = i.Body.Name
	order.Currency = i.Body.Currency.Code
	order.ID = partyID
	err = tx.Order.WithContext(ctx).Save(&order)
	if err != nil {
		return res, err
	}
	s.createPurchaseOrderLine(ctx, tx, i, &order)
	res.Entity = dto.OrderDtoFromModel(&order)
	err = tx.Commit()
	if err != nil {
		return res, err
	}
	return res, nil
}

func (s *purchaseRepository) createPurchaseOrderLine(ctx context.Context, tx *query.QueryTx,
	d *dto.CreatePurchaseOrderRequest, order *model.Order) error {
	orderLines := make([]*model.ItemLine, len(d.Body.Lines))
	for i, line := range d.Body.Lines {
		itemPrice, err := tx.ItemPrice.Where(s.conn.Q.ItemPrice.ID.Eq(line.ItemID)).
			First()
		if err != nil {
			return err
		}
		orderLine := model.ItemLine{}
		orderLine.ItemID = itemPrice.ID
		orderLine.PartyID = &order.ID
		orderLine.Quantity = line.Quantity
		orderLine.Rate = itemPrice.Rate
		orderLines[i] = &orderLine
	}
	err := tx.ItemLine.WithContext(ctx).CreateInBatches(orderLines, len(orderLines))
	return err
}
