// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// IMicroservice is interface for centralized service management
type IMicroservice interface {
	Start() error
	Stop()
	Cleanup() error

	// Consumer Services
	Consume(servers string, topic string, groupID string, h ServiceHandleFunc) error
	// ConsumeBatch(servers string, topic string, groupID string, h ServiceHandleFunc) error
}

// Microservice is the centralized service management
type Microservice struct {
	exitChannel chan bool
}

// ServiceHandleFunc is the handler for each Microservice
type ServiceHandleFunc func(ctx IContext) error

// NewMicroservice is the constructor function of Microservice
func NewMicroservice() *Microservice {
	return &Microservice{}
}

// Consume register service endpoint for Consumer service
func (ms *Microservice) Consume(servers string, topic string, groupID string, readTimeout time.Duration, h ServiceHandleFunc) error {

	go func() {
		c, err := ms.newKafkaConsumer(servers, groupID)
		if err != nil {
			return
		}

		defer c.Close()

		for {
			if readTimeout <= 0 {
				// readtimeout -1 indicates no timeout
				readTimeout = -1
			}

			msg, err := c.ReadMessage(readTimeout)
			if err != nil {
				kafkaErr, ok := err.(kafka.Error)
				if ok {
					if kafkaErr.Code() == kafka.ErrTimedOut {
						if readTimeout == -1 {
							// No timeout just continue to read message again
							continue
						}
					}
				}
				return
			}

			// Execute Handler
			h(NewConsumerContext(ms, string(msg.Value)))
		}
	}()

	return nil
}

// Start start all registered services
func (ms *Microservice) Start() error {
	// There are 2 ways to exit from Microservices
	// 1. The SigTerm can be send from outside program such as from k8s
	// 2. Send true to ms.exitChannel
	osQuit := make(chan os.Signal, 1)
	ms.exitChannel = make(chan bool, 1)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)
	exit := false
	for {
		if exit {
			break
		}
		select {
		case <-osQuit:
			// if exitHTTP != nil {
			// 	exitHTTP <- true
			// }
			// for i := 0; i < scN; i++ {
			// 	exitSC <- true
			// }
			exit = true
		case <-ms.exitChannel:
			// if exitHTTP != nil {
			// 	exitHTTP <- true
			// }
			// for i := 0; i < scN; i++ {
			// 	exitSC <- true
			// }
			exit = true
		}
	}

	return nil
}

// Stop stop the services
func (ms *Microservice) Stop() {
	if ms.exitChannel == nil {
		return
	}
	ms.exitChannel <- true
}

// Cleanup clean resources up from every registered services before exit
func (ms *Microservice) Cleanup() error {
	return nil
}

// newKafkaConsumer create new Kafka consumer
func (ms *Microservice) newKafkaConsumer(servers string, groupID string) (*kafka.Consumer, error) {
	// Configurations
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	config := &kafka.ConfigMap{

		// Alias for metadata.broker.list: Initial list of brokers as a CSV list of broker host or host:port.
		// The application may also use rd_kafka_brokers_add() to add brokers during runtime.
		"bootstrap.servers": servers,

		// Client group id string. All clients sharing the same group.id belong to the same group.
		"group.id": groupID,

		// Action to take when there is no initial offset in offset store or the desired offset is out of range:
		// 'smallest','earliest' - automatically reset the offset to the smallest offset,
		// 'largest','latest' - automatically reset the offset to the largest offset,
		// 'error' - trigger an error which is retrieved by consuming messages and checking 'message->err'.
		"auto.offset.reset": "latest",

		// Protocol used to communicate with brokers.
		// plaintext, ssl, sasl_plaintext, sasl_ssl
		"security.protocol": "plaintext",

		// Client group session and failure detection timeout.
		// The consumer sends periodic heartbeats (heartbeat.interval.ms) to indicate its liveness to the broker.
		// If no hearts are received by the broker for a group member within the session timeout,
		// the broker will remove the consumer from the group and trigger a rebalance.
		// The allowed range is configured with the broker configuration properties group.min.session.timeout.ms
		// and group.max.session.timeout.ms. Also see max.poll.interval.ms.
		// 10000ms 	= 10s (default)
		// 20000ms  = 20s
		"session.timeout.ms": 10,

		// Automatically and periodically commit offsets in the background.
		// Note: setting this to false does not prevent the consumer from fetching previously committed start offsets.
		// To circumvent this behaviour set specific start offsets per partition in the call to assign().
		"enable.auto.commit": true,

		// The frequency in milliseconds that the consumer offsets are committed (written) to offset storage. (0 = disable).
		// default = 5000ms (5s)
		// 5s is too large, it might cause double process message easily, so we reduce this to 200ms (if we turn on enable.auto.commit)
		"auto.commit.interval.ms": 5000,

		// Automatically store offset of last message provided to application.
		// The offset store is an in-memory store of the next offset to (auto-)commit for each partition
		// and cs.Commit() <- offset-less commit
		"enable.auto.offset.store": true,

		// Minimum number of messages per topic+partition librdkafka tries to maintain in the local consumer queue.
		// default = 100000
		"queued.min.messages": 100000,

		// Maximum number of kilobytes per topic+partition in the local consumer queue.
		// This value may be overshot by fetch.message.max.bytes.
		// This property has higher priority than queued.min.messages.
		// default = 1048576 kbytes = 1 Gigabytes
		// 100000 kbytes = 100 Megabytes
		"queued.max.messages.kbytes": 100000,

		// Initial maximum number of bytes per topic+partition to request when fetching messages from the broker.
		// If the client encounters a message larger than this value it will gradually try to increase it
		// until the entire message can be fetched.
		// default = 1048576
		// 1048576 bytes = 1 Megabytes
		"fetch.message.max.bytes": 1048576,

		// Maximum amount of data the broker shall return for a Fetch request.
		// Messages are fetched in batches by the consumer and if the first message batch in the first non-empty partition of the Fetch request is larger than this value,
		// then the message batch will still be returned to ensure the consumer can make progress.
		// The maximum message batch size accepted by the broker is defined via message.max.bytes (broker config) or
		// max.message.bytes (broker topic config).
		// fetch.max.bytes is automatically adjusted upwards to be at least message.max.bytes (consumer config).
		// default = 52428800
		// 52428800 bytes = 52 Megabytes
		"fetch.max.bytes": 52428800,

		// Maximum Kafka protocol request message size.
		// default = 1000000
		// 1000000 bytes = 1 Megabytes
		"message.max.bytes": 1000000,

		// How long to postpone the next fetch request for a topic+partition in case of a fetch error.
		// 500 ms = 0.5 sec
		"fetch.error.backoff.ms": 500,

		// Enable TCP keep-alives (SO_KEEPALIVE) on broker sockets
		"socket.keepalive.enable": true,

		// Default timeout for network requests.
		// Producer: ProduceRequests will use the lesser value of socket.timeout.ms
		//           and remaining message.timeout.ms for the first message in the batch.
		// Consumer: FetchRequests will use fetch.wait.max.ms + socket.timeout.ms.
		// Admin: Admin requests will use socket.timeout.ms or explicitly set rd_kafka_AdminOptions_set_operation_timeout() value.
		// Default = 60000ms = 60s
		// 300000 = 5m
		"socket.timeout.ms": 300000,

		// Minimum number of bytes the broker responds with.
		// If fetch.wait.max.ms expires the accumulated data will be sent to the client regardless of this setting.
		// This property allows a consumer to specify the minimum amount of data that it wants to receive from the broker when fetching records.
		// If a broker receives a request for records from a consumer but the new records amount to fewer bytes than fetch.min.bytes,
		// the broker will wait until more messages are available before sending the records back to the consumer.
		// This reduces the load on both the consumer and the broker as they have to handle fewer back-and-forth messages
		// in cases where the topics don’t have much new activity (or for lower activity hours of the day).
		// You will want to set this parameter higher than the default if the consumer is using too much CPU when there isn’t much data available,
		// or reduce load on the brokers when you have large number of consumers.
		// For example, suppose the value is set to 6 bytes and the timeout on a poll is set to 100ms.
		// If there are 5 bytes available and no further bytes come in before the 100ms expire, the poll returns with nothing.
		// Default = 1bytes
		// 10 = 10bytes
		"fetch.min.bytes": 10,

		// Maximum time the broker may wait to fill the response with fetch.min.bytes.
		// By setting fetch.min.bytes, you tell Kafka to wait until it has enough data to send before responding to the consumer.
		// fetch.max.wait.ms lets you control how long to wait. By default, Kafka will wait up to 500 ms.
		// This results in up to 500 ms of extra latency in case there is not enough data flowing to the Kafka topic to satisfy the minimum amount of data to return.
		// If you want to limit the potential latency (usually due to SLAs controlling the maximum latency of the application), you can set fetch.max.wait.ms to a lower value.
		// If you set fetch.max.wait.ms to 100 ms and fetch.min.bytes to 1 MB, Kafka will receive a fetch request from the consumer
		// and will respond with data either when it has 1 MB of data to return or after 100 ms, whichever happens first.
		// Default = 100ms =  0.1s
		// 200 = 0.2s
		"fetch.wait.max.ms": 100,

		// The maximum delay between invocations of poll() when using consumer group management.
		// This places an upper bound on the amount of time that the consumer can be idle before fetching more records.
		// If poll() is not called before expiration of this timeout, then the consumer is considered failed and the group will
		// rebalance in order to reassign the partitions to another member. For consumers using a non-null group.instance.id which reach this timeout,
		// partitions will not be immediately reassigned. Instead, the consumer will stop sending heartbeats and partitions will be reassigned
		// after expiration of session.timeout.ms. This mirrors the behavior of a static consumer which has shutdown.
		// 300000ms = 5m (default)
		// 600000ms = 10m
		// We don't want to bother with timeout that cause rebalance, so we set max.poll.interval.ms to very high value,
		// so we rely on monitoring (prometheus) to monitor consumer that use too much time to process it messages
		"max.poll.interval.ms": 300000,

		// Name of partition assignment strategy to use when elected group leader assigns partitions to group members.
		// default = range
		// possible value = range,roundrobin
		"partition.assignment.strategy": "range",
	}

	kc, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return kc, err
}
