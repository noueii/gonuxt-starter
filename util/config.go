package util

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	DbDriver             string
	DbUser               string
	DbPassword           string
	DbHost               string
	DbPort               string
	DbName               string
	DbSSLConnection      bool
	DbURL                string
	DBMigrationsLocation string
	HTTPHost             string
	HTTPPort             string
	HTTPAddr             string
	GRPCHost             string
	GRPCPort             string
	GRPCAddr             string
	TokenSymmetricKey    string
	TokenAccessDuration  time.Duration
	TokenRefreshDuration time.Duration
	OAuth                OAuth
	CORSAllowedOrigins   []string
	Environment          ENV
}

type OAuth struct {
	Google oauth2.Config
}

type ENV string

const (
	Production  ENV = "prod"
	Development ENV = "dev"
	Test        ENV = "test"
)

func LoadConfig(env ENV, fp ...string) (*Config, error) {
	path := ""
	if len(fp) > 0 {
		path = fp[0]
	}
	filePath := fmt.Sprintf("%s%s.env", path, env)
	if env == Test {
		filePath = "../" + filePath
		fmt.Println(filePath)
	}

	envars, err := godotenv.Read(filePath)

	if err != nil {
		return nil, err
	}

	dbDriver, ok := envars["DB_DRIVER"]
	if !ok || len(dbDriver) == 0 {
		return nil, fmt.Errorf("%s environment variable 'DB_DRIVER' not found", env)
	}

	dbUser, ok := envars["DB_USER"]
	if !ok || len(dbUser) == 0 {
		return nil, fmt.Errorf("%s environment variable 'DB_USER' not found", env)
	}

	dbPassword, ok := envars["DB_PASSWORD"]
	if !ok || len(dbPassword) == 0 {
		return nil, fmt.Errorf("%s environment variable 'DB_PASSWORD' not found", env)
	}

	dbHost, ok := envars["DB_HOST"]
	if !ok || len(dbHost) == 0 {
		return nil, fmt.Errorf("%s environment variable 'DB_HOST' not found", env)
	}

	dbPort, ok := envars["DB_PORT"]
	if !ok || len(dbPort) == 0 {
		return nil, fmt.Errorf("%s environment variable 'DB_PORT' not found", env)
	}

	dbName, ok := envars["DB_NAME"]
	if !ok || len(dbName) == 0 {
		return nil, fmt.Errorf("%s environment variable 'DB_NAME' not found", env)
	}

	dbSSLConnectionString, ok := envars["DB_SSL_ENABLE"]
	if !ok || len(dbSSLConnectionString) == 0 {
		return nil, fmt.Errorf("%s environment variable 'DB_SSL_ENABLE' not found", env)
	}

	dbSSLConnection, err := strconv.ParseBool(dbSSLConnectionString)

	dbURL := fmt.Sprintf("%s://%s:%s@%s:%s/%s", dbDriver, dbUser, dbPassword, dbHost, dbPort, dbName)

	if err != nil {
		return nil, err
	}

	if !dbSSLConnection {
		dbURL += "?sslmode=disable"
	}

	dbMigrationsLocation, ok := envars["DB_MIGRATIONS_LOCATION"]
	if !ok || len(dbMigrationsLocation) == 0 {
		return nil, fmt.Errorf("%s environment valiable 'DB_MIGRATIONS_LOCATION' not found", env)
	}

	httpHost, ok := envars["HTTP_HOST"]
	if !ok || len(httpHost) == 0 {
		return nil, fmt.Errorf("%s environment variable 'HTTP_HOST' not found", env)
	}

	httpPort, ok := envars["HTTP_PORT"]
	if !ok || len(httpPort) == 0 {
		return nil, fmt.Errorf("%s environment variable 'HTTP_PORT' not found", env)
	}

	httpAddr := fmt.Sprintf("%s:%s", httpHost, httpPort)

	tokenSymmetricKey, ok := envars["TOKEN_SYMMETRIC_KEY"]
	if !ok {
		return nil, fmt.Errorf("%s environment variable 'TOKEN_SYMMETRIC_KEY' not found", env)
	}

	if len(tokenSymmetricKey) != 32 {
		return nil, fmt.Errorf("%s environment variable 'TOKEN_SYMMETRIC_KEY' has invalid length", env)
	}

	tokenAccessDurationString, ok := envars["TOKEN_ACCESS_DURATION"]
	if !ok || len(tokenAccessDurationString) == 0 {
		return nil, fmt.Errorf("%s environment variable 'TOKEN_ACCESS_DURATION' not found", env)
	}

	tokenAccessDuration, err := time.ParseDuration(tokenAccessDurationString)
	if err != nil {
		return nil, err
	}

	tokenRefreshDurationString, ok := envars["TOKEN_REFRESH_DURATION"]
	if !ok || len(tokenRefreshDurationString) == 0 {
		return nil, fmt.Errorf("%s environment variable 'TOKEN_REFRESH_DURATION' not found", env)
	}

	tokenRefreshDuration, err := time.ParseDuration(tokenRefreshDurationString)
	if err != nil {
		return nil, err
	}

	grpcHost, ok := envars["GRPC_HOST"]
	if !ok || len(grpcHost) == 0 {
		return nil, fmt.Errorf("%s environment variable 'GRPC_HOST' not found", env)
	}

	grpcPort, ok := envars["GRPC_PORT"]
	if !ok || len(grpcPort) == 0 {
		return nil, fmt.Errorf("%s environment variable 'GRPC_PORT' not found", env)
	}

	grpcAddr := fmt.Sprintf("%s:%s", grpcHost, grpcPort)

	corsAllowedOriginsString, ok := envars["CORS_ORIGINS"]

	if !ok || len(corsAllowedOriginsString) == 0 {
		return nil, fmt.Errorf("%s environment variable 'CORS_ORIGINS' not found", env)
	}

	corsAllowedOrigins := strings.Split(corsAllowedOriginsString, ",")

	googleOAuthClientID, ok := envars["GOOGLE_OAUTH_CLIENT_ID"]

	if !ok || len(googleOAuthClientID) == 0 {
		return nil, fmt.Errorf("%s environment variable 'GOOGLE_OAUTH_CLIENT_ID' not found", env)
	}

	googleOAuthClientSecret, ok := envars["GOOGLE_OAUTH_CLIENT_SECRET"]

	if !ok || len(googleOAuthClientSecret) == 0 {
		return nil, fmt.Errorf("%s environment variable 'GOOGLE_OAUTH_CLIENT_SECRET' not found", env)
	}

	googleOauthConfig := &oauth2.Config{
		ClientID:     googleOAuthClientID,
		ClientSecret: googleOAuthClientSecret,
		RedirectURL:  "http://localhost:8080/v1/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	oAuth := OAuth{
		Google: *googleOauthConfig,
	}

	return &Config{
		DbDriver:             dbDriver,
		DbUser:               dbUser,
		DbPassword:           dbPassword,
		DbHost:               dbHost,
		DbPort:               dbPort,
		DbName:               dbName,
		DbSSLConnection:      dbSSLConnection,
		DbURL:                dbURL,
		DBMigrationsLocation: dbMigrationsLocation,
		HTTPHost:             httpHost,
		HTTPPort:             httpPort,
		HTTPAddr:             httpAddr,
		GRPCHost:             grpcHost,
		GRPCPort:             grpcPort,
		GRPCAddr:             grpcAddr,
		TokenSymmetricKey:    tokenSymmetricKey,
		TokenAccessDuration:  tokenAccessDuration,
		TokenRefreshDuration: tokenRefreshDuration,
		CORSAllowedOrigins:   corsAllowedOrigins,
		Environment:          env,
		OAuth:                oAuth,
	}, nil
}

func LoadEnv() (ENV, error) {
	envs, err := godotenv.Read("app.env")

	if err != nil {
		return "", err
	}

	env, ok := envs["ENVIRONMENT"]
	if !ok {
		return "", fmt.Errorf("environment variable 'ENVIRONMENT' not found")
	}

	if env != string(Production) && env != string(Development) && env != string(Test) {
		return "", fmt.Errorf("invalid environment variable 'ENVIRONMENT'")
	}

	return ENV(env), nil
}
