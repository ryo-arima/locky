package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/ryo-arima/locky/pkg/entity/response"
)

type CommonRepository interface {
	Login(request request.LoginRequest) (response response.LoginResponse)
	RefreshToken(refreshToken string) (response response.RefreshTokenResponse)
	Logout(accessToken string) (response response.CommonResponse)
	ValidateToken(accessToken string) (response response.CommonResponse)
	GetUserInfo(accessToken string) (response response.CommonResponse)
}

type commonRepository struct {
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
func (rcvr commonRepository) Login(loginRequest request.LoginRequest) (response response.LoginResponse) {
	// Updated to match server router: POST /v1/share/common/auth/tokens
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens"

	// Prepare the HTTP request
	jsonData, err := json.Marshal(loginRequest)
	if err != nil {
		fmt.Println("Error marshaling login request:", err)
		response.Code = "CLIENT_AUTH_LOGIN_001"
		response.Message = "Failed to marshal login request"
		return response
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		response.Code = "CLIENT_AUTH_LOGIN_002"
		response.Message = "Failed to create HTTP request"
		return response
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		response.Code = "CLIENT_AUTH_LOGIN_003"
		response.Message = "Failed to send HTTP request"
		return response
	}
	defer resp.Body.Close()

	// Decode the response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		// try read raw for debug
		_, _ = io.Copy(io.Discard, resp.Body)
		fmt.Println("Error decoding response:", err)
		response.Code = "CLIENT_AUTH_LOGIN_004"
		response.Message = "Failed to decode response"
		return response
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Login successful! Token Pair: %+v\n", response.TokenPair)
		if response.TokenPair != nil {
			os.Setenv("LOCKY_ACCESS_TOKEN", response.TokenPair.AccessToken)
			os.Setenv("LOCKY_REFRESH_TOKEN", response.TokenPair.RefreshToken)
			// Save token (profile determination already saved individually on controller side; this is redundant save)
			saveTokenPair(response.TokenPair.AccessToken, response.TokenPair.RefreshToken)
		}
	} else {
		fmt.Printf("Login failed: %s\n", response.Message)
	}

	return response
}

// RefreshToken refreshes the access token using refresh token
func (rcvr commonRepository) RefreshToken(refreshToken string) (response response.RefreshTokenResponse) {
	// Updated to match server router: POST /v1/share/common/auth/tokens/refresh
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens/refresh"

	requestData := map[string]string{
		"refresh_token": refreshToken,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling refresh request:", err)
		response.Code = "CLIENT_AUTH_REFRESH_001"
		response.Message = "Failed to marshal refresh request"
		return response
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		response.Code = "CLIENT_AUTH_REFRESH_002"
		response.Message = "Failed to create HTTP request"
		return response
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		response.Code = "CLIENT_AUTH_REFRESH_003"
		response.Message = "Failed to send HTTP request"
		return response
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		response.Code = "CLIENT_AUTH_REFRESH_004"
		response.Message = "Failed to decode response"
		return response
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Token refresh successful! New Token Pair: %+v\n", response.TokenPair)
		if response.TokenPair != nil {
			os.Setenv("LOCKY_ACCESS_TOKEN", response.TokenPair.AccessToken)
			os.Setenv("LOCKY_REFRESH_TOKEN", response.TokenPair.RefreshToken)
			// Save token (profile determination already saved individually on controller side; this is redundant save)
			saveTokenPair(response.TokenPair.AccessToken, response.TokenPair.RefreshToken)
		}
	} else {
		fmt.Printf("Token refresh failed: %s\n", response.Message)
	}

	return response
}

// Logout performs user logout
func (rcvr commonRepository) Logout(accessToken string) (response response.CommonResponse) {
	// Updated to match server router: DELETE /v1/share/common/auth/tokens
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens"

	req, err := http.NewRequest("DELETE", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		response.Code = "CLIENT_AUTH_LOGOUT_001"
		response.Message = "Failed to create HTTP request"
		return response
	}

	if accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+accessToken)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		response.Code = "CLIENT_AUTH_LOGOUT_002"
		response.Message = "Failed to send HTTP request"
		return response
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		response.Code = "CLIENT_AUTH_LOGOUT_003"
		response.Message = "Failed to decode response"
		return response
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Logout successful!")
	} else {
		fmt.Printf("Logout failed: %s\n", response.Message)
	}

	return response
}

// ValidateToken validates an access token
func (rcvr commonRepository) ValidateToken(accessToken string) (response response.CommonResponse) {
	// Updated to match server router: GET /v1/share/common/auth/tokens/validate
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens/validate"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		response.Code = "CLIENT_AUTH_VALIDATE_001"
		response.Message = "Failed to create HTTP request"
		return response
	}
	if accessToken == "" {
		// try file system
		accessToken = loadAccessTokenFromFiles()
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		response.Code = "CLIENT_AUTH_VALIDATE_002"
		response.Message = "Failed to send HTTP request"
		return response
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		response.Code = "CLIENT_AUTH_VALIDATE_003"
		response.Message = "Failed to decode response"
		return response
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("Token is valid: %+v\n", response.Commons)
	} else {
		fmt.Printf("Token validation failed: %s\n", response.Message)
	}

	return response
}

// GetUserInfo retrieves user information using access token
func (rcvr commonRepository) GetUserInfo(accessToken string) (response response.CommonResponse) {
	// Updated to match server router: GET /v1/share/common/auth/tokens/user
	endpoint := rcvr.BaseConfig.YamlConfig.Application.Client.ServerEndpoint + "/v1/share/common/auth/tokens/user"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		response.Code = "CLIENT_AUTH_USERINFO_001"
		response.Message = "Failed to create HTTP request"
		return response
	}
	if accessToken == "" {
		accessToken = loadAccessTokenFromFiles()
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		response.Code = "CLIENT_AUTH_USERINFO_002"
		response.Message = "Failed to send HTTP request"
		return response
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding response:", err)
		response.Code = "CLIENT_AUTH_USERINFO_003"
		response.Message = "Failed to decode response"
		return response
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Printf("User info retrieved: %+v\n", response.Commons)
	} else {
		fmt.Printf("Failed to get user info: %s\n", response.Message)
	}

	return response
}

func NewCommonRepository(baseConfig config.BaseConfig) CommonRepository {
	return &commonRepository{BaseConfig: baseConfig}
}
