package usecase

import (
	"fmt"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/repository/internal"
	"github.com/ryo-arima/locky/pkg/client/repository/private"
	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type User interface {
	Bootstrap(request request.User, format string) string
	GetInternal(request request.User, format string) string
	GetPrivate(request request.User, format string) string
	CreatePublic(request request.User, format string) string
	CreatePrivate(request request.User, format string) string
	UpdateInternal(request request.User, format string) string
	UpdatePrivate(request request.User, format string) string
	DeleteInternal(request request.User, format string) string
	DeletePrivate(request request.User, format string) string
}

type user struct {
	repo repository.User
}

func NewUser(conf config.BaseConfig) User {
	return &user{repo: repository.NewUser(conf)}
}

func (u *user) Bootstrap(req request.User, format string) string {
	resp := u.repo.BootstrapUserForDB(req)
	return Format(format, resp)
}

func (u *user) GetInternal(req request.User, format string) string {
	resp := u.repo.GetUserForInternal(req)
	return Format(format, resp)
}

func (u *user) GetPrivate(req request.User, format string) string {
	resp := u.repo.GetUserForPrivate(req)
	return Format(format, resp)
}

func (u *user) CreatePublic(req request.User, format string) string {
	resp := u.repo.CreateUserForPublic(req)
	return Format(format, resp)
}

func (u *user) CreatePrivate(req request.User, format string) string {
	resp := u.repo.CreateUserForPrivate(req)
	return Format(format, resp)
}

func (u *user) UpdateInternal(req request.User, format string) string {
	resp := u.repo.UpdateUserForInternal(req)
	return Format(format, resp)
}

func (u *user) UpdatePrivate(req request.User, format string) string {
	resp := u.repo.UpdateUserForPrivate(req)
	return Format(format, resp)
}

func (u *user) DeleteInternal(req request.User, format string) string {
	resp := u.repo.DeleteUserForInternal(req)
	return Format(format, resp)
}

func (u *user) DeletePrivate(req request.User, format string) string {
	resp := u.repo.DeleteUserForPrivate(req)
	return Format(format, resp)
}

// usersTableString renders Users as a table string.
func usersTableString(res response.Users) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"ID", "UUID", "EMAIL", "NAME"}, "\t"))
	for _, u := range res.Users {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", u.ID, u.UUID, u.Email, u.Name)
	}
	w.Flush()
	return buf.String()
}
