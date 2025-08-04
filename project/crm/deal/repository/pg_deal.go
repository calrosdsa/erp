package deal_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/internal/domain"
	"erp/pkg/db"

	"strings"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gen/helper"
)

type DealRepository interface {
	CreateDeal(tx *query.QueryTx, req *common.RequestContext, d dto.DealData) (res model.Deal, err error)
	EditDeal(tx *query.QueryTx, req *common.RequestContext, d dto.DealData) (err error)
	GetDeal(req *common.RequestContext, d dto.RequestEntity) (res dto.DealDetailDto, err error)
	GetDeals(req *common.RequestContext, d dto.DealsRequest) (res []dto.DealDto, err error)
	DealTransition(tx *query.QueryTx, req *common.RequestContext, d dto.EntityTransitionData) (err error)
	
}

type dealRepo struct {
	query     helpers.QueryHelper
	convertor helpers.ConvertorHelper
	Q         *query.Query
}

func NewDealRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) DealRepository {
	return &dealRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		query:     helpers.Query,
	}
}

func (r *dealRepo) DealTransition(tx *query.QueryTx, req *common.RequestContext, d dto.EntityTransitionData) (err error) {
	if d.SourceStageID == d.DestinationStageID {
		if d.SourceIndex < d.DestinationIndex {
			if err = r.updateDealIndexes(tx, d.SourceStageID, d.SourceIndex, false); err != nil {
				return
			}
		} else {
			if err = r.updateDealIndexes(tx, d.DestinationStageID, d.DestinationIndex, true); err != nil {
				return
			}
		}
		_, err = tx.Deal.Where(
			tx.Deal.ID.Eq(d.ID),
		).UpdateSimple(tx.Deal.Index.Value(d.DestinationIndex))
		if err != nil {
			return
		}

	} else {
		if err = r.updateDealIndexes(tx, d.SourceStageID, d.SourceIndex, false); err != nil {
			return
		}
		if err = r.updateDealIndexes(tx, d.DestinationStageID, d.DestinationIndex, true); err != nil {
			return
		}
		_, err = tx.Deal.Where(
			tx.Deal.ID.Eq(d.ID),
		).UpdateSimple(
			tx.Deal.Index.Value(d.DestinationIndex),
			tx.Deal.StageID.Value(d.DestinationStageID),
		)
		if err != nil {
			return
		}
	}
	

	return
}

func (r *dealRepo) updateDealIndexes(tx *query.QueryTx, stageID int64, index int32, isAdd bool) (err error) {
	var expr field.AssignExpr
	var cond gen.Condition
	if isAdd {
		expr = tx.Deal.Index.Add(1)
		cond = tx.Deal.Index.Gte(index)
	} else {
		expr = tx.Deal.Index.Sub(1)
		cond = tx.Deal.Index.Gt(index)
	}
	// query := `update from deals set updated_at = now(), index = index`
	_, err = tx.Deal.Where(
		tx.Deal.StageID.Eq(stageID),
		cond,
	).UpdateSimple(expr)
	return
}

func (r *dealRepo) CreateDeal(tx *query.QueryTx, req *common.RequestContext, d dto.DealData) (res model.Deal, err error) {
	id, err := tx.Deal.InsertParty(proto.PartyType_deal.String())
	if err != nil {
		return
	}
	fields := d.Fields
	res.ID = id
	res.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}
	// Update indexes before saving the deal to prevent modifying indexes for newly created entries.
	err = r.updateDealIndexes(tx, res.StageID, 0, true)
	if err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Deal.Save(&res)
	if err != nil {
		return
	}
	if err = tx.Deal.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}
	err = r.upsertParticipants(tx, d.Participants, res.ID)

	return
}
func (r *dealRepo) EditDeal(tx *query.QueryTx, req *common.RequestContext, d dto.DealData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Deal.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Deal{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Deal.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = r.upsertParticipants(tx, d.Participants, d.ID)
	return
}

func (r *dealRepo) upsertParticipants(tx *query.QueryTx, participants []dto.ParticipantData, dealID int64) (err error) {
	for _, participant := range participants {
		switch participant.Action {
		case string(domain.CREATE):
			dealParticipant := model.DealParticipant{
				DealID:    dealID,
				ProfileID: participant.ID,
			}
			err = tx.DealParticipant.Save(&dealParticipant)
		case string(domain.DELETE):
			_, err = tx.DealParticipant.Where(
				tx.DealParticipant.DealID.Eq(dealID),
				tx.DealParticipant.ProfileID.Eq(participant.ID),
			).Delete()
		}
	}
	return
}
func (r *dealRepo) GetDeal(req *common.RequestContext, d dto.RequestEntity) (res dto.DealDetailDto, err error) {
	e := r.Q.Deal
	stageQ := r.Q.Stage
	responsibleQ := r.Q.Profile
	dealID := r.convertor.StrtoInt(d.ID)
	customerQ := r.Q.Customer
	err = e.WithContext(req.Ctx).Select(
		e.ID, e.UUID, e.Name, e.CreatedAt, e.Amount, e.DealType, e.Currency, e.Source,
		e.SourceInformation, e.StartDate,
		e.StageID, stageQ.Name.As("stage"), stageQ.Index.As("stage_index"),
		e.ResponsibleID, responsibleQ.GivenName.As("responsible_given_name"),
		responsibleQ.FamilyName.As("responsible_family_name"), responsibleQ.UUID.As("responsible_uuid"),
		e.CustomerID,customerQ.Name.As("customer"),
	).
		Join(stageQ, stageQ.ID.EqCol(e.StageID)).
		Join(responsibleQ, responsibleQ.ID.EqCol(e.ResponsibleID)).
		LeftJoin(customerQ,customerQ.ID.EqCol(e.CustomerID)).
		Where(
			e.CompanyID.Eq(req.ActiveCompany.ID),
			e.ID.Eq(dealID),
		).Scan(&res.Deal)
	if err != nil {
		return res, err
	}
	res.Participants,err = r.getParticipants(req,res.Deal.ID)
	return
}

func (r *dealRepo) getParticipants(req *common.RequestContext, dealID int64) (res []dto.ProfileDto, err error) {
	err = r.Q.DealParticipant.UnderlyingDB().Raw(`
		select 
		p.id,p.uuid,p.given_name,p.family_name,
		CONCAT(p.given_name, ' ', p.family_name) as full_name,
		p.email_address,p.phone_number
		from deal_participants as e
		inner join profiles as p on p.id = e.profile_id
		where e.deal_id = ? and e.deleted_at is null
	`,dealID).Scan(&res).Error
	return 
}

func (r *dealRepo) GetDeals(req *common.RequestContext, d dto.DealsRequest) (res []dto.DealDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Deal
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.dealsQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *dealRepo) dealsQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
			e.id,e.uuid,e.created_at,e.name,e.amount,e.currency,
			e.deal_type,e.source,e.source_information,e.start_date,
			e.stage_id,
			e.responsible_id,r.given_name as responsible_given_name,r.family_name as responsible_family_name,
			r.uuid as responsible_uuid,
			e.customer_id,c.name as customer
			from deals as e 
			inner join profiles as r on r.id = e.responsible_id
			left join customers as c on c.id = e.customer_id
		`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}
