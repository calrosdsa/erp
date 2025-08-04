package addressservice

// import (
// 	"context"
// 	"erp/api/common"
// 	"erp/api/dto"
// 	"erp/internal/app/connection"
// 	"erp/internal/app/entity"
// 	"erp/internal/app/service/helpers"
// 	"erp/pkg/logger"
// 	"time"

// 	"gorm.io/gorm"
// )

// type AddressService struct {
// 	conn    *connection.Connection
// 	timeout time.Duration
// 	emitLog helpers.EmitLog
// }

// func NewAddressService(
// 	conn *connection.Connection,
// 	timeout time.Duration,
// 	helpers *helpers.Helpers,
// ) *AddressService {
// 	return &AddressService{
// 		conn:    conn,
// 		timeout: timeout,
// 		emitLog: helpers.Logger.EmitLog("address-service"),
// 	}
// }
// func (s *AddressService) CreatePartyAddress(req *common.RequestContext, i *dto.CreateAddressPartyRequest) (entity.PartyAddress, error) {
// 	ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
// 	defer cancel()
// 	tx := s.conn.Db.Begin()
// 	err := tx.Error
// 	defer func() {
// 		if err != nil {
// 			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreatePartyAddress"))
// 			tx.Rollback()
// 		}
// 	}()
// 	var partyAddress entity.PartyAddress
// 	address, err := s.createAddress(tx, i.Body.Address)
// 	if err != nil {
// 		return partyAddress, err
// 	}
// 	partyAddress.AddressID = address.ID
// 	partyAddress.Address = address
// 	r := i.Body.AddressParty
// 	partyAddress.IsBillingAddress = r.IsBillingAddress
// 	partyAddress.IsShippingAddress = r.IsShippingAddress
// 	partyAddress.PartyID = r.PartyID
// 	err = tx.WithContext(ctx).Save(&partyAddress).Error
// 	if err := tx.Commit().Error; err != nil {
// 		return partyAddress, err
// 	}
// 	return partyAddress, nil
// }

// func (s *AddressService) createAddress(tx *gorm.DB, i dto.AddressDto) (entity.Address, error) {
// 	// ctx, cancel := context.WithTimeout(req.Ctx, s.timeout)
// 	// defer cancel()
// 	var address entity.Address
// 	r := i
// 	address.FullName = r.FullName
// 	address.CountryCode = r.CountryCode
// 	address.Company = r.Company
// 	address.City = r.City
// 	address.Province = r.Province
// 	address.StreetLine1 = r.StreetLine1
// 	address.StreetLine2 = r.StreetLine2
// 	address.PostalCode = r.PostalCode
// 	address.PhoneNumber = r.PhoneNumber
// 	address.IdentificationNumber = r.IdentificationNumber
// 	err := tx.Save(&address).Error
// 	defer func() {
// 		if err != nil {
// 			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateAddress"))
// 		}
// 	}()
// 	return address, err
// }
