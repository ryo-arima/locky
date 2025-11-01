package response

import (
	"time"

	"github.com/ryo-arima/locky/pkg/entity/model"
)

type CommonResponse struct {
	Code    string
	Message string
	Commons []Common
}

type Common struct {
	ID        uint
	UUID      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}

// LoginResponse represents successful login response
type LoginResponse struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	TokenPair *model.TokenPair `json:"token_pair,omitempty"`
	User      *User            `json:"user,omitempty"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	TokenPair *model.TokenPair `json:"token_pair,omitempty"`
}
