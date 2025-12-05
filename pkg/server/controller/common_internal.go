package controller

import (
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type CommonInternal interface {
}

type commonInternal struct {
	CommonRepository repository.Common
}

func NewCommonInternal(commonRepository repository.Common) CommonInternal {
	return &commonInternal{CommonRepository: commonRepository}
}
