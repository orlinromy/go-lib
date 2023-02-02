package rmq

import "time"

type PublisherConfig struct {
	Enabled               bool          `json:"enabled" mapstructure:"enabled"`
	Name                  string        `json:"name" mapstructure:"name"`
	RoutingKey            string        `json:"routing_key" mapstructure:"routing_key"`
	Mandatory             bool          `json:"mandatory" mapstructure:"mandatory"`
	Immediate             bool          `json:"immediate" mapstructure:"immediate"`
	AutoGenerateMessageID bool          `json:"auto_generate_message_id" mapstructure:"auto_generate_message_id"`
	PublisherConfirmed    bool          `json:"publisher_confirmed" mapstructure:"publisher_confirmed"`
	Timeout               time.Duration `json:"timeout" mapstructure:"timeout"` // second
	NoWait                bool          `json:"no_wait" mapstructure:"no_wait"`
}

type ConsumerConfig struct {
	Queue     string                 `json:"queue" mapstructure:"queue"`
	Enabled   bool                   `json:"enabled" mapstructure:"enabled"`
	Name      string                 `json:"name" mapstructure:"name"`
	AutoAck   bool                   `json:"auto_ack" mapstructure:"auto_ack"`
	Exclusive bool                   `json:"exclusive" mapstructure:"exclusive"`
	NoLocal   bool                   `json:"no_local" mapstructure:"no_local"`
	NoWait    bool                   `json:"no_wait" mapstructure:"no_wait"`
	Args      map[string]interface{} `json:"args" mapstructure:"args"`

	// Fair dispatch
	EnabledPrefetch bool `json:"enabled_prefetch" mapstructure:"enabled_prefetch"`
	PrefetchCount   int  `json:"prefetch_count" mapstructure:"prefetch_count"`
	PrefetchSize    int  `json:"prefetch_size" mapstructure:"prefetch_size"`
	Global          bool `json:"global" mapstructure:"global"`
}

type MessageRetryConfig struct {
	// retry
	Enabled           bool `json:"enabled" mapstructure:"enabled"`
	HandleDeadMessage bool `json:"handle_dead_message" mapstructure:"handle_dead_message"`
	RetryCountLimit   int  `json:"retry_count_limit" mapstructure:"retry_count_limit"`
}

type QueueManagerConfig struct {
	Enabled         bool      `json:"enabled" mapstructure:"enabled"`
	ConnURIs        []string  `json:"conn_uris" mapstructure:"conn_uris"`
	VirtualHost     string    `json:"virtual_host" mapstructure:"virtual_host"`
	AutoReconnect   bool      `json:"auto_reconnect" mapstructure:"auto_reconnect"`
	EnablePublisher bool      `json:"enable_publisher" mapstructure:"enable_publisher"`
	EnableConsumer  bool      `json:"enable_consumer" mapstructure:"enable_consumer"`
	Reconnect       Reconnect `json:"reconnect" mapstructure:"reconnect"`
}

type QueueConfig struct {
	Name       string                 `json:"name" mapstructure:"name"`
	Durable    bool                   `json:"durable" mapstructure:"durable"`
	AutoDelete bool                   `json:"auto_delete" mapstructure:"auto_delete"`
	Exclusive  bool                   `json:"exclusive" mapstructure:"exclusive"`
	NoWait     bool                   `json:"no_wait" mapstructure:"no_wait"`
	Args       map[string]interface{} `json:"args" mapstructure:"args"`
}

type QueueBindConfig struct {
	Queue      string                 `json:"queue" mapstructure:"queue"`
	Exchange   string                 `json:"exchange" mapstructure:"exchange"`
	BindingKey string                 `json:"binding_key" mapstructure:"binding_key"`
	NoWait     bool                   `json:"no_wait" mapstructure:"no_wait"`
	Args       map[string]interface{} `json:"args" mapstructure:"args"`
}

type ExchangeConfig struct {
	Exchange     string                 `json:"exchange" mapstructure:"exchange"`
	ExchangeType string                 `json:"exchange_type" mapstructure:"exchange_type"`
	Durable      bool                   `json:"durable" mapstructure:"durable"`
	AutoDelete   bool                   `json:"auto_delete" mapstructure:"auto_delete"`
	Exclusive    bool                   `json:"exclusive" mapstructure:"exclusive"`
	NoWait       bool                   `json:"no_wait" mapstructure:"no_wait"`
	Args         map[string]interface{} `json:"args" mapstructure:"args"`
}

type Reconnect struct {
	Interval   time.Duration `json:"interval" mapstructure:"interval"` // in milli sec
	MaxAttempt int           `json:"max_attempt" mapstructure:"max_attempt"`
}
