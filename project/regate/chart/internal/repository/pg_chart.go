package chart_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/proto"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gorm"
)

type ChartRepository interface {
	GetIncomeChartData(req *common.RequestContext, i dto.ChartDataBody) (
		res []dto.ChartDataDto, err error)
	GetIncomeAvgChartData(req *common.RequestContext, i dto.ChartDataBody) (
		res []dto.ChartDataDto, err error)
	GetBookingHoursChartData(req *common.RequestContext, i dto.ChartDataBody) (
		res []dto.ChartDataDto, err error)

	GetBookingHoursAvgChartData(req *common.RequestContext, i dto.ChartDataBody) (
		res []dto.ChartDataDto, err error)
}

type chartRepository struct {
	conn db.Connection
	DB   *gorm.DB
}

func NewChartRepository(
	conn db.Connection,
) ChartRepository {
	return &chartRepository{
		conn: conn,
		DB:   conn.GetDB(),
	}
}

func (r *chartRepository) GetIncomeChartData(req *common.RequestContext, i dto.ChartDataBody) (res []dto.ChartDataDto, err error) {
	// Prepare the query using named parameters for better readability
	query := `
	SELECT date_trunc($3, d.date) AS name,
		COALESCE(SUM(rc.paid_amount) FILTER (WHERE rc.type = $5), 0) AS value,
		COALESCE(SUM(rc.paid_amount) FILTER (WHERE rc.type = $6), 0) AS value2
	FROM (
		SELECT date_trunc('day', ($1::date - offs)) AS date
		FROM generate_series(0, ($1::date - $2::date), 1) AS offs
	) d
	LEFT JOIN r_booking_slots rc
	ON d.date = date_trunc('day', rc.datetime AT TIME ZONE 'America/La_Paz') AND rc.company_id = $4
	GROUP BY date_trunc($3, d.date)
	ORDER BY date_trunc($3, d.date);
	`

	// Execute the query with proper parameters
	err = r.DB.Raw(query, i.EndDate, i.StartDate, i.TimeUnit, req.ActiveCompany.ID,
		proto.BookingType_LOCAL_BOOKING.String(),
		proto.BookingType_EVENT_BOOKING.String()).Scan(&res).Error

	// Return the result and any potential error
	return res, err
}
func (r *chartRepository) GetIncomeAvgChartData(req *common.RequestContext, i dto.ChartDataBody) (
	res []dto.ChartDataDto, err error) {
	timeUnit, generateCount, startStep := r.timeUnitFromAvgData(i.TimeUnit)
	query := fmt.Sprintf(`select val AS name,
		coalesce(sum(rc.paid_amount) filter (where rc.type = $5),0) AS value,
		coalesce(sum(rc.paid_amount) filter (where rc.type = $6),0) AS value2
			from  generate_series($7,$4,1) as val
		left join r_booking_slots as rc on EXTRACT(%s FROM rc.datetime AT TIME ZONE 'America/La_Paz') = val
		and  rc.company_id = $3
		and  rc.datetime::date >= $2 AND rc.datetime::date <= $1
		GROUP BY val order by val;`, timeUnit)
	err = r.DB.Raw(query, i.EndDate, i.StartDate, req.ActiveCompany.ID,
		generateCount, proto.BookingType_LOCAL_BOOKING.String(),
		proto.BookingType_EVENT_BOOKING.String(), startStep).Scan(&res).Error
	return res, err
}

func (r *chartRepository) GetBookingHoursChartData(req *common.RequestContext, i dto.ChartDataBody) (
	res []dto.ChartDataDto, err error) {
	query := `
		SELECT date_trunc($3, d.date) AS name,
			COALESCE(count(rc.id) FILTER (WHERE rc.type = $5), 0) AS value,
			COALESCE(count(rc.id) FILTER (WHERE rc.type = $6), 0) AS value2
		FROM (
			SELECT date_trunc('day', ($1::date - offs)) AS date
			FROM generate_series(0, ($1::date - $2::date), 1) AS offs
		) d
		LEFT JOIN r_booking_slots rc
		ON d.date = date_trunc('day', rc.datetime AT TIME ZONE 'America/La_Paz') AND rc.company_id = $4
		GROUP BY date_trunc($3, d.date)
		ORDER BY date_trunc($3, d.date);
		`

	// Execute the query with proper parameters
	err = r.DB.Raw(query, i.EndDate, i.StartDate, i.TimeUnit, req.ActiveCompany.ID,
		proto.BookingType_LOCAL_BOOKING.String(),
		proto.BookingType_EVENT_BOOKING.String()).Scan(&res).Error

	// Return the result and any potential error
	return res, err
}

func (r *chartRepository) GetBookingHoursAvgChartData(req *common.RequestContext, i dto.ChartDataBody) (
	res []dto.ChartDataDto, err error) {
	timeUnit, generateCount, startStep := r.timeUnitFromAvgData(i.TimeUnit)
	query := fmt.Sprintf(`select val AS name,
		coalesce(count(rc.id) filter (where rc.type = $5),0) AS value,
		coalesce(count(rc.id) filter (where rc.type = $6),0) AS value2
			from  generate_series($7,$4,1) as val
		left join r_booking_slots as rc on EXTRACT(%s FROM rc.datetime AT TIME ZONE 'America/La_Paz') = val
		and  rc.company_id = $3
		and  rc.datetime::date >= $2 AND rc.datetime::date <= $1
		GROUP BY val order by val;`, timeUnit)
	err = r.DB.Raw(query, i.EndDate, i.StartDate, req.ActiveCompany.ID,
		generateCount, proto.BookingType_LOCAL_BOOKING.String(),
		proto.BookingType_EVENT_BOOKING.String(), startStep).Scan(&res).Error
	return res, err
}

func (r *chartRepository) timeUnitFromAvgData(timeUnit string) (string, int, int) {
	switch timeUnit {
	case proto.TimeUnit_hour.String():
		return "hour", 23, 0
	case proto.TimeUnit_day.String():
		return "dow", 6, 0
	case proto.TimeUnit_month.String():
		return "month", 12, 1
	default:
		return "dow", 6, 0
	}
}
