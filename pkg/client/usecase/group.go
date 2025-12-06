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
	internal internal.Group
	private  private.Group
}

func NewGroup(conf config.BaseConfig) Group {
	return &group{
		internal: internal.NewGroup(conf),
		private:  private.NewGroup(conf),
	}
}

func (u *group) Bootstrap(req request.Group, format string) string {
	resp := u.internal.BootstrapGroupForDB(req)
	return Format(format, resp)
}

func (u *group) GetInternal(req request.Group, format string) string {
	resp := u.internal.GetGroupForInternal(req)
	return Format(format, resp)
}
func (u *group) GetPrivate(req request.Group, format string) string {
	resp := u.private.GetGroupForPrivate(req)
	return Format(format, resp)
}
func (u *group) CreateInternal(req request.Group, format string) string {
	resp := u.internal.CreateGroupForInternal(req)
	return Format(format, resp)
}
func (u *group) CreatePrivate(req request.Group, format string) string {
	resp := u.private.CreateGroupForPrivate(req)
	return Format(format, resp)
}
func (u *group) UpdateInternal(req request.Group, format string) string {
	resp := u.internal.UpdateGroupForInternal(req)
	return Format(format, resp)
}
func (u *group) UpdatePrivate(req request.Group, format string) string {
	resp := u.private.UpdateGroupForPrivate(req)
	return Format(format, resp)
}
func (u *group) DeleteInternal(req request.Group, format string) string {
	resp := u.internal.DeleteGroupForInternal(req)
	return Format(format, resp)
}
func (u *group) DeletePrivate(req request.Group, format string) string {
	resp := u.private.DeleteGroupForPrivate(req)
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
