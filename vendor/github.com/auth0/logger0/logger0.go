package logger0

import "context"

type Publisher interface {
	Publish(ctx context.Context, rec *LogRecord) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, req SubscribeRequest) (<-chan SubscriberEvent, error)
}

type SubscriberEvent struct {
	LogRecord *LogRecord
	Error     error
}

type SubscribeRequest struct {
	Type    LogRecord_Type
	Tenant  string
	Filters map[string]string
}
