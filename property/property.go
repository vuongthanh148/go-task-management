package property

import (
	"context"
	"fmt"
	"kn-assignment/internal/log"
	"time"

	"cloud.google.com/go/profiler"
	"github.com/kelseyhightower/envconfig"
)

func Init(ctx context.Context) {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf(ctx, "read env error : %s", err.Error())
	}
	setPostgresConnString()
}

func setPostgresConnString() {
	postgres := cfg.Postgres
	cfg.PostgresConfig.ConnString = fmt.Sprintf(cfg.PostgresConfig.ConnUri, postgres.Host, postgres.Port, postgres.Database, postgres.User, postgres.Password)
}

func Get() config {
	return cfg
}

var cfg config

type config struct {
	Server         serviceProperties
	Gin            gin
	ProfilerConfig profilerConfig
	Secret         secretConfig
	Postgres       postgres
	PostgresConfig PostgresConfig
}

type serviceProperties struct {
	ServerProperties
	PostgresConnectionURI  string `envconfig:"POSTGRES_CONNECTION_URI"`
	PostgresProjectId      string `envconfig:"POSTGRES_PROJECT_ID"`
	PostgresPasswordSecret string `envconfig:"POSTGRES_PASSWORD_SECRET"`
	DataServerAuthHost     string `envconfig:"DATA_SERVER_AUTH_HOST"`
	DataServerAuthUsername string `envconfig:"DATA_SERVER_AUTH_USERNAME"`
	DataServerAuthPassword string `envconfig:"DATA_SERVER_AUTH_PASSWORD"`
	DataServerHost         string `envconfig:"DATA_SERVER_HOST"`
	DataServerApiKey       string `envconfig:"DATA_SERVER_API_KEY"`

	DataProjectId string `envconfig:"DATA_PROJECT_ID"`
	DataSecret    string `envconfig:"DATA_SECRET"`

	MaxGoroutineDB int `envconfig:"LIMIT_GOROUTINE_DB_CONNECT" default:"100"`
}
type gin struct {
	Mode string `envconfig:"GIN_MODE" default:"debug"`
}

type profilerConfig struct {
	Cfg profiler.Config
}

type postgres struct {
	User     string `envconfig:"POSTGRES_USER" default:"user"`
	Host     string `envconfig:"POSTGRES_HOST" default:"localhost"`
	Port     string `envconfig:"POSTGRES_PORT" default:"5432"`
	Password string `envconfig:"POSTGRES_PASSWORD" default:"password"`
	Database string `envconfig:"POSTGRES_DATABASE" default:"taskdb"`
}

type secretConfig struct {
	PostgresPasswordSecret string `envconfig:"POSTGRES_PASSWORD_SECRET"`
	JWTSecretKey           string `envconfig:"JWT_SECRET_KEY"`
}

type PostgresConfig struct {
	// ConnUri: "host=localhost port=5430 database=profile user=postgres password=xxx"
	ConnUri    string `envconfig:"POSTGRES_CONN_URI" env:"POSTGRES_CONN_URI" default:"host=%s port=%s database=%s user=%s password=%s"`
	ConnString string
	// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	MaxConnLifetime time.Duration `envconfig:"POSTGRES_MAX_CONN_LIFETIME" env:"POSTGRES_MAX_CONN_LIFETIME" default:"1h"`

	// MaxConnLifetimeJitter is the duration after MaxConnLifetime to randomly decide to close a connection.
	// This helps prevent all connections from being closed at the exact same time, starving the pool.
	// MaxConnLifetimeJitter time.Duration `envconfig:"POSTGRES_MAX_CONN_LIFETIME_JITTER" env:"POSTGRES_MAX_CONN_LIFETIME_JITTER"`

	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	MaxConnIdleTime time.Duration `envconfig:"POSTGRES_MAX_CONN_IDLE_TIME" env:"POSTGRES_MAX_CONN_IDLE_TIME" default:"30m"`

	// MaxConns is the maximum size of the pool. The default is the greater of 4 or runtime.NumCPU().
	MaxConns int32 `envconfig:"POSTGRES_MAX_CONNS" env:"POSTGRES_MAX_CONNS" default:"4"`

	// MinConns is the minimum size of the pool. After connection closes, the pool might dip below MinConns. A low
	// number of MinConns might mean the pool is empty after MaxConnLifetime until the health check has a chance
	// to create new connections.
	MinConns int32 `envconfig:"POSTGRES_MIN_CONNS" env:"POSTGRES_MIN_CONNS" default:"0"`
}

type ServerProperties struct {
	DebugMode            bool   `envconfig:"DEBUG_MODE" long:"debug-mode" description:"turn on/off debug mode (default: false)" env:"DEBUG_MODE"`
	PrintConsoleFormat   bool   `envconfig:"CONSOLE_FORMAT" long:"print-console-format" description:"log to print console format or not (default: false)" env:"CONSOLE_FORMAT"`
	ShutdownTimeout      int64  `envconfig:"SHUTDOWN_TIMEOUT" long:"shutdown-timeout" description:"graceful shutdown timeout" env:"SHUTDOWN_TIMEOUT" default:"300"`
	Port                 string `envconfig:"PORT" long:"port" description:"server running port" env:"PORT" `
	ProjectID            string `envconfig:"GOOGLE_CLOUD_PROJECT" long:"project-id" description:"Google project id" env:"GOOGLE_CLOUD_PROJECT"`
	ServiceName          string `envconfig:"SERVICE_NAME" long:"service-name" description:"Service name" env:"SERVICE_NAME"`
	ServiceDescription   string `envconfig:"SERVICE_DESCRIPTION" long:"service-description" description:"Service description" env:"SERVICE_DESCRIPTION" default:""`
	RunLocal             bool   `envconfig:"RUN_LOCAL" long:"run-local" description:"Is service running on local (default: false)" env:"RUN_LOCAL"`
	LogIgnorePaths       string `envconfig:"LOG_IGNORE_PATHS" long:"log-ignore-paths" description:"url path to ignore logging (full path without host)" env:"LOG_IGNORE_PATHS"`
	ApiDocs              bool   `envconfig:"API_DOCS" long:"api-docs" description:"expose api docs url (default: false)" env:"API_DOCS"`
	ApiDocsSchema        string `envconfig:"API_DOCS_SCHEMA" long:"api-docs-schema" description:"Api docs schema" env:"API_DOCS_SCHEMA" default:"http"`
	ApiDocsVersion       string `envconfig:"API_DOCS_VERSION" long:"api-docs-version" description:"Api docs version" env:"API_DOCS_VERSION" default:"v0.0.1"`
	LogClientIgnorePaths string `envconfig:"LOG_CLIENT_IGNORE_PATHS" long:"log-client-ignore-paths" description:"url path to ignore client logging (full path without host)" env:"LOG_CLIENT_IGNORE_PATHS"`
	Host                 string `envconfig:"HOST" long:"host" description:"Host" env:"HOST" default:"localhost"`
	GinMode              string `envconfig:"GIN_MODE" long:"gin-mode" description:"Gin mode" env:"GIN_MODE"`
	ClientLogMasking     bool   `envconfig:"CLIENT_LOG_MASKING" long:"client-log-masking" description:"Client log masking" env:"CLIENT_LOG_MASKING"`

	AccessTokenExpiry  time.Duration `envconfig:"ACCESS_TOKEN_TIME" long:"access-token-time" description:"Access token expiry time" env:"ACCESS_TOKEN_TIME" default:"15m"`
	RefreshTokenExpiry time.Duration `envconfig:"REFRESH_TOKEN_TIME" long:"refresh-token-time" description:"Refresh token expiry time" env:"REFRESH_TOKEN_TIME" default:"168h"`
}
