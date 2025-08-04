package price_list_repo

import (
	"context"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
)

type PriceListEventRepo interface {
	CreateDefaultPriceList(ctx context.Context,tx *query.QueryTx,d model.CompanyDefault)(err error)
}

type priceListEventRepo struct {

}

func NewPriceListEventRepo(

) PriceListEventRepo{
	return &priceListEventRepo{}
}

func (r *priceListEventRepo) CreateDefaultPriceList(ctx context.Context,tx *query.QueryTx,d model.CompanyDefault)(err error){
	sellingPriceListID,err := tx.PriceList.WithContext(ctx).InsertParty(proto.PartyType_priceList.String())
	if err != nil {
		return
	}	

	buyingPriceListID,err := tx.PriceList.WithContext(ctx).InsertParty(proto.PartyType_priceList.String())
	if err != nil {
		return
	}	
	sellingPriceList := model.PriceList{
		ID: sellingPriceListID,
		Name: "Precio de Venta Predeterminado",
		Currency: d.Currency,
		IsSelling: true,
		CompanyID: d.CompanyID,

	}
	buyingPriceList := model.PriceList{
		ID: buyingPriceListID,
		Name: "Precio de Compra Predeterminado",
		Currency: d.Currency,
		IsBuying: true,
		CompanyID: d.CompanyID,
	}
	err = tx.PriceList.WithContext(ctx).CreateInBatches([]*model.PriceList{&buyingPriceList,&sellingPriceList},2)

	return
}
