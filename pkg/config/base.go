package config

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// LoggerConfig represents logger configuration
// Logger implementation is provided by server/share or client/share packages
type LoggerConfig struct {
	Component    string `json:"component" yaml:"component"`
	Service      string `json:"service" yaml:"service"`
	Level        string `json:"level" yaml:"level"`
	Structured   bool   `json:"structured" yaml:"structured"`
	EnableCaller bool   `json:"enable_caller" yaml:"enable_caller"`
	Output       string `json:"output" yaml:"output"`
}

type BaseConfig struct {
	DBConnection *gorm.DB
	YamlConfig   YamlConfig
	Logger       interface{} // Logger implementation from server/share or client/share
}

type YamlConfig struct {
	Application Application  `yaml:"Application"`
	MySQL       MySQL        `yaml:"MySQL"`
	Redis       Redis        `yaml:"Redis"`
	Logger      LoggerConfig `yaml:"Logger"`
}

type IntOrString int

// UnmarshalYAML: receive number or string and convert to number. Non-numeric strings return 0 with warning log.
func (ios *IntOrString) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("invalid yaml node for IntOrString")
	}
	s := value.Value
	if n, err := strconv.Atoi(s); err == nil {
		*ios = IntOrString(n)
		return nil
	}
	log.Printf("Redis db value '%s' is not numeric. Defaulting to 0.", s)
	*ios = 0
	return nil
}

type Redis struct {
	Host string      `yaml:"host"`
	Port int         `yaml:"port"`
	User string      `yaml:"user"`
	Pass string      `yaml:"pass"`
	DB   IntOrString `yaml:"db"`
}

// RedisConfig defines Redis-related server configurations
type RedisConfig struct {
	JWTCache bool `yaml:"jwt_cache"` // Enable JWT token caching in Redis
	CacheTTL int  `yaml:"cache_ttl"` // JWT cache TTL in seconds (0 = use token expiry)
}

// Server defines server-related configurations
type Server struct {
	Admin     Admin       `yaml:"admin"`
	JWTSecret string      `yaml:"jwt_secret"`
	LogLevel  string      `yaml:"log_level"` // debug / info / warn / error
	Redis     RedisConfig `yaml:"redis"`     // Redis-related configurations
}

type Mail struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	UseTLS   bool   `yaml:"use_tls"`
}

type Common struct {
}

type Client struct {
	ServerEndpoint string `yaml:"ServerEndpoint"`
	UserEmail      string `yaml:"UserEmail"`
	UserPassword   string `yaml:"UserPassword"`
}

type Application struct {
	Common Common `yaml:"Common"`
	Server Server `yaml:"Server"`
	Client Client `yaml:"Client"`
	Mail   Mail   `yaml:"Mail"`
}

type Admin struct {
	Emails []string `yaml:"emails"`
}

type MySQL struct {
	Host string `yaml:"host"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
	Port string `yaml:"port"`
	Db   string `yaml:"db"`
}

// NewBaseConfig: creates a new BaseConfig instance with configuration loaded from app.yaml or Secrets Manager
func NewBaseConfig() *BaseConfig {
	return NewBaseConfigWithContext(context.Background())
}

// NewClientConfig: creates a new BaseConfig instance for client with custom default config file path
// If CONFIG_FILE is not set, it defaults to etc/app.yaml (same as server)
// Usage: CONFIG_FILE=etc/client.yaml go run cmd/client/*/main.go
func NewClientConfig() *BaseConfig {
	return NewBaseConfigWithContext(context.Background())
}

// NewBaseConfigWithContext: creates a new BaseConfig instance with configuration loaded from app.yaml or Secrets Manager
func NewBaseConfigWithContext(ctx context.Context) *BaseConfig {
	var config YamlConfig
	var configSource string

	// Determine configuration source
	useSecretsManager := os.Getenv("USE_SECRETSMANAGER") == "true"
	if useSecretsManager {
		configSource = "secretsmanager"
	} else {
		configSource = "localfile"
	}

	// Load configuration based on source
	switch configSource {
	case "secretsmanager":
		secretID, useLocal := GetConfigFromEnv()
		if secretID == "" {
			log.Println("USE_SECRETSMANAGER is true but SECRET_ID is not set, falling back to file-based config")
			// Fall through to localfile case
			configSource = "localfile"
		} else {
			configPtr, err := LoadConfigFromSecretsManager(ctx, secretID, useLocal)
			if err != nil {
				log.Printf("Failed to load config from Secrets Manager: %v, falling back to file-based config", err)
				// Fall through to localfile case
				configSource = "localfile"
			} else {
				log.Println("Successfully loaded configuration from Secrets Manager")
				config = *configPtr
				// Skip to initialization
				goto initializeLogger
			}
		}
		fallthrough

	case "localfile":
		configFilePath := os.Getenv("CONFIG_FILE")
		if configFilePath == "" {
			configFilePath = "etc/app.yaml"
		}

		yamlFile, err := os.Open(configFilePath)
		if err != nil {
			log.Fatalf("Failed to open config file %s: %v", configFilePath, err)
		}
		defer yamlFile.Close()

		byteData, err := io.ReadAll(yamlFile)
		if err != nil {
			log.Fatalf("Failed to read config file %s: %v", configFilePath, err)
		}

		err = yaml.Unmarshal(byteData, &config)
		if err != nil {
			log.Fatalf("Failed to unmarshal YAML from %s: %v", configFilePath, err)
		}
		log.Printf("Successfully loaded configuration from file (%s)", configFilePath)

	default:
		log.Fatalf("Invalid configuration source: %s", configSource)
	}

initializeLogger:
	// Initialize logger config with default values if not configured
	if config.Logger.Component == "" {
		config.Logger.Component = "locky"
	}
	if config.Logger.Service == "" {
		config.Logger.Service = "locky-server"
	}
	if config.Logger.Level == "" {
		config.Logger.Level = "INFO"
	}
	if config.Logger.Output == "" {
		config.Logger.Output = "stdout"
	}

	baseConfig := &BaseConfig{
		YamlConfig:   config,
		DBConnection: nil,
	}

	return baseConfig
}

// NewBaseConfigFromSource: creates a new BaseConfig instance based on CONFIG_SOURCE environment variable
// Valid CONFIG_SOURCE values: "secretsmanager", "localfile" (default)
func NewBaseConfigFromSource(ctx context.Context) *BaseConfig {
	configSource := os.Getenv("CONFIG_SOURCE")

	switch configSource {
	case "secretsmanager":
		log.Println("CONFIG_SOURCE=secretsmanager: Using AWS Secrets Manager for configuration")
		os.Setenv("USE_SECRETSMANAGER", "true")
		return NewBaseConfigWithContext(ctx)
	case "localfile", "":
		if configSource == "" {
			log.Println("CONFIG_SOURCE not set, using local file for configuration (default)")
		} else {
			log.Println("CONFIG_SOURCE=localfile: Using local file for configuration")
		}
		os.Setenv("USE_SECRETSMANAGER", "false")
		return NewBaseConfigWithContext(ctx)
	default:
		log.Fatalf("Invalid CONFIG_SOURCE: %s. Valid values are 'secretsmanager' or 'localfile'", configSource)
		return nil
	}
}

// ConnectDB: connect to MySQL only when needed (safe to call multiple times)
func (bc *BaseConfig) ConnectDB() error {
	if bc.DBConnection != nil {
		return nil
	}
	db := NewDBConnection(bc.YamlConfig)
	if db == nil {
		return fmt.Errorf("failed to connect database")
	}
	bc.DBConnection = db
	return nil
}

func NewDBConnection(conf YamlConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=skip-verify", conf.MySQL.User, conf.MySQL.Pass, conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Db)

	log.Printf("[C-NDBC-1] Attempting database connection to %s:%s/%s", conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Db)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[C-NDBC-3] Failed to connect to database %s:%s/%s: %v", conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Db, err)
		return nil
	}

	log.Printf("[C-NDBC-2] Database connection established to %s:%s/%s", conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Db)
	return db
}
