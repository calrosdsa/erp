package helpers

import (
	"context"
	"erp/internal/app/config"
	"erp/pkg/logger"
	"fmt"

	permify_payload "buf.build/gen/go/permifyco/permify/protocolbuffers/go/base/v1"
	permify_grpc "github.com/Permify/permify-go/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PermifyHelper struct {
	configService *config.ConfigService
	client        *permify_grpc.Client
	emitLog       EmitLog
}

func NewPermifyHelper(
	configService *config.ConfigService,
	logger *LoggerHelper,
) *PermifyHelper {
	permifyConfig := configService.GetPermifyConfig()
	client, err := permify_grpc.NewClient(
		permify_grpc.Config{
			Endpoint: permifyConfig.Endpoint,
		},
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		// log.Panic(err)
	}
	helper := &PermifyHelper{
		client:  client,
		emitLog: logger.EmitLog("permify"),
	}
	// helper.CreateNewTenant(&permify_payload.TenantCreateRequest{
	// 	Id: "EA96361B",
	// 	Name: "My organization",
	// })
// 	helper.WriteSchema(&permify_payload.SchemaWriteRequest{
// 		TenantId: "t1",
// 		Schema: `
// 		entity user {}

// entity role {
//     relation assignee @user
// }

// entity template {
//     relation view @role#assignee
//     relation delete @role#assignee
//     relation edit @role#assignee
//     relation create @role#assignee
// }
// 		`,
// 	})

	// helper.WriteRelationships(&permify_payload.RelationshipWriteRequest{
	// 	TenantId: "t2",
	// 	Metadata: &permify_payload.RelationshipWriteRequestMetadata {
	// 		SchemaVersion: "cr9qmbkldtpvng6hk0vg",
	// 	},
	// 	Tuples: [] * permify_payload.Tuple {
	// 		{
	// 			Entity: & permify_payload.Entity {
	// 				Type: "role",
	// 				Id: "admin",
	// 			},
	// 			Relation: "assignee",
	// 			Subject: & permify_payload.Subject {
	// 				Type: "user",
	// 				Id: "1",
	// 			},
	// 		},
	// 	},
	// })

	// helper.WriteRelationships(&permify_payload.RelationshipWriteRequest{
	// 	TenantId: "t2",
	// 	Metadata: &permify_payload.RelationshipWriteRequestMetadata {
	// 		SchemaVersion: "cr9qmbkldtpvng6hk0vg",
	// 	},
	// 	Tuples: [] * permify_payload.Tuple {
	// 		{
	// 			Entity: & permify_payload.Entity {
	// 				Type: "company",
	// 				Id: "main",
	// 			},
	// 			Relation: "edit",
	// 			Subject: & permify_payload.Subject {
	// 				Type: "role",
	// 				Id: "admin",
	// 				Relation: "assignee",
	// 			},
	// 		},
	// 	},
	// })

	// helper.Check(&permify_payload.PermissionCheckRequest{
	// 	TenantId: "t2",
	// 	Metadata: &permify_payload.PermissionCheckRequestMetadata{
	// 		// SnapToken: "twoAAAAAAAA=", // rr --> relationship write response
	// 		SchemaVersion: "cr9r4akldtpvng6hk10g", // sr --> schema write response
	// 		Depth:         10,
	// 	},
	// 	Entity: &permify_payload.Entity{
	// 		Type: "company",
	// 		Id:   "main",
	// 	},
	// 	Permission: "view",
	// 	Subject: &permify_payload.Subject{
	// 		Type:     "role",
	// 		Id:       "member",
	// 		Relation: "assignee",
	// 	},
	// })

	return helper
}
func (h *PermifyHelper) WriteRelationships(ctx context.Context,payload *permify_payload.RelationshipWriteRequest) error {
	_, err := h.client.Data.WriteRelationships(ctx, payload)
	if err != nil {
		h.emitLog.Err(err, logger.OptionsLog.WithMethod("WriteRelationships"))
	}
	//USE rr.SnapToken
	// fmt.Println(rr)
	return err
}

func (h *PermifyHelper)DeleteRelationships(ctx context.Context,payload *permify_payload.RelationshipDeleteRequest) error {
	_,err := h.client.Data.DeleteRelationships(ctx,payload)
	if err != nil {
		h.emitLog.Err(err, logger.OptionsLog.WithMethod("WriteRelationships"))
	}
	return err
}

func (h *PermifyHelper) Check(ctx context.Context,payload *permify_payload.PermissionCheckRequest) (bool) {
	// return true	
	cr, err := h.client.Permission.Check(ctx, payload)
	if cr == nil {
		fmt.Println("NIL DATA",payload)
		return false
	}
	if err != nil {
		return false
	}
	if cr.Can == permify_payload.CheckResult_CHECK_RESULT_ALLOWED {
		return true
	} else {
		return false
	}
}

func (h *PermifyHelper) CreateNewTenant(payload *permify_payload.TenantCreateRequest) (
	*permify_payload.TenantCreateResponse, error) {
	ct, err := h.client.Tenancy.Create(context.Background(), payload)
	if err != nil {
		h.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateNewTenant"))
	}
	return ct, err
}

func (h *PermifyHelper) WriteSchema(payload *permify_payload.SchemaWriteRequest) (
	*permify_payload.SchemaWriteResponse, error) {
	sr, err := h.client.Schema.Write(context.Background(), payload)
	if err != nil {
		h.emitLog.Err(err, logger.OptionsLog.WithMethod("WriteSchema"))
	}
	return sr, err
}

