package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

var (
	// ErrNatsConn is returned when an error occurs in connecting to NATS.
	ErrNatsConn = errors.New("error opening nats connection")
	// ErrNatsJetstream is returned when an error occurs in setting up the NATS Jetstream context.
	ErrNatsJetstream = errors.New("error in nats jetstream")
	// ErrNatsJetstreamAddStream os returned when an attempt to add a NATS Jetstream fails.
	ErrNatsJetstreamAddStream = errors.New("error adding stream to nats jetstream")
)

// NatsJetStream wraps the NATs JetStream connector to implement the StreamBroker interface.
type NatsJetStream struct {
	nats.JetStreamContext
	conn           *nats.Conn
	appName        string
	credsFile      string
	natsURL        string
	streamName     string
	streamPrefix   string
	streamSubjects []string
	basicAuthUser  string
	basicAuthPass  string
	connectTimeout time.Duration
}

// Open connects to the NATS Jetstream.
func (n *NatsJetStream) Open() error {
	opts := []nats.Option{
		nats.Name(n.appName),
		nats.Timeout(n.connectTimeout),
		nats.RetryOnFailedConnect(true),
	}

	if n.basicAuthUser != "" {
		opts = append(opts, nats.UserInfo(n.basicAuthUser, n.basicAuthPass))
	} else {
		opts = append(opts, nats.UserCredentials(n.credsFile))
	}

	conn, err := nats.Connect(n.natsURL, opts...)
	if err != nil {
		return errors.Wrap(ErrNatsConn, err.Error())
	}

	js, err := conn.JetStream()
	if err != nil {
		return errors.Wrap(ErrNatsJetstream, err.Error())
	}

	// TODO: this should be created outside of serverservice,
	// since we have to create consumers as well
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     n.streamName,
		Subjects: n.streamSubjects,
	})
	if err != nil {
		return errors.Wrap(ErrNatsJetstreamAddStream, err.Error())
	}

	n.JetStreamContext = js
	n.conn = conn

	return nil
}

// PublishAsyncWithContext publishes an event onto the NATS Jetstream.
func (n *NatsJetStream) PublishAsyncWithContext(ctx context.Context, objType ObjectType, eventType EventType, data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	msg := newEventStreamMessage(n.appName, eventType, objType)
	msg.AdditionalData = map[string]interface{}{"data": b}
	msg.EventType = string(eventType)

	subject := fmt.Sprintf("%s.%s.%s", n.streamPrefix, objType, eventType)
	if _, err := n.PublishAsync(subject, b); err != nil {
		return err
	}

	return nil
}

// Close closes the NATS Jetstream connection.
func (n *NatsJetStream) Close() {
	if n.conn != nil {
		n.conn.Close()
	}
}
