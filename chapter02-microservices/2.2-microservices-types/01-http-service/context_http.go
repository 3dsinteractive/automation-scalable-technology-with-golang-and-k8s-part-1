package main

import (
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
)

// HTTPContext implement IContext it is context for HTTP
type HTTPContext struct {
	ms *Microservice
	c  echo.Context
}

// NewHTTPContext is the constructor function for HTTPContext
func NewHTTPContext(ms *Microservice, c echo.Context) *HTTPContext {
	return &HTTPContext{
		ms: ms,
		c:  c,
	}
}

// Log will log a message
func (ctx *HTTPContext) Log(message string) {
	fmt.Println("HTTP: ", message)
}

// Param return parameter by name
func (ctx *HTTPContext) Param(name string) string {
	return ctx.c.Param(name)
}

// QueryParam return query param
func (ctx *HTTPContext) QueryParam(name string) string {
	return ctx.c.QueryParam(name)
}

// ReadInput read the request body and return it as string
func (ctx *HTTPContext) ReadInput() string {
	body, err := ioutil.ReadAll(ctx.c.Request().Body)
	if err != nil {
		return ""
	}
	return string(body)
}

// Response return response to client
func (ctx *HTTPContext) Response(responseCode int, responseData interface{}) {
	ctx.c.JSON(responseCode, responseData)
}
