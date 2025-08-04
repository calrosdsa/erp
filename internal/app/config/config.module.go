package config

type ConfigModule struct{
	ConfigService *ConfigService
}

func Init(config *AppConfig) *ConfigModule{
	configService :=  NewConfigService(config)
	return &ConfigModule{
		ConfigService: configService,
	}
}