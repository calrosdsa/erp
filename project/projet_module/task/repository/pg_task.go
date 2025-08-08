package task_repo

import (
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/gen/proto"
	"erp/internal/app/service/helpers"
	"erp/pkg/db"

	"strings"

	"gorm.io/gen/field"
	"gorm.io/gen/helper"
	"gorm.io/gen"
)

type TaskRepository interface {
	CreateTask(tx *query.QueryTx, req *common.RequestContext, d dto.TaskData) (res model.Task, err error)
	EditTask(tx *query.QueryTx, req *common.RequestContext, d dto.TaskData) (err error)
	GetTask(req *common.RequestContext, d dto.RequestEntity) (res dto.TaskDetailDto, err error)
	GetTasks(req *common.RequestContext, d dto.TasksRequest) (res []dto.TaskDto, err error)
	TaskTransition(tx *query.QueryTx, req *common.RequestContext, d dto.EntityTransitionData) (err error)
}

type taskRepo struct {
	query     helpers.QueryHelper
	convertor helpers.ConvertorHelper
	Q         *query.Query
}

func NewTaskRepository(
	db db.Connection,
	helpers *helpers.Helpers,
) TaskRepository {
	return &taskRepo{
		Q:         db.GetQ(),
		convertor: helpers.Convertor,
		query:     helpers.Query,
	}
}


func (r *taskRepo) TaskTransition(tx *query.QueryTx, req *common.RequestContext, d dto.EntityTransitionData) (err error) {
	if d.SourceStageID == d.DestinationStageID {
		if d.SourceIndex < d.DestinationIndex {
			if err = r.updateTaskIndexes(tx, int32(d.SourceStageID), d.SourceIndex, false); err != nil {
				return
			}
		} else {
			if err = r.updateTaskIndexes(tx, int32(d.DestinationStageID), d.DestinationIndex, true); err != nil {
				return
			}
		}
		_, err = tx.Task.Where(
			tx.Task.ID.Eq(d.ID),
		).UpdateSimple(tx.Task.Index.Value(d.DestinationIndex))
		if err != nil {
			return
		}

	} else {
		if err = r.updateTaskIndexes(tx, int32(d.SourceStageID), d.SourceIndex, false); err != nil {
			return
		}
		if err = r.updateTaskIndexes(tx, int32(d.DestinationStageID), d.DestinationIndex, true); err != nil {
			return
		}
		_, err = tx.Task.Where(
			tx.Task.ID.Eq(d.ID),
		).UpdateSimple(
			tx.Task.Index.Value(d.DestinationIndex),
			tx.Task.StageID.Value(int32(d.DestinationStageID)),
		)
		if err != nil {
			return
		}
	}
	

	return
}

func (r *taskRepo) updateTaskIndexes(tx *query.QueryTx, stageID int32, index int32, isAdd bool) (err error) {
	var expr field.AssignExpr
	var cond gen.Condition
	if isAdd {
		expr = tx.Task.Index.Add(1)
		cond = tx.Task.Index.Gte(index)
	} else {
		expr = tx.Task.Index.Sub(1)
		cond = tx.Task.Index.Gt(index)
	}
	_, err = tx.Task.Where(
		tx.Task.StageID.Eq(stageID),
		cond,
	).UpdateSimple(expr)
	return
}

func (r *taskRepo) CreateTask(tx *query.QueryTx, req *common.RequestContext, d dto.TaskData) (res model.Task, err error) {
	id, err := tx.Task.InsertParty(proto.PartyType_task.String())
	if err != nil {
		return
	}
	fields := d.Fields
	res.ID = id
	res.CompanyID = req.ActiveCompany.ID
	if err = r.convertor.CopyStructData(fields, &res); err != nil {
		return
	}

	err = r.updateTaskIndexes(tx, res.StageID, 0, true)
	if err != nil {
		return
	}

	err = tx.WithContext(req.Ctx).Task.Save(&res)
	if err != nil {
		return
	}
	if err = tx.Task.InsertActivity(id, req.Profile.ID, proto.ActivityType_CREATE.String(), nil); err != nil {
		return
	}

	return
}

func (r *taskRepo) EditTask(tx *query.QueryTx, req *common.RequestContext, d dto.TaskData) (err error) {
	data, err := r.convertor.DataMap(d.Fields)
	if err != nil {
		return
	}
	err = tx.Task.UnderlyingDB().WithContext(req.Ctx).Model(
		&model.Task{ID: d.ID},
	).Updates(data).Error
	if err != nil {
		return
	}
	err = tx.Task.InsertActivity(d.ID, req.Profile.ID, proto.ActivityType_EDIT.String(), nil)
	if err != nil {
		return
	}
	return
}

func (r *taskRepo) GetTask(req *common.RequestContext, d dto.RequestEntity) (res dto.TaskDetailDto, err error) {
	e := r.Q.Task
	projectQ := r.Q.Project
	assigneeQ := r.Q.Profile
	taskID := r.convertor.StrtoInt(d.ID)
	stageQ := r.Q.Stage
	err = e.WithContext(req.Ctx).Select(
		e.ID, e.UUID, e.Title, e.CreatedAt, e.Priority, e.DueDate,
		e.Description, e.ProjectID, projectQ.Name.As("project"),
		e.Assignee.As("assignee_id"), assigneeQ.GivenName.As("assignee_given_name"),
		assigneeQ.FamilyName.As("assignee_family_name"),
		e.StageID, stageQ.Name.As("stage"), stageQ.Index.As("stage_index"),
	).Join(projectQ, projectQ.ID.EqCol(e.ProjectID)).
		LeftJoin(assigneeQ, assigneeQ.ID.EqCol(e.Assignee)).
		LeftJoin(stageQ, stageQ.ID.EqCol(e.StageID)).
		Where(
			e.CompanyID.Eq(req.ActiveCompany.ID),
			e.ID.Eq(taskID),
		).Scan(&res.Task)
	if err != nil {
		return res, err
	}
	return
}

func (r *taskRepo) GetTasks(req *common.RequestContext, d dto.TasksRequest) (res []dto.TaskDto, err error) {
	var (
		generateSQL strings.Builder
	)
	builder := r.Q.WithContext(req.Ctx).Task
	queryData := r.convertor.GenerateQueryMap(d)
	params := r.tasksQuery(req, queryData, &generateSQL)
	err = builder.UnderlyingDB().Raw(generateSQL.String(), params...).Scan(&res).Error
	return
}

func (r *taskRepo) tasksQuery(req *common.RequestContext, d map[string]string, generateSQL *strings.Builder,
) (params []interface{}) {
	var (
		whereSQL strings.Builder
	)
	generateSQL.WriteString(`SELECT 
		e.id,e.uuid,e.created_at,e.title,e.priority,e.due_date,e.description,
		e.project_id,e.stage_id,s.name as stage,
		e.assignee as assignee_id,a.given_name as assignee_given_name,a.family_name as assignee_family_name
		from tasks as e 
		left join profiles as a on a.id = e.assignee
		left join stages as s on s.id = e.stage_id
	`)
	whereSQL.WriteString(` e.deleted_at is null and e.company_id = ? `)
	params = append(params, req.ActiveCompany.ID)
	columnFilters := []string{}
	r.query.FilterBuilder(&whereSQL, &params, d, columnFilters...)

	helper.JoinWhereBuilder(generateSQL, whereSQL)

	r.query.OrderAndLimitBuilder(generateSQL, d)
	return
}
