package dto

import "erp/gen/db/model"

// type RequestDataPartyReference struct {
// 	UUID      string `json:"uuid"`
// 	PartyType string `json:"partyType"`
// }

type RequestPartyReference struct {
	AuthParams
	OptionalQueryParams
	PartyType string `path:"party_type"`
}

type RequestAddPartyReference struct {
	AuthParams
	Body struct {
		PartyID     int64 `json:"party_id"`
		ReferenceID int64 `json:"reference_id"`
	}
}

type PartyConnections struct {
	Connections int32  `json:"connections"`
	PartyType   string `json:"party_type"`
	PartyID     int64  `json:"party_id"`
	Entity      string `json:"entity"`
	Href        string `json:"href"`
}

type PartyTypeDto struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type PartyDto struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

type PartyReferenceDto struct {
	UUID      string `json:"uuid"`
	PartyType string `json:"code"`
	Reference string `json:"reference"`
	Name      string `json:"name"`
}

func PartyTypeDtoFromModel(m *model.PartyType) PartyTypeDto {
	r := PartyTypeDto{}
	r.Name = m.Name
	r.Code = m.Code
	return r
}
