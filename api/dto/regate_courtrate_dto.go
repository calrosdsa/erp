package dto

type (
	UpdateCourtRatesRequest struct {
		AuthParams
		Body UpdateCourtRatesBody
	}

	UpdateCourtRatesBody struct {
		CourtRates []CourtRateData `json:"court_rate" required:"true"`
		CourtUUID  string          `json:"court_uuid" required:"true"`
		Action     string          `json:"action" required:"true"`
	}

	CourtRateData struct {
		Rate    float64 `json:"rate" required:"true"`
		Enabled bool    `json:"enabled" required:"true"`
		DayWeek int32   `json:"day_week" required:"true"`
		Time    string  `json:"time" required:"true"`
	}

	CourtRateDto struct {
		Rate     float64 `json:"rate"`
		Enabled  bool    `json:"enabled"`
		DayWeek  int32   `json:"day_week"`
		Time     string  `json:"time"`
		Currency string  `json:"currency"`
	}
)
