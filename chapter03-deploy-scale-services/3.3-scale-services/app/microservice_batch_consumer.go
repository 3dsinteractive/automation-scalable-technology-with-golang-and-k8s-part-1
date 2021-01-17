// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func (ms *Microservice) consumeBatch(
	servers string,
	topic string,
	groupID string,
	readTimeout time.Duration,
	batchSize int,
	batchTimeout time.Duration,
	h ServiceHandleFunc) error {

	// Batch Filler
	fill := func(b *Batch, payload interface{}) error {
		p := payload.(string)
		b.Add(p)
		return nil
	}

	// Batch Executer
	exec := func(b *Batch) error {
		messages := make([]string, 0)
		for {
			item := b.Read()
			if item == nil {
				break
			}
			message := item.(string)
			messages = append(messages, message)
		}

		if len(messages) == 0 {
			return nil
		}

		// Execute Handler
		h(NewBatchConsumerContext(ms, messages))
		return nil
	}

	// Payloads loader
	payload := make(chan interface{})
	quit := make(chan bool, 1)

	go func() {

		c, err := ms.newKafkaConsumer(servers, groupID)
		if err != nil {
			quit <- true
			return
		}

		// This will close kafka consumer before exit function
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
				return
			}

			message := string(msg.Value)
			payload <- message
		}
	}()

	go func() {
		// Gracefull shutdown routine
		osQuit := make(chan os.Signal, 1)
		signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)

		select {
		case <-quit:
			close(payload)
		case <-osQuit:
			close(payload)
		}
	}()

	// Error listener
	errc := make(chan error)
	defer close(errc)
	go func() {
		for err := range errc {
			if err != nil {
				ms.Log("BatchConsumer", err.Error())
			}
		}
	}()

	be := NewBatchEvent(batchSize, batchTimeout, fill, exec, payload, errc)
	be.Start()

	return nil
}

// ConsumeBatch register service endpoint for Batch Consumer service
func (ms *Microservice) ConsumeBatch(
	servers string,
	topic string,
	groupID string,
	readTimeout time.Duration,
	batchSize int,
	batchTimeout time.Duration,
	h ServiceHandleFunc) error {

	go ms.consumeBatch(servers, topic, groupID, readTimeout, batchSize, batchTimeout, h)
	return nil
}
