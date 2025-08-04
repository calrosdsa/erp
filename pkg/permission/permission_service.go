package permission

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/domain"
	"erp/pkg/db"
	"erp/pkg/logger"
	"fmt"
	"strconv"
	permify_payload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	"gorm.io/gorm"
)

var schemaVersion = "crbnh8d7s45s73bo60v0"
var tenatID = "1"

type permissionService struct {
	conn *gorm.DB
	permifyClient *PermifyClient
	emitLog       logger.EmitLog
}

func NewPermissionService(
	db db.Connection,
	logger logger.Logger,
)PermissionService{
	permifyClient := NewPermifyClient(logger)
	return &permissionService{
		conn: db.GetDB(),
		emitLog: logger.EmitLog("permify-service"),
		permifyClient: permifyClient,
	}
}


func (s *permissionService) CheckIfUserIsCompanyAdmin(ctx context.Context, req *common.RequestContext,
	companyID int64, permission domain.Permission) bool {
	return s.CheckCustomPermission(ctx, req,
		domain.PermissionOpts.WithEntityType(domain.COMPANY_ENTITY),
		domain.PermissionOpts.WithEntityID(strconv.Itoa(int(companyID))),
		domain.PermissionOpts.WithPermission(domain.CREATE_USER),
		domain.PermissionOpts.WithSubjectType(domain.USER_ENTITY),
		domain.PermissionOpts.WithSubjectId(strconv.Itoa(int(req.User.ID))),
	)
}

func (s *permissionService) CheckCustomPermission(ctx context.Context, req *common.RequestContext, permissionOpts ...domain.PermissionOpt) bool {
	opts := domain.PermissionOpts.Apply(permissionOpts...)
	fmt.Println("OPTS", opts)
	parentCompany, err := s.GetParentCompany(req, s.conn)
	if err != nil {
		return false
	}
	payload := &permify_payload.PermissionCheckRequest{
		TenantId: strconv.Itoa(int(parentCompany.ID)),
		Metadata: &permify_payload.PermissionCheckRequestMetadata{
			// SnapToken: "twoAAAAAAAA=", // rr --> relationship write response
			SchemaVersion: schemaVersion, // sr --> schema write response
			Depth:         10,
		},
		Entity: &permify_payload.Entity{
			Type: string(opts.EntityType),
			Id:   opts.EntityID,
		},
		Permission: string(opts.Permission),
		Subject: &permify_payload.Subject{
			Type: string(opts.SubjectType),
			Id:   opts.SubjectID,
		},
	}
	if opts.SubjectRelation != "" {
		payload.Subject.Relation = opts.SubjectRelation
	}
	return s.permifyClient.Check(ctx, payload)
}

func (s *permissionService) CheckPermission(ctx context.Context, req *common.RequestContext,
	entity domain.EntityTemplate, actionName domain.ActionType) bool {
	parentCompany, err := s.GetParentCompany(req, s.conn)
	if err != nil {
		return false
	}
	return s.permifyClient.Check(ctx, &permify_payload.PermissionCheckRequest{
		TenantId: strconv.Itoa(int(parentCompany.ID)),
		Metadata: &permify_payload.PermissionCheckRequestMetadata{
			// SnapToken: "twoAAAAAAAA=", // rr --> relationship write response
			SchemaVersion: schemaVersion, // sr --> schema write response
			Depth:         10,
		},
		Entity: &permify_payload.Entity{
			Type: "template",
			Id:   entity.Name,
		},
		Permission: string(actionName),
		Subject: &permify_payload.Subject{
			Type:     "role",
			Id:       strconv.Itoa(int(req.Role.ID)),
			Relation: "assignee",
		},
	})
}

func (s *permissionService) WriteTemplateAuthData(req *common.RequestContext, roleID int64,
	entityName string, actionName string) error {

	parentCompany, err := s.GetParentCompany(req, s.conn)
	if err != nil {
		return err
	}
	// fmt.Println("PARENT COMPANY CODE",parentCompany.ID)
	s.permifyClient.WriteRelationships(req.Ctx, &permify_payload.RelationshipWriteRequest{
		TenantId: strconv.Itoa(int(parentCompany.ID)),
		Metadata: &permify_payload.RelationshipWriteRequestMetadata{
			SchemaVersion: schemaVersion,
		},
		Tuples: []*permify_payload.Tuple{
			{
				Entity: &permify_payload.Entity{
					Type: "template",
					Id:   entityName,
				},
				Relation: actionName,
				Subject: &permify_payload.Subject{
					Type:     "role",
					Id:       strconv.Itoa(int(roleID)),
					Relation: "assignee",
				},
			},
		},
	})
	return err
}

func (s *permissionService) DeleteTemplateAuthData(req *common.RequestContext, roleID int64,
	entityName string, actionName string) error {
	parentCompany, err := s.GetParentCompany(req, s.conn)
	if err != nil {
		return err
	}
	s.permifyClient.DeleteRelationships(req.Ctx, &permify_payload.RelationshipDeleteRequest{
		TenantId: strconv.Itoa(int(parentCompany.ID)),
		Filter: &permify_payload.TupleFilter{
			Entity: &permify_payload.EntityFilter{
				Type: "template",
				Ids:  []string{entityName},
			},
			Relation: actionName,
			Subject: &permify_payload.SubjectFilter{
				Type:     "role",
				Ids:      []string{strconv.Itoa(int(roleID))},
				Relation: "assignee",
			},
		},
	})
	return err
}

func (s *permissionService) GetParentCompany(req *common.RequestContext, tx *gorm.DB) (model.Company, error) {
	var (
		res model.Company
		err error
	)
	defer func(){
		if err != nil {
			s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetParentCompany"))
		}
	}()
	if req.ActiveCompany.ParentID == nil {
		return req.ActiveCompany,nil
	}
	query := `WITH RECURSIVE companies_cte(id,parent_id) as (
			SELECT 
				companies.id, 
				companies.parent_id
				FROM companies
				WHERE companies.id = ?
			UNION ALL 
			SELECT 
				companies.id, 
				companies.parent_id
				FROM companies_cte,
					companies
				WHERE companies.id = companies_cte.parent_id
		)
		SELECT *
		FROM companies_cte where parent_id is NULL;`
	if err := s.conn.Raw(query,req.ActiveCompany.ID).WithContext(req.Ctx).Scan(&res).Error; err != nil {
		return res, err
	}
	return res, nil
}


func (s *permissionService) GetActions(ctx context.Context,entityID int64)[]dto.ActionDto{
	var actions []model.Action
	err := s.conn.WithContext(ctx).Where(&model.Action{EntityID: entityID}).Find(&actions).Error
	if err != nil {
		fmt.Println("FAIL TO GET ACTIONS")
	}
	actionDtos := make([]dto.ActionDto, len(actions))
	for i, action := range actions {
		actionDtos[i] = dto.ActionDtoFromModel(&action)
	}

	return actionDtos
}