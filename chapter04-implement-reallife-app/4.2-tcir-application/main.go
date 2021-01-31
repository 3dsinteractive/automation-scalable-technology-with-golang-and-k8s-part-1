// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"encoding/json"
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
		startBatchPTaskWorkerNode(ms, cfg)
	case "external-api":
		start3rdPartyMockAPI(ms, cfg)
	}

	ms.Start()
}

func startRegisterAPI(ms *Microservice, cfg IConfig) {
	ms.AsyncPOST("/api/citizen", cfg.CacheServer(), cfg.MQServers(), func(ctx IContext) error {
		// 1. Read Input (Not using it right now, just for example)
		input := ctx.ReadInput()
		ctx.Log("POST: /api/citizen " + input)

		// 2. Generate citizenID and send it to MQ
		//    The citizen id should be received from client, but for code to be easy to read, we just create it
		citizenID := randString()
		citizen := map[string]interface{}{
			"citizen_id": citizenID,
		}
		prod := ctx.Producer(cfg.MQServers())
		err := prod.SendMessage(cfg.CitizenRegisteredTopic(), "", citizen)
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
	topic := cfg.CitizenRegisteredTopic()
	groupID := "mail-consumer"
	timeout := time.Duration(-1)

	// 1. Create topic "citizen registered" if not exists
	mq := NewMQ(cfg.MQServers(), ms)
	mq.CreateTopicR(topic, 5, 1, time.Hour*24*30)

	// 2. Start consumer to consume message from "citizen registered" topic
	ms.Consume(cfg.MQServers(), topic, groupID, timeout, func(ctx IContext) error {
		msg := ctx.ReadInput()

		// 3. Parse input to citizen object
		citizen := &Citizen{}
		err := json.Unmarshal([]byte(msg), &citizen)
		if err != nil {
			ctx.Log(err.Error())
			return err
		}

		// 4. Call validation API (Response AVG at 1 second, so we set timeout at 5 seconds)
		req := ctx.Requester("", 5*time.Second)
		validationResStr, err := req.Post(cfg.CitizenValidationAPI(),
			map[string]string{"citizen_id": citizen.CitizenID})
		if err != nil {
			ctx.Log(err.Error())
			return err
		}

		validationRes := map[string]interface{}{}
		err = json.Unmarshal([]byte(validationResStr), &validationRes)
		if err != nil {
			ctx.Log(err.Error())
			return err
		}

		isValid, _ := validationRes["status"]
		if isValid != "ok" {
			// 5. Send Email to citizen to reject register if validation is not OK
			//    We just log to console, but for the real code, this should send the email
			ctx.Log("Mail rejection has sent to " + citizen.CitizenID)
			return nil
		}

		// 6. Send Email to citizen to confirm validation
		//    We just log to console, but for the real code, this should send the email
		ctx.Log("Mail confirmation has sent to " + citizen.CitizenID)

		// 7. Produce message to topic "citizen confirmed"
		prod := ctx.Producer(cfg.MQServers())
		err = prod.SendMessage(cfg.CitizenConfirmedTopic(), "", citizen)
		if err != nil {
			ctx.Log(err.Error())
			return err
		}

		return nil
	})
}

func startBatchScheduler(ms *Microservice, cfg IConfig) {
	ms.Schedule(time.Hour, func(ctx IContext) error {
		// 1. Batch Scheduler will run during 00.00 - 00.59
		nowH := ctx.Now().Hour()
		if nowH != 0 {
			return nil
		}

		// 2. Will start PTask to execute all workers
		//    This run only 1 time a day, to make sure it will run, use 30 secs timeout
		rqt := ctx.Requester("", 30*time.Second)
		res, err := rqt.Post(cfg.BatchDeliverAPI(),
			map[string]string{"task_id": "batch_deliver", "worker_count": "5"})
		if err != nil {
			ctx.Log(err.Error())
			return err
		}
		ctx.Log(res)

		return nil
	})
}

func startBatchPTaskAPI(ms *Microservice, cfg IConfig) {
	ms.PTaskEndpoint("/ptask/delivery", cfg.CacheServer(), cfg.MQServers())
}

func startBatchPTaskWorkerNode(ms *Microservice, cfg IConfig) {
	ms.PTaskWorkerNode("/ptask/delivery",
		cfg.CacheServer(),
		cfg.MQServers(),
		func(ctx IContext) error {

			newMS := NewMicroservice()
			newMS.ConsumeBatch(
				cfg.MQServers(),
				cfg.CitizenConfirmedTopic(),
				"deliver-consumer",
				5*time.Minute, // Read Timeout
				5,             // Batch Size
				5*time.Second, // Batch Timeout
				func(newCtx IContext) error {
					inputs := newCtx.ReadInputs()
					for _, input := range inputs {
						newCtx.Log("Deliver to " + input)
					}
					return nil
				})
			newMS.Start()

			ctx.Response(http.StatusOK, map[string]interface{}{"status": "success"})

			return nil
		})
}

func start3rdPartyMockAPI(ms *Microservice, cfg IConfig) {
	ms.POST("/3rd-party/validate", func(ctx IContext) error {
		time.Sleep(1 * time.Second)
		status := map[string]interface{}{
			"status": "ok",
		}
		ctx.Response(http.StatusOK, status)
		return nil
	})

	ms.POST("/3rd-party/delivery", func(ctx IContext) error {
		time.Sleep(1 * time.Second)
		status := map[string]interface{}{
			"status": "ok",
		}
		ctx.Response(http.StatusOK, status)
		return nil
	})
}
