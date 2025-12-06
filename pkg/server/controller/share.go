package controller

import (
	"github.com/ryo-arima/locky/pkg/global"
	"github.com/ryo-arima/locky/pkg/server/share"
)

// Local aliases for cleaner logging code
var (
	INFO  = share.GetServerLogger().INFO
	DEBUG = share.GetServerLogger().DEBUG
	WARN  = share.GetServerLogger().WARN
	ERROR = share.GetServerLogger().ERROR
)

// Local MCode definitions
var (
	SRNRSR1 = global.SRNRSR1
	SRNRSR2 = global.SRNRSR2
	Mcode   = global.Mcode
)
