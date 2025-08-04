package pluginservice

import (
	"context"
	"erp/api/common"
	"erp/api/dto"
	"erp/internal/app/config"
	"erp/internal/app/connection"
	"erp/internal/app/entity"
	"fmt"
	"time"
)

type PluginService struct {
	configService *config.ConfigService
	timeout       *time.Duration
	conn          *connection.Connection
}

func NewPluginService(configService *config.ConfigService, timeout *time.Duration,
	conn *connection.Connection) *PluginService {
	return &PluginService{
		configService: configService,
		timeout:       timeout,
		conn:          conn,
	}
}

func (s *PluginService) GetPlugins() []config.PluginApp {
	return s.configService.GetPlugins()
}

func (s *PluginService) GetPlugin(req *common.RequestContext, d *dto.PluginDetailRequest) (entity.CompanyPlugins, error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	var companyPlugin entity.CompanyPlugins
	err := s.conn.Db.WithContext(ctx).First(&companyPlugin, "company_id = $1 and plugin = $2", req.ActiveCompany.ID, d.Plugin).Error
	err = s.conn.Db.Raw("select company_id,plugin, pgp_sym_decrypt(credentials::bytea, '!qazSymKeyXsw2') as credentials from company_plugins where company_id = $1 and plugin = $2",req.ActiveCompany.ID,d.Plugin).
	Scan(&companyPlugin).Error
	if err != nil {
		return entity.CompanyPlugins{}, err
	}
	return companyPlugin, err
}

func (s *PluginService) AddPlugin(req *common.RequestContext, d *dto.AddPluginRequest) (entity.CompanyPlugins, error) {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	
	var companyPlugin entity.CompanyPlugins
	companyPlugin.CompanyID = int(req.ActiveCompany.ID)
	companyPlugin.Plugin = d.Body.Plugin
	// err := s.conn.Db.WithContext(ctx).Save(&companyPlugin).Error
	err := s.conn.Db.WithContext(ctx).Exec("insert into company_plugins(company_id,plugin,credentials) values($1,$2,pgp_sym_encrypt($3,'!qazSymKeyXsw2'))",
    companyPlugin.CompanyID,companyPlugin.Plugin,companyPlugin.Credentials).Error
	fmt.Println("ERR ADD PLUGIN",err)
	return companyPlugin, err
}

func (s *PluginService) UpdatePluginCredentials(req *common.RequestContext, d *dto.UpdateCredentialsPluginRequest) error {
	ctx, cancel := context.WithTimeout(req.Ctx, *s.timeout)
	defer cancel()
	cryptoPass := s.configService.GetDbConfig().CryptoPass
	// fmt.Println(d.Body,company.ID,d.Plugin)
	err := s.conn.Db.WithContext(ctx).Exec("update company_plugins set credentials = pgp_sym_encrypt($1,$4) where company_id = $2 and plugin = $3",
		d.Body.Credentials, req.ActiveCompany.ID, d.Plugin,cryptoPass).Error
	if err != nil {
		return err
	}
	return nil
}

// func (s *PluginService) getCompany(ctx context.Context, uuid string) (entity.Company, error) {
// 	var company entity.Company
// 	err := s.conn.Db.WithContext(ctx).First(&company, "uuid = $1", uuid).Error
	
// 	if err != nil {
// 		return company, err
// 	}
// 	return company, err
// }
