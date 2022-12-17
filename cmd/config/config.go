package config

import (
	"time"

	apiconfig "github.com/cfabrica46/api-config"
)

func configEntries() []apiconfig.ConfigEntry {
	return []apiconfig.ConfigEntry{
		{
			VariableName: "port",
			Description:  "Puerto a utilizar",
			Shortcut:     "p",
			DefaultValue: ":8080",
		},
		{
			VariableName: "timeout",
			Description:  "timeout por defecto ",
			DefaultValue: 30,
		},
		{
			VariableName: "uri_prefix",
			Description:  "Prefijo de URL con version",
			DefaultValue: "",
		},
		{
			VariableName: "database_user",
			Description:  "Usuario DB",
			DefaultValue: "cfabrica46",
		},
		{
			VariableName: "database_pass",
			Description:  "Password DB",
			DefaultValue: "abcd",
		},
		{
			VariableName: "database_host",
			Description:  "Direccion IPV4",
			DefaultValue: "localhost",
		},
		{
			VariableName: "database_port",
			Description:  "PORT TCP",
			DefaultValue: "5432",
		},
		{
			VariableName: "database_name",
			Description:  "Name DB",
			DefaultValue: "gokit_app_storage",
		},
	}
}

type APIConfig struct {
	*apiconfig.CfgBase
	DBConfig DBConfig
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func GetAPIConfig() (*APIConfig, error) {
	typeResolver := apiconfig.NewVariableTypeResolver()
	flagConfigurator := apiconfig.NewFlagConfigurator(typeResolver)
	configurator := apiconfig.NewConfigurator(flagConfigurator, typeResolver)

	cfg, err := configurator.Configure(configEntries())
	if err != nil {
		return nil, err
	}

	return &APIConfig{
		CfgBase: &apiconfig.CfgBase{
			Port:      cfg["port"].(string),
			Timeout:   time.Duration(cfg["timeout"].(int)) * time.Second,
			URIPrefix: cfg["uri_prefix"].(string),
		},
		DBConfig: DBConfig{
			Host:     cfg["database_host"].(string),
			Port:     cfg["database_port"].(string),
			User:     cfg["database_user"].(string),
			Password: cfg["database_pass"].(string),
			DBName:   cfg["database_name"].(string),
		},
	}, nil
}
