// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"context"
	"time"

	"github.com/labstack/echo"
)

// GET register service endpoint for HTTP GET
func (ms *Microservice) GET(path string, h ServiceHandleFunc) {
	ms.echo.GET(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// POST register service endpoint for HTTP POST
func (ms *Microservice) POST(path string, h ServiceHandleFunc) {
	ms.echo.POST(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// PUT register service endpoint for HTTP PUT
func (ms *Microservice) PUT(path string, h ServiceHandleFunc) {
	ms.echo.PUT(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// PATCH register service endpoint for HTTP PATCH
func (ms *Microservice) PATCH(path string, h ServiceHandleFunc) {
	ms.echo.PATCH(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// DELETE register service endpoint for HTTP DELETE
func (ms *Microservice) DELETE(path string, h ServiceHandleFunc) {
	ms.echo.DELETE(path, func(c echo.Context) error {
		return h(NewHTTPContext(ms, c))
	})
}

// startHTTP will start HTTP service, this function will block thread
func (ms *Microservice) startHTTP(exitChannel chan bool) error {
	// Caller can exit by sending value to exitChannel
	go func() {
		<-exitChannel
		ms.stopHTTP()
	}()
	return ms.echo.Start(":8080")
}

func (ms *Microservice) stopHTTP() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ms.echo.Shutdown(ctx)
}
