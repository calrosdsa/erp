package sellingservice

// import (
// 	// "erp/api/common"
// 	"context"
// 	"erp/api/common"
// 	"erp/api/dto"
// 	"erp/gen/db/model"
// 	"erp/internal/app/config"
// 	sellingconfig "erp/internal/app/config/selling_config"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	"erp/internal/app/plugin"
// 	"erp/internal/app/service/helpers"
// 	"fmt"
// 	"time"
// )

// type SalesOrderService struct {
// 	conn            *connection.Connection
// 	timeout         *time.Duration
// 	currencyHelper  *helpers.CurrencyHelper
// 	generatorHelper *helpers.GeneratorHelper
// 	plugins         *plugin.PluginModule
// }

// func NewSalesOrderService(conn *connection.Connection, timeout *time.Duration, helpers *helpers.Helpers,
// 	plugins *plugin.PluginModule) *SalesOrderService {
// 	return &SalesOrderService{
// 		conn:            conn,
// 		timeout:         timeout,
// 		currencyHelper:  helpers.Currency,
// 		generatorHelper: helpers.Generator,
// 		plugins:         plugins,
// 	}
// }

// func (s *SalesOrderService) CreateSaleOrderService(req *common.RequestContext, d *dto.CreateSalesOrderBody) (
// 	res dto.EntityResponse[entity.SalesOrder],err error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
// 	defer cancel()
// 	tx := s.conn.Db.Begin()
// 	defer func() {
		
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()
// 	if tx.Error != nil {
// 		return
// 	}
// 	var salesOrder entity.SalesOrder

// 	salesOrder.PartyID = d.ClientID
// 	salesOrder.CompanyID = uint(req.ActiveCompany.ID)
// 	salesOrder.OrderType = d.OrderType

// 	if d.BillingAddressID != nil {
// 		salesOrder.BillingAddressID = d.BillingAddressID
// 	}
// 	if d.ShippingAddressID != nil {
// 		salesOrder.ShippingAddressID = d.ShippingAddressID
// 	}
	

// 	disableOrder := d.DeleteAt.Valid

// 	if disableOrder {
// 		salesOrder.DeletedAt = d.DeleteAt
// 	}

// 	salesOrder.Code = s.conn.GenerateCode(ctx,tx,&entity.SalesOrder{},req.ActiveCompany.ID)
// 	salesOrder.DeliveryDate = d.DeliveryDate

// 	tx.WithContext(ctx).Save(&salesOrder)

// 	salesOrderPlugins := make([]entity.SalesOrderPlugin, len(d.Plugins))
// 	for i, plugin := range d.Plugins {
// 		var salesOrderPlugin entity.SalesOrderPlugin
// 		salesOrderPlugin.Plugin = plugin
// 		salesOrderPlugin.SalesOrderID = salesOrder.ID
// 		salesOrderPlugins[i] = salesOrderPlugin
// 	}
// 	salesOrder.SalesOrderPlugin = salesOrderPlugins
// 	tx.WithContext(ctx).Save(&salesOrder)

// 	fmt.Println("NEW ORDER ", salesOrder)
// 	orderLines := make([]entity.SalesItemLine, len(d.SalesItemLines))
// 	for i, item := range d.SalesItemLines {
// 		orderLineItem := entity.SalesItemLine{}
// 		var nItem model.ItemPrice
// 		nItem.CompanyID = req.ActiveCompany.ID
// 		nItem.ID = item.ItemPriceID
// 		err = s.conn.GetEntity(ctx, tx, &nItem)
// 		if err != nil {
// 			return 
// 		}
// 		orderLineItem.ItemPrice = nItem

// 		orderLineItem.ItemPriceID = nItem.ID
// 		orderLineItem.Currency = item.Currency
// 		orderLineItem.Rate = s.currencyHelper.FloatToInt(item.Rate)
// 		orderLineItem.ItemQuantity = item.ItemQuanitity
// 		orderLineItem.SalesOrderID = salesOrder.ID
// 		orderLines[i] = orderLineItem
// 	}
// 	err = tx.WithContext(ctx).CreateInBatches(orderLines, len(orderLines)).Error
// 	if err != nil {
// 		return
// 	}

// 	err = tx.Commit().Error
// 	if err != nil {
// 		return
// 	}
// 	res.Body.Result = salesOrder
	
// 	for _,plugin := range d.Plugins {
// 		if plugin == config.PLUGIN_SQUARE && disableOrder {
// 			strategyResp,err1 := s.getSalesOrderStrategy(plugin).CreateSalesOrder(req,&salesOrder,orderLines)
// 			res.Body.Result.Data = strategyResp
// 			if err1 != nil {
// 				err = err1
// 				return
// 			}
// 		}
// 	}
// 	return 
// }

// func (s *SalesOrderService) GetSalesOrders(req *common.RequestContext, d *dto.RequestPaginationData) (dto.PaginationResult[[]entity.SalesOrder], error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
// 	defer cancel()
// 	var result dto.PaginationResult[[]entity.SalesOrder]

// 	queryBuilder := s.conn.Db.WithContext(ctx).Table("sales_orders").
// 		Where("company_id = ?", req.ActiveCompany.ID)

// 	if d.Query != "" {
// 		queryBuilder = queryBuilder.Where("name ILIKE ?", "%"+d.Query+"%")
// 	}

// 	if req.Session.Role == entity.ROLE_CLIENT {
// 		queryBuilder.Where("party_id = ?", req.GetClientID())
// 	}

// 	err := queryBuilder.
// 		Scopes(s.conn.Paginate(req.Params)).
// 		Find(&result.Results).Error
// 	if err != nil {
// 		return result, err
// 	}

// 	err = queryBuilder.
// 		Count(&result.Total).Error
// 	return result, err
// }

// func (s *SalesOrderService) GetSalesOrderDetail(req *common.RequestContext, d *dto.RequestSalesOrderDetail) (
// 	res dto.ResponseSalesOrderDetail, err error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
// 	defer cancel()
// 	queryBuilder := s.conn.Db.WithContext(ctx).Where("company_id = ? and code = ?", req.ActiveCompany.ID, d.OrderCode)

// 	if req.Session.Role == entity.ROLE_CLIENT {
// 		queryBuilder.Where("client_id = ?", req.GetClientID())
// 	}
// 	err = queryBuilder.Preload("SalesOrderPlugin").Preload("Company").First(&res.Body.SalesOrder).Error

// 	err = s.conn.Db.WithContext(ctx).Preload("ItemPrice.Item").
// 		Preload("ItemPrice.Tax").Find(&res.Body.SalesOrderItems, "sales_order_id = ?", res.Body.SalesOrder.ID).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	for i,plugin := range res.Body.SalesOrder.SalesOrderPlugin {
// 		dataString,err := s.getSalesOrderStrategy(plugin.Plugin).GetSalesOrderDetail(req,&res.Body.SalesOrder,res.Body.SalesOrderItems)
// 		if err != nil {
// 			fmt.Println("FAIL TO GET DATA PLUGIN",err,plugin.Plugin)
// 		}
// 		// fmt.Println("DATA STRING")
// 		plugin.Data = dataString
// 		res.Body.SalesOrder.SalesOrderPlugin[i] =  plugin
// 	}
// 	err = nil
// 	return
// }

// func (s *SalesOrderService) getSalesOrderStrategy(plugin string) sellingconfig.SalesOrderStrategy {
// 	switch plugin {
// 	case config.PLUGIN_SQUARE:
// 		return s.plugins.GetPlugin(config.PLUGIN_SQUARE).SalesOrderStrategy
// 	}
// 	return &sellingconfig.DefaultSalesOrderStrategy{}
// }
