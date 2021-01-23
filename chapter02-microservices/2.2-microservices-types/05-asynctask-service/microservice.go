// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
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
	Log(message string)

	// HTTP Services
	GET(path string, h ServiceHandleFunc)
	POST(path string, h ServiceHandleFunc)
	PUT(path string, h ServiceHandleFunc)
	PATCH(path string, h ServiceHandleFunc)
	DELETE(path string, h ServiceHandleFunc)

	// Consumer Services
	Consume(servers string, topic string, groupID string, h ServiceHandleFunc) error

	// AsyncTask Services
	AsyncPOST(path string, mqServers string, h ServiceHandleFunc)
	AsyncPUT(path string, mqServers string, h ServiceHandleFunc)
}

// Microservice is the centralized service management
type Microservice struct {
	echo        *echo.Echo
	exitChannel chan bool
}

// ServiceHandleFunc is the handler for each Microservice
type ServiceHandleFunc func(ctx IContext) error

// NewMicroservice is the constructor function of Microservice
func NewMicroservice() *Microservice {
	return &Microservice{
		echo: echo.New(),
	}
}

// startAsyncTaskConsumer read async task message from message queue and execute with handler
func (ms *Microservice) startAsyncTaskConsumer(path string, cacheServer string, mqServers string, h ServiceHandleFunc) {
	topic := escapeName(path)
	mq := NewMQ(mqServers, ms)
	err := mq.CreateTopicR(topic, 5, 1, time.Hour*24*30) // retain message for 30 days
	if err != nil {
		ms.Log("ATASK", err.Error())
	}

	ms.Consume(mqServers, topic, "atask", -1, func(ctx IContext) error {
		message := map[string]interface{}{}
		err := json.Unmarshal([]byte(ctx.ReadInput()), &message)
		if err != nil {
			return err
		}
		ref, _ := message["ref"].(string)
		input, _ := message["input"].(string)
		return h(NewAsyncTaskContext(ms, cacheServer, ref, input))
	})
}

// handleAsyncTaskRequest accept async task request and send it to message queue
func (ms *Microservice) handleAsyncTaskRequest(path string, cacherServer string, mqServers string, ctx IContext) error {
	topic := escapeName(path)

	// 1. Read Input
	input := ctx.ReadInput()

	// 2. Generate REF
	ref := fmt.Sprintf("atask-%s", randString())

	// 3. Set Status in Cache
	cacher := ctx.Cacher(cacherServer)
	status := map[string]interface{}{
		"status": "processing",
	}
	expire := time.Minute * 30
	cacher.Set(ref, status, expire)

	// 4. Send Message to MQ
	prod := ctx.Producer(mqServers)
	message := map[string]interface{}{
		"ref":   ref,
		"input": input,
	}
	prod.SendMessage(topic, "", message)

	// 5. Response REF
	res := map[string]string{
		"ref": ref,
	}
	ctx.Response(http.StatusOK, res)
	return nil
}

func (ms *Microservice) handleAsyncTaskResponse(path string, cacheServer string, ctx IContext) error {
	// 1. ReadInput (REF from query string)
	ref := ctx.QueryParam("ref")

	// 2. Read Status from Cache
	cacher := ctx.Cacher(cacheServer)
	statusJS, err := cacher.Get(ref)
	if err != nil {
		return err
	}

	// 3. Return Status
	status := map[string]interface{}{}
	err = json.Unmarshal([]byte(statusJS), &status)
	if err != nil {
		return err
	}
	ctx.Response(http.StatusOK, status)
	return nil
}

// AsyncPOST register async task service for HTTP POST
func (ms *Microservice) AsyncPOST(path string, cacheServer string, mqServers string, h ServiceHandleFunc) {
	ms.startAsyncTaskConsumer(path, cacheServer, mqServers, h)
	ms.GET(path, func(ctx IContext) error {
		return ms.handleAsyncTaskResponse(path, cacheServer, ctx)
	})
	ms.POST(path, func(ctx IContext) error {
		return ms.handleAsyncTaskRequest(path, cacheServer, mqServers, ctx)
	})
}

// AsyncPUT register async task service for HTTP PUT
func (ms *Microservice) AsyncPUT(path string, cacheServer string, mqServers string, h ServiceHandleFunc) {
	ms.startAsyncTaskConsumer(path, cacheServer, mqServers, h)
	ms.GET(path, func(ctx IContext) error {
		return ms.handleAsyncTaskResponse(path, cacheServer, ctx)
	})
	ms.PUT(path, func(ctx IContext) error {
		return ms.handleAsyncTaskRequest(path, cacheServer, mqServers, ctx)
	})
}

// GET register service endpoint for HTTP GET
func (ms *Microservice) GET(path string, h ServiceHandleFunc) {
	ms.echo.GET(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// POST register service endpoint for HTTP POST
func (ms *Microservice) POST(path string, h ServiceHandleFunc) {
	ms.echo.POST(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// PUT register service endpoint for HTTP PUT
func (ms *Microservice) PUT(path string, h ServiceHandleFunc) {
	ms.echo.PUT(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// PATCH register service endpoint for HTTP PATCH
func (ms *Microservice) PATCH(path string, h ServiceHandleFunc) {
	ms.echo.PATCH(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// DELETE register service endpoint for HTTP DELETE
func (ms *Microservice) DELETE(path string, h ServiceHandleFunc) {
	ms.echo.DELETE(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// startHTTP will start HTTP service, this function will block thread
func (ms *Microservice) startHTTP(exitChannel chan bool) error {
	// Caller can exit by sending value to exitChannel
	go func() {
		<-exitChannel
		ms.stopHTTP()
	}()
	return ms.echo.Start(":8080")
}

func (ms *Microservice) stopHTTP() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ms.echo.Shutdown(ctx)
}

func (ms *Microservice) consumeSingle(servers string, topic string, groupID string, readTimeout time.Duration, h ServiceHandleFunc) {
	c, err := ms.newKafkaConsumer(servers, groupID)
	if err != nil {
		return
	}

	defer c.Close()

	c.Subscribe(topic, nil)

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
			ms.Log("Consumer", err.Error())
			return
		}

		// Execute Handler
		h(NewConsumerContext(ms, string(msg.Value)))
	}
}

// Consume register service endpoint for Consumer service
func (ms *Microservice) Consume(servers string, topic string, groupID string, readTimeout time.Duration, h ServiceHandleFunc) error {
	go ms.consumeSingle(servers, topic, groupID, readTimeout, h)
	return nil
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
			// for i := 0; i < scN; i++ {
			// 	exitSC <- true
			// }
			exit = true
		case <-ms.exitChannel:
			// Exit from HTTP as well
			if exitHTTP != nil {
				exitHTTP <- true
			}
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

// Log log message to console
func (ms *Microservice) Log(tag string, message string) {
	fmt.Println(tag+": ", message)
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
