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

	serviceID := cfg.ServiceID()

	switch serviceID {
	case "register-api":
		startRegisterAPI(ms, cfg)
	case "mail-consumer":
		startMailConsumer(ms, cfg)
	case "batch-scheduler":
		startBatchScheduler(ms, cfg)
	case "batch-ptask-api":
		startBatchPTaskAPI(ms, cfg)
	case "batch-ptask-worker":
		startBatchPTaskWorker(ms, cfg)
	}

	ms.Start()
}

func startRegisterAPI(ms *Microservice, cfg IConfig) {
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

func startMailConsumer(ms *Microservice, cfg IConfig) {
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

func startBatchScheduler(ms *Microservice, cfg IConfig) {

}

func startBatchPTaskAPI(ms *Microservice, cfg IConfig) {

}

func startBatchPTaskWorker(ms *Microservice, cfg IConfig) {

}
