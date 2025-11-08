package main

import (
	"github.com/ryo-arima/locky/pkg/client"
	"github.com/ryo-arima/locky/pkg/config"
)

func main() {
	conf := config.NewClientConfig()
	// bootstrap コマンド等でDBが必要になる直前に usecase 内で ConnectDB を呼ぶ設計に変更も可
	_ = conf.ConnectDB() // ここで接続 (必要でなければ削除可)
	client.ClientForAdminUser(*conf)
}
