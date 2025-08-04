package uomservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/connection"
	"erp/internal/app/domain"
	"fmt"
	"time"

	"gorm.io/gen/field"
)

type UOMService struct {
	conn *connection.Connection
	timeout time.Duration	
}

func NewUOMService(conn *connection.Connection,timeout time.Duration) *UOMService{
	return &UOMService{
		conn: conn,
		timeout: timeout,
	}
}

func (s *UOMService) GetUnitOfMeasures(req *common.RequestContext,d *dto.UOMsRequest)([]dto.UOMDto,error){
	ctx,cancel := context.WithTimeout(req.Ctx,s.timeout)
	defer cancel()
	var (
		res []dto.UOMDto
	)
	fmt.Println("LANGUAGE CODE",req.LanguageCode)
	
	uomTranslation := s.conn.Q.UnitOfMeasureTranslation
	uom := s.conn.Q.UnitOfMeasure

	uomTranslationCond := []field.Expr{}
	uomTranslationCond = append(uomTranslationCond, uomTranslation.BaseID.EqCol(uom.ID))
	if d.Query != "" {
		uomTranslationCond = append(uomTranslationCond, uomTranslation.Name.Eq(d.Query))
	}
	uomTranslationCond = append(uomTranslationCond, uomTranslation.LanguageCode.Eq(string(req.LanguageCode)))

	err := s.conn.Q.UnitOfMeasure.WithContext(ctx).Select(uom.ID,uomTranslation.Name,uom.Code).
	Join(uomTranslation,uomTranslationCond...).Limit(domain.DEFAULT_LIMIT).Scan(&res)
	return res,err
}