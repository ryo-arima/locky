package usecase

import (
	"fmt"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Group interface {
	Bootstrap(request request.Group, format string) string
	GetInternal(request request.Group, format string) string
	GetPrivate(request request.Group, format string) string
	CreateInternal(request request.Group, format string) string
	CreatePrivate(request request.Group, format string) string
	UpdateInternal(request request.Group, format string) string
	UpdatePrivate(request request.Group, format string) string
	DeleteInternal(request request.Group, format string) string
	DeletePrivate(request request.Group, format string) string
}

type group struct {
	internalRepo repository.GroupInternal
	privateRepo  repository.GroupPrivate
}

func NewGroup(conf config.BaseConfig) Group {
	return &group{
		internalRepo: repository.NewGroupInternal(conf),
		privateRepo:  repository.NewGroupPrivate(conf),
	}
}

func (u *group) Bootstrap(req request.Group, format string) string {
	resp := u.internalRepo.BootstrapGroupForDB(req)
	return Format(format, resp)
}

func (u *group) GetInternal(req request.Group, format string) string {
	resp := u.internalRepo.GetGroup(req)
	return Format(format, resp)
}
func (u *group) GetPrivate(req request.Group, format string) string {
	resp := u.privateRepo.GetGroup(req)
	return Format(format, resp)
}
func (u *group) CreateInternal(req request.Group, format string) string {
	resp := u.internalRepo.CreateGroup(req)
	return Format(format, resp)
}
func (u *group) CreatePrivate(req request.Group, format string) string {
	resp := u.privateRepo.CreateGroup(req)
	return Format(format, resp)
}
func (u *group) UpdateInternal(req request.Group, format string) string {
	resp := u.internalRepo.UpdateGroup(req)
	return Format(format, resp)
}
func (u *group) UpdatePrivate(req request.Group, format string) string {
	resp := u.privateRepo.UpdateGroup(req)
	return Format(format, resp)
}
func (u *group) DeleteInternal(req request.Group, format string) string {
	resp := u.internalRepo.DeleteGroup(req)
	return Format(format, resp)
}
func (u *group) DeletePrivate(req request.Group, format string) string {
	resp := u.privateRepo.DeleteGroup(req)
	return Format(format, resp)
}

// groupsTableString renders Groups as a table string.
func groupsTableString(res response.Groups) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"ID", "UUID", "NAME"}, "\t"))
	for _, g := range res.Groups {
		fmt.Fprintf(w, "%d\t%s\t%s\n", g.ID, g.UUID, g.Name)
	}
	w.Flush()
	return buf.String()
}
