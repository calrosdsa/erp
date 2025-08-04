package db

import (
	"context"
	"erp/gen/db/query"

	"gorm.io/gorm"
)

type Connection interface {
	GetDB() *gorm.DB
	GetQ() *query.Query
	GetDbHelper() DbHelper

	// GetPartyByUUID(req *common.RequestContext,partyType string,uuid string) *int64
	GetPartyByUUID(ctx context.Context, companyID int64, partyType string, uuid string) *int64

	GenerateCode(ctx context.Context, entity interface{}, companyId int64) (code string)
	ValidateCode(ctx context.Context, entity interface{}, companyId int64, code string) (err error)

	InsertReference(ctx context.Context, tx *query.QueryTx, partyID int64, referenceId int64) error
}

type DbHelper interface {
	InsertReferences(ctx context.Context, tx *query.QueryTx, partyID int64, 
	references []*int64,args ...interface{}) error
	DeleteReferences(ctx context.Context, tx *query.QueryTx, partyID int64) error
	ValidateName(db *gorm.DB, entity interface{}) error
	GenerateCode(db *gorm.DB, entity interface{}, companyId int64) (code string)}
