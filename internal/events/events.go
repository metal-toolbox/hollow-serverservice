// Package events provides types and methods to interact with a messaging stream broker.
package events

import (
	"context"
	"fmt"
	"time"

	"go.infratographer.com/x/pubsubx"
)

type (
	// ObjectType is the kind of the serverservice object included in an event.
	ObjectType string

	// EventType is a type identifying the Event kind that has occurred on a serverservice objects.
	EventType string
)

const (
	// URN is the uniform resource name, the last two fields are the ObjectType, ActionKind and the
	URN = "urn:hollow:serverservice:%s:%s/"

	// Create action kind identifies objects that were created.
	Create EventType = "create"

	// Update action kind identifies objects that were updated.
	Update EventType = "update"

	// Delete action kind identifies objects that were removed.
	Delete EventType = "delete"
)

// StreamBroker provides methods to interact with the event stream.
type StreamBroker interface {
	Open() error
	// PublishWithContext publishes the message to the message broker in an async manner.
	PublishAsyncWithContext(ctx context.Context, objType ObjectType, eventType EventType, data interface{}) error
	Close()
}

// NewStreamBroker returns a StreamBorker instance.
func NewStreamBroker(
	appName,
	credsFile,
	natsURL,
	streamName,
	streamPrefix string,
	streamSubjects []string,
	basicAuthUser,
	basicAuthPass string,
	connectTimeout time.Duration,
) StreamBroker {
	return &NatsJetStream{
		appName:        appName,
		credsFile:      credsFile,
		natsURL:        natsURL,
		streamName:     streamName,
		streamPrefix:   streamPrefix,
		streamSubjects: streamSubjects,
		basicAuthUser:  basicAuthUser,
		basicAuthPass:  basicAuthPass,
		connectTimeout: connectTimeout,
	}
}

func newURN(eventType EventType, objType ObjectType) string {
	return fmt.Sprintf(URN, eventType, objType)
}

func newEventStreamMessage(appName string, eventType EventType, objType ObjectType) *pubsubx.Message {
	return &pubsubx.Message{
		EventType:  string(eventType),
		SubjectURN: newURN(eventType, objType),
		Timestamp:  time.Now().UTC(),
		Source:     appName,
	}
}
