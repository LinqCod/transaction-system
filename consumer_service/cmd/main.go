package main

import (
	"context"
	"encoding/json"
	"github.com/linqcod/transaction-system/consumer_service/internal/model"
	"github.com/linqcod/transaction-system/consumer_service/internal/repository"
	"github.com/linqcod/transaction-system/consumer_service/pkg/config"
	"github.com/linqcod/transaction-system/consumer_service/pkg/currencyapi"
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
			if err != nil {
				transaction.Status = model.StatusError
				logTransaction(transaction, logger)
				logger.Errorf("error while getting account by card number: %v", err)
				continue
			}

			currencyCoefficient, err := currencyapi.ConvertCurrencies(transaction.Currency, "RUB")
			if err != nil {
				transaction.Status = model.StatusError
				logTransaction(transaction, logger)
				logger.Errorf("error while converting currencies: %v", err)
				continue
			}

			transaction.Amount *= currencyCoefficient

			if transaction.Type == model.InvoiceTransactionType {
				newFrozenBalance := account.FrozenBalance - transaction.Amount

				if newFrozenBalance < 0 {
					transaction.Status = model.StatusError
					logTransaction(transaction, logger)
					logger.Errorln("error: account frozen balance cannot be less than zero")
					continue
				}

				if err = accountRepository.UpdateAccountFrozenBalance(transaction.CardNumber, newFrozenBalance); err != nil {
					transaction.Status = model.StatusError
					logTransaction(transaction, logger)
					logger.Errorln("error while updating account frozen balance: %v", err)
					continue
				}
			}

			if transaction.Type == model.WithdrawTransactionType {
				transaction.Amount *= -1
			}

			newBalance := account.Balance + transaction.Amount

			if newBalance < 0 {
				transaction.Status = model.StatusError
				logTransaction(transaction, logger)
				logger.Errorln("error: account balance cannot be less than zero")
				continue
			}

			if err = accountRepository.UpdateAccountBalance(transaction.CardNumber, newBalance); err != nil {
				transaction.Status = model.StatusError
				logTransaction(transaction, logger)
				logger.Errorln("error while updating account balance: %v", err)
				continue
			}

			transaction.Status = model.StatusSuccess
			logTransaction(transaction, logger)
		}
	}
}

func logTransaction(transaction model.Transaction, logger *zap.SugaredLogger) {
	logger.Infof("Transaction: Card: %s, Currency: %s, Amount: %f, Type: %s, Status:%s",
		transaction.CardNumber,
		transaction.Currency,
		transaction.Amount,
		transaction.Type,
		transaction.Status,
	)
}
