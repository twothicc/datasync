package env

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/twothicc/common-go/logger"
	"go.uber.org/zap"
)

type envConfigs struct {
	Port string
	Env  string
}

var EnvConfigs = &envConfigs{}

func Init(ctx context.Context) {
	err := godotenv.Load()
	if err != nil {
		logger.WithContext(ctx).Error("fail to init env", zap.Error(err))
	}

	EnvConfigs.Port = os.Getenv(PORT)
	EnvConfigs.Env = os.Getenv(ENV)
}

// IsTest - Indicates if environment is test or production
func IsTest() bool {
	if EnvConfigs.Env == TEST {
		return true
	}

	return false
}
