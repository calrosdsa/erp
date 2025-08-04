package dto

import "time"

type (
	CreatePianoFormRequest struct {
		AcceptLanguageHeader
		Body CreatePianoBody
	}
	CreatePianoBody struct {
		FirstName      string    `json:"first-name"`
		LastName       string    `json:"last-name"`
		Email          string    `json:"email"`
		Phone          string    `json:"phone"`
		RentPiano      string      `json:"rent-piano"`
		MovingDate     string `json:"moving-date"`
		PianoType      string    `json:"piano-type"`
		StairsPickup   string      `json:"stairs-pickup"`
		PickupFlights  string    `json:"pickup-flights"`
		PickupStreet   string    `json:"pickup-street"`
		PickupCity     string    `json:"pickup-city"`
		PickupState    string    `json:"pickup-state"`
		PickupZip      string    `json:"pickup-zip"`
		StairsDropoff  string      `json:"stairs-dropoff"`
		DropoffFlights string    `json:"dropoff-flights"`
		DropoffStreet  string    `json:"dropoff-street"`
		DropoffCity    string    `json:"dropoff-city"`
		DropoffState   string    `json:"dropoff-state"`
		DropoffZip     string    `json:"dropoff-zip"`
	}

	PianoExportRequest struct {
		AuthParams
		Body struct {
			FromDate time.Time `json:"from_date"`
			ToDate   time.Time `json:"to_date"`
		}
	}
)
