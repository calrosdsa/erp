package db

import (
	"context"
	"crypto/rand"
	"erp/internal/domain"
	"erp/pkg/logger"
	"fmt"

	"gorm.io/gorm"
	
	"erp/gen/db/query"
	"erp/gen/db/model"

)

type dbHelper struct {
	emitLog logger.EmitLog
}

func NewDbHelper(
	logger logger.Logger,
) DbHelper {
	return &dbHelper{
		emitLog: logger.EmitLog("db-helper"),
	}
}

func (c  *dbHelper)DeleteReferences(ctx context.Context, tx *query.QueryTx, partyID int64) error{
	_, err := tx.PartyReference.WithContext(ctx).Unscoped().Where(
		tx.PartyReference.ReferenceID.Eq(partyID),
	).Delete()
	return err
}


func (c *dbHelper) InsertReferences(ctx context.Context, tx *query.QueryTx, partyID int64, references []*int64,
	args ...interface{}) (err error) {
	var partyReferences []*model.PartyReference
	var (
		isEdit bool
	)
	if len(args) == 1 {
		isEdit,_ =args[0].(bool)
	}
	ids := make([]int64,len(references))
	for _, reference := range references {
		if reference != nil {
			partyReference := &model.PartyReference{
				PartyID:     *reference,
				ReferenceID: partyID,
			}
			partyReferences = append(partyReferences, partyReference)
			ids = append(ids,  *reference)
		}
	}
	if isEdit {
		_,err = tx.PartyReference.WithContext(ctx).Unscoped().Where(
			tx.PartyReference.PartyID.In(ids...),
		).Delete()
		if err != nil {
			return
		}
	}
	err = tx.PartyReference.WithContext(ctx).CreateInBatches(partyReferences, len(partyReferences))
	return
}
func (c *dbHelper) GenerateCode(db *gorm.DB, entity interface{}, companyId int64) (code string) {
	for {
		code = c.generateCode()
		// Check if the code already exists in the database
		var count int64
		err := db.Model(entity).
			Where("company_id = ? AND code = ?", companyId, code).
			Count(&count).Error
		if err != nil {
			c.emitLog.Err(err, logger.OptionsLog.WithMethod("GenerateCode"))
			break
		}
		// If the code is unique, break the loop
		if count == 0 {
			break
		}
	}
	return
}

func (r *dbHelper) ValidateName(db *gorm.DB, entity interface{}) error {
	var count int64
	err := db.Model(entity).Where(entity).
			Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return domain.ERROR_NAME_TAKEN
	}
	return nil
}

func (c *dbHelper) generateCode() string {
	n := 4
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%X", b)
}
