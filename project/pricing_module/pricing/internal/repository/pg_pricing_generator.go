package pricing_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/query"
	"erp/pkg/db"
	// w
)

type PricingGeneratorRepo interface {
	GeneratePo(tx *query.QueryTx,req *common.RequestContext, d *dto.PricingDataRequest) (err error)
	GetPriceListID(tx *query.QueryTx,req *common.RequestContext, name string, isBuying bool,
		isSelling bool) (int64, error)
	GetItemID(tx *query.QueryTx,req *common.RequestContext, partNumber string) (int64, error)
	GetItemPrice(tx *query.QueryTx,req *common.RequestContext, itemID int64, priceList int64) (
		itemPriceID int64, err error)
}

type pricingGeneratorRepo struct {
	Q *query.Query
}

func NewPricingGeneratorRepo(
	conn db.Connection,
) PricingGeneratorRepo {
	return &pricingGeneratorRepo{
		Q: conn.GetQ(),
	}
}

func (r *pricingGeneratorRepo) GeneratePo(tx *query.QueryTx,req *common.RequestContext, d *dto.PricingDataRequest) (err error) {
	return nil
}

func (r *pricingGeneratorRepo) GetItemID(tx *query.QueryTx,req *common.RequestContext, partNumber string) (int64, error) {
	item, err := tx.Item.WithContext(req.Ctx).Select(tx.Item.ID).Where(
		tx.Item.CompanyID.Eq(req.ActiveCompany.ID),
		tx.Item.Code.Eq(partNumber),
	).First()
	if err != nil {
		return 0, err
	}
	return item.ID, err

}

func (r *pricingGeneratorRepo) GetItemPrice(tx *query.QueryTx,req *common.RequestContext, itemID int64, priceList int64) (
	itemPriceID int64, err error) {
	itemPrice, err := tx.ItemPrice.WithContext(req.Ctx).Select(tx.ItemPrice.ID).Where(
		tx.ItemPrice.CompanyID.Eq(req.ActiveCompany.ID),
		tx.ItemPrice.ItemID.Eq(itemID),
		tx.ItemPrice.PriceListID.Eq(priceList),
	).First()
	if err != nil {
		return 0, err
	}
	return itemPrice.ID, err

}

func (r *pricingGeneratorRepo) GetPriceListID(tx *query.QueryTx,req *common.RequestContext, name string, isBuying bool,
	isSelling bool) (int64, error) {
	priceList, err := tx.PriceList.WithContext(req.Ctx).Select(tx.PriceList.ID).Where(
		tx.PriceList.CompanyID.Eq(req.ActiveCompany.ID),
		tx.PriceList.Name.Eq(name),
		tx.PriceList.IsSelling.Is(isSelling),
		tx.PriceList.IsBuying.Is(isBuying),
	).First()
	if err != nil {
		return 0, err
	}
	return priceList.ID, err

}

func (r *pricingGeneratorRepo) generateItems(req *common.RequestContext, lines []dto.PricingData) (err error) {
	// items  := make([]*model.Item,len(lines))
	// for items, line := range lines {
	// 	item := &model.Item{}
	// 	item.UnitOfMeasureID = domain.DEFAULT_UOM

	// 	// do something
	// }
	return nil
}
