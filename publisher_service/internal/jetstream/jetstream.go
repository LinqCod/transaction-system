package jetstream

import (
	"github.com/nats-io/nats.go"
	"log"
)

const (
	streamName     = "TRANSACTIONS"
	streamSubjects = "TRANSACTIONS.*"

	SubjectName = "TRANSACTIONS.CREATED"
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

func CreateStream(js nats.JetStreamContext) error {
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
