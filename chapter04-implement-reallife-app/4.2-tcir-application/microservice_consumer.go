// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

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
			ms.Stop()
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
