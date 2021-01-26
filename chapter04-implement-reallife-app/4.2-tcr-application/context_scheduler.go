// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// SchedulerContext implement IContext it is context for Consumer
type SchedulerContext struct {
	ms *Microservice
}

// NewSchedulerContext is the constructor function for SchedulerContext
func NewSchedulerContext(ms *Microservice) *SchedulerContext {
	return &SchedulerContext{
		ms: ms,
	}
}

// Log will log a message
func (ctx *SchedulerContext) Log(message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("Scheduler:", fns[len(fns)-1], line, message)
}

// Param return parameter by name (empty in scheduler)
func (ctx *SchedulerContext) Param(name string) string {
	return ""
}

// QueryParam return empty in scheduler
func (ctx *SchedulerContext) QueryParam(name string) string {
	return ""
}

// ReadInput return message (return empty in scheduler)
func (ctx *SchedulerContext) ReadInput() string {
	return ""
}

// ReadInputs return messages in batch (return nil in scheduler)
func (ctx *SchedulerContext) ReadInputs() []string {
	return nil
}

// Response return response to client
func (ctx *SchedulerContext) Response(responseCode int, responseData interface{}) {
	return
}

// Now return now
func (ctx *SchedulerContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *SchedulerContext) Cacher(server string) ICacher {
	return ctx.ms.getCacher(server)
}

// Producer return producer
func (ctx *SchedulerContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *SchedulerContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}

// Requester return Requester
func (ctx *SchedulerContext) Requester(baseURL string, timeout time.Duration) IRequester {
	return NewRequester(baseURL, timeout, ctx.ms)
}
