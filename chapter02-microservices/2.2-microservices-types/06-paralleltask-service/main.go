// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"math/rand"
	"net/http"
	"time"
)

func main() {
	ms := NewMicroservice()

	cacheServer := "localhost:6379"
	mqServers := "localhost:9094"

	// 1. Start PTask endpoint
	ms.PTaskEndpoint("/citizen/batch", cacheServer, mqServers)

	// 2. Start 2 worker nodes
	for i := 0; i < 2; i++ {
		ms.PTaskWorkerNode("/citizen/batch", cacheServer, mqServers, func(ctx IContext) error {
			ctx.Log(ctx.ReadInput())
			resStr := randString()
			res := map[string]interface{}{
				"result": resStr,
			}
			n := rand.Intn(5)
			time.Sleep(time.Duration(n) * time.Second)
			ctx.Response(http.StatusOK, res)
			return nil
		})
	}

	defer ms.Cleanup()
	ms.Start()
}
