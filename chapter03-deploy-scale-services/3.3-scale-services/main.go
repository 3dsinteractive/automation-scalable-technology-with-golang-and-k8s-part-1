// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"net/http"
	"time"
)

func main() {
	cfg := NewConfig()

	ms := NewMicroservice()
	ms.RegisterLivenessProbeEndpoint("/healthz")

	// 1. Read SERVICE_ID via environment variable to run each service in separate deployments
	serviceID := cfg.ServiceID()

	// 2. Finally, these 2 services will run in separate deployments,
	//    so we could control how many replicas each services will be
	switch serviceID {
	case "register-api":
		startHTTP(ms, cfg)
	case "mail-consumer":
		startConsumer(ms, cfg)
	}

	ms.Start()
}

func startHTTP(ms *Microservice, cfg IConfig) {
	ms.POST("/citizen", func(ctx IContext) error {

		// 1. Read Input (Not using it right now, just for example)
		input := ctx.ReadInput()
		ctx.Log("POST: /citizen " + input)

		// 2. Generate citizenID and send it to MQ
		citizenID := randString()
		citizen := map[string]interface{}{
			"citizen_id": citizenID,
		}
		prod := ctx.Producer(cfg.MQServers())
		err := prod.SendMessage("when-citizen-has-registered", "", citizen)
		if err != nil {
			ctx.Log(err.Error())
			return err
		}

		// 3. Response citizenID
		status := map[string]interface{}{
			"status":     "success",
			"citizen_id": citizenID,
		}
		ctx.Response(http.StatusOK, status)
		return nil
	})
}

func startConsumer(ms *Microservice, cfg IConfig) {
	topic := "when-citizen-has-registered"
	groupID := "mail-consumer"
	timeout := time.Duration(-1)

	mq := NewMQ(cfg.MQServers(), ms)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*30)
	ms.Consume(cfg.MQServers(), topic, groupID, timeout, func(ctx IContext) error {
		msg := ctx.ReadInput()
		ctx.Log("Mail has sent to " + msg)
		return nil
	})
}
