package response

import (
	"time"

	"github.com/ryo-arima/locky/pkg/entity/model"
)

type Commons struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Commons []Common `json:"commons,omitempty"`
}

type Common struct {
	ID        uint       `json:"id"`
	UUID      string     `json:"uuid"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// Login represents successful login response
type Login struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	TokenPair *model.TokenPair `json:"token_pair,omitempty"`
	User      *User            `json:"user,omitempty"`
}

// RefreshToken represents refresh token response
type RefreshToken struct {
	Code      string           `json:"code"`
	Message   string           `json:"message"`
	TokenPair *model.TokenPair `json:"token_pair,omitempty"`
}

// ValidateToken represents token validation response
type ValidateToken struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
