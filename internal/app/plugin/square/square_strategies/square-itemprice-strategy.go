package squarestrategies

// import (
// 	"context"
// 	"erp/api/common"
// 	_strategy "erp/internal/app/config/stock_config"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	entitysquare "erp/internal/app/plugin/square/entitiy_square"
// 	squareservice "erp/internal/app/plugin/square/square_service"
// 	squaretypes "erp/internal/app/plugin/square/square_types"
// 	"erp/internal/app/service/helpers"
// 	"erp/pkg/logger"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"time"
// )

// var (
// 	_ _strategy.ItemPriceStrategy = (*SquareItemPriceStrategy)(nil)
// )

// type SquareItemPriceStrategy struct {
// 	conn          *connection.Connection
// 	squareService *squareservice.SquareService
// 	timeout time.Duration
// 	emitLog helpers.EmitLog
// }

// func NewSquareItemPriceStrategy(conn *connection.Connection, squareService *squareservice.SquareService,
// 	timeout *time.Duration,
// 	helpers *helpers.Helpers,
// 	) _strategy.ItemPriceStrategy {
// 	return &SquareItemPriceStrategy{
// 		conn:          conn,
// 		squareService: squareService,
// 		timeout: *timeout,
// 		emitLog: helpers.Logger.EmitLog("square-itemprice-strategy"),
// 	}
// }

// func (s *SquareItemPriceStrategy) GetItemPrice(ctx context.Context, req *common.RequestContext, d *entity.ItemPrice) (res string, err error) {
// 	credentials ,err:= s.squareService.GetCredentials(req)
// 	if err != nil {
// 		s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetItemPrice_GetCredentials"))
// 		// fmt.Println(err)
// 		return
// 	}
// 	squareObject,err := s.getSquareObject(ctx,d.ID)
// 	if err != nil {
// 		s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetItemPrice_getSquareObject"))
// 		return 
// 	}
// 	res,err = s.retrieveCatalogObject(&squareObject,&credentials)
// 	if err != nil {
// 		s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetItemPrice_retrieveCatalogObject"))
// 	}
// 	return
// }

// func (s *SquareItemPriceStrategy) retrieveCatalogObject(squareObject *entitysquare.SquareObject,
// 	credentials *squaretypes.SquareCredentials)(
// 	res string,err error,
// ){
// 	url := fmt.Sprintf("%scatalog/object/%s", s.squareService.GetBaseUrl(),squareObject.ObjectVariationId)
// 	req,err := http.NewRequest("GET",url,nil)
// 	if err != nil {
// 		return
// 	}
// 	req.Header.Add("Content-Type","application/json")
// 	req.Header.Add("Authorization",fmt.Sprintf("Bearer %s",credentials.AccessToken))
// 	req.Header.Add("Square-Version", credentials.ApiVersion)
// 	client := &http.Client{
// 		Timeout: s.timeout,
// 	}
// 	resp,err := client.Do(req)
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	data,err := io.ReadAll(resp.Body)
// 	// err = json.NewDecoder(resp.Body).Decode(&res)
// 	if err != nil {
// 		return
// 	}
// 	return string(data),nil
// }

// func (s *SquareItemPriceStrategy) getSquareObject(ctx context.Context, itemPriceid uint) (res entitysquare.SquareObject, err error) {
// 	err = s.conn.Db.WithContext(ctx).Where(&entitysquare.SquareObject{ItemPriceID: uint(itemPriceid)}).First(&res).Error
// 	return
// }
