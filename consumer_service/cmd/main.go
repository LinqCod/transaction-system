package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/linqcod/transaction-system/consumer_service/internal/model"
	"github.com/linqcod/transaction-system/consumer_service/internal/repository"
	"github.com/linqcod/transaction-system/consumer_service/pkg/config"
	"github.com/linqcod/transaction-system/consumer_service/pkg/database"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

const (
	subjectName = "TRANSACTIONS.CREATED"
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

	accountRepository := repository.NewAccountRepository(context.Background(), db)

	// nats jetstream connection
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		logger.Fatalf("error while connecting to nats jetstream: %v", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		logger.Fatalf("error while creating jetsream: %v", err)
	}

	sub, _ := js.PullSubscribe(subjectName, "transaction-processing", nats.PullMaxWaiting(100))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msgs, _ := sub.Fetch(10, nats.Context(ctx))
		for _, msg := range msgs {
			msg.Ack()

			var transaction model.Transaction
			err = json.Unmarshal(msg.Data, &transaction)
			if err != nil {
				logger.Fatal(err)
			}

			account, err := accountRepository.GetAccount(transaction.CardNumber)
			if errors.Is(err, sql.ErrNoRows) {
				transaction.Status = model.StatusError
				//TODO: error because of nil account
			} else if err != nil {
				logger.Errorf("error while getting account by card number: %v", err)
				//TODO: how
				continue
			}

			switch transaction.Type {
			case model.InvoiceTransactionType:
				newBalance := account.Balance + transaction.Amount
				if err = accountRepository.UpdateAccountBalance(transaction.CardNumber, newBalance); err != nil {
					logger.Errorf("error while updating account balance: %v", err)
					continue
				}
				transaction.Status = model.StatusSuccess
			case model.WithdrawTransactionType:
				newBalance := account.Balance - transaction.Amount
				if newBalance < 0 {
					transaction.Status = model.StatusError
				} else if err = accountRepository.UpdateAccountBalance(transaction.CardNumber, newBalance); err != nil {
					logger.Errorf("error while updating account balance: %v", err)
					continue
				} else {
					transaction.Status = model.StatusSuccess
				}
			}

			logger.Infof("Final: Card: %s, Currency: %s, Amount: %f, Type: %s, Status:%s", transaction.CardNumber, transaction.Currency, transaction.Amount, transaction.Type, transaction.Status)
		}
	}
}
