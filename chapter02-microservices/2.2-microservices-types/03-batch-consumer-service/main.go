// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"time"
)

func main() {
	ms := NewMicroservice()

	servers := "localhost:9094"
	topic := "when-citizen-has-registered-batch-" + randString()
	groupID := "validation-consumer"
	timeout := time.Duration(-1)
	batchSize := 3
	batchTimeout := time.Second * 5

	ms.ConsumeBatch(servers, topic, groupID, timeout, batchSize, batchTimeout,
		func(ctx IContext) error {
			msgs := ctx.ReadInputs()
			ctx.Log("Begin Batch")
			for _, msg := range msgs {
				ctx.Log(msg)
			}
			ctx.Log("End Batch")
			return nil
		})

	prod := NewProducer(servers, ms)
	go func() {
		for i := 0; i < 10; i++ {
			prod.SendMessage(topic, "", map[string]interface{}{"message_id": i})
			time.Sleep(time.Second)
		}

		// Wait for last batch then exit program
		time.Sleep(5 * time.Second)
		ms.Stop()
	}()

	defer ms.Cleanup()
	ms.Start()
}
