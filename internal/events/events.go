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
	// urnFormatString is the format for the uniform resource name.
	//
	// The string is to be formatted as "urn:<namespace>:<ObjectType>:<object UUID>"
	urnFormatString = "urn:%s:%s:%s"

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
	PublishAsyncWithContext(ctx context.Context, objType ObjectType, eventType EventType, objID string, obj interface{}) error
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
	streamUrnNs,
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
		streamUrnNs:    streamUrnNs,
		basicAuthUser:  basicAuthUser,
		basicAuthPass:  basicAuthPass,
		connectTimeout: connectTimeout,
	}
}

func newURN(namespace string, objType ObjectType, objID string) string {
	return fmt.Sprintf(urnFormatString, namespace, objType, objID)
}

func newEventStreamMessage(appName, urnNamespace string, eventType EventType, objType ObjectType, objID string) *pubsubx.Message {
	return &pubsubx.Message{
		EventType:  string(eventType),
		ActorURN:   "", // To be filled in with the data from the client request JWT.
		SubjectURN: newURN(urnNamespace, objType, objID),
		Timestamp:  time.Now().UTC(),
		Source:     appName,
	}
}
