// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	ms := NewMicroservice()

	startHTTP(ms)
	startConsumer(ms)
	startBatchConsumer(ms)
	exitSC := startScheduler(ms)
	startAsyncTask(ms)
	startParallelTask(ms)

	defer func() {
		exitSC <- true
	}()

	ms.Start()
}

func startHTTP(ms *Microservice) {
	ms.POST("/citizen", func(ctx IContext) error {
		ctx.Log("POST: /citizen")
		status := map[string]interface{}{
			"status": "success",
		}
		ctx.Response(http.StatusOK, status)
		return nil
	})

	ms.GET("/citizen/:id", func(ctx IContext) error {
		id := ctx.Param("id")
		page := ctx.QueryParam("page")
		ctx.Log("GET: /citizen/" + id)
		citizen := map[string]interface{}{
			"id":   id,
			"page": page,
		}
		ctx.Response(http.StatusOK, citizen)
		return nil
	})

	ms.PUT("/citizen/:id", func(ctx IContext) error {
		id := ctx.Param("id")
		ctx.Log("PUT: /citizen/" + id)
		citizen := map[string]interface{}{
			"id": id,
		}
		ctx.Response(http.StatusOK, citizen)
		return nil
	})

	ms.DELETE("/citizen/:id", func(ctx IContext) error {
		id := ctx.Param("id")
		ctx.Log("DELETE: /citizen/" + id)
		status := map[string]interface{}{
			"status": "success",
		}
		ctx.Response(http.StatusOK, status)
		return nil
	})
}

func startConsumer(ms *Microservice) {
	servers := "localhost:9094"
	topic := "when-citizen-has-registered-" + randString()
	groupID := "validation-consumer"
	timeout := time.Duration(-1)

	ms.Consume(servers, topic, groupID, timeout, func(ctx IContext) error {
		msg := ctx.ReadInput()
		ctx.Log(msg)
		return nil
	})

	prod := ms.getProducer(servers)
	go func() {
		for i := 0; i < 10; i++ {
			prod.SendMessage(topic, "", map[string]interface{}{"message_id": i})
			time.Sleep(time.Second)
		}

		// Exit program
		ms.Stop()
	}()
}

func startBatchConsumer(ms *Microservice) {
	mqServers := "localhost:9094"
	topic := "when-citizen-has-registered-batch-" + randString()
	groupID := "validation-consumer"
	timeout := time.Duration(-1)
	batchSize := 3
	batchTimeout := time.Second * 5

	ms.ConsumeBatch(mqServers, topic, groupID, timeout, batchSize, batchTimeout,
		func(ctx IContext) error {
			msgs := ctx.ReadInputs()
			ctx.Log("Begin Batch")
			for _, msg := range msgs {
				ctx.Log(msg)
			}
			ctx.Log("End Batch")
			return nil
		})

	prod := ms.getProducer(mqServers)
	go func() {
		for i := 0; i < 10; i++ {
			prod.SendMessage(topic, "", map[string]interface{}{"message_id": i})
			time.Sleep(time.Second)
		}

		// Wait for last batch then exit program
		time.Sleep(5 * time.Second)
		ms.Stop()
	}()
}

func startScheduler(ms *Microservice) chan bool {
	timer := 1 * time.Second
	exitScheduler := ms.Schedule(timer, func(ctx IContext) error {
		now := time.Now()
		ctx.Log(fmt.Sprintf("Tick at %s", now.Format("15:04:05")))
		return nil
	})
	return exitScheduler
}

func startAsyncTask(ms *Microservice) {
	cacheServer := "localhost:6379"
	mqServers := "localhost:9094"
	ms.AsyncPOST("/citizen/register", cacheServer, mqServers, func(ctx IContext) error {
		ctx.Log(ctx.ReadInput())
		res := map[string]interface{}{
			"id": "123",
		}
		ctx.Response(http.StatusOK, res)
		return nil
	})
}

func startParallelTask(ms *Microservice) {
	cacheServer := "localhost:6379"
	mqServers := "localhost:9094"
	ms.PTaskEndpoint("/citizen/batch", cacheServer, mqServers)

	// Start 3 workers
	for i := 0; i < 3; i++ {
		ms.PTaskWorkerNode("/citizen/batch", cacheServer, mqServers, func(ctx IContext) error {
			ctx.Log(ctx.ReadInput())
			res := map[string]interface{}{
				"id": "123",
			}
			n := rand.Intn(5)
			time.Sleep(time.Duration(n) * time.Second)
			ctx.Response(http.StatusOK, res)
			return nil
		})
	}
}
