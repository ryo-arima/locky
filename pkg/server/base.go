package server

import (
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/middleware"
)

func Main(conf config.BaseConfig) {
	conf.Logger.INFO(middleware.ToConfigMCode(middleware.SM1), "Starting locky server on port 8000", nil)
	router := InitRouter(conf)
	conf.Logger.INFO(middleware.ToConfigMCode(middleware.SM3), "Server is ready", nil)
	router.Run(":8000")
}
