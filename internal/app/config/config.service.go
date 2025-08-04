package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type ConfigService struct {
	config *AppConfig
}
var ErrEnvVarEmpty = errors.New("getenv: environment variable empty")
func NewConfigService(config *AppConfig) *ConfigService {
	fmt.Println(config.EmailOptions)
	if config.EmailOptions.Processor == nil {
		config.EmailOptions.Processor = &ProcessorOptions{QueueSize: 2000,NumWorkers: 50,}
	}
	if config.EmailOptions.Transport == nil {
		host,_:= getenvStr("EMAIL_HOST")
		port,_ := getenvInt("EMAIL_PORT")
		emailUser,_ := getenvStr("EMAIL_USER")
		emailPass,_ := getenvStr("EMAIL_PASS")
		config.EmailOptions.Transport = &TransportOptions{
			Host: host,
			Port: port,
			Auth: Auth{
				User: emailUser,
				Pass:emailPass,
			},
		}
	}
	fmt.Println(*config.EmailOptions.Processor)
	fmt.Println(config.EmailOptions)
	c:= &ConfigService{
		config:config,
	}


	return c
}


func (c *ConfigService) GetClientConfig() ClientConfig {
	return c.config.ClientConfig
}

func (c *ConfigService) GetEmailOptions() EmailOptions{
	fmt.Println(c.config.EmailOptions)
	
	// if c.config.EmailOptions == nil {
	// 	c.config.EmailOptions = &EmailOptions{
	// 		Processor: ProcessorOptions{
	// 			QueueSize: 2000,
	// 			NumWorkers: 50,
	// 		},
	// 	}
	// }
	return c.config.EmailOptions
}


func (c *ConfigService) GetTimeoutAPICall() time.Duration {
	return time.Duration(2) * time.Minute
}

func (c *ConfigService) GetApiOptions() ApiOptions {
	return c.config.ApiOptions
}

func (c *ConfigService) GetDbConfig() DbConfig {
	return c.config.DbConfig
}

func (c *ConfigService) GetPlugins() []PluginApp {
	return c.config.Plugins
}

func (c *ConfigService) GetPermifyConfig() PermifyAuthorization {
	return c.config.PermifyAuthorization
}

func(c *ConfigService)GetDefaultLanguage() string {
	return c.config.DefaultLanguage
}

func getenvStr(key string) (string, error) {
    v := os.Getenv(key)
    if v == "" {
        return v, ErrEnvVarEmpty
    }
    return v, nil
}

func getenvInt(key string) (int, error) {
    s, err := getenvStr(key)
    if err != nil {
        return 0, err
    }
    v, err := strconv.Atoi(s)
    if err != nil {
        return 0, err
    }
    return v, nil
}
