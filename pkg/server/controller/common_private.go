package controller

import (
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type CommonPrivate interface {
}

type commonPrivate struct {
	CommonRepository repository.Common
}

func NewCommonPrivate(commonRepository repository.Common) CommonPrivate {
	return &commonPrivate{CommonRepository: commonRepository}
}
