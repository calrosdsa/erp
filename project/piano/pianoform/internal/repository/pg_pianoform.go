package pianoform_repo

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
	piano_domain "erp/project/piano/intenal/domain"
	"fmt"
	"strings"
	"time"

	"gorm.io/gen/field"
)

type PianoFormRepository interface {
	CreatePianoForm(ctx context.Context, i dto.CreatePianoBody) (err error)
	GetPianoForms(req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]*model.PianoForm], error)
	GetPianoForm(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ResultEntity[*model.PianoForm], error)
	ExportData(req *common.RequestContext, i *dto.PianoExportRequest) (
		res []*model.PianoForm, err error)
}

type pianoFormRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewPianoReposotory(
	conn db.Connection,
	helpers *helpers.Helpers,
) PianoFormRepository {
	return &pianoFormRepository{
		Q:         conn.GetQ(),
		convertor: helpers.Convertor,
	}
}

func (r *pianoFormRepository) ExportData(req *common.RequestContext, i *dto.PianoExportRequest) (
	res []*model.PianoForm, err error) {
	res, err = r.Q.PianoForm.WithContext(req.Ctx).Where(
		r.Q.PianoForm.CompanyID.Eq(req.ActiveCompany.ID),
		r.Q.PianoForm.CreatedAt.Between(i.Body.FromDate, i.Body.ToDate),
	).Limit(domain.MAX_LIMIT).Find()
	return
}

func (r *pianoFormRepository) GetPianoForms(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]*model.PianoForm], err error) {
	var (
		order field.Expr
	)
	builder := r.Q.PianoForm.WithContext(req.Ctx).Where(
		r.Q.PianoForm.CompanyID.Eq(req.ActiveCompany.ID),
	)
	res.Total, err = builder.Count()
	if err != nil {
		return
	}

	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
	orderCol, ok := r.Q.PianoForm.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
	if ok {
		fmt.Println("ORDER", i.Order)
		if strings.ToUpper(i.Order) == "ASC" {
			fmt.Println("ORDER", "ASC")
			order = orderCol.Asc()
		} else {
			fmt.Println("ORDER", "DESC")
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Results, err = builder.Find()
	return
}
func (r *pianoFormRepository) GetPianoForm(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[*model.PianoForm], err error) {
	res.Entity, err = r.Q.PianoForm.WithContext(req.Ctx).Where(
		r.Q.PianoForm.CompanyID.Eq(req.ActiveCompany.ID),
		r.Q.PianoForm.ID.Eq(r.convertor.StrtoInt(i.ID)),
	).First()

	return
}

func (r *pianoFormRepository) CreatePianoForm(ctx context.Context, i dto.CreatePianoBody) (err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	partyID, err := tx.PianoForm.InsertParty(proto.PianoPartyType_pianoForms.String())
	if err != nil {
		return
	}
    // Parse the date string
    layout := "2006-01-02T15:04"
    t, err := time.Parse(layout, i.MovingDate)
    if err != nil {
        fmt.Println("Error parsing date:", err)
        return
    }
	pianoForm := model.PianoForm{
		ID:             partyID,
		FirstName:      i.FirstName,
		LastName:       i.LastName,
		PhoneNumber:    i.Phone,
		Email:          i.Email,
		RentPiano:      i.RentPiano == "1",
		MovingDate:     t, // Assuming you need to parse the date from string
		PianoType:      i.PianoType,
		StairsPickup:   i.StairsPickup == "1",
		PickupFlights:  i.PickupFlights,
		PickupStreet:   i.PickupStreet,
		PickupCity:     i.PickupCity,
		PickupState:    i.PickupState,
		PickupZip:      i.PickupZip,
		StairsDropoff:  i.StairsDropoff == "1",
		DropoffFlights: i.DropoffFlights,
		DropoffStreet:  i.DropoffStreet,
		DropoffCity:    i.DropoffCity,
		DropoffState:   i.DropoffState,
		DropoffZip:     i.DropoffZip,
		CompanyID:      piano_domain.DEFAULT_COMPANY,
	}

	err = tx.PianoForm.Save(&pianoForm)
	if err != nil {
		return
	}
	err = tx.Commit()
	return err
}
