package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"

	"github.com/ryo-arima/locky/pkg/client/repository/share"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type Common interface {
	Login(request request.Login) response.Login
	RefreshToken(refreshToken string) response.RefreshToken
	Logout(accessToken string) response.Commons
	ValidateToken(accessToken string) response.ValidateToken
	GetUserInfo(accessToken string) response.Commons
}

type common struct {
	repo share.Common
}

func NewCommon(conf config.BaseConfig) Common {
	return &common{repo: repository.NewCommon(conf)}
}

func (u *common) Login(req request.Login) response.Login {
	return u.repo.Login(req)
}
func (u *common) RefreshToken(refreshToken string) response.RefreshToken {
	return u.repo.RefreshToken(refreshToken)
}
func (u *common) Logout(accessToken string) response.Commons {
	return u.repo.Logout(accessToken)
}
func (u *common) ValidateToken(accessToken string) response.ValidateToken {
	return u.repo.ValidateToken(accessToken)
}
func (u *common) GetUserInfo(accessToken string) response.Commons {
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
	case response.Users:
		return usersTableString(data)
	case *response.Users:
		return usersTableString(*data)
	case response.Groups:
		return groupsTableString(data)
	case *response.Groups:
		return groupsTableString(*data)
	case response.Members:
		return membersTableString(data)
	case *response.Members:
		return membersTableString(*data)
	case response.Roles:
		return share.RolesTableStringAlias(data)
	case *response.Roles:
		return share.RolesTableStringAlias(*data)
	case response.Login:
		return loginTableString(data)
	case *response.Login:
		return loginTableString(*data)
	case response.RefreshToken:
		return refreshTableString(data)
	case *response.RefreshToken:
		return refreshTableString(*data)
	case response.Commons:
		return commonTableString(data)
	case *response.Commons:
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

func loginTableString(res response.Login) string {
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

func refreshTableString(res response.RefreshToken) string {
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

func commonTableString(res response.Commons) string {
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
