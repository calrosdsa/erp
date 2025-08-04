package court_repo

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"

	// "fmt"

	"gorm.io/gorm"
)

type CourtRateRepository interface {
	UpdateCourtRatesSchedule(req *common.RequestContext, i dto.UpdateCourtRatesBody) (err error)
	GetCourtRates(req *common.RequestContext, i *dto.RequestEntity) (
		res []dto.CourtRateDto, err error)
}

type courtRateRepository struct {
	conn           db.Connection
	Q              *query.Query
	currencyHelper helpers.CurrencyHelper
	convertor      helpers.ConvertorHelper
}

func NewCourtRateRepository(
	conn db.Connection,
	helpers *helpers.Helpers,
) CourtRateRepository {
	return &courtRateRepository{
		conn:           conn,
		Q:              conn.GetQ(),
		currencyHelper: helpers.Currency,
		convertor:      helpers.Convertor,
	}
}

func (r *courtRateRepository) GetCourtRates(req *common.RequestContext, i *dto.RequestEntity) (
	res []dto.CourtRateDto, err error) {
	courtRateQ := r.Q.CourtRate
	courtID := r.convertor.StrtoInt(i.ID)
	court, err := r.Q.Court.WithContext(req.Ctx).Select(r.Q.Court.ID).Where(
		r.Q.Court.ID.Eq(courtID),
		r.Q.Court.CompanyID.Eq(req.ActiveCompany.ID)).First()
	if err != nil {
		return
	}
	err = r.Q.CourtRate.WithContext(req.Ctx).Select(
		courtRateQ.Rate, courtRateQ.Enabled, courtRateQ.Currency, courtRateQ.DayWeek, courtRateQ.Time,
		courtRateQ.Currency,
	).Where(courtRateQ.CourtID.Eq(court.ID)).Scan(&res)
	if err != nil {
		return
	}

	return
}

func (r *courtRateRepository) UpdateCourtRatesSchedule(req *common.RequestContext, d dto.UpdateCourtRatesBody) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	court, err := tx.WithContext(req.Ctx).Court.
		Select(r.Q.Court.ID).Where(
		r.Q.Court.UUID.Eq(d.CourtUUID),
		r.Q.Court.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}
	var courtRates []*model.CourtRate
	for _, courtRateData := range d.CourtRates {
		m := &model.CourtRate{}
		m.CompanyID = req.ActiveCompany.ID
		m.CourtID = court.ID
		m.Enabled = courtRateData.Enabled
		m.Time = courtRateData.Time
		m.Currency = domain.DEFAULT_CURRENCY
		m.DayWeek = courtRateData.DayWeek
		m.Rate = r.currencyHelper.FloatToInt(courtRateData.Rate)
		c, err1 := r.getCourtRate(req.Ctx, tx, m)
		err = err1
		if err != nil {
			return err
		}
		switch d.Action {
		case proto.ActivityType_CREATE.String():
			if c == 0 {
				courtRates = append(courtRates, m)
			} 
		case proto.ActivityType_EDIT.String():
			if c == 0 {
				continue
			}
			_, err = tx.CourtRate.WithContext(req.Ctx).Where(
				tx.CourtRate.ID.Eq(c),
			).Updates(model.CourtRate{
				Enabled: courtRateData.Enabled,
				Rate:    r.currencyHelper.FloatToInt(courtRateData.Rate),
			})
			if err != nil {
				return
			}
		case proto.ActivityType_DELETE.String():
			if c == 0 {
				continue
			}
			_, err = tx.CourtRate.WithContext(req.Ctx).Where(
				tx.CourtRate.ID.Eq(c),
			).Delete()
			if err != nil {
				return
			}
		}

		// fmt.Println("INTERVAL ID", c, courtRateData.Enabled)
		// if c == 0 {
		// 	courtRates = append(courtRates, m)
		// } else {
		// 	_, err = tx.CourtRate.WithContext(req.Ctx).Where(
		// 		tx.CourtRate.ID.Eq(c),
		// 	).Updates(model.CourtRate{
		// 		Enabled: courtRateData.Enabled,
		// 		Rate:    int32(r.currencyHelper.FloatToInt(courtRateData.Rate)),
		// 	})
		// 	if err != nil {
		// 		return
		// 	}

		// }
		// courtRates[i] = m
	}
	err = tx.CourtRate.WithContext(req.Ctx).CreateInBatches(courtRates, len(courtRates))
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r *courtRateRepository) getCourtRate(ctx context.Context, tx *query.QueryTx, courtRate *model.CourtRate) (res int32, err error) {
	c, err := tx.CourtRate.WithContext(ctx).Select(tx.CourtRate.ID).Where(
		tx.CourtRate.CourtID.Eq(courtRate.CourtID),
		tx.CourtRate.DayWeek.Eq(courtRate.DayWeek),
		tx.CourtRate.Currency.Eq(courtRate.Currency),
		tx.CourtRate.Time.Eq(courtRate.Time),
	).First()
	if err == gorm.ErrRecordNotFound {
		return res, nil
	}
	if err != nil {
		return
	}
	return c.ID, err
}
