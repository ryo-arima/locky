package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/model"
	"github.com/ryo-arima/locky/pkg/server/repository"
)

// MockUserRepository implements repository.UserRepository for testing
type MockUserRepository struct {
	Users           []model.Users
	GetListFunc     func(ctx context.Context, filter repository.UserQueryFilter) ([]model.Users, error)
	GetByIDFunc     func(ctx context.Context, id uint) (*model.Users, error)
	CreateFunc      func(ctx context.Context, user *model.Users) error
	UpdateFunc      func(ctx context.Context, user *model.Users) error
	DeleteFunc      func(ctx context.Context, id uint) error
	CountFunc       func(ctx context.Context, filter repository.UserQueryFilter) (int64, error)
}

func (m *MockUserRepository) GetUserList(ctx context.Context, filter repository.UserQueryFilter) ([]model.Users, error) {
	if m.GetListFunc != nil {
		return m.GetListFunc(ctx, filter)
	}
	result := []model.Users{}
	for _, u := range m.Users {
		if filter.Email != nil && u.Email != *filter.Email {
			continue
		}
		if filter.Name != nil && u.Name != *filter.Name {
			continue
		}
		if filter.UUID != nil && u.UUID != *filter.UUID {
			continue
		}
		result = append(result, u)
	}
	return result, nil
}

func (m *MockUserRepository) GetUserByID(ctx context.Context, id uint) (*model.Users, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	for _, u := range m.Users {
		if u.ID == id {
			return &u, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user *model.Users) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	user.ID = uint(len(m.Users) + 1)
	m.Users = append(m.Users, *user)
	return nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user *model.Users) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	for i, u := range m.Users {
		if u.ID == user.ID {
			m.Users[i] = *user
			return nil
		}
	}
	return nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	for i, u := range m.Users {
		if u.ID == id {
			m.Users = append(m.Users[:i], m.Users[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockUserRepository) CountUsers(ctx context.Context, filter repository.UserQueryFilter) (int64, error) {
	if m.CountFunc != nil {
		return m.CountFunc(ctx, filter)
	}
	return int64(len(m.Users)), nil
}

// MockGroupRepository implements repository.GroupRepository for testing
type MockGroupRepository struct {
	Groups      []model.Groups
	GetListFunc func(ctx context.Context, filter repository.GroupQueryFilter) ([]model.Groups, error)
	CreateFunc  func(ctx context.Context, group *model.Groups) error
	UpdateFunc  func(ctx context.Context, group *model.Groups) error
	DeleteFunc  func(ctx context.Context, id uint) error
}

func (m *MockGroupRepository) GetGroupList(ctx context.Context, filter repository.GroupQueryFilter) ([]model.Groups, error) {
	if m.GetListFunc != nil {
		return m.GetListFunc(ctx, filter)
	}
	return m.Groups, nil
}

func (m *MockGroupRepository) GetGroupByID(ctx context.Context, id uint) (*model.Groups, error) {
	for _, g := range m.Groups {
		if g.ID == id {
			return &g, nil
		}
	}
	return nil, nil
}

func (m *MockGroupRepository) CreateGroup(ctx context.Context, group *model.Groups) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, group)
	}
	group.ID = uint(len(m.Groups) + 1)
	m.Groups = append(m.Groups, *group)
	return nil
}

func (m *MockGroupRepository) UpdateGroup(ctx context.Context, group *model.Groups) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, group)
	}
	for i, g := range m.Groups {
		if g.ID == group.ID {
			m.Groups[i] = *group
			return nil
		}
	}
	return nil
}

func (m *MockGroupRepository) DeleteGroup(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	for i, g := range m.Groups {
		if g.ID == id {
			m.Groups = append(m.Groups[:i], m.Groups[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockGroupRepository) CountGroups(ctx context.Context, filter repository.GroupQueryFilter) (int64, error) {
	return int64(len(m.Groups)), nil
}

// MockMemberRepository implements repository.MemberRepository for testing
type MockMemberRepository struct {
	Members     []model.Members
	GetListFunc func(ctx context.Context, filter repository.MemberQueryFilter) ([]model.Members, error)
	CreateFunc  func(ctx context.Context, member *model.Members) error
	DeleteFunc  func(ctx context.Context, id uint) error
}

func (m *MockMemberRepository) GetMemberList(ctx context.Context, filter repository.MemberQueryFilter) ([]model.Members, error) {
	if m.GetListFunc != nil {
		return m.GetListFunc(ctx, filter)
	}
	return m.Members, nil
}

func (m *MockMemberRepository) GetMemberByID(ctx context.Context, id uint) (*model.Members, error) {
	for _, mem := range m.Members {
		if mem.ID == id {
			return &mem, nil
		}
	}
	return nil, nil
}

func (m *MockMemberRepository) CreateMember(ctx context.Context, member *model.Members) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, member)
	}
	member.ID = uint(len(m.Members) + 1)
	m.Members = append(m.Members, *member)
	return nil
}

func (m *MockMemberRepository) DeleteMember(ctx context.Context, id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	for i, mem := range m.Members {
		if mem.ID == id {
			m.Members = append(m.Members[:i], m.Members[i+1:]...)
			return nil
		}
	}
	return nil
}

func (m *MockMemberRepository) CountMembers(ctx context.Context, filter repository.MemberQueryFilter) (int64, error) {
	return int64(len(m.Members)), nil
}

// MockCommonRepository implements repository.CommonRepository for testing
type MockCommonRepository struct {
	JWTSecret          string
	InvalidatedTokens  map[string]bool
	GenerateTokenFunc  func(email, role, uuid string) (string, string, error)
	ValidateTokenFunc  func(tokenString string) (*model.JWTClaims, error)
	InvalidateFunc     func(ctx context.Context, jti string, expiration time.Duration) error
	IsInvalidatedFunc  func(ctx context.Context, jti string) (bool, error)
	VerifyPasswordFunc func(hashedPassword, password string) error
}

func (m *MockCommonRepository) GenerateJWTToken(claims model.JWTClaims) (string, error) {
	return "mock-access-token-" + claims.Email, nil
}

func (m *MockCommonRepository) ValidateJWTToken(tokenString string) (*model.JWTClaims, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(tokenString)
	}
	return &model.JWTClaims{
		Email: "test@example.com",
		Role:  "user",
		UUID:  "test-uuid",
		Jti:   "test-jti",
	}, nil
}

func (m *MockCommonRepository) ParseTokenUnverified(tokenString string) (*model.JWTClaims, error) {
	return &model.JWTClaims{
		Email: "test@example.com",
		Role:  "user",
		UUID:  "test-uuid",
		Jti:   "test-jti",
	}, nil
}

func (m *MockCommonRepository) InvalidateToken(ctx context.Context, tokenString string) error {
	if m.InvalidateFunc != nil {
		return m.InvalidateFunc(ctx, tokenString, 0)
	}
	if m.InvalidatedTokens == nil {
		m.InvalidatedTokens = make(map[string]bool)
	}
	m.InvalidatedTokens[tokenString] = true
	return nil
}

func (m *MockCommonRepository) IsTokenInvalidated(ctx context.Context, jti string) (bool, error) {
	if m.IsInvalidatedFunc != nil {
		return m.IsInvalidatedFunc(ctx, jti)
	}
	if m.InvalidatedTokens == nil {
		return false, nil
	}
	return m.InvalidatedTokens[jti], nil
}

func (m *MockCommonRepository) VerifyPassword(hashedPassword, password string) error {
	if m.VerifyPasswordFunc != nil {
		return m.VerifyPasswordFunc(hashedPassword, password)
	}
	// Simple mock: just compare strings
	if hashedPassword == password {
		return nil
	}
	return fmt.Errorf("password mismatch")
}

func (m *MockCommonRepository) GenerateTokenPair(userID uint, userUUID, email, name, role string) (*model.TokenPair, error) {
	return &model.TokenPair{
		AccessToken:  "mock-access-token",
		RefreshToken: "mock-refresh-token",
	}, nil
}

func (m *MockCommonRepository) GetBaseConfig() config.BaseConfig {
	return config.BaseConfig{}
}

func (m *MockCommonRepository) GenerateJWTSecret() (string, error) {
	return "mock-jwt-secret", nil
}

func (m *MockCommonRepository) ValidateJWTSecretStrength(secret string) error {
	return nil
}

func (m *MockCommonRepository) HashPassword(password string) (string, error) {
	return "hashed-" + password, nil
}

func (m *MockCommonRepository) ValidatePasswordStrength(password string) error {
	return nil
}

func (m *MockCommonRepository) DeleteTokenCache(token string) {
	// Mock implementation
}

func (m *MockCommonRepository) SendEmail(ctx context.Context, to, subject, body string, isHTML bool) error {
	return nil
}

func (m *MockCommonRepository) SendWelcomeEmail(ctx context.Context, to, name string) error {
	return nil
}

func (m *MockCommonRepository) SendPasswordResetEmail(ctx context.Context, to, name, resetURL string) error {
	return nil
}
