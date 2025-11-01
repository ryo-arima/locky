package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	yaml "gopkg.in/yaml.v3"
)

type BaseConfig struct {
	DBConnection *gorm.DB
	YamlConfig   YamlConfig
}

type YamlConfig struct {
	Application Application `yaml:"Application"`
	MySQL       MySQL       `yaml:"MySQL"`
	Redis       Redis       `yaml:"Redis"`
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

type Server struct {
	Admin     Admin  `yaml:"admin"`
	JWTSecret string `yaml:"jwt_secret"`
	LogLevel  string `yaml:"log_level"` // Added: debug / info / warn / error
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

func NewBaseConfig() BaseConfig {
	buf1, err := os.ReadFile("etc/app.yaml")
	if err != nil {
		panic(err)
	}

	var d1 YamlConfig
	err = yaml.Unmarshal(buf1, &d1)
	if err != nil {
		panic(err)
	}

	// Do not connect to DB here, use lazy connection
	baseConfig := &BaseConfig{YamlConfig: d1}
	return *baseConfig
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Printf("DSN (without password): %s:***@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=skip-verify",
			conf.MySQL.User, conf.MySQL.Host, conf.MySQL.Port, conf.MySQL.Db)
		return nil
	}
	log.Println("Successfully connected to database")
	return db
}
