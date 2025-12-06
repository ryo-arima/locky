package server

import (
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/global"
	"github.com/ryo-arima/locky/pkg/server/share"
)

func Main(conf config.BaseConfig) {
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		logger.INFO(global.SSM1, "Starting locky server on port 8000")
	}
	router := InitRouter(conf)
	if logger, ok := conf.Logger.(share.LoggerInterface); ok {
		logger.INFO(global.SSM3, "Server is ready")
	}
	router.Run(":8000")
}
