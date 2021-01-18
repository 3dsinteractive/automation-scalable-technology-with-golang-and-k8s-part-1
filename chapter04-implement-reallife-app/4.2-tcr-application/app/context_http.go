package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

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
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("HTTP:", fns[len(fns)-1], line, message)
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

// ReadInputs return nil in HTTP Context
func (ctx *HTTPContext) ReadInputs() []string {
	return nil
}

// Response return response to client
func (ctx *HTTPContext) Response(responseCode int, responseData interface{}) {
	ctx.c.JSON(responseCode, responseData)
}

// Now return now
func (ctx *HTTPContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *HTTPContext) Cacher(server string) ICacher {
	return ctx.ms.getCacher(server)
}

// Producer return producer
func (ctx *HTTPContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *HTTPContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}

// Requester return Requester
func (ctx *HTTPContext) Requester(baseURL string, timeout time.Duration) IRequester {
	return NewRequester(baseURL, timeout, ctx.ms)
}
