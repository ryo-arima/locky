package controller

import (
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type CommonControllerForPrivate interface {
}

type commonControllerForPrivate struct {
	CommonRepository repository.CommonRepository
}

func NewCommonControllerForPrivate(commonRepository repository.CommonRepository) CommonControllerForPrivate {
	return &commonControllerForPrivate{CommonRepository: commonRepository}
}
