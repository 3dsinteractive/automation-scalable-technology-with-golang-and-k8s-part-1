// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// IMQ is interface to manage Kafka
type IMQ interface {
	CreateTopic(topic string, partitions int, replications int) error
	CreateTopicR(topic string, partitions int, replications int, retentionPeriod time.Duration) error
}

// MQ is message queue
type MQ struct {
	ms      *Microservice
	servers string
}

// NewMQ return new MQ
func NewMQ(servers string, ms *Microservice) *MQ {
	return &MQ{
		ms:      ms,
		servers: servers,
	}
}

func (q *MQ) getAdminClient() (*kafka.AdminClient, error) {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{"bootstrap.servers": q.servers})
	if err != nil {
		return nil, err
	}
	return admin, nil
}

// CreateTopicR create topic with retention period
func (q *MQ) CreateTopicR(topic string, partitions int, replications int, retentionPeriod time.Duration) error {
	return q.createTopic(topic, partitions, replications, retentionPeriod)
}

// CreateTopic create the topic
func (q *MQ) CreateTopic(topic string, partitions int, replications int) error {
	return q.createTopic(topic, partitions, replications, 0)
}

func (q *MQ) createTopic(topic string, partitions int, replications int, retentionPeriod time.Duration) error {
	if retentionPeriod <= 0 {
		retentionPeriod = 7 * (time.Hour * 24) // default = 7 days (Message will keep 7 days)
	}

	admin, err := q.getAdminClient()
	if err != nil {
		return err
	}

	defer admin.Close()

	// Operation timeout for create topic = 5 minutes
	timeout, err := time.ParseDuration("5m")
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	retentionPeriodMillisec := fmt.Sprintf("%d", int64(retentionPeriod/time.Millisecond))

	var results []kafka.TopicResult
	if timeout > 0 {
		results, err = admin.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{{
				Topic:             topic,
				NumPartitions:     partitions,
				ReplicationFactor: replications,
				Config: map[string]string{
					"retention.ms": retentionPeriodMillisec,
				}}},
			kafka.SetAdminOperationTimeout(timeout))
	} else {
		results, err = admin.CreateTopics(
			ctx,
			[]kafka.TopicSpecification{{
				Topic:             topic,
				NumPartitions:     partitions,
				ReplicationFactor: replications,
				Config: map[string]string{
					"retention.ms": retentionPeriodMillisec,
				}}})
	}
	if err != nil {
		return err
	}

	for _, result := range results {
		q.ms.Log("MQ", result.String())
	}

	return nil
}
