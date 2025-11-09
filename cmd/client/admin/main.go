package main

import (
	"github.com/ryo-arima/locky/pkg/client"
	"github.com/ryo-arima/locky/pkg/config"
)

func main() {
	conf := config.NewClientConfig()
	client.ClientForAdminUser(*conf)
}
