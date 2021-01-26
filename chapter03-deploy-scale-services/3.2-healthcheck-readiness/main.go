// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"net/http"
	"time"
)

func main() {
	ms := NewMicroservice()

	cacheServer := "redis:6379"
	cacher := ms.getCacher(cacheServer)
	cacher.Set("key1", "value1", time.Duration(-1))

	// 1. Healthcheck endpoint will register to /healthz
	ms.RegisterLivenessProbeEndpoint("/healthz")

	// 2. Start application services
	startHTTP(ms)

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
