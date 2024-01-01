package lib

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(cfg *Config) *zap.Logger {
	dev, err := cfg.Get(Development)
	if err != nil {
		dev = "false"
	}

	var logger *zap.Logger
	if dev == "true" {
		logger, err = zap.NewDevelopment()

	} else {
		conf := zap.NewProductionConfig()
		conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		logger, err = conf.Build()
	}
	if err != nil {
		log.Fatalf("Error creating Logger: %s", err.Error())
	}

	return logger
}
