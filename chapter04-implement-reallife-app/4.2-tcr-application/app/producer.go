// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// IProducer is interface for producer
type IProducer interface {
	// SendMessage will send message to the partition
	SendMessage(topic string, key string, message interface{}) error
	// Close the producer
	Close() error
}

// Producer implement IProducer, is the service to send message to Kafka
type Producer struct {
	ms      *Microservice
	servers string
	prod    *kafka.Producer
}

// NewProducer return new instance of Producer
func NewProducer(servers string, ms *Microservice) *Producer {
	return &Producer{
		ms:      ms,
		servers: servers,
	}
}

func (p *Producer) getProducer() *kafka.Producer {
	if p.prod == nil {
		prod, _ := p.newKafkaProducer(p.servers)
		p.prod = prod
	}
	return p.prod
}

// SendMessage send message to topic synchronously
func (p *Producer) SendMessage(topic string, key string, message interface{}) error {

	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	var keyBytes []byte
	if len(key) > 0 {
		keyBytes = []byte(key)
	}

	// Send Message Synchrounously
	deliveryChan := make(chan kafka.Event)

	prod := p.getProducer()
	err = prod.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(messageJSON),
		Key:            keyBytes,
	}, deliveryChan)
	if err != nil {
		return err
	}

	<-deliveryChan
	close(deliveryChan)

	return nil
}

// Close the producer
func (p *Producer) Close() error {
	if p.prod == nil {
		return nil
	}

	prod := p.prod
	prod.Flush(5000) // 5s for flush message in queue
	prod.Close()

	p.ms.Log("PROD", "Close successfully")

	return nil
}

func (p *Producer) newKafkaProducer(servers string) (*kafka.Producer, error) {

	// Configurations
	// https://github.com/edenhill/librdkafka/blob/master/CONFIGURATION.md
	return kafka.NewProducer(&kafka.ConfigMap{

		// Alias for metadata.broker.list: Initial list of brokers as a CSV list of broker host or host:port.
		// The application may also use rd_kafka_brokers_add() to add brokers during runtime.
		"bootstrap.servers": servers,

		// Protocol used to communicate with brokers.
		// plaintext, ssl, sasl_plaintext, sasl_ssl
		"security.protocol": "plaintext",

		// Maximum number of messages allowed on the producer queue. This queue is shared by all topics and partitions.
		// default is 100000 messages
		// our default = 1000000 messages
		"queue.buffering.max.messages": 1000000,

		// Maximum total message size sum allowed on the producer queue.
		// This queue is shared by all topics and partitions.
		// This property has higher priority than queue.buffering.max.messages.
		// default is 1048576 kbytes = 1 Gigabytes
		// --
		// Below is the explaination from Kafka consumer page https://docs.confluent.io/current/installation/configuration/producer-configs.html
		// This setting should correspond roughly to the total memory the producer will use, but is not a hard bound
		// since not all memory the producer uses is used for buffering.
		// Some additional memory will be used for compression (if compression is enabled) as well as for maintaining in-flight requests.
		// --
		// librdkafka default is 1048576 kbytes = 1 Gigabytes
		//  10000 kbytes =  10 Megabytes
		//  50000 kbytes =  50 Megabytes
		// 100000 kbytes = 100 Megabytes <-- our default value
		"queue.buffering.max.kbytes": 100000,

		// Delay in milliseconds to wait for messages in the producer queue to accumulate before constructing message batches (MessageSets) to transmit to brokers.
		// A higher value allows larger and more effective (less overhead, improved compression) batches of messages to accumulate at the expense of increased message delivery latency.
		// 0 = send immediately
		"queue.buffering.max.ms": 0,
		"compression.codec":      "snappy",

		// Maximum number of messages batched in one MessageSet.
		// The total MessageSet size is also limited by message.max.bytes.
		// ---
		// Q: Why are my messages not being sent directly?
		// A: Because the number of messages in the queue has not reached batch.num.messages yet (and not timed out yet queue.buffering.max.ms)
		// ---
		// Q: Why are the MessageSets sent to the broker not bigger?
		// A: Because it is limited by batch.num.messages
		// default value = 10K (for async)
		"batch.num.messages": 10000,

		// Maximum Kafka protocol request message size.
		// Due to differing framing overhead between protocol versions the producer is unable to reliably enforce a strict max message limit
		// at produce time and may exceed the maximum size by one message in protocol ProduceRequests,
		// the broker will enforce the the topic's max.message.bytes limit (see Apache Kafka documentation).
		// Default = 1000000
		// 1000000 bytes = 1 Megabytes
		// 10000000 bytes = 10 Megabytes
		"message.max.bytes": 10000000,

		// How many times to retry sending a failing Message. Note: retrying may cause reordering unless enable.idempotence is set to true.
		// ---
		// Reference for detail belows, https://blog.newrelic.com/engineering/kafka-best-practices/
		// The right value will depend on your application;
		// for applications where data-loss cannot be tolerated, set this to max value 10000000
		// This guards against situations where the broker leading the partition isn’t able to respond to a produce request right away.
		// default = 2
		// MAX = 10000000
		// When we set message.send.max.retries 10000000, so it depends on message.timeout.ms (set below) before we lost the message
		"message.send.max.retries": 10000000,

		// The backoff time in milliseconds before retrying a protocol request.
		// 100 ms = 0.1 sec
		"retry.backoff.ms": 100,

		// delivery.timeout.ms is alias for message.timeout.ms
		// This value is only enforced locally and limits the time a produced message waits for successful delivery.
		// This is the maximum time librdkafka may use to deliver a message (including retries).
		// Delivery error occurs when either the retry count or the message timeout are exceeded.
		// default 300000 = 5m
		// 900000 = 15m
		// 3600000 = 1h
		// 10800000 = 3h
		// 21600000 = 6h
		// 43200000 = 12h
		"message.timeout.ms": 43200000,

		// Default timeout for network requests.
		// Producer: ProduceRequests will use the lesser value of socket.timeout.ms and remaining message.timeout.ms for the **first message in the batch.
		// Consumer: FetchRequests will use fetch.wait.max.ms + socket.timeout.ms.
		// Admin: Admin requests will use socket.timeout.ms or explicitly set rd_kafka_AdminOptions_set_operation_timeout() value.
		// Default = 60000ms = 60s
		// 600000 = 10m
		// 300000 = 5m
		"socket.timeout.ms": 300000,

		// This field indicates the number of acknowledgements the leader broker must receive from ISR brokers before responding to the request:
		// 0=Broker does not send any response/ack to client,
		// -1 or all=Broker will block until message is committed by all in sync replicas (ISRs).
		// If there are less than min.insync.replicas (broker configuration) in the ISR set the produce request will fail.
		// The message will be resent up to message.send.max.retries times before reporting a failure back to the application.
		"request.required.acks": -1,

		// The ack timeout of the producer request in milliseconds.
		// This value is only enforced by the broker and relies on request.required.acks being != 0.
		// The configuration controls the maximum amount of time the client will wait for the response of a request.
		// If the response is not received before the timeout elapses the client will resend the request if necessary
		// or fail the request if retries are exhausted. This should be larger than replica.lag.time.max.ms (a broker configuration)
		// to reduce the possibility of message duplication due to unnecessary producer retries.
		// 30000 = 30s (default)
		// 60000 = 60s
		"request.timeout.ms": 60000,

		// Enable TCP keep-alives (SO_KEEPALIVE) on broker sockets
		"socket.keepalive.enable": true,

		// https://www.cloudkarafka.com/blog/2019-04-10-apache-kafka-idempotent-producer-avoiding-message-duplication.html
		// When set to true, the producer will ensure that messages are successfully produced exactly once and in the original produce order.
		// The following configuration properties are adjusted automatically (if not modified by the user) when idempotence is enabled:
		// max.in.flight.requests.per.connection=5 (must be less than or equal to 5),
		// retries=INT32_MAX (must be greater than 0),
		// acks=all, queuing.strategy=fifo. Producer instantation will fail if user-supplied configuration is incompatible.
		"enable.idempotence": true,
	})
}
