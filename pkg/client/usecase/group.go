package usecase

import (
	"fmt"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type GroupUsecase interface {
	Bootstrap(request request.GroupRequest, format string) string
	GetInternal(request request.GroupRequest, format string) string
	GetPrivate(request request.GroupRequest, format string) string
	CreateInternal(request request.GroupRequest, format string) string
	CreatePrivate(request request.GroupRequest, format string) string
	UpdateInternal(request request.GroupRequest, format string) string
	UpdatePrivate(request request.GroupRequest, format string) string
	DeleteInternal(request request.GroupRequest, format string) string
	DeletePrivate(request request.GroupRequest, format string) string
}

type groupUsecase struct {
	repo repository.GroupRepository
}

func NewGroupUsecase(conf config.BaseConfig) GroupUsecase {
	return &groupUsecase{repo: repository.NewGroupRepository(conf)}
}

func (u *groupUsecase) Bootstrap(req request.GroupRequest, format string) string {
	resp := u.repo.BootstrapGroupForDB(req)
	return Format(format, resp)
}

func (u *groupUsecase) GetInternal(req request.GroupRequest, format string) string {
	resp := u.repo.GetGroupForInternal(req)
	return Format(format, resp)
}
func (u *groupUsecase) GetPrivate(req request.GroupRequest, format string) string {
	resp := u.repo.GetGroupForPrivate(req)
	return Format(format, resp)
}
func (u *groupUsecase) CreateInternal(req request.GroupRequest, format string) string {
	resp := u.repo.CreateGroupForInternal(req)
	return Format(format, resp)
}
func (u *groupUsecase) CreatePrivate(req request.GroupRequest, format string) string {
	resp := u.repo.CreateGroupForPrivate(req)
	return Format(format, resp)
}
func (u *groupUsecase) UpdateInternal(req request.GroupRequest, format string) string {
	resp := u.repo.UpdateGroupForInternal(req)
	return Format(format, resp)
}
func (u *groupUsecase) UpdatePrivate(req request.GroupRequest, format string) string {
	resp := u.repo.UpdateGroupForPrivate(req)
	return Format(format, resp)
}
func (u *groupUsecase) DeleteInternal(req request.GroupRequest, format string) string {
	resp := u.repo.DeleteGroupForInternal(req)
	return Format(format, resp)
}
func (u *groupUsecase) DeletePrivate(req request.GroupRequest, format string) string {
	resp := u.repo.DeleteGroupForPrivate(req)
	return Format(format, resp)
}

// groupsTableString renders GroupResponse as a table string.
func groupsTableString(res response.GroupResponse) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"ID", "UUID", "NAME"}, "\t"))
	for _, g := range res.Groups {
		fmt.Fprintf(w, "%d\t%s\t%s\n", g.ID, g.UUID, g.Name)
	}
	w.Flush()
	return buf.String()
}
