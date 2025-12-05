package server

import (
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/server/share"
)

func Main(conf config.BaseConfig) {
	conf.Logger.INFO(share.ToConfigMCode(share.SM1), "Starting locky server on port 8000", nil)
	router := InitRouter(conf)
	conf.Logger.INFO(share.ToConfigMCode(share.SM3), "Server is ready", nil)
	router.Run(":8000")
}
