package pianoform_ucase

import (
	"bytes"
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/service/helpers"
	"erp/internal/domain/repository"
	"erp/pkg/exporter"
	"erp/pkg/logger"
	pianoform_repo "erp/project/piano/pianoform/internal/repository"
)

type PianoFormUseCase interface {
	CreatePianoForm(ctx context.Context, i dto.CreatePianoBody) (err error)
	GetPianoForms(req *common.RequestContext, i *dto.RequestPaginationData) (
		dto.PaginationResult[[]*model.PianoForm], error)
	GetPianoForm(req *common.RequestContext, i *dto.RequestEntity) (
		dto.ResultEntity[*model.PianoForm], error)
	ExportData(req *common.RequestContext, i *dto.PianoExportRequest) (res *bytes.Buffer, err error)
}

type pianoFormUseCase struct {
	emiLog        logger.EmitLog
	pianoFormRepo pianoform_repo.PianoFormRepository
	core          repository.CoreService
	excelExporter exporter.ExcelExporter
}

func NewPianoUseCase(
	logger logger.Logger,
	pianoFormRepo pianoform_repo.PianoFormRepository,
	core repository.CoreService,
	helpers *helpers.Helpers,
) PianoFormUseCase {
	return &pianoFormUseCase{
		emiLog:        logger.EmitLog("pianoform-usecase"),
		pianoFormRepo: pianoFormRepo,
		core:          core,
		excelExporter: helpers.ExcelExporter,
	}
}

func (u *pianoFormUseCase) ExportData(req *common.RequestContext, i *dto.PianoExportRequest) (res *bytes.Buffer, err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("ExportData"))
		}
	}()
	exportData, err := u.pianoFormRepo.ExportData(req, i)
	data := make([][]interface{}, len(exportData))
	for idx, d := range exportData {
		//COmplete
		data[idx] = []interface{}{
			d.FirstName,
			d.LastName,
			d.Email,
			d.PhoneNumber,
			d.MovingDate,
			d.PianoType,
			d.PickupFlights,
			d.PickupStreet,
			d.PickupCity,
			d.PickupState,
			d.PickupZip,
			d.DropoffFlights,
			d.DropoffStreet,
			d.DropoffCity,
			d.DropoffState,
			d.DropoffZip,
			d.CompanyID,
			d.RentPiano,
			d.StairsDropoff,
			d.StairsPickup,
		}
	}

	headers := []interface{}{"First Name", "Last Name", "Email", "Phone Number",
		"Moving Date", "Piano Type", "Pickup Flights", "Pickup Street",
		"Pickup City", "Pickup State", "Pickup Zip", "Dropoff Flights",
		"Dropoff Street", "Dropoff City", "Dropoff State", "Dropoff Zip",
		"Company ID", "Rent Piano", "Stairs Dropoff", "Stairs Pickup"}
	res, err = u.excelExporter.Export("Sheet", headers, data)
	return
}

func (u *pianoFormUseCase) GetPianoForm(req *common.RequestContext, i *dto.RequestEntity) (
	res dto.ResultEntity[*model.PianoForm], err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("GetPianoForms"))
		}
	}()
	res, err = u.pianoFormRepo.GetPianoForm(req, i)
	if err != nil {
		return
	}
	res.Activities = u.core.GerActivitiesByPartyID(req, res.Entity.ID)
	return
}

func (u *pianoFormUseCase) GetPianoForms(req *common.RequestContext, i *dto.RequestPaginationData) (
	res dto.PaginationResult[[]*model.PianoForm], err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("GetPianoForms"))
		}
	}()
	res, err = u.pianoFormRepo.GetPianoForms(req, i)
	return
}

func (u *pianoFormUseCase) CreatePianoForm(ctx context.Context, i dto.CreatePianoBody) (err error) {
	defer func() {
		if err != nil {
			u.emiLog.Err(err, logger.OptionsLog.WithMethod("CreatePianoForm"))
		}
	}()
	err = u.pianoFormRepo.CreatePianoForm(ctx, i)
	return err
}
