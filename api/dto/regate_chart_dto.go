package dto


type (
	ChartDataRequest struct {
		Chart string `path:"chart"`
		Body  ChartDataBody
	}
	ChartDashboardDataRequest struct {
		Body  ChartDataBody
	}
	ChartDataBody struct {
		StartDate string `json:"start_date" required:"true"`
		EndDate   string `json:"end_date" required:"true"`

		TimeUnit string `json:"time_unit" required:"false"`
	}

	ChartDataDto struct {
		Name   string   `json:"name"`
		Value  float64  `json:"value"`
		Value2 *float64 `json:"value2"`
	}

	ChartDashboardData struct {
		Income []ChartDataDto `json:"income"`
		IncomeAvg []ChartDataDto `json:"income_avg"`
		BookingHours []ChartDataDto `json:"booking_hours"`
		BookingHoursAvg []ChartDataDto `json:"booking_hours_avg"`
	}
)

