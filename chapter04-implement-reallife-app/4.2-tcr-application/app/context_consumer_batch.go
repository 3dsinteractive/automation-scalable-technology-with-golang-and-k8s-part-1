// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// BatchConsumerContext implement IContext it is context for Consumer
type BatchConsumerContext struct {
	ms       *Microservice
	messages []string
}

// NewBatchConsumerContext is the constructor function for BatchConsumerContext
func NewBatchConsumerContext(ms *Microservice, messages []string) *BatchConsumerContext {
	return &BatchConsumerContext{
		ms:       ms,
		messages: messages,
	}
}

// Log will log a message
func (ctx *BatchConsumerContext) Log(message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("Batch Consumer:", fns[len(fns)-1], line, message)
}

// Param return parameter by name (empty in case of Consumer)
func (ctx *BatchConsumerContext) Param(name string) string {
	return ""
}

// QueryParam return empty in consumer batch
func (ctx *BatchConsumerContext) QueryParam(name string) string {
	return ""
}

// ReadInput return message (return empty in batch consumer)
func (ctx *BatchConsumerContext) ReadInput() string {
	return ""
}

// ReadInputs return messages in batch
func (ctx *BatchConsumerContext) ReadInputs() []string {
	return ctx.messages
}

// Response return response to client
func (ctx *BatchConsumerContext) Response(responseCode int, responseData interface{}) {
	return
}

// Now return now
func (ctx *BatchConsumerContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *BatchConsumerContext) Cacher(server string) ICacher {
	return ctx.ms.getCacher(server)
}

// Producer return producer
func (ctx *BatchConsumerContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *BatchConsumerContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}

// Requester return Requester
func (ctx *BatchConsumerContext) Requester(baseURL string, timeout time.Duration) IRequester {
	return NewRequester(baseURL, timeout, ctx.ms)
}
