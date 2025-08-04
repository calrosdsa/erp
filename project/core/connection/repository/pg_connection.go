package connection_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"
	// "fmt"

	"github.com/samber/lo"
)

type ConnectionRepository interface {
	GetConnectionsEntity(req *common.RequestContext, d dto.RequestEntity) (res []dto.ConnectionDto, err error)
}

type connectionRepo struct {
	Q         *query.Query
	convertor helpers.ConvertorHelper
}

func NewConnectionRepo(
	db db.Connection,
	helpers *helpers.Helpers,
) ConnectionRepository {
	return &connectionRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
	}
}

func (r *connectionRepo) GetConnectionsEntity(req *common.RequestContext, d dto.RequestEntity) (
	connections []dto.ConnectionDto, err error) {
	//Convert entityId from string to integer
	entityID := r.convertor.StrtoInt(d.ID)
	e := r.Q.Connection
	refEntQ := r.Q.Entity
	// var connections []dto.ConnectionDto
	err = e.WithContext(req.Ctx).Select(
		e.ID, e.SectionName, e.EntityID,
		refEntQ.Name.As("entity_name"),
		refEntQ.Href.As("entity_href"),
		refEntQ.HasModal.As("entity_has_modal"),
	).
		LeftJoin(refEntQ, refEntQ.ID.EqCol(e.ReferenceEntityID)).
		Where(
			e.EntityID.Eq(entityID),
		).Scan(&connections)
	if err != nil {
		return
	}
	// fmt.Println("CONNECTIONS",connections,"ENTITY ID",entityID)
	//Map id connections records to
	connectionEntityIds := lo.Map(connections, func(x dto.ConnectionDto, index int) int64 {
		return int64(x.ID)
	})

	partyRefQ := r.Q.PartyReference

	partyRefernces, err := partyRefQ.Select(
		partyRefQ.ReferenceID,
	).Where(
		partyRefQ.ReferenceID.In(connectionEntityIds...),
	).Find()
	if err != nil {
		return
	}
	// Group party references by their connection ID for quick lookup
	groupReferences := lo.GroupBy(partyRefernces, func(x *model.PartyReference) int64 {
		return x.ReferenceID
	})
	for i, reference := range connections {
		if val, exist := groupReferences[int64(reference.ID)]; exist {
			connections[i].Count = len(val)
		}
	}
	return
}

    // func(r *connectionRepo) 
