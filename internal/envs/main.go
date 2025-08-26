package envs

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type EnvironmentType string

const (
	DEVELOP    EnvironmentType = "develop"
	STAGE      EnvironmentType = "stage"
	PRODUCTION EnvironmentType = "production"
)

type key int

const (
	APP_ENV = iota + 1
	APP_URL
	SERVER_NAME
	SERVER_PORT
	POSTGRES_USERNAME
	POSTGRES_PASSWORD
	POSTGRES_HOST
	POSTGRES_PORT
	POSTGRES_DATABASE
	ACCESS_CLIENT_ID
	MAIL_HOST
	MAIL_PORT
	MAIL_USERNAME
	MAIL_PASSWORD
	RABBITMQ_USERNAME
	RABBITMQ_PASSWORD
	RABBITMQ_HOST
	RABBITMQ_VHOST
	RABBITMQ_PORT
)

func (k key) String() string {
	switch k {
	case APP_ENV:
		return "APP_ENV"
	case APP_URL:
		return "APP_URL"
	case SERVER_NAME:
		return "SERVER_NAME"
	case SERVER_PORT:
		return "SERVER_PORT"
	case POSTGRES_USERNAME:
		return "POSTGRES_USERNAME"
	case POSTGRES_PASSWORD:
		return "POSTGRES_PASSWORD"
	case POSTGRES_HOST:
		return "POSTGRES_HOST"
	case POSTGRES_PORT:
		return "POSTGRES_PORT"
	case POSTGRES_DATABASE:
		return "POSTGRES_DATABASE"
	case ACCESS_CLIENT_ID:
		return "ACCESS_CLIENT_ID"
	case MAIL_HOST:
		return "MAIL_HOST"
	case MAIL_PORT:
		return "MAIL_PORT"
	case MAIL_USERNAME:
		return "MAIL_USERNAME"
	case MAIL_PASSWORD:
		return "MAIL_PASSWORD"
	case RABBITMQ_USERNAME:
		return "RABBITMQ_USERNAME"
	case RABBITMQ_PASSWORD:
		return "RABBITMQ_PASSWORD"
	case RABBITMQ_HOST:
		return "RABBITMQ_HOST"
	case RABBITMQ_VHOST:
		return "RABBITMQ_VHOST"
	case RABBITMQ_PORT:
		return "RABBITMQ_PORT"
	default:
		return "Unknown"
	}
}

func Get(key key) (value string) {
	value = os.Getenv(key.String())
	return value
}

func GetInt(key key) (value int) {
	strValue := os.Getenv(key.String())
	value, _ = strconv.Atoi(strValue)
	return
}

func checkEnvVariables() error {
	missingVars := []string{}

	for _, k := range []key{
		APP_ENV,
		APP_URL,
		SERVER_NAME,
		SERVER_PORT,
		POSTGRES_USERNAME,
		POSTGRES_PASSWORD,
		POSTGRES_HOST,
		POSTGRES_PORT,
		POSTGRES_DATABASE,
		ACCESS_CLIENT_ID,
		MAIL_HOST,
		MAIL_PORT,
		MAIL_USERNAME,
		MAIL_PASSWORD,
		RABBITMQ_USERNAME,
		RABBITMQ_PASSWORD,
		RABBITMQ_HOST,
		RABBITMQ_VHOST,
		RABBITMQ_PORT,
	} {
		if os.Getenv(k.String()) == "" {
			missingVars = append(missingVars, k.String())
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("the following environment variables are missing: %v", missingVars)
	}

	return nil
}

func IsDev() bool {
	return os.Getenv(key.String(APP_ENV)) == string(DEVELOP)
}

func IsStage() bool {
	return os.Getenv(key.String(APP_ENV)) == string(STAGE)
}

func IsProd() bool {
	return os.Getenv(key.String(APP_ENV)) == string(PRODUCTION)
}

func Load() {
	_ = godotenv.Load()

	err := checkEnvVariables()

	if err != nil {
		panic(err)
	}
}
