package main

import (
	"github.com/linqcod/transaction-system/consumer_service/internal/jetstream"
	"github.com/linqcod/transaction-system/consumer_service/pkg/config"
	"github.com/linqcod/transaction-system/consumer_service/pkg/database"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func init() {
	config.LoadConfig(".env")
}

func main() {
	// init zap logger
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	baseLogger, err := loggerConfig.Build()
	if err != nil {
		log.Fatalf("error while building zap logger: %v", err)
	}

	logger := baseLogger.Sugar()

	// init db connection
	db, err := database.InitDB()
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	// nats server connection
	js, err := jetstream.Connect()
	if err != nil {
		logger.Fatalf("error while connecting to nats server: %v", err)
	}

	if err = jetstream.Subscribe(js); err != nil {
		logger.Fatalf("error while subscribing to subject: %v", err)
	}
}
