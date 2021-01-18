// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// AsyncTaskContext implement IContext it is context for Consumer
type AsyncTaskContext struct {
	ms          *Microservice
	cacheServer string
	ref         string
	input       string
}

// NewAsyncTaskContext is the constructor function for AsyncTaskContext
func NewAsyncTaskContext(ms *Microservice, cacheServer string, ref string, input string) *AsyncTaskContext {
	return &AsyncTaskContext{
		ms:          ms,
		cacheServer: cacheServer,
		ref:         ref,
		input:       input,
	}
}

// Log will log a message
func (ctx *AsyncTaskContext) Log(message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("ATASK:", fns[len(fns)-1], line, message)
}

// Param return parameter by name (empty in AsyncTask)
func (ctx *AsyncTaskContext) Param(name string) string {
	return ""
}

// QueryParam return empty in async task
func (ctx *AsyncTaskContext) QueryParam(name string) string {
	return ""
}

// ReadInput return message (return empty in AsyncTask)
func (ctx *AsyncTaskContext) ReadInput() string {
	return ctx.input
}

// ReadInputs return messages in batch (return nil in AsyncTask)
func (ctx *AsyncTaskContext) ReadInputs() []string {
	return nil
}

// Response return response to client
func (ctx *AsyncTaskContext) Response(responseCode int, responseData interface{}) {
	cacher := ctx.Cacher(ctx.cacheServer)
	res := map[string]interface{}{
		"status": "success",
		"code":   responseCode,
		"data":   responseData,
	}
	cacher.Set(ctx.ref, res, 30*time.Minute)
}

// Now return now
func (ctx *AsyncTaskContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *AsyncTaskContext) Cacher(server string) ICacher {
	return ctx.ms.getCacher(server)
}

// Producer return producer
func (ctx *AsyncTaskContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *AsyncTaskContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}

// Requester return Requester
func (ctx *AsyncTaskContext) Requester(baseURL string, timeout time.Duration) IRequester {
	return NewRequester(baseURL, timeout, ctx.ms)
}
