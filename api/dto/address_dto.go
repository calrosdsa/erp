package dto


type (
	RequestAddresses struct {
		DefaultListParams
		Title string `query:"title" required:"false"`
	}

	

)
// type CreateAddressPartyRequest struct {
// 	Body struct {
// 		Address      AddressDto      `json:"address" required:"true"`
// 		AddressParty AddressPartyDto `json:"addressParty" required:"true"`
// 	}
// }

// type AddressPartyDto struct {
// 	IsShippingAddress bool `json:"isShippingAddress" required:"true"`
// 	IsBillingAddress  bool `json:"isBillingAddress" required:"true"`
// 	PartyID           uint `json:"partyId" required:"true"`
// }

// type AddressDto struct {
// 	FullName             string `json:"fullName" required:"true"`
// 	Company              string `json:"company" required:"false"`
// 	StreetLine1          string `json:"streetLine" required:"false"`
// 	StreetLine2          string `json:"streetLine2" required:"false"`
// 	City                 string `json:"city" required:"true"`
// 	Province             string `json:"province" required:"false"`
// 	PostalCode           string `json:"postalCode" required:"true"`
// 	PhoneNumber          string `json:"phoneNumber" required:"false"`
// 	CountryCode          string `json:"countryCode" required:"flase"`
// 	IdentificationNumber string `json:"identificationNumber" required:"false"`
// }
