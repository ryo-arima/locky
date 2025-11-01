package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"

	"github.com/ryo-arima/locky/pkg/client/repository"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type CommonUsecase interface {
	Login(request request.LoginRequest) response.LoginResponse
	RefreshToken(refreshToken string) response.RefreshTokenResponse
	Logout(accessToken string) response.CommonResponse
	ValidateToken(accessToken string) response.CommonResponse
	GetUserInfo(accessToken string) response.CommonResponse
}

type commonUsecase struct {
	repo repository.CommonRepository
}

func NewCommonUsecase(conf config.BaseConfig) CommonUsecase {
	return &commonUsecase{repo: repository.NewCommonRepository(conf)}
}

func (u *commonUsecase) Login(req request.LoginRequest) response.LoginResponse {
	return u.repo.Login(req)
}
func (u *commonUsecase) RefreshToken(refreshToken string) response.RefreshTokenResponse {
	return u.repo.RefreshToken(refreshToken)
}
func (u *commonUsecase) Logout(accessToken string) response.CommonResponse {
	return u.repo.Logout(accessToken)
}
func (u *commonUsecase) ValidateToken(accessToken string) response.CommonResponse {
	return u.repo.ValidateToken(accessToken)
}
func (u *commonUsecase) GetUserInfo(accessToken string) response.CommonResponse {
	return u.repo.GetUserInfo(accessToken)
}

// Format formats the given value into table, json, or yaml and returns it as string.
func Format(format string, v interface{}) string {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "json":
		b, _ := json.MarshalIndent(v, "", "  ")
		return string(b) + "\n"
	case "yaml":
		b, _ := yaml.Marshal(v)
		return string(b)
	default:
		return tableString(v)
	}
}

func tableString(v interface{}) string {
	switch data := v.(type) {
	case response.UserResponse:
		return usersTableString(data)
	case *response.UserResponse:
		return usersTableString(*data)
	case response.GroupResponse:
		return groupsTableString(data)
	case *response.GroupResponse:
		return groupsTableString(*data)
	case response.MemberResponse:
		return membersTableString(data)
	case *response.MemberResponse:
		return membersTableString(*data)
	case response.RoleResponse:
		return repository.RolesTableStringAlias(data)
	case *response.RoleResponse:
		return repository.RolesTableStringAlias(*data)
	case response.LoginResponse:
		return loginTableString(data)
	case *response.LoginResponse:
		return loginTableString(*data)
	case response.RefreshTokenResponse:
		return refreshTableString(data)
	case *response.RefreshTokenResponse:
		return refreshTableString(*data)
	case response.CommonResponse:
		return commonTableString(data)
	case *response.CommonResponse:
		return commonTableString(*data)
	default:
		b, _ := json.Marshal(data)
		return string(b) + "\n"
	}
}

func newTabWriterBuf() (*tabwriter.Writer, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 2, 4, 2, ' ', 0)
	return w, buf
}

func loginTableString(res response.LoginResponse) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"FIELD", "VALUE"}, "\t"))
	fmt.Fprintf(w, "Code\t%s\n", res.Code)
	fmt.Fprintf(w, "Message\t%s\n", res.Message)
	if res.User != nil {
		fmt.Fprintf(w, "User\t%s (%s)\n", res.User.Name, res.User.Email)
	}
	if res.TokenPair != nil {
		fmt.Fprintf(w, "AccessToken\t%s\n", res.TokenPair.AccessToken)
		fmt.Fprintf(w, "RefreshToken\t%s\n", res.TokenPair.RefreshToken)
		fmt.Fprintf(w, "TokenType\t%s\n", res.TokenPair.TokenType)
		fmt.Fprintf(w, "ExpiresIn\t%d\n", res.TokenPair.ExpiresIn)
	}
	w.Flush()
	return buf.String()
}

func refreshTableString(res response.RefreshTokenResponse) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"FIELD", "VALUE"}, "\t"))
	fmt.Fprintf(w, "Code\t%s\n", res.Code)
	fmt.Fprintf(w, "Message\t%s\n", res.Message)
	if res.TokenPair != nil {
		fmt.Fprintf(w, "AccessToken\t%s\n", res.TokenPair.AccessToken)
		fmt.Fprintf(w, "RefreshToken\t%s\n", res.TokenPair.RefreshToken)
		fmt.Fprintf(w, "TokenType\t%s\n", res.TokenPair.TokenType)
		fmt.Fprintf(w, "ExpiresIn\t%d\n", res.TokenPair.ExpiresIn)
	}
	w.Flush()
	return buf.String()
}

func commonTableString(res response.CommonResponse) string {
	w, buf := newTabWriterBuf()
	fmt.Fprintln(w, strings.Join([]string{"CODE", "MESSAGE"}, "\t"))
	fmt.Fprintf(w, "%s\t%s\n", res.Code, res.Message)
	if len(res.Commons) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, strings.Join([]string{"ID", "UUID", "CREATED_AT"}, "\t"))
		for _, c := range res.Commons {
			created := ""
			if c.CreatedAt != nil {
				created = c.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", c.ID, c.UUID, created)
		}
	}
	w.Flush()
	return buf.String()
}
