package controller

import (
	"github.com/ryo-arima/locky/pkg/server/share"
)

type CommonPrivate interface {
}

type commonPrivate struct {
	CommonRepository share.Common
}

func NewCommonPrivate(commonRepository share.Common) CommonPrivate {
	return &commonPrivate{CommonRepository: commonRepository}
}
