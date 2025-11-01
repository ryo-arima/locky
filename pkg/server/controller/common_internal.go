package controller

import (
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type CommonControllerForInternal interface {
}

type commonControllerForInternal struct {
	CommonRepository repository.CommonRepository
}

func NewCommonControllerForInternal(commonRepository repository.CommonRepository) CommonControllerForInternal {
	return &commonControllerForInternal{CommonRepository: commonRepository}
}
