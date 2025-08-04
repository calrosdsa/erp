package module_repo

import "erp/internal/domain"

func DefaultSectionEntities() map[string][]SectionEntity {

	sectionEntities := make(map[string][]SectionEntity)
	sectionEntities["Accounting"] = []SectionEntity{
		{SectionName: "Empresa", EntityID: domain.COMPANY.ID},
		{SectionName: "Empresa", EntityID: domain.LEDGER.ID},

		{SectionName: "Pagos", EntityID: domain.PAYMENT.ID},
		{SectionName: "Pagos", EntityID: domain.JOURNAL_ENTRY.ID},

		{SectionName: "Ajustes", EntityID: domain.CURRENCY_EXCHANGE.ID},
		{SectionName: "Ajustes", EntityID: domain.CHARGES_TEMPLATE.ID},

		{SectionName: "Dimensión contable", EntityID: domain.PROJECT.ID},
		{SectionName: "Dimensión contable", EntityID: domain.COST_CENTER.ID},

		{SectionName: "Reportes", EntityID: domain.GENERAL_LEDGER.ID},
		{SectionName: "Reportes", EntityID: domain.INCOME_STATEMENT.ID},
		{SectionName: "Reportes", EntityID: domain.CASH_FLOW.ID},
		{SectionName: "Reportes", EntityID: domain.BALANCE_SHEET.ID},

	}

	sectionEntities["Payables"] = []SectionEntity{
		{SectionName: "Facturación", EntityID: domain.SUPPLIER.ID},
		{SectionName: "Facturación", EntityID: domain.PURCHASE_INVOICE.ID},
		{SectionName: "Facturación", EntityID: domain.PURCHASE_RECORD.ID},

		{SectionName: "Pagos", EntityID: domain.PAYMENT.ID},
		{SectionName: "Pagos", EntityID: domain.JOURNAL_ENTRY.ID},

		{SectionName: "Reportes", EntityID: domain.ACCOUNT_PAYABLE.ID},
		{SectionName: "Reportes", EntityID: domain.ACCOUNT_PAYABLE_SUMARY.ID},	
	}

	sectionEntities["Receivables"] = []SectionEntity{
		{SectionName: "Facturación", EntityID: domain.CUSTOMER.ID},
		{SectionName: "Facturación", EntityID: domain.SALE_INVOICE.ID},
		{SectionName: "Facturación", EntityID: domain.SALES_RECORD.ID},

		{SectionName: "Pagos", EntityID: domain.PAYMENT.ID},
		{SectionName: "Pagos", EntityID: domain.JOURNAL_ENTRY.ID},

		{SectionName: "Reportes", EntityID: domain.ACCOUNT_RECEIVABLE.ID},
		{SectionName: "Reportes", EntityID: domain.ACCOUNT_RECEIVABLE_SUMARY.ID},	
	}

	sectionEntities["Buying"] = []SectionEntity{
		{SectionName: "Compra", EntityID: domain.PURCHASE_ORDER.ID},
		{SectionName: "Compra", EntityID: domain.PURCHASE_RECEIPT.ID},
		{SectionName: "Compra", EntityID: domain.SUPPLIER_QUOTATION.ID},

		{SectionName: "Artículo y Precio", EntityID: domain.ITEM.ID},
		{SectionName: "Artículo y Precio", EntityID: domain.ITEM_GROUP.ID},
		{SectionName: "Artículo y Precio", EntityID: domain.ITEM_PRICE.ID},
		{SectionName: "Artículo y Precio", EntityID: domain.PRICE_LIST.ID},

		{SectionName: "Proveedor", EntityID: domain.SUPPLIER.ID},
		{SectionName: "Proveedor", EntityID: domain.SUPPLIER_GROUP.ID},
		{SectionName: "Proveedor", EntityID: domain.CONTACT.ID},
		{SectionName: "Proveedor", EntityID: domain.ADDRESS.ID},
		
	}

	sectionEntities["Selling"] = []SectionEntity{
		{SectionName: "Venta", EntityID: domain.SALE_ORDER.ID},
		{SectionName: "Venta", EntityID: domain.DELIVERY_NOTE.ID},
		{SectionName: "Venta", EntityID: domain.QUOTATION.ID},

		{SectionName: "Artículo y Precio", EntityID: domain.ITEM.ID},
		{SectionName: "Artículo y Precio", EntityID: domain.ITEM_GROUP.ID},
		{SectionName: "Artículo y Precio", EntityID: domain.ITEM_PRICE.ID},
		{SectionName: "Artículo y Precio", EntityID: domain.PRICE_LIST.ID},

		{SectionName: "Cliente", EntityID: domain.CUSTOMER.ID},
		{SectionName: "Cliente", EntityID: domain.CUSTOMER_GROUP.ID},
		{SectionName: "Cliente", EntityID: domain.CONTACT.ID},
		{SectionName: "Cliente", EntityID: domain.ADDRESS.ID},


		
	}

	sectionEntities["Stock"] = []SectionEntity{
		{SectionName: "Catálogo de Artículos", EntityID: domain.ITEM.ID},
		{SectionName: "Catálogo de Artículos", EntityID: domain.ITEM_GROUP.ID},

		{SectionName: "Transacciones de Inventario", EntityID: domain.STOCK_ENTRY.ID},
		{SectionName: "Transacciones de Inventario", EntityID: domain.PURCHASE_RECEIPT.ID},

		{SectionName: "Reportes", EntityID: domain.STOCK_LEDGER.ID},
		{SectionName: "Reportes", EntityID: domain.STOCK_BALANCE.ID},
		{SectionName: "Reportes", EntityID: domain.SERIALNO_RESUME.ID},

		{SectionName: "Número de Serie y Lote", EntityID: domain.SERIAL_NO.ID},
		{SectionName: "Número de Serie y Lote", EntityID: domain.BATCH_BUNDLE.ID},

		{SectionName: "Precios", EntityID: domain.ITEM_PRICE.ID},
		{SectionName: "Precios", EntityID: domain.PRICE_LIST.ID},

		{SectionName: "Ajuste", EntityID: domain.WAREHOUSE.ID},
	}

	sectionEntities["Setting"] = []SectionEntity{
		{SectionName: "Autorización",EntityID: domain.ROLE.ID},
		{SectionName: "Autorización",EntityID: domain.USER.ID},
		{SectionName: "App",EntityID: domain.MODULE.ID},
	}


	return sectionEntities
}
