package repository

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	BaseConfig config.BaseConfig
}

// helper: token file paths
func tokenDirs() []string {
	return []string{
		filepath.Join("etc", ".locky", "client", "admin"),
		filepath.Join("etc", ".locky", "client", "app"),
	}
}

func readFirstExisting(paths []string) string {
	for _, p := range paths {
		b, err := os.ReadFile(p)
		if err == nil && len(b) > 0 {
			return strings.TrimSpace(string(b))
		}
	}
	return ""
}

func loadAccessTokenFromFiles() string {
	var candidates []string
	for _, d := range tokenDirs() {
		candidates = append(candidates, filepath.Join(d, "access_token"))
	}
	return readFirstExisting(candidates)
}

func saveTokenPair(access, refresh string) {
	for _, d := range tokenDirs() {
		_ = os.MkdirAll(d, 0o755)
		if access != "" {
			_ = os.WriteFile(filepath.Join(d, "access_token"), []byte(access), 0o600)
		}
		if refresh != "" {
			_ = os.WriteFile(filepath.Join(d, "refresh_token"), []byte(refresh), 0o600)
		}
	}
}

// Login performs user authentication and returns JWT tokens
func (rcvr *common) Login(loginRequest request.Login) response.Login {
	var result response.Login
	// Updated to match server router: POST /v1/share/common/auth/tokens
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens"

	// Prepare the HTTP request
	jsonData, err := json.Marshal(loginRequest)
	if err != nil {
		result.Code = "CLIENT_AUTH_LOGIN_001"
		result.Message = "Failed to marshal login request"
		return result
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Code = "CLIENT_AUTH_LOGIN_002"
		result.Message = "Failed to create HTTP request"
		return result
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Code = "CLIENT_AUTH_LOGIN_003"
		result.Message = "Failed to send HTTP request"
		return result
	}
	defer resp.Body.Close()

	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		// try read raw for debug
		_, _ = io.Copy(io.Discard, resp.Body)
		result.Code = "CLIENT_AUTH_LOGIN_004"
		result.Message = "Failed to decode response"
		return result
	}

	// Save token on success (no output here - handled by controller layer)
	if resp.StatusCode == http.StatusOK && result.TokenPair != nil {
		os.Setenv("LOCKY_ACCESS_TOKEN", result.TokenPair.AccessToken)
		os.Setenv("LOCKY_REFRESH_TOKEN", result.TokenPair.RefreshToken)
		// Save token (profile determination already saved individually on controller side; this is redundant save)
		saveTokenPair(result.TokenPair.AccessToken, result.TokenPair.RefreshToken)
	}

	return result
}

// RefreshToken refreshes the access token using refresh token
func (rcvr *common) RefreshToken(refreshToken string) response.RefreshToken {
	// Updated to match server router: POST /v1/share/common/auth/tokens/refresh
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens/refresh"

	requestData := map[string]string{
		"refresh_token": refreshToken,
	}

	var result response.RefreshToken
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		result.Code = "CLIENT_AUTH_REFRESH_001"
		result.Message = "Failed to marshal refresh request"
		return result
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		result.Code = "CLIENT_AUTH_REFRESH_002"
		result.Message = "Failed to create HTTP request"
		return result
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Code = "CLIENT_AUTH_REFRESH_003"
		result.Message = "Failed to send HTTP request"
		return result
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "CLIENT_AUTH_REFRESH_004"
		result.Message = "Failed to decode response"
		return result
	}

	if resp.StatusCode == http.StatusOK {
		if result.TokenPair != nil {
			os.Setenv("LOCKY_ACCESS_TOKEN", result.TokenPair.AccessToken)
			os.Setenv("LOCKY_REFRESH_TOKEN", result.TokenPair.RefreshToken)
			// Save token (profile determination already saved individually on controller side; this is redundant save)
			saveTokenPair(result.TokenPair.AccessToken, result.TokenPair.RefreshToken)
		}
	}

	return result
}

// Logout performs user logout
func (rcvr *common) Logout(accessToken string) response.Commons {
	// Updated to match server router: DELETE /v1/share/common/auth/tokens
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens"

	var result response.Commons
	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		result.Code = "CLIENT_AUTH_LOGOUT_001"
		result.Message = "Failed to create HTTP request"
		return result
	}

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Code = "CLIENT_AUTH_LOGOUT_002"
		result.Message = "Failed to send HTTP request"
		return result
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "CLIENT_AUTH_LOGOUT_003"
		result.Message = "Failed to decode response"
		return result
	}

	return result
}

// ValidateToken validates an access token
func (rcvr *common) ValidateToken(accessToken string) response.ValidateToken {
	// Updated to match server router: GET /v1/share/common/auth/tokens/validate
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens/validate"

	var result response.ValidateToken
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		result.Code = "CLIENT_AUTH_VALIDATE_001"
		result.Message = "Failed to create HTTP request"
		return result
	}
	if accessToken == "" {
		// try file system
		accessToken = loadAccessTokenFromFiles()
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Code = "CLIENT_AUTH_VALIDATE_002"
		result.Message = "Failed to send HTTP request"
		return result
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "CLIENT_AUTH_VALIDATE_003"
		result.Message = "Failed to decode response"
		return result
	}

	return result
}

// GetUserInfo retrieves user information using access token
func (rcvr *common) GetUserInfo(accessToken string) response.Commons {
	// Updated to match server router: GET /v1/share/common/auth/tokens/user
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens/user"

	var result response.Commons
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		result.Code = "CLIENT_AUTH_USERINFO_001"
		result.Message = "Failed to create HTTP request"
		return result
	}
	if accessToken == "" {
		accessToken = loadAccessTokenFromFiles()
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		result.Code = "CLIENT_AUTH_USERINFO_002"
		result.Message = "Failed to send HTTP request"
		return result
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		result.Code = "CLIENT_AUTH_USERINFO_003"
		result.Message = "Failed to decode response"
		return result
	}

	return result
}

func NewCommon(baseConfig config.BaseConfig) Common {
	return &common{BaseConfig: baseConfig}
}
