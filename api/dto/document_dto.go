package dto

type (
	EditAddressAndContactRequest struct {
		Body AddressAndContactData
	}
	EditDocTermRequest struct {
		Body DocTermsData
	}
	EditDocAccountingRequest struct {
		Body DocAccountingData
	}

	AddressAndContactData struct {
		ID        int64                   `json:"id" required:"true"`
		PartyType string                  `json:"party_type" required:"false"`
		Fields    AddressAndContactFields `json:"fields"`
	}
	AddressAndContactFields struct {
		ShippingAddressID *int64 `json:"shipping_address_id" required:"false"`
		BillingAddressID  *int64 `json:"billing_address_id" required:"false"`
		PartyAddressID    *int64 `json:"party_address_id" required:"false"`
		ContactID         *int64 `json:"contact_id" required:"false"`
	}

	AddressAndContact struct {
		ShippingAddress AddressDto `json:"shipping_address"`
		BillingAddress  AddressDto `json:"billing_address"`
		PartyAddress    AddressDto `json:"party_address"`
		PartyContact    ContactDto `json:"party_contact"`
		// Contact Addre
	}

	AddressAndContactDto struct {
		ShippingAddressID   *int64  `json:"shipping_address_id"`
		ShippingAddressUUID *string `json:"shipping_address_uuid"`
		ShippingAddress     *string `json:"shipping_address"`
		ContactID           *int64  `json:"contact_id"`
		ContactUUID         *string `json:"contact_uuid"`
		Contact             *string `json:"contact"`
		BillingAddressID    *int64  `json:"billing_address_id"`
		BillingAddressUUID  *string `json:"billing_address_uuid"`
		BillingAddress      *string `json:"billing_address"`
		PartyAddressID      *int64  `json:"party_address_id"`
		PartyAddressUUID    *string `json:"party_address_uuid"`
		PartyAddress        *string `json:"party_address"`
	}

	DocAccountingData struct {
		DocID        int64               `json:"doc_id"`
		DocPartyType string              `json:"doc_party_type"`
		Fields       DocAccountingFields `json:"fields"`
	}

	DocAccountingFields struct {
		DebitAccountID  *int64 `json:"debit_account_id" required:"false"`
		CreditAccountID *int64 `json:"credit_account_id" required:"false"`
	}

	DocAccountingDto struct {
		DebitAccountID    *int64  `json:"debit_account_id"`
		DebitAccountUUID  *string `json:"debit_account_uuid"`
		DebitAccount      *string `json:"debit_account"`
		CreditAccountID   *int64  `json:"credit_account_id"`
		CreditAccountUUID *string `json:"credit_account_uuid"`
		CreditAccount     *string `json:"credit_account"`
	}

	DocTermsData struct {
		DocID        int64          `json:"doc_id"`
		DocPartyType string         `json:"doc_party_type"`
		Fields       DocTermsFields `json:"fields"`
	}
	DocTermsFields struct {
		TermsAndConditionID   *int64 `json:"terms_and_condition_id" required:"false"`
		PaymentTermTemplateID *int64 `json:"payment_term_template_id" required:"false"`
	}

	DocTermsDto struct {
		TermsAndConditionID     *int64  `json:"terms_and_condition_id"`
		TermsAndCondition       *string `json:"terms_and_condition"`
		TermsAndConditionUUID   *string `json:"terms_and_condition_uuid"`
		PaymentTermTemplateID   *int64  `json:"payment_term_template_id"`
		PaymentTermTemplate     *string `json:"payment_term_template"`
		PaymentTermTemplateUUID *string `json:"payment_term_template_uuid"`
	}
)
