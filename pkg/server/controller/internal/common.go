package internal

import (
	"github.com/ryo-arima/locky/pkg/server/repository"
)

type CommonController interface {
}

type commonController struct {
	CommonRepository repository.CommonRepository
}

func NewCommonController(commonRepository repository.CommonRepository) CommonController {
	return &commonController{CommonRepository: commonRepository}
}
