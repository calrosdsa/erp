package clientservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/gen/db/model"
	"erp/internal/app/config"
	clientconfig "erp/internal/app/config/client_config"
	"erp/internal/app/connection"
	"erp/internal/app/entity"
	"erp/internal/app/event-bus/event"
	"erp/internal/app/plugin"
	"erp/internal/app/service/helpers"
	// partyservice "erp/internal/app/service/services/party_service"
	userservice "erp/internal/app/service/services/user_service"
	"erp/pkg/logger"
	"fmt"

	// "sync"
	"time"

	"gorm.io/gorm"
)

type ClientService struct {
	conn          *connection.Connection
	timeout       *time.Duration
	plugin        *plugin.PluginModule
	configService *config.ConfigService
	userService   *userservice.UserService
	// partyService  *partyservice.PartyService
	eventHelper   *helpers.EventHelper
	emitLog       helpers.EmitLog
	cache         *helpers.CacheHelper
}

func NewClientService(conn *connection.Connection, timeout *time.Duration, plugin *plugin.PluginModule,
	configService *config.ConfigService, userService *userservice.UserService, helpers *helpers.Helpers,
	// partyService *partyservice.PartyService,
) *ClientService {
	return &ClientService{
		conn:          conn,
		timeout:       timeout,
		plugin:        plugin,
		configService: configService,
		userService:   userService,
		eventHelper:   helpers.Event,
		emitLog:       helpers.Logger.EmitLog("client-service"),
		cache:         helpers.Cache,
		// partyService:  partyService,
	}
}

func (s *ClientService) EditClient(req *common.RequestContext, i *dto.EditClientRequest) error {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var err error
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("EditClient"))
		}
	}()
	var client entity.Client
	data := i.Body
	client.GivenName = data.GivenName
	client.FamilyName = data.FamilyName
	client.OrganizationName = data.OrganizationName
	client.PhoneNumber = data.PhoneNumber.Number
	client.CountryCode = data.PhoneNumber.CountryCode
	err = s.conn.Db.WithContext(ctx).Where(&entity.Client{
		Uuid: i.UserSessionUuid,
	}).Updates(&client).Error
	return err
}

func (s *ClientService) GetClientProfile(req *common.RequestContext) (entity.Client, error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var err error
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("GetClientProfile"))
		}
	}()
	var client entity.Client
	client.Uuid = req.Profile.UUID
	err = s.conn.Db.WithContext(ctx).Preload("ClientKeyValueData").Where(&client).
		First(&client).Error
	return client, err
}

func (s *ClientService) CreateCustomer(req *common.RequestContext, d *dto.CreateClientRequest) (res entity.Client, err error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	tx := s.conn.Db.Begin()
	defer func() {
		if err != nil {
			s.emitLog.Err(err, logger.OptionsLog.WithMethod("CreateCustomer"))
			tx.Rollback()
		}
	}()
	if err = tx.Error; err != nil {
		return
	}
	user, err := s.createClientUser(ctx, tx, req, d)
	if err != nil {
		return
	}

	
	fmt.Println("USERID", user.ID)
	clientRequestData := d.Body.ClientRequestDto
	var client entity.Client
	// client.UserID = user.ID
	client.CompanyID = uint(req.ActiveCompany.ID)
	err = s.getClientCompany(tx,&client)
	
	fmt.Println("CLIENT COMPANY",client)
	// partyClient, err := s.partyService.CreateClientParty(ctx, tx,client)
	// if err != nil {
	// 	return
	// }
	// fmt.Println("CLIENT COMPANY",partyClient)
	// client.ID = uint(partyClient.ID)	
	client.FamilyName = clientRequestData.FamilyName
	client.GivenName = clientRequestData.GivenName
	client.Code = s.conn.GenerateCode(ctx, tx, &entity.Client{}, req.ActiveCompany.ID)
	client.EmailAddress = clientRequestData.EmailAddress
	client.OrganizationName = clientRequestData.CompanyName
	client.PhoneNumber = clientRequestData.PhoneNumber
	// client.UserID = user.ID
	disableClient := clientRequestData.DeleteAt.Valid
	if disableClient {
		fmt.Println("DISABLING CLIENT")
		client.DeletedAt = clientRequestData.DeleteAt
	}
	// client.User = user
	client.CountryCode = clientRequestData.Country.Code
	client.CompanyID = uint(req.ActiveCompany.ID)
	// client.Company = req.ActiveCompany

	// err = tx.WithContext(ctx).Where(entity.Client{EmailAddress: client.EmailAddress,CompanyID: req.ActiveCompany.ID}).Assign(client).
	// 	FirstOrCreate(&client).Error
	err = tx.WithContext(ctx).Save(&client).Error
	if err != nil {
		return
	}
	//Insert client to m2m clients

	err = s.insertUserClients(ctx, tx, int(user.ID), int(client.ID))
	if err != nil {
		return
	}

	err = s.insertClientKeyValues(ctx, tx, client.ID, clientRequestData.KeyValues)
	if err != nil {
		return
	}

	err = tx.Commit().Error
	if err != nil {
		return
	}
	// var wg sync.WaitGroup
	for _, plugin := range clientRequestData.Plugins {
		if plugin.Plugin == config.PLUGIN_SQUARE {
			strategy := s.getClientStrategy(plugin.Plugin)
			err = strategy.CreateCustomer(req, &client, clientRequestData.Metadata)
			if err != nil {
				fmt.Println("FAIL TO CREATE SQUARE CUSTOMER", err)
				return client, err
			}
		} else {
			// wg.Add(1)
			go func(req *common.RequestContext, client *entity.Client, companyPlugin *entity.CompanyPlugins) {
				strategy := s.getClientStrategy(companyPlugin.Plugin)
				err := strategy.CreateCustomer(req, client, clientRequestData.Metadata)
				if err != nil {
					fmt.Println("FAIL TO CREATE CUSTOMER", err)
				}
				// wg.Done()
			}(req, &client, &plugin)
		}
	}

	if !disableClient {
		s.SendCredentialEmail(*req, client)
	}

	// wg.Wait()

	return client, nil
}

func (s *ClientService) getClientCompany(tx *gorm.DB,client *entity.Client)(error){
	err := tx.Unscoped().Where(&entity.Client{CompanyID: client.CompanyID,UserID: client.UserID}).First(&client).Error
	return err
}

func (s *ClientService) SendCredentialEmail(req common.RequestContext, client entity.Client) {
	go func() {
		defer func() {
			fmt.Println("CLOSING GOROUTINE")
		}()
		s.eventHelper.Publish(event.NOTIFICATION_EVENT, event.NotificationData{
			NotificationEventType: event.NOTIFY_NEW_CLIENT,
			Data: event.NotificationPayload{
				Payload:        client,
				RequestContext: req,
			},
		},
		)
	}()
}

func (s *ClientService) insertUserClients(ctx context.Context, tx *gorm.DB, userId int, clientId int) (err error) {
	err = tx.WithContext(ctx).Exec(`insert into user_clients(user_id,client_id)values(?,?)
	 on conflict (user_id,client_id) do nothing`,
		userId, clientId).Error
	return
}

func (s *ClientService) insertClientKeyValues(ctx context.Context, tx *gorm.DB, clientId uint,
	keyValues []entity.ClientKeyValueData) (err error) {
	if len(keyValues) == 0 {
		return
	}
	for i, keyValue := range keyValues {
		keyValue.BaseID = clientId
		keyValues[i] = keyValue
	}
	err = tx.WithContext(ctx).CreateInBatches(keyValues, len(keyValues)).Error
	return
}

func (s *ClientService) createClientUser(ctx context.Context, tx *gorm.DB, req *common.RequestContext,
	d *dto.CreateClientRequest) (res model.User, err error) {
	res.Identifier = d.Body.ClientRequestDto.EmailAddress
	role, err := s.getClientRole(ctx)
	if err != nil {
		return
	}
	// err = s.userService.InsertUser(ctx, tx, &res)
	// if err != nil {
	// 	fmt.Println("FAIL TO INSERT USER")
	// 	return
	// }
	fmt.Println("USER", res)
	// err = s.userService.InsertUserCompany(ctx, tx, res.ID, req.ActiveCompany.ID)
	// if err != nil {
	// 	fmt.Println("FAIL TO INSERT COMPANY")
	// 	return
	// }
	err = s.userService.InsertUserRole(ctx, tx, res.ID, role.ID)
	if err != nil {
		fmt.Println("FAILT TO INSERT ROLE")
		return
	}
	return
}

// func (s *ClientService) GetClient(ctx context.Context,code string) (client entity.Client, err error) {
// 	ctx, cancel := context.WithTimeout(ctx, *s.timeout)
// 	defer cancel()
// 	defer func(){
// 		if err != nil {
// 			s.emitLog.Err(err,logger.OptionsLog.WithMethod("GetClient"))

// 		}
// 	}()
// 	err = s.getClientFromCache(ctx, &client)
// 	if err != nil {
// 		err = s.conn.Db.WithContext(ctx).Preload("Company").
// 			Preload("ClientKeyValueData").Where(&entity.Client{Code: code}).First(&client).Error
// 		if err != nil {
// 			return
// 		}

// 		err = s.cache.Set(ctx, s.getClientKeyCache(code), client)
// 		if err != nil {
// 			return
// 		}
// 	}
// 	return
// }

// func (s *ClientService) getClientKeyCache(id string) string {
// 	return fmt.Sprintf("client-%d", id)
// }
// func (s *ClientService) getClientFromCache(ctx context.Context, client *entity.Client) (err error) {
// 	key := s.getClientKeyCache(client.ID)
// 	err = s.cache.Get(ctx, key, &client)
// 	if err != nil {
// 		fmt.Println("USER NOT PRESENT ON CACHE")
// 	}
// 	return
// }

func (s *ClientService) getClientRole(ctx context.Context) (model.Role, error) {
	var role model.Role
	err := s.conn.Db.WithContext(ctx).First(&role, "code='client'").Error
	return role, err
}

func (s *ClientService) getClientStrategy(plugin string) clientconfig.ClientStrategy {
	switch plugin {
	case config.PLUGIN_SQUARE:
		return s.plugin.GetPlugin(config.PLUGIN_SQUARE).ClientStrategy
	}
	return clientconfig.DefaultItemStrategy{}
}
