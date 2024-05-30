package config

import (
	"main.go/config/app_config"
	"main.go/config/db_config"
)

func InitConfig() {
	app_config.InitAppConfig()
	db_config.InitDbConfig()
}
