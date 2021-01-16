// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import "fmt"

// AsyncTaskContext implement IContext it is context for Consumer
type AsyncTaskContext struct {
	ms *Microservice
}

// NewAsyncTaskContext is the constructor function for AsyncTaskContext
func NewAsyncTaskContext(ms *Microservice) *AsyncTaskContext {
	return &AsyncTaskContext{
		ms: ms,
	}
}

// Log will log a message
func (ctx *AsyncTaskContext) Log(message string) {
	fmt.Println("AsyncTask: ", message)
}

// Param return parameter by name (empty in AsyncTask)
func (ctx *AsyncTaskContext) Param(name string) string {
	return ""
}

// ReadInput return message (return empty in AsyncTask)
func (ctx *AsyncTaskContext) ReadInput() string {
	return ""
}

// ReadInputs return messages in batch (return nil in AsyncTask)
func (ctx *AsyncTaskContext) ReadInputs() []string {
	return nil
}

// Response return response to client
func (ctx *AsyncTaskContext) Response(responseCode int, responseData interface{}) {
	return
}

// Cacher return cacher
func (ctx *AsyncTaskContext) Cacher(server string) ICacher {
	return NewCacher(server)
}

// Producer return producer
func (ctx *AsyncTaskContext) Producer(servers string) IProducer {
	return NewProducer(servers, ctx.ms)
}
