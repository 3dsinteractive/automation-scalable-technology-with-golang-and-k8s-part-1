// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
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

// Now return time.Now
func (ctx *SchedulerContext) Now() time.Time {
	return time.Now()
}

// Log will log a message
func (ctx *SchedulerContext) Log(message string) {
	fmt.Println("Scheduler: ", message)
}

// Param return parameter by name (empty in scheduler)
func (ctx *SchedulerContext) Param(name string) string {
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
