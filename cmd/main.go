package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/linqcod/transaction-system/cmd/api"
	"github.com/linqcod/transaction-system/internal/model"
	"github.com/linqcod/transaction-system/pkg/config"
	"github.com/linqcod/transaction-system/pkg/database"
	_ "github.com/linqcod/transaction-system/pkg/database"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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
			log.Println("transaction-review service")
			log.Printf("Card: %s, Currency: %s, Amount: %f, Status:%s\n", transaction.Card, transaction.Currency, transaction.Amount, transaction.Status)
		}
	}

	return nil
}

//func main() {
//	// Connect to NATS server
//	nc, err := nats.Connect("nats://nats:4222")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
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

	//init routing
	router := api.InitRouter(context.Background(), logger, db)

	// init server
	serverAddr := fmt.Sprintf(":%s", viper.GetString("SERVER_PORT"))
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// graceful shutdown
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Fatalf("error while trying to shutdown http server: %v", err)
		}
		close(stopped)
	}()

	logger.Infof("Starting HTTP server on %s", serverAddr)

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-stopped

	log.Printf("Have a nice day :)")
}
