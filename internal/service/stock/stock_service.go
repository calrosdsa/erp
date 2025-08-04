package stock_service

import (
	"erp/api/common"
	"erp/gen/db/model"
	"erp/gen/db/query"
	"erp/internal/domain/repository"
	"erp/pkg/db"
	"erp/pkg/logger"
)

type stockService struct {
	emitLog logger.EmitLog
	Q       *query.Query
}

func NewStockService(
	logger logger.Logger,
	conn db.Connection,
) repository.StockService {
	return &stockService{
		Q:       conn.GetQ(),
		emitLog: logger.EmitLog("stock-service"),
	}
}



func (s *stockService) GetStockDefault(req *common.RequestContext) (res model.StockDefault, err error) {
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetStockDefault"))
		}
	}()

	stockDefault, err := s.Q.StockDefault.WithContext(req.Ctx).Where(
		s.Q.StockDefault.CompanyID.Eq(req.ActiveCompany.ID),
	).First()
	if err != nil {
		return
	}
	return *stockDefault, nil
}
