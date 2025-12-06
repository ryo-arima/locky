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

type Member interface {
	Bootstrap(request request.Member, format string) string
	GetInternal(request request.Member, format string) string
	GetPrivate(request request.Member, format string) string
	CreateInternal(request request.Member, format string) string
	CreatePrivate(request request.Member, format string) string
	UpdateInternal(request request.Member, format string) string
	UpdatePrivate(request request.Member, format string) string
	DeleteInternal(request request.Member, format string) string
	DeletePrivate(request request.Member, format string) string
}

type member struct {
	repo repository.Member
}

func NewMember(conf config.BaseConfig) Member {
	return &member{repo: repository.NewMember(conf)}
}

func (u *member) Bootstrap(req request.Member, format string) string {
	resp := u.repo.BootstrapMemberForDB(req)
	return Format(format, resp)
}
func (u *member) GetInternal(req request.Member, format string) string {
	resp := u.repo.GetMemberForInternal(req)
	return Format(format, resp)
}
func (u *member) GetPrivate(req request.Member, format string) string {
	resp := u.repo.GetMemberForPrivate(req)
	return Format(format, resp)
}
func (u *member) CreateInternal(req request.Member, format string) string {
	resp := u.repo.CreateMemberForInternal(req)
	return Format(format, resp)
}
func (u *member) CreatePrivate(req request.Member, format string) string {
	resp := u.repo.CreateMemberForPrivate(req)
	return Format(format, resp)
}
func (u *member) UpdateInternal(req request.Member, format string) string {
	resp := u.repo.UpdateMemberForInternal(req)
	return Format(format, resp)
}
func (u *member) UpdatePrivate(req request.Member, format string) string {
	resp := u.repo.UpdateMemberForPrivate(req)
	return Format(format, resp)
}
func (u *member) DeleteInternal(req request.Member, format string) string {
	resp := u.repo.DeleteMemberForInternal(req)
	return Format(format, resp)
}
func (u *member) DeletePrivate(req request.Member, format string) string {
	resp := u.repo.DeleteMemberForPrivate(req)
	return Format(format, resp)
}

// membersTableString renders Members as a table string.
func membersTableString(res response.Members) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"ID", "UUID", "USER_UUID", "GROUP_UUID", "ROLE"}, "\t"))
	for _, m := range res.Members {
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n", m.ID, m.UUID, m.UserUUID, m.GroupUUID, m.Role)
	}
	w.Flush()
	return buf.String()
}
