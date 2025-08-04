package pricing_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/di"
	"erp/pkg/logger"
	order_usecase "erp/project/order/usecase"
	pricing_repo "erp/project/pricing_module/pricing/internal/repository"
	quotation_ucase "erp/project/quotation/usecase"
	item_ucase "erp/project/stock/item/usecase"
	itemprice_ucase "erp/project/stock/itemprice/usecase"
	price_list_ucase "erp/project/stock/price_list/usecase"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PricingGeneratorUcase interface {
	GeneratePo(req *common.RequestContext, d *dto.PricingDataRequest) (err error)
	GenerateQuotation(req *common.RequestContext, d *dto.PricingDataRequest) (err error)
}

type pricingGeneratorUcase struct {
	c                    di.Container
	emitLog              logger.EmitLog
	priceListUcase       price_list_ucase.PriceListUseCase
	itemUcase            item_ucase.ItemUseCase
	itemPriceUcase       itemprice_ucase.ItemPriceUseCase
	pricingGeneratorRepo pricing_repo.PricingGeneratorRepo
	orderUcase           order_usecase.OrderUseCase
	quotationUcase       quotation_ucase.QuotationUseCase
	currency             helpers.CurrencyHelper
}

func NewPricingGeneratorUcase(
	c di.Container,
	logger logger.Logger,
	pricingGeneratorRepo pricing_repo.PricingGeneratorRepo,
	helpers *helpers.Helpers,
) PricingGeneratorUcase {
	return &pricingGeneratorUcase{
		c:                    c,
		priceListUcase:       c.Get(domain.PriceListUseCase).(price_list_ucase.PriceListUseCase),
		itemUcase:            c.Get(domain.ItemUseCase).(item_ucase.ItemUseCase),
		itemPriceUcase:       c.Get(domain.ItemPriceUseCase).(itemprice_ucase.ItemPriceUseCase),
		orderUcase:           c.Get(domain.OrderUseCase).(order_usecase.OrderUseCase),
		quotationUcase:       c.Get(domain.QuotationUseCase).(quotation_ucase.QuotationUseCase),
		emitLog:              logger.EmitLog("pricing-generator-ucase"),
		pricingGeneratorRepo: pricingGeneratorRepo,
		currency:             helpers.Currency,
	}
}

func (u *pricingGeneratorUcase) GenerateQuotation(req *common.RequestContext, d *dto.PricingDataRequest) (err error) {
	tx := u.c.Get(domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GeneratePo"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	currency := string(common.CurrencyCodeBOB)
	priceListID, err := u.getPriceList(tx, req, d, false, true, currency)
	if err != nil {
		return
	}
	var lineItemsData []dto.LineItemData
	for _, line := range d.Body.PricingLineItems {
		if line.SupplierID == nil {
			continue
		}
		itemID, err := u.getItem(tx, req, line)
		if err != nil {
			return err
		}
		if itemID == 0 {
			continue
		}
		_, err = u.getItemPrice(tx, req, itemID, priceListID, *line.PrecioUnitarioTc)
		if err != nil {
			return err
		}
		lineItemData := dto.LineItemData{
			ItemID:          itemID,
			Rate:            *line.PrecioUnitarioTc,
			Quantity:        *line.Quantity,
			UnitOfMeasureID: domain.DEFAULT_UOM,
			LineType:        proto.ItemLineType_QUOTATION_LINE_ITEM.String(),
		}
		lineItemsData = append(lineItemsData, lineItemData)
	}
	quotationBody := dto.QuotationBody{
		Quotation: dto.QuotationData{
			Fields: dto.QuotationFields{
				PartyID:      *d.Body.Pricing.CustomerID,
				PostingDate:  time.Now(),
				PostingTime:  time.Now().Format(time.TimeOnly),
				Tz:           domain.DEFAULT_TZ,
				Currency:     currency,
				ValidTill:    time.Now().AddDate(0, 0, 7),
				ProjectID:    d.Body.Pricing.ProjectID,
				CostCenterID: d.Body.Pricing.CostCenterID,
				PriceListID:  &priceListID,
			},
			QuotationPartyType: proto.PartyType_salesQuotation.String(),
			References:         []*int64{&d.Body.Pricing.ID},
		},
		CreateItemLines: dto.CreateItemLines{
			Lines: lineItemsData,
		},
	}
	_, err = u.quotationUcase.CreateQuotationTx(tx, req, quotationBody)
	return
}

func (u *pricingGeneratorUcase) GeneratePo(req *common.RequestContext, d *dto.PricingDataRequest) (err error) {
	// ctx := u.c.Scoped(req.Ctx)
	tx := u.c.Get(domain.DbKey).(*query.Query).Begin()
	defer func(tx *query.QueryTx) {
		if err != nil {
			u.emitLog.Err(err, logger.OptionsLog.WithMethod("GeneratePo"))
		}
		err = u.closeTx(tx, err)
		fmt.Println("TRANSACTION ERROR", err)
	}(tx)
	currency := string(common.CurrencyCodeUSD)
	priceListID, err := u.getPriceList(tx, req, d, true, false, currency)
	if err != nil {
		return
	}
	partyLineItems := make(map[int64][]dto.LineItemData)
	for _, line := range d.Body.PricingLineItems {
		if line.SupplierID == nil {
			continue
		}
		itemID, err := u.getItem(tx, req, line)
		if err != nil {
			return err
		}
		fmt.Println("ITEM ID", itemID)
		if itemID == 0 {
			continue
		}
		_, err = u.getItemPrice(tx, req, itemID, priceListID, *line.CostAlm)
		if err != nil {
			return err
		}
		// fmt.Println("ITEM PRICE ID", itemPriceID)
		lineItemData := dto.LineItemData{
			ItemID:          itemID,
			Rate:            *line.CostAlm,
			Quantity:        *line.Quantity,
			UnitOfMeasureID: domain.DEFAULT_UOM,
			LineType:        proto.ItemLineType_ITEM_LINE_ORDER.String(),
		}
		if val, ok := partyLineItems[*line.SupplierID]; ok {
			partyLineItems[*line.SupplierID] = append(val, lineItemData)
		} else {
			partyLineItems[*line.SupplierID] = []dto.LineItemData{lineItemData}
		}
	}
	fmt.Println("map supplier", partyLineItems)

	for partyID, lines := range partyLineItems {
		var totalAmount float64
		for _, line := range lines {
			totalAmount += float64(line.Quantity) * line.Rate
		}
		fmt.Println("TOTAL AMOUNT", totalAmount)
		orderData := dto.OrderBody{
			Order: dto.OrderData{
				Fields: dto.OrderFields{
					PartyID:      partyID,
					PostingDate:  time.Now(),
					PostingTime:  time.Now().Format(time.TimeOnly),
					Tz:           domain.DEFAULT_TZ,
					Currency:     currency,
					ProjectID:    d.Body.Pricing.ProjectID,
					CostCenterID: d.Body.Pricing.CostCenterID,
					PriceListID:  &priceListID,
				},
				OrderPartyType: proto.PartyType_purchaseOrder.String(),
				TotalAmount:    totalAmount,
				References:     []*int64{&d.Body.Pricing.ID},
			},
			CreateItemLines: dto.CreateItemLines{
				Lines: lines,
			},
		}
		_, err = u.orderUcase.CreateOrderTx(tx, req, orderData)
		if err != nil {
			return err
		}
	}

	return err
}

func (u *pricingGeneratorUcase) getItemPrice(tx *query.QueryTx, req *common.RequestContext,
	itemID int64, priceListID int64, rate float64) (itemPriceID int64, err error) {
	itemPriceID, err = u.pricingGeneratorRepo.GetItemPrice(tx, req, itemID, priceListID)
	if err == gorm.ErrRecordNotFound {
		itemPrice, err := u.itemPriceUcase.CreateItemPriceTx(tx, req, dto.ItemPriceData{
			Fields: dto.ItemPriceFields{
				ItemID:          itemID,
				PriceListID:     priceListID,
				Rate:            u.currency.FloatToInt(rate),
				ItemQuantity:    1,
				UnitOfMeasureID: domain.DEFAULT_UOM,
			},
		})
		return itemPrice.ID, err
	}
	return
}

func (u *pricingGeneratorUcase) getPriceList(tx *query.QueryTx, req *common.RequestContext, d *dto.PricingDataRequest,
	isBuying bool, isSelling bool, currency string) (priceListID int64, err error) {
	complement := ""
	if isBuying {
		complement = "Compra"
	}
	if isSelling {
		complement = "Venta"
	}
	name := d.Body.Pricing.Code + " " + complement
	priceListID, err = u.pricingGeneratorRepo.GetPriceListID(tx, req, name, isBuying, isSelling)
	if err == gorm.ErrRecordNotFound {
		priceList, err := u.priceListUcase.CreatePriceListTx(tx, req, &dto.CreatePriceListRequest{
			Body: dto.PriceListData{
				Name:      name,
				IsBuying:  isBuying,
				IsSelling: isSelling,
				Currency:  currency,
			},
		})
		return priceList.ID, err
	}
	return
}

func (u *pricingGeneratorUcase) getItem(tx *query.QueryTx, req *common.RequestContext, d dto.PricingLineItemData,
) (itemID int64, err error) {
	if d.PartNumber == nil {
		return
	}
	itemID, err = u.pricingGeneratorRepo.GetItemID(tx, req, *d.PartNumber)
	if err == gorm.ErrRecordNotFound {
		item, err := u.itemUcase.CreateItemTx(tx, req, dto.ItemData{
			Fields: dto.ItemFields{
				Name:            *d.PartNumber,
				Code:              d.PartNumber,
				MaintainStock:   true,
				UnitOfMeasureID: domain.DEFAULT_UOM,
				Description:     d.Description,
			},
		})
		return item.ID, err
	}
	return
}

func (s *pricingGeneratorUcase) closeTx(tx *query.QueryTx, err error) error {
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
