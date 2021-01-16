package main

import (
	"github.com/labstack/echo"
)

// IMicroservice is interface for centralized service management
type IMicroservice interface {
	Start() error
	Cleanup() error

	// HTTP Services
	GET(path string, h ServiceHandleFunc)
	POST(path string, h ServiceHandleFunc)
	PUT(path string, h ServiceHandleFunc)
	PATCH(path string, h ServiceHandleFunc)
	DELETE(path string, h ServiceHandleFunc)
}

// Microservice is the centralized service management
type Microservice struct {
	echo *echo.Echo
}

// ServiceHandleFunc is the handler for each Microservice
type ServiceHandleFunc func(ctx IContext) error

// NewMicroservice is the constructor function of Microservice
func NewMicroservice() *Microservice {
	return &Microservice{
		echo: echo.New(),
	}
}

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

func (ms *Microservice) startHTTP() error {
	return ms.echo.Start(":8080")
}

// Start start all registered services
func (ms *Microservice) Start() error {
	// Start HTTP Services
	err := ms.startHTTP()
	return err
}

// Cleanup clean resources up from every registered services before exit
func (ms *Microservice) Cleanup() error {
	return nil
}
