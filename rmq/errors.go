package rmq

type queueError struct {
	Code    int
	Message string
}

func (e *queueError) Error() string {
	return e.Message
}

const (
	QueueManagerDisabled   int = 001
	ConnIsNotOpened        int = 002
	InvalidExchangeConfig  int = 003
	InvalidPublisherConfig int = 004
	InvalidPublishArgs     int = 005
	IdGenerationFailed     int = 006
	InvalidQueueConfig     int = 005
	InvalidQueueBindConfig int = 006
	InvalidConsumerConfig  int = 006
	RetryCountExceeded     int = 007
)

var (
	queueManagerIsDisabledError = &queueError{
		Code:    QueueManagerDisabled,
		Message: "queue manager is disabled.",
	}
	connIsNotOpened = &queueError{
		Code:    ConnIsNotOpened,
		Message: "connection is not open",
	}
	invalidExchangeConfig = &queueError{
		Code:    InvalidExchangeConfig,
		Message: "missing or invalid exchange configuration",
	}
	invalidPublisherConfig = &queueError{
		Code:    InvalidPublisherConfig,
		Message: "missing or invalid publisher configuration",
	}
	invalidPublishArgs = &queueError{
		Code:    InvalidPublishArgs,
		Message: "invalid publish args",
	}
	IDGenerationFailed = &queueError{
		Code:    IdGenerationFailed,
		Message: "failed to generate message ID",
	}
	invalidQueueConfig = &queueError{
		Code:    InvalidQueueConfig,
		Message: "missing or invalid queue configuration",
	}
	invalidQueueBindConfig = &queueError{
		Code:    InvalidQueueBindConfig,
		Message: "missing or invalid queue bind configuration",
	}
	invalidConsumerConfig = &queueError{
		Code:    InvalidConsumerConfig,
		Message: "missing or invalid consumer configuration",
	}
	retryCountExceeded = &queueError{
		Code:    RetryCountExceeded,
		Message: "retry count exceeded for msg",
	}
)

func IsQueueManagerDisabled(err error) bool {
	queueErr, ok := err.(*queueError)
	if !ok {
		return ok
	}
	return queueErr.Code == QueueManagerDisabled
}
