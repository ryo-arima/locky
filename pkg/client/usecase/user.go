package usecase

import (
	"fmt"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/repository"
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
	internalRepo repository.UserInternal
	privateRepo  repository.UserPrivate
	publicRepo   repository.UserPublic
}

func NewUser(conf config.BaseConfig) User {
	return &user{
		internalRepo: repository.NewUserInternal(conf),
		privateRepo:  repository.NewUserPrivate(conf),
		publicRepo:   repository.NewUserPublic(conf),
	}
}

func (u *user) Bootstrap(req request.User, format string) string {
	resp := u.internalRepo.BootstrapUserForDB(req)
	return Format(format, resp)
}

func (u *user) GetInternal(req request.User, format string) string {
	resp := u.internalRepo.GetUser(req)
	return Format(format, resp)
}

func (u *user) GetPrivate(req request.User, format string) string {
	resp := u.privateRepo.GetUser(req)
	return Format(format, resp)
}

func (u *user) CreatePublic(req request.User, format string) string {
	resp := u.publicRepo.CreateUser(req)
	return Format(format, resp)
}

func (u *user) CreatePrivate(req request.User, format string) string {
	resp := u.privateRepo.CreateUser(req)
	return Format(format, resp)
}

func (u *user) UpdateInternal(req request.User, format string) string {
	resp := u.internalRepo.UpdateUser(req)
	return Format(format, resp)
}

func (u *user) UpdatePrivate(req request.User, format string) string {
	resp := u.privateRepo.UpdateUser(req)
	return Format(format, resp)
}

func (u *user) DeleteInternal(req request.User, format string) string {
	resp := u.internalRepo.DeleteUser(req)
	return Format(format, resp)
}

func (u *user) DeletePrivate(req request.User, format string) string {
	resp := u.privateRepo.DeleteUser(req)
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
