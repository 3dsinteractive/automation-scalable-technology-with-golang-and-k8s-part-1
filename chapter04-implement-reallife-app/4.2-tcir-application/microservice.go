// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/echo"
)

// IMicroservice is interface for centralized service management
type IMicroservice interface {
	Start() error
	Stop()
	Cleanup() error
	Log(tag string, message string)

	// HTTP Services
	GET(path string, h ServiceHandleFunc)
	POST(path string, h ServiceHandleFunc)
	PUT(path string, h ServiceHandleFunc)
	PATCH(path string, h ServiceHandleFunc)
	DELETE(path string, h ServiceHandleFunc)

	// Consumer Services
	Consume(servers string, topic string, groupID string, readTimeout time.Duration,
		h ServiceHandleFunc) error

	// Batch Consumer Services
	ConsumeBatch(servers string, topic string, groupID string, readTimeout time.Duration,
		batchSize int, batchTimeout time.Duration, h ServiceHandleFunc) error

	// Scheduler Services
	Schedule(timer time.Duration, h ServiceHandleFunc) chan bool /*exit channel*/

	// AsyncTask Services
	AsyncPOST(path string, cacheServer string, mqServers string, h ServiceHandleFunc)
	AsyncPUT(path string, cacheServer string, mqServers string, h ServiceHandleFunc)

	// ParallelTask Services
	PTaskWorkerNode(path string, cacheServer string, mqServers string, h ServiceHandleFunc)
	PTaskEndpoint(path string, cacheServer string, mqServers string)

	// Healthcheck
	RegisterLivenessProbeEndpoint(path string)
}

// Microservice is the centralized service management
type Microservice struct {
	echo        *echo.Echo
	exitChannel chan bool
	prod        IProducer
	cacher      ICacher
}

// ServiceHandleFunc is the handler for each Microservice
type ServiceHandleFunc func(ctx IContext) error

// NewMicroservice is the constructor function of Microservice
func NewMicroservice() *Microservice {
	return &Microservice{
		echo: echo.New(),
	}
}

func (ms *Microservice) getProducer(mqServers string) IProducer {
	if ms.prod == nil {
		ms.prod = NewProducer(mqServers, ms)
	}
	return ms.prod
}

func (ms *Microservice) getCacher(cacheServer string) ICacher {
	if ms.cacher == nil {
		ms.cacher = NewCacher(cacheServer, ms)
	}
	return ms.cacher
}

// Start start all registered services
func (ms *Microservice) Start() error {

	httpN := len(ms.echo.Routes())
	var exitHTTP chan bool
	if httpN > 0 {
		exitHTTP = make(chan bool, 1)
		go func() {
			ms.startHTTP(exitHTTP)
		}()
	}

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
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		case <-ms.exitChannel:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
			exit = true
		}
	}

	ms.Cleanup()
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
	ms.Log("MS", "Start cleanup")
	if ms.cacher != nil {
		ms.cacher.Close()
	}
	if ms.prod != nil {
		ms.prod.Close()
	}
	return nil
}

// Log log message to console
func (ms *Microservice) Log(tag string, message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println(tag+":", fns[len(fns)-1], line, message)
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
		"auto.offset.reset": "earliest",

		// Protocol used to communicate with brokers.
		// plaintext, ssl, sasl_plaintext, sasl_ssl
		"security.protocol": "plaintext",

		// Automatically and periodically commit offsets in the background.
		// Note: setting this to false does not prevent the consumer from fetching previously committed start offsets.
		// To circumvent this behaviour set specific start offsets per partition in the call to assign().
		"enable.auto.commit": true,

		// The frequency in milliseconds that the consumer offsets are committed (written) to offset storage. (0 = disable).
		// default = 5000ms (5s)
		// 5s is too large, it might cause double process message easily, so we reduce this to 200ms (if we turn on enable.auto.commit)
		"auto.commit.interval.ms": 500,

		// Automatically store offset of last message provided to application.
		// The offset store is an in-memory store of the next offset to (auto-)commit for each partition
		// and cs.Commit() <- offset-less commit
		"enable.auto.offset.store": true,

		// Enable TCP keep-alives (SO_KEEPALIVE) on broker sockets
		"socket.keepalive.enable": true,
	}

	kc, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return kc, err
}
