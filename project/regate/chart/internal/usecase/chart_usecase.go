package chart_ucase

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/proto"
	"erp/pkg/logger"
	chart_repo "erp/project/regate/chart/internal/repository"
	"fmt"
)

type ChartUseCase interface {
	GetChartData(req *common.RequestContext, i *dto.ChartDataRequest) (
		res []dto.ChartDataDto, err error)

	GetDashboardData(req *common.RequestContext, i *dto.ChartDashboardDataRequest) (
		res dto.ChartDashboardData, err error)
}

type chartUseCase struct {
	emitLog   logger.EmitLog
	chartRepo chart_repo.ChartRepository
}

func NewChartUseCase(
	logger logger.Logger,
	chartRepo chart_repo.ChartRepository,
) ChartUseCase {
	return &chartUseCase{
		emitLog:   logger.EmitLog("chart-usecase"),
		chartRepo: chartRepo,
	}
}

func (r *chartUseCase) GetDashboardData(req *common.RequestContext, i *dto.ChartDashboardDataRequest) (
	res dto.ChartDashboardData, err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetChartData"))
		}
	}()
	res.Income,err = r.chartRepo.GetIncomeChartData(req,i.Body)
	if err != nil {
		return
	}
	res.IncomeAvg,err = r.chartRepo.GetIncomeAvgChartData(req,i.Body)
	if err != nil {
		return
	}
	res.BookingHours,err = r.chartRepo.GetBookingHoursChartData(req,i.Body)
	if err != nil {
		return
	}
	i.Body.TimeUnit = proto.TimeUnit_hour.String()
	res.BookingHoursAvg,err = r.chartRepo.GetBookingHoursAvgChartData(req,i.Body)
	if err != nil {
		return
	}
	return 
}


func (r *chartUseCase) GetChartData(req *common.RequestContext, i *dto.ChartDataRequest) (
	res []dto.ChartDataDto, err error) {
	defer func() {
		if err != nil {
			r.emitLog.Err(err, logger.OptionsLog.WithMethod("GetChartData"))
		}
	}()
	fmt.Println("chart", i.Chart)
	switch i.Chart {
	case proto.ChartType_INCOME.String():
		res, err = r.chartRepo.GetIncomeChartData(req, i.Body)
	case proto.ChartType_INCOME_AVG.String():
		res, err = r.chartRepo.GetIncomeAvgChartData(req, i.Body)
	case proto.ChartType_BOOKING_HOUR.String():
		res, err = r.chartRepo.GetBookingHoursChartData(req, i.Body)
	case proto.ChartType_BOOKING_HOUR_AVG.String():
		res, err = r.chartRepo.GetBookingHoursAvgChartData(req, i.Body)
	
	}

	return
}
