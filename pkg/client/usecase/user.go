package usecase

import (
	"fmt"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type UserUsecase interface {
	Bootstrap(request request.UserRequest, format string) string
	GetInternal(request request.UserRequest, format string) string
	GetPrivate(request request.UserRequest, format string) string
	CreatePublic(request request.UserRequest, format string) string
	CreatePrivate(request request.UserRequest, format string) string
	UpdateInternal(request request.UserRequest, format string) string
	UpdatePrivate(request request.UserRequest, format string) string
	DeleteInternal(request request.UserRequest, format string) string
	DeletePrivate(request request.UserRequest, format string) string
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(conf config.BaseConfig) UserUsecase {
	return &userUsecase{repo: repository.NewUserRepository(conf)}
}

func (u *userUsecase) Bootstrap(req request.UserRequest, format string) string {
	resp := u.repo.BootstrapUserForDB(req)
	return Format(format, resp)
}

func (u *userUsecase) GetInternal(req request.UserRequest, format string) string {
	resp := u.repo.GetUserForInternal(req)
	return Format(format, resp)
}

func (u *userUsecase) GetPrivate(req request.UserRequest, format string) string {
	resp := u.repo.GetUserForPrivate(req)
	return Format(format, resp)
}

func (u *userUsecase) CreatePublic(req request.UserRequest, format string) string {
	resp := u.repo.CreateUserForPublic(req)
	return Format(format, resp)
}

func (u *userUsecase) CreatePrivate(req request.UserRequest, format string) string {
	resp := u.repo.CreateUserForPrivate(req)
	return Format(format, resp)
}

func (u *userUsecase) UpdateInternal(req request.UserRequest, format string) string {
	resp := u.repo.UpdateUserForInternal(req)
	return Format(format, resp)
}

func (u *userUsecase) UpdatePrivate(req request.UserRequest, format string) string {
	resp := u.repo.UpdateUserForPrivate(req)
	return Format(format, resp)
}

func (u *userUsecase) DeleteInternal(req request.UserRequest, format string) string {
	resp := u.repo.DeleteUserForInternal(req)
	return Format(format, resp)
}

func (u *userUsecase) DeletePrivate(req request.UserRequest, format string) string {
	resp := u.repo.DeleteUserForPrivate(req)
	return Format(format, resp)
}

// usersTableString renders UserResponse as a table string.
func usersTableString(res response.UserResponse) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"ID", "UUID", "EMAIL", "NAME"}, "\t"))
	for _, u := range res.Users {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", u.ID, u.UUID, u.Email, u.Name)
	}
	w.Flush()
	return buf.String()
}
