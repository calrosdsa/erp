package group_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type GroupRepository interface {
	GetGroupDescendents(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		res dto.ResultEntity[[]dto.GroupHierarchyDto], err error)
	GetGroup(req *common.RequestContext, d *dto.RequestEntityWithParty) (
		res dto.ResultEntity[dto.GroupDto], err error)
	GetGroups(req *common.RequestContext, d *dto.RequestPaginationPartyData) (
		res dto.PaginationResult[[]dto.GroupDto], err error)
	CreateGroup(req *common.RequestContext, i *dto.CreateGroupRequest) (err error)
	EditGroup(req *common.RequestContext, d *dto.EditGroupRequest) (err error)
	GetTreeView(req *common.RequestContext, d *dto.RequestDataWithPartyType) (
		res []dto.TreeEntryDto, err error)
}

type groupRepository struct {
	Q         *query.Query
	DB        *gorm.DB
	convertor helpers.ConvertorHelper
}

func NewGroupRepository(
	conn db.Connection,

	helpers *helpers.Helpers,
) GroupRepository {
	return &groupRepository{
		Q:         conn.GetQ(),
		DB:        conn.GetDB(),
		convertor: helpers.Convertor,
	}
}
func (r *groupRepository) GetTreeView(req *common.RequestContext, d *dto.RequestDataWithPartyType) (
	res []dto.TreeEntryDto, err error) {
	query := `WITH RECURSIVE data_cte AS (
		SELECT 
			id, 
			uuid,
			parent_id,
			is_group,
			name
		FROM groups
		WHERE parent_id IS NULL and company_id = ?
		UNION ALL 
		SELECT 
			l.id, 
			l.uuid,
			l.parent_id,
			l.is_group,
			l.name
		FROM groups l
		INNER JOIN data_cte d
			ON l.parent_id = d.id  
	)
	SELECT * FROM data_cte;`
	err = r.Q.Group.UnderlyingDB().WithContext(req.Ctx).Raw(query, req.ActiveCompany.ID).Scan(&res).Error
	return
}

func (r *groupRepository) EditGroup(req *common.RequestContext, d *dto.EditGroupRequest) (
	err error) {
	tx := r.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	groupQ := tx.Group
	var columns []field.AssignExpr

	columns = append(columns, groupQ.Name.Value(d.Body.Name)) // groupQ..Value(d.Body.),

	_, err = tx.Group.WithContext(req.Ctx).Where(
		groupQ.ID.Eq(d.Body.ID), groupQ.CompanyID.Eq(req.ActiveCompany.ID),
	).UpdateSimple(columns...)
	if err != nil {
		return
	}
	err = tx.Group.InsertActivity(d.Body.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (s *groupRepository) GetGroupDescendents(req *common.RequestContext, d *dto.RequestEntityWithParty) (
	res dto.ResultEntity[[]dto.GroupHierarchyDto], err error) {
	query := `
    WITH RECURSIVE groups_cte(id,parent_id,uuid,parent_uuid,name,is_group,enabled,depth) as (
	SELECT 
    	groups.id, 
    	groups.parent_id,
    	groups.uuid,
		(select uuid from groups where id = groups.parent_id),
    	groups.name,
    	groups.is_group, 
    	groups.enabled, 
    	0
        FROM groups
        WHERE uuid = ?
	UNION ALL 
	SELECT 
	    groups.id, 
    	groups.parent_id, 
		groups.uuid,
    	(select uuid from groups where id = groups_cte.id),
    	groups.name, 
    	groups.is_group, 
    	groups.enabled, 
    	depth +1
        FROM groups_cte,
             groups
        WHERE groups.parent_id = groups_cte.id
	)
	SELECT g.uuid,g.parent_uuid,g.name,g.is_group,g.enabled,g.depth
	FROM groups_cte as g;
    `

	if err := s.DB.Raw(query, d.ID).WithContext(req.Ctx).Scan(&res.Entity).Error; err != nil {
		return res, err
	}
	return res, err
}

func (s *groupRepository) GetGroup(req *common.RequestContext, d *dto.RequestEntityWithParty) (
	res dto.ResultEntity[dto.GroupDto], err error) {
	id := s.convertor.StrtoInt(d.ID)
	groupQ := s.Q.Group
	group, err := groupQ.WithContext(req.Ctx).Where(
		groupQ.CompanyID.Eq(req.ActiveCompany.ID),
		groupQ.ID.Eq(id),
	).First()
	if err != nil {
		return res, err
	}
	res.Entity = dto.GroupDtoFromModel(group)
	return res, err
}

func (r *groupRepository) GetGroups(req *common.RequestContext, i *dto.RequestPaginationPartyData) (
	res dto.PaginationResult[[]dto.GroupDto], err error) {
	var (
		conds []gen.Condition
		order field.Expr
	)
	groupQ := r.Q.Group
	builder := r.Q.WithContext(req.Ctx).Group

	//ADDING CONDITIONS
	conds = append(conds, groupQ.CompanyID.Eq(req.ActiveCompany.ID))
	if i.IsGroup != "" {
		conds = append(conds, groupQ.IsGroup.Is(r.convertor.StrToBool(i.IsGroup)))
	}

	limit, offset := r.convertor.ToPaginationParams(i.Page, i.Size)
	orderCol, ok := r.Q.Payment.GetFieldByName(i.OrderColumn) // maybe orderColStr == "id"
	if ok {
		if i.Order == "ASC" {
			order = orderCol.Asc()
		} else {
			order = orderCol.Desc()
		}
		builder = builder.Order(order)
		// User doesn't contains orderColStr
	}
	partyQ := r.Q.Party
	builder = builder.Select(
		groupQ.ID, groupQ.UUID, groupQ.Name, groupQ.CreatedAt,
		groupQ.IsGroup, groupQ.Ordinal,
	).
		Join(partyQ, partyQ.ID.EqCol(groupQ.ID), partyQ.PartyTypeCode.Eq(i.PartyType)).
		Where(conds...)
	total, err := builder.Count()
	if err != nil {
		return res, err
	}
	err = builder.Limit(limit).Offset(offset).Scan(&res.Results)
	res.Total = total
	return
}

func (s *groupRepository) CreateGroup(req *common.RequestContext, i *dto.CreateGroupRequest) (err error) {
	tx := s.Q.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	var group model.Group
	// supplierG := s.Q.Group
	partyID, err := tx.Group.InsertParty(i.Body.PartyTypeCode)
	if err != nil {
		return err
	}
	group.ID = partyID
	group.CompanyID = int64(req.ActiveCompany.ID)
	group.Name = i.Body.Name
	// group.IsGroup = i.Body.IsGroup
	// group.Enabled = i.Body.Enabled
	group.Ordinal = 0

	// if i.Body.ParentID != nil && *i.Body.ParentID != 0 {
	// 	parentGroup, err := tx.Group.WithContext(req.Ctx).
	// 		Select(supplierG.ID).
	// 		Where(
	// 			supplierG.ID.Eq(*i.Body.ParentID),
	// 			supplierG.CompanyID.Eq(int64(req.ActiveCompany.ID)),
	// 		).First()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	group.ParentID = &parentGroup.ID
	// 	group.Ordinal = parentGroup.Ordinal + 1
	// }

	err = tx.Group.Save(&group)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// get entity that the group is related
// func (h *groupRepository) GetEntityGroup(partyCode string) *domain.EntityTemplate {
// 	switch partyCode {
// 	case domain.PARTY_SUPPLIER_GROUP:
// 		return &domain.SUPPLIER
// 	default:
// 		return nil
// 	}
// }
