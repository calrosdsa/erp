package integrationservice

// import (
// 	"context"
// 	"erp/api/dto"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	"erp/internal/app/service/helpers"
// 	"time"
// )

// type TecluMobilityService struct {
// 	conn    *connection.Connection
// 	timeout time.Duration
// 	currencyHelper *helpers.CurrencyHelper
// }

// func NewTecluMobilityService(
// 	conn *connection.Connection,
// 	timeout time.Duration,
// 	helpers *helpers.Helpers,
// ) *TecluMobilityService {
// 	return &TecluMobilityService{
// 		conn:    conn,
// 		timeout: timeout,
// 		currencyHelper: helpers.Currency,
// 	}
// }

// func (s *TecluMobilityService) GetItemPrice(ctx context.Context, d *dto.TecluMobilityRequestItemPrice) (
// 	res entity.ItemPrice, err error,
// ) {
// 	ctx, cancel := context.WithTimeout(ctx, s.timeout)
// 	defer cancel()
// 	var item entity.Item
// 	companyId := int64(8)
// 	itemGroupId := uint(4)
// 	err = s.conn.Db.WithContext(ctx).Where(&entity.Item{Name: d.Body.ItemCode, CompanyID: companyId, ItemGroupID: itemGroupId}).First(&item).Error
// 	switch d.Body.Type {
// 	case "ITEM":
// 		err = s.conn.Db.WithContext(ctx).Where(&entity.ItemPrice{ItemID: item.ID, CompanyID: companyId}).First(&res).Error
// 		return
// 	case "SERVICE":
// 		err = s.conn.Db.WithContext(ctx).Where(&entity.ItemPrice{ItemID: item.ID,
// 			CompanyID: companyId, ItemQuantity: d.Body.BillingPeriod}).First(&res).Error
// 		return
// 	}
// 	return
// }


// func (s *TecluMobilityService) ToSalesOrderLineDto(ctx context.Context,d dto.OrderData)([]dto.SalesItemLineDto,error){
// 	ctx,cancel := context.WithTimeout(ctx,s.timeout)
// 	defer cancel()
// 	res := make([]dto.SalesItemLineDto,len(d.OrderLine))
// 	for i,line := range d.OrderLine{
// 		var itemPrice entity.ItemPrice
// 		var saleItemLineDto dto.SalesItemLineDto
// 		err := s.conn.Db.WithContext(ctx).Where(&entity.ItemPrice{Base: entity.Base{
// 			ID: line.ItemPriceId,
// 		}}).Preload("ItemPriceList").First(&itemPrice).Error
// 		if err != nil {
// 			return res,err
// 		}
// 		saleItemLineDto.Currency = itemPrice.ItemPriceList.Currency
// 		saleItemLineDto.ItemPriceID = line.ItemPriceId
// 		saleItemLineDto.ItemQuanitity = line.Quantity
// 		saleItemLineDto.Rate = s.currencyHelper.IntToFloat(itemPrice.Rate)
// 		res[i] = saleItemLineDto
// 	}
// 	return res,nil
// }