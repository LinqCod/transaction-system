package jetstream

import (
	"encoding/json"
	"github.com/linqcod/transaction-system/consumer_service/internal/model"
	"github.com/nats-io/nats.go"
	"log"
)

const (
	subjectName = "TRANSACTIONS.CREATED"
)

func Connect() (nats.JetStreamContext, error) {
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return js, nil
}

func Subscribe(js nats.JetStreamContext) error {
	if _, err := js.Subscribe(subjectName, func(msg *nats.Msg) {
		msg.Ack()

		var transaction model.Transaction
		err := json.Unmarshal(msg.Data, &transaction)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Got: Card: %s, Currency: %s, Amount: %f, Type: %s, Status:%s", transaction.CardNumber, transaction.Currency, transaction.Amount, transaction.Type, transaction.Status)
	}); err != nil {
		return err
	}

	return nil
}
