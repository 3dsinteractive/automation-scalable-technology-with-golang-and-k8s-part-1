// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import "net/http"

func main() {
	ms := NewMicroservice()

	cacheServer := "localhost:6379"
	mqServers := "localhost:9094"
	ms.PTask("/citizen/batch", cacheServer, mqServers)

	// Start 3 workers
	for i := 0; i < 3; i++ {
		ms.PTaskWorker("/citizen/batch", cacheServer, mqServers, func(ctx IContext) error {
			ctx.Log(ctx.ReadInput())
			res := map[string]interface{}{
				"id": "123",
			}
			ctx.Response(http.StatusOK, res)
			return nil
		})
	}

	defer ms.Cleanup()
	ms.Start()
}
