package main

import (
	"context"
	"encoding/json"
	"github.com/linqcod/transaction-system/consumer_service/internal/model"
	"github.com/linqcod/transaction-system/consumer_service/pkg/config"
	"github.com/linqcod/transaction-system/consumer_service/pkg/database"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"strconv"
	"time"
)

func init() {
	config.LoadConfig(".env")
}

const (
	streamName     = "TRANSACTIONS"
	streamSubjects = "TRANSACTIONS.*"

	subjectName = "TRANSACTIONS.CREATED"
)

func createStream(js nats.JetStreamContext) error {
	stream, err := js.StreamInfo(streamName)
	if err != nil {
		log.Println(err)
	}

	if stream == nil {
		log.Printf("creating stream %q and subjects %q", streamName, streamSubjects)
		_, err := js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func createTransaction(js nats.JetStreamContext) error {
	var transaction model.Transaction

	for i := 1; i <= 10; i++ {
		transaction = model.Transaction{
			Card:     "Card-" + strconv.Itoa(i),
			Currency: "RUB",
			Amount:   float64(i),
			Status:   model.StatusCreated,
		}
		transactionJSON, err := json.Marshal(transaction)
		if err != nil {
			return err
		}

		_, err = js.Publish(subjectName, transactionJSON)
		if err != nil {
			return err
		}
		log.Printf("Transaction for Card: %s has been published\n", transaction.Card)
	}

	return nil
}

func subscribe(js nats.JetStreamContext) error {
	sub, _ := js.PullSubscribe(subjectName, "transaction-review", nats.PullMaxWaiting(128))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		msgs, _ := sub.Fetch(10, nats.Context(ctx))
		for _, msg := range msgs {
			msg.Ack()
			var transaction model.Transaction
			if err := json.Unmarshal(msg.Data, &transaction); err != nil {
				log.Fatal(err)
			}
			log.Println("transaction-review publisher_service")
			log.Printf("Card: %s, Currency: %s, Amount: %f, Status:%s\n", transaction.Card, transaction.Currency, transaction.Amount, transaction.Status)
		}
	}

	return nil
}

//func main() {
//	// Connect to NATS server
//	nc, err := jetstream.Connect("jetstream://jetstream:4222")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	js, err := nc.JetStream(jetstream.PublishAsyncMaxPending(256))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	if err = createStream(js); err != nil {
//		log.Fatal(err)
//	}
//
//	go func() {
//		if err = createTransaction(js); err != nil {
//			log.Fatal(err)
//		}
//	}()
//
//	if err = subscribe(js); err != nil {
//		log.Fatal(err)
//	}
//}

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

	// TODO: subscribe
}
