// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import "fmt"

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
	fmt.Println("Batch Consumer: ", message)
}

// Param return parameter by name (empty in case of Consumer)
func (ctx *BatchConsumerContext) Param(name string) string {
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
