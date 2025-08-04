package journal_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	"fmt"

	"gorm.io/gen"
	"gorm.io/gen/field"
)

type JournalRepository interface {
	CreateJournalEntry(req *common.RequestContext, d dto.JournalEntryData) (res dto.JournalEntryDto, err error)
	EditJournalEntry(req *common.RequestContext, d dto.JournalEntryData) (err error)
	GetJournalEntry(req *common.RequestContext, d *dto.RequestEntity) (
		res dto.ResultEntity[dto.JournalEntryDetailDto], err error)
	GetJournalEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
		res dto.PaginationResult[[]dto.JournalEntryDto], err error)
	UpdateStatus(req *common.RequestContext, tx *query.QueryTx,
		d dto.UpdateStatusWithEvent,nextState string) (res model.JournalEntry, lines []*model.JournalEntryLine, err error)
}

type journalRepository struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
	dbHelper  db.DbHelper
	currency  helpers.CurrencyHelper
}

func NewJournalRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) JournalRepository {
	return &journalRepository{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		dbHelper:  db.GetDbHelper(),
		currency:  helpers.Currency,
	}
}

func (r *journalRepository) UpdateStatus(req *common.RequestContext, tx *query.QueryTx,
	d dto.UpdateStatusWithEvent,nextState string) (res model.JournalEntry, line []*model.JournalEntryLine, err error) {
	journalEntryQ := tx.JournalEntry
	id := r.convertor.StrtoInt(d.Body.PartyID)
	fmt.Println("JOURNAL ENTRY ID", id)
	_, err = tx.JournalEntry.WithContext(req.Ctx).Where(
		journalEntryQ.CompanyID.Eq(req.ActiveCompany.ID),
		journalEntryQ.Status.Eq(d.Body.CurrentState),
		journalEntryQ.ID.Eq(id),
	).UpdateSimple(journalEntryQ.Status.Value(nextState))
	if err != nil {
		return
	}
	journalEntry, err := tx.JournalEntry.WithContext(req.Ctx).Where(
		journalEntryQ.CompanyID.Eq(req.ActiveCompany.ID),
		journalEntryQ.ID.Eq(id),
	).First()
	if err != nil {
		return
	}
	lineItems, err := tx.JournalEntryLine.WithContext(req.Ctx).Where(
		tx.JournalEntryLine.JournalEntryID.Eq(journalEntry.ID),
	).Find()
	if err != nil {
		return
	}

	return *journalEntry, lineItems, err
}

func(r *journalRepository) EditJournalEntry(req *common.RequestContext,d dto.JournalEntryData)(err error){
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.JournalEntry.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.JournalEntry{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}

	err = r.processJournalEntryLines(tx, req, d, d.ID)
	if err != nil {
		return
	}

	err = tx.JournalEntry.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err =tx.Commit()
	return
}

func (r *journalRepository) CreateJournalEntry(req *common.RequestContext, d dto.JournalEntryData) (
	res dto.JournalEntryDto, err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var journalEntry model.JournalEntry
	journalEntry.Code = r.dbHelper.GenerateCode(r.Q.JournalEntry.UnderlyingDB(), model.JournalEntry{}, req.ActiveCompany.ID)
	id, err := tx.Address.InsertParty(proto.PartyType_workspace.String())
	if err != nil {
		return
	}
	fields := d.Fields
	journalEntry.ID = id
	journalEntry.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &journalEntry); err != nil {
		return
	}
	
	err = tx.WithContext(req.Ctx).JournalEntry.Save(&journalEntry)
	if err != nil {
		return
	}

	err = r.processJournalEntryLines(tx, req, d, id)
	if err != nil {
		return
	}
	
	if err = tx.Address.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	res = dto.JournalEntryDtoFromModel(&journalEntry)
	err = tx.Commit()
	
	return
}

func (r *journalRepository) processJournalEntryLines(tx *query.QueryTx, req *common.RequestContext,
	d dto.JournalEntryData, journalEntryID int64) (err error) {
	e := tx.JournalEntryLine
	_,err = e.WithContext(req.Ctx).Unscoped().Where(
		e.JournalEntryID.Eq(journalEntryID),
	).Delete()
	if err != nil {
		return
	}
	entryLines := make([]*model.JournalEntryLine, len(d.EntryLines))
	for i, line := range d.EntryLines {
		entryLine := &model.JournalEntryLine{}
		entryLine.Credit = r.currency.FloatToInt(line.Credit)
		entryLine.Debit = r.currency.FloatToInt(line.Debit)
		entryLine.Currency = req.CompanyDefaults.Currency
		entryLine.JournalEntryID = journalEntryID
		entryLine.LedgerID = line.LedgerID
		entryLine.ProjectID = line.ProjectID
		entryLine.CostCenterID = line.CostCenterID
		entryLines[i] = entryLine
	}
	err = tx.JournalEntryLine.WithContext(req.Ctx).CreateInBatches(entryLines, len(entryLines))
	return
}

func (r *journalRepository) GetJournalEntry(req *common.RequestContext, d *dto.RequestEntity) (
	res dto.ResultEntity[dto.JournalEntryDetailDto], err error) {
	journalEntryQ := r.Q.JournalEntry
	err = journalEntryQ.WithContext(req.Ctx).Select(
		journalEntryQ.ID, journalEntryQ.Code, journalEntryQ.Status,
		journalEntryQ.PostingDate, journalEntryQ.EntryType,
	).
		Where(journalEntryQ.CompanyID.Eq(req.ActiveCompany.ID),
			journalEntryQ.Code.Eq(d.ID)).
		Scan(&res.Entity.JournalEntry)
	if err != nil {
		return
	}
	res.Entity.JournalEntryLines, err = r.getJournalEntryLines(req, res.Entity.JournalEntry.ID)

	return
}

func (r *journalRepository) getJournalEntryLines(req *common.RequestContext, journalEntryID int64) (
	res []dto.JournalEntryLineDto, err error) {
	lineQ := r.Q.JournalEntryLine
	ledgerQ := r.Q.Ledger
	projectQ := r.Q.Project
	costCenterQ := r.Q.CostCenter
	err = r.Q.JournalEntryLine.WithContext(req.Ctx).Select(
		lineQ.Debit, lineQ.Credit, lineQ.Currency,
		ledgerQ.Name.As("account"), ledgerQ.ID.As("account_id"),
		costCenterQ.Name.As("cost_center"), costCenterQ.ID.As("cost_center_id"),
		projectQ.Name.As("project"), projectQ.ID.As("project_id"),
	).Where(
		lineQ.JournalEntryID.Eq(journalEntryID),
	).
		Join(ledgerQ, ledgerQ.ID.EqCol(lineQ.LedgerID)).
		LeftJoin(projectQ, projectQ.ID.EqCol(lineQ.ProjectID)).
		LeftJoin(costCenterQ, costCenterQ.ID.EqCol(lineQ.CostCenterID)).
		Scan(&res)
	return
}

func (r *journalRepository) GetJournalEntries(req *common.RequestContext, d *dto.RequestPaginationData) (
	res dto.PaginationResult[[]dto.JournalEntryDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	journalEntryQ := r.Q.JournalEntry
	builder := r.Q.WithContext(req.Ctx).JournalEntry

	limit, offset := r.convertor.ToPaginationParams(d.Page, d.Size)
	orderCol, ok := r.Q.JournalEntry.GetFieldByName(d.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if d.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
	}
	//ADDING CONDITIONS
	conds = append(conds, journalEntryQ.CompanyID.Eq(req.ActiveCompany.ID))

	builder = builder.Select(
		journalEntryQ.ID, journalEntryQ.Code, journalEntryQ.Status,
		journalEntryQ.PostingDate, journalEntryQ.EntryType,
	).Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}
