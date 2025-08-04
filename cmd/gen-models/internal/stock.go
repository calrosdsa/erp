package internal

import (
	"gorm.io/gen"
)

func StockModels(g *gen.Generator) []interface{} {
	item := g.GenerateModel("items")
	itemInventorySetting := g.GenerateModel("item_inventory_settings")
	itemLine := g.GenerateModel("item_lines")
	itemLineReceipt := g.GenerateModel("item_line_receipts")
	itemLineStockEntry := g.GenerateModel("item_line_stock_entries")
	deliveryLineItems := g.GenerateModel("delivery_line_items")
	invoicedItemLine := g.GenerateModel("invoiced_item_lines")
	itemPrice := g.GenerateModel("item_prices")

	stockSettings := g.GenerateModel("stock_settings")
	stockDefault := g.GenerateModel("stock_defaults")
	stockEntries := g.GenerateModel("stock_entries")

	stockTx := g.GenerateModel("stock_transactions")
	priceList := g.GenerateModel("price_lists")

	serialNo := g.GenerateModel("serial_nos")
	serialNoTransaction := g.GenerateModel("serial_no_transactions")
	batchBundle := g.GenerateModel("batch_bundles")

	return []interface{}{
		item,
		itemLine,
		itemLineReceipt,
		deliveryLineItems,
		itemLineStockEntry,
		invoicedItemLine,
		itemPrice,
		stockSettings,
		stockDefault,
		stockTx,
		stockEntries,
		priceList,
		serialNo,
		serialNoTransaction,
		batchBundle,
		itemInventorySetting,
	}
}
