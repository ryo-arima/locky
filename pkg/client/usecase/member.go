package usecase

import (
	"fmt"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/repository"
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
	internalRepo repository.MemberInternal
	privateRepo  repository.MemberPrivate
}

func NewMember(conf config.BaseConfig) Member {
	return &member{
		internalRepo: repository.NewMemberInternal(conf),
		privateRepo:  repository.NewMemberPrivate(conf),
	}
}

func (u *member) Bootstrap(req request.Member, format string) string {
	resp := u.internalRepo.BootstrapMemberForDB(req)
	return Format(format, resp)
}
func (u *member) GetInternal(req request.Member, format string) string {
	resp := u.internalRepo.GetMember(req)
	return Format(format, resp)
}
func (u *member) GetPrivate(req request.Member, format string) string {
	resp := u.privateRepo.GetMember(req)
	return Format(format, resp)
}
func (u *member) CreateInternal(req request.Member, format string) string {
	resp := u.internalRepo.CreateMember(req)
	return Format(format, resp)
}
func (u *member) CreatePrivate(req request.Member, format string) string {
	resp := u.privateRepo.CreateMember(req)
	return Format(format, resp)
}
func (u *member) UpdateInternal(req request.Member, format string) string {
	resp := u.internalRepo.UpdateMember(req)
	return Format(format, resp)
}
func (u *member) UpdatePrivate(req request.Member, format string) string {
	resp := u.privateRepo.UpdateMember(req)
	return Format(format, resp)
}
func (u *member) DeleteInternal(req request.Member, format string) string {
	resp := u.internalRepo.DeleteMember(req)
	return Format(format, resp)
}
func (u *member) DeletePrivate(req request.Member, format string) string {
	resp := u.privateRepo.DeleteMember(req)
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
