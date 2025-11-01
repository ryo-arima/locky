package controller

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryo-arima/locky/pkg/client/usecase"
	"github.com/ryo-arima/locky/pkg/config"
	"github.com/ryo-arima/locky/pkg/entity/request"
	"github.com/spf13/cobra"
)

// global output format (table/json/yaml)
var outputFormat = "table"

// SetOutputFormat sets global output format
func SetOutputFormat(format string) {
	format = strings.ToLower(strings.TrimSpace(format))
	switch format {
	case "table", "json", "yaml":
		outputFormat = format
	default:
		outputFormat = "table"
	}
}

// GetOutputFormat returns current output format
func GetOutputFormat() string { return outputFormat }

// PrintMessage prints message as per current format via usecase formatter
func PrintMessage(msg string) {
	type message struct {
		Message string `json:"message" yaml:"message"`
	}
	fmt.Print(usecase.Format(GetOutputFormat(), message{Message: msg}))
}

// saveTokenPairToFiles writes access & refresh tokens under etc/.locky/client/{profile}/
func saveTokenPairToFiles(profile string, access, refresh string) {
	if profile == "" {
		profile = "app"
	}
	baseDir := filepath.Join("etc", ".locky", "client", profile)
	_ = os.MkdirAll(baseDir, 0o755)
	if access != "" {
		_ = os.WriteFile(filepath.Join(baseDir, "access_token"), []byte(access), 0o600)
	}
	if refresh != "" {
		_ = os.WriteFile(filepath.Join(baseDir, "refresh_token"), []byte(refresh), 0o600)
	}
}

// loadRefreshTokenFromFiles tries admin then app directory
func loadRefreshTokenFromFiles() string {
	candidates := []string{
		filepath.Join("etc", ".locky", "client", "admin", "refresh_token"),
		filepath.Join("etc", ".locky", "client", "app", "refresh_token"),
	}
	for _, p := range candidates {
		b, err := os.ReadFile(p)
		if err == nil && len(b) > 0 {
			return strings.TrimSpace(string(b))
		}
	}
	return ""
}

// isAdminEmail checks email against config admin emails
func isAdminEmail(conf config.BaseConfig, email string) bool {
	for _, a := range conf.YamlConfig.Application.Server.Admin.Emails {
		if strings.EqualFold(a, email) {
			return true
		}
	}
	return false
}

// InitLoginCmd creates a login command
func InitCommonLoginCmd(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewCommonUsecase(conf)
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "login with email and password",
		Long:  "authenticate user and receive JWT tokens",
		Run: func(cmd *cobra.Command, args []string) {
			email, err := cmd.Flags().GetString("email")
			if err != nil {
				log.Fatal(err)
			}
			password, err := cmd.Flags().GetString("password")
			if err != nil {
				log.Fatal(err)
			}

			if email == "" || password == "" {
				PrintMessage("Email and password are required")
				return
			}

			loginResponse := uc.Login(request.LoginRequest{
				Email:    email,
				Password: password,
			})

			// Save tokens to environment variables or files for later use
			if loginResponse.TokenPair != nil {
				os.Setenv("LOCKY_ACCESS_TOKEN", loginResponse.TokenPair.AccessToken)
				os.Setenv("LOCKY_REFRESH_TOKEN", loginResponse.TokenPair.RefreshToken)
				profile := "app"
				if loginResponse.User != nil && isAdminEmail(conf, loginResponse.User.Email) {
					profile = "admin"
				}
				saveTokenPairToFiles(profile, loginResponse.TokenPair.AccessToken, loginResponse.TokenPair.RefreshToken)
			}

			fmt.Print(usecase.Format(GetOutputFormat(), loginResponse))
		},
	}
	loginCmd.Flags().StringP("email", "e", "", "user email")
	loginCmd.Flags().StringP("password", "p", "", "user password")
	loginCmd.MarkFlagRequired("email")
	loginCmd.MarkFlagRequired("password")
	return loginCmd
}

// InitRefreshTokenCmd creates a refresh token command
func InitCommonRefreshTokenCmd(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewCommonUsecase(conf)
	refreshCmd := &cobra.Command{
		Use:   "refresh",
		Short: "refresh access token using refresh token",
		Long:  "refresh the access token using the stored refresh token",
		Run: func(cmd *cobra.Command, args []string) {
			refreshToken, err := cmd.Flags().GetString("refresh-token")
			if err != nil {
				log.Fatal(err)
			}

			// If no refresh token provided via flag, try environment variable
			if refreshToken == "" {
				refreshToken = os.Getenv("LOCKY_REFRESH_TOKEN")
			}

			// If still empty, try file system
			if refreshToken == "" {
				refreshToken = loadRefreshTokenFromFiles()
			}

			if refreshToken == "" {
				PrintMessage("Refresh token is required. Provide via --refresh-token flag, env var, or login first")
				return
			}

			refreshResponse := uc.RefreshToken(refreshToken)

			// Update environment variables & files
			if refreshResponse.TokenPair != nil {
				os.Setenv("LOCKY_ACCESS_TOKEN", refreshResponse.TokenPair.AccessToken)
				os.Setenv("LOCKY_REFRESH_TOKEN", refreshResponse.TokenPair.RefreshToken)
				// Cannot know which profile with only token; keep both for convenience
				saveTokenPairToFiles("admin", refreshResponse.TokenPair.AccessToken, refreshResponse.TokenPair.RefreshToken)
				saveTokenPairToFiles("app", refreshResponse.TokenPair.AccessToken, refreshResponse.TokenPair.RefreshToken)
			}

			fmt.Print(usecase.Format(GetOutputFormat(), refreshResponse))
		},
	}
	refreshCmd.Flags().StringP("refresh-token", "r", "", "refresh token")
	return refreshCmd
}

// InitLogoutCmd creates a logout command
func InitCommonLogoutCmd(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewCommonUsecase(conf)
	logoutCmd := &cobra.Command{
		Use:   "logout",
		Short: "logout and invalidate tokens",
		Long:  "logout user and optionally invalidate tokens",
		Run: func(cmd *cobra.Command, args []string) {
			accessToken, err := cmd.Flags().GetString("access-token")
			if err != nil {
				log.Fatal(err)
			}

			// If no access token provided via flag, try environment variable
			if accessToken == "" {
				accessToken = os.Getenv("LOCKY_ACCESS_TOKEN")
			}

			logoutResponse := uc.Logout(accessToken)

			// Clear environment variables on success
			if logoutResponse.Code == "" || logoutResponse.Code == "SERVER_CONTROLLER_LOGOUT_SUCCESS" {
				os.Unsetenv("LOCKY_ACCESS_TOKEN")
				os.Unsetenv("LOCKY_REFRESH_TOKEN")
			}

			fmt.Print(usecase.Format(GetOutputFormat(), logoutResponse))
		},
	}
	logoutCmd.Flags().StringP("access-token", "a", "", "access token")
	return logoutCmd
}

// InitValidateTokenCmd creates a validate token command
func InitCommonValidateTokenCmd(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewCommonUsecase(conf)
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "validate access token",
		Long:  "validate the current access token",
		Run: func(cmd *cobra.Command, args []string) {
			accessToken, err := cmd.Flags().GetString("access-token")
			if err != nil {
				log.Fatal(err)
			}

			// If no access token provided via flag, try environment variable
			if accessToken == "" {
				accessToken = os.Getenv("LOCKY_ACCESS_TOKEN")
			}

			if accessToken == "" {
				PrintMessage("Access token is required. Provide via --access-token flag or LOCKY_ACCESS_TOKEN environment variable")
				return
			}

			validateResponse := uc.ValidateToken(accessToken)
			fmt.Print(usecase.Format(GetOutputFormat(), validateResponse))
		},
	}
	validateCmd.Flags().StringP("access-token", "a", "", "access token")
	return validateCmd
}

// InitGetUserInfoCmd creates a get user info command
func InitCommonUserInfoCmd(conf config.BaseConfig) *cobra.Command {
	uc := usecase.NewCommonUsecase(conf)
	userInfoCmd := &cobra.Command{
		Use:   "userinfo",
		Short: "get user information",
		Long:  "get user information using access token",
		Run: func(cmd *cobra.Command, args []string) {
			accessToken, err := cmd.Flags().GetString("access-token")
			if err != nil {
				log.Fatal(err)
			}

			// If no access token provided via flag, try environment variable
			if accessToken == "" {
				accessToken = os.Getenv("LOCKY_ACCESS_TOKEN")
			}

			if accessToken == "" {
				PrintMessage("Access token is required. Provide via --access-token flag or LOCKY_ACCESS_TOKEN environment variable")
				return
			}

			userInfoResponse := uc.GetUserInfo(accessToken)
			fmt.Print(usecase.Format(GetOutputFormat(), userInfoResponse))
		},
	}
	userInfoCmd.Flags().StringP("access-token", "a", "", "access token")
	return userInfoCmd
}
