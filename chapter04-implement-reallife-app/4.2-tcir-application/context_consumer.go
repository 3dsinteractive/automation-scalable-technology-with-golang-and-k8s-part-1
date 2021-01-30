// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// ConsumerContext implement IContext it is context for Consumer
type ConsumerContext struct {
	ms      *Microservice
	message string
}

// NewConsumerContext is the constructor function for ConsumerContext
func NewConsumerContext(ms *Microservice, message string) *ConsumerContext {
	return &ConsumerContext{
		ms:      ms,
		message: message,
	}
}

// Log will log a message
func (ctx *ConsumerContext) Log(message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("Consumer:", fns[len(fns)-1], line, message)
}

// Param return parameter by name (empty in case of Consumer)
func (ctx *ConsumerContext) Param(name string) string {
	return ""
}

// QueryParam return empty in consumer
func (ctx *ConsumerContext) QueryParam(name string) string {
	return ""
}

// ReadInput return message
func (ctx *ConsumerContext) ReadInput() string {
	return ctx.message
}

// ReadInputs return nil in case Consumer
func (ctx *ConsumerContext) ReadInputs() []string {
	return nil
}

// Response return response to client
func (ctx *ConsumerContext) Response(responseCode int, responseData interface{}) {
	return
}

// Now return now
func (ctx *ConsumerContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *ConsumerContext) Cacher(server string) ICacher {
	return ctx.ms.getCacher(server)
}

// Producer return producer
func (ctx *ConsumerContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *ConsumerContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}

// Requester return Requester
func (ctx *ConsumerContext) Requester(baseURL string, timeout time.Duration) IRequester {
	return NewRequester(baseURL, timeout, ctx.ms)
}
