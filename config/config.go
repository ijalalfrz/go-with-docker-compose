package config

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// Config is an app configuration.
type Config struct {
	Redis struct {
		// Options *redis.Options
	}
	BasicAuth struct {
		Username string
		Password string
	}
	Application struct {
		Port string
		Name string
	}
	Crypto struct {
		Secret string
		IV     string
	}
	Logger struct {
		Formatter logrus.Formatter
	}
	Mariadb struct {
		Driver             string
		Host               string
		Port               string
		Username           string
		Password           string
		Database           string
		DSN                string
		MaxOpenConnections int
		MaxIdleConnections int
	}
}

// Load will load the configuration.
func Load() *Config {
	cfg := new(Config)
	cfg.redis()
	cfg.basicAuth()
	cfg.crypto()
	cfg.logFormatter()
	cfg.app()
	cfg.mariadb()
	return cfg
}

func (cfg *Config) redis() {
	// host := os.Getenv("REDIS_HOST")
	// password := os.Getenv("REDIS_PASSWORD")
	// db, _ := strconv.ParseInt(os.Getenv("REDIS_DATABASE"), 10, 64)

	// options := &redis.Options{
	// 	Addr:     host,
	// 	Password: password,
	// 	DB:       int(db),
	// }

	// cfg.Redis.Options = options
}

func (cfg *Config) basicAuth() {
	username := os.Getenv("BASIC_AUTH_USERNAME")
	password := os.Getenv("BASIC_AUTH_PASSWORD")

	cfg.BasicAuth.Username = username
	cfg.BasicAuth.Password = password
}

func (cfg *Config) crypto() {
	secret := os.Getenv("AES_SECRET")
	iv := os.Getenv("AES_IV")

	cfg.Crypto.IV = iv
	cfg.Crypto.Secret = secret
}

func (cfg *Config) logFormatter() {
	formatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			s := strings.Split(f.Function, ".")
			funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return funcname, filename
		},
	}

	cfg.Logger.Formatter = formatter
}

func (cfg *Config) app() {
	appName := os.Getenv("APP_NAME")
	port := os.Getenv("PORT")

	cfg.Application.Port = port
	cfg.Application.Name = appName
}

func (cfg *Config) mariadb() {
	host := os.Getenv("MARIADB_HOST")
	port := os.Getenv("MARIADB_PORT")
	username := os.Getenv("MARIADB_USERNAME")
	password := os.Getenv("MARIADB_PASSWORD")
	database := os.Getenv("MARIADB_DATABASE")
	maxOpenConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_OPEN_CONNECTIONS"), 10, 64)
	maxIdleConnections, _ := strconv.ParseInt(os.Getenv("MARIADB_MAX_IDLE_CONNECTIONS"), 10, 64)

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	connVal := url.Values{}
	connVal.Add("parseTime", "1")
	connVal.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", dbConnectionString, connVal.Encode())

	cfg.Mariadb.Driver = "mysql"
	cfg.Mariadb.Host = host
	cfg.Mariadb.Port = port
	cfg.Mariadb.Username = username
	cfg.Mariadb.Password = password
	cfg.Mariadb.Database = database
	cfg.Mariadb.DSN = dsn
	cfg.Mariadb.MaxOpenConnections = int(maxOpenConnections)
	cfg.Mariadb.MaxIdleConnections = int(maxIdleConnections)
}
