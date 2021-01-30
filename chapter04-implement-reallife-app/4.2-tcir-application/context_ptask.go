// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"
	"time"
)

// PTaskContext implement IContext it is context for ParallelTask
type PTaskContext struct {
	ms          *Microservice
	cacheServer string
	taskID      string
	workerID    string
	input       string
}

// NewPTaskContext is the constructor function for PTaskContext
func NewPTaskContext(ms *Microservice, cacheServer string, taskID string, workerID string, input string) *PTaskContext {
	return &PTaskContext{
		ms:          ms,
		cacheServer: cacheServer,
		taskID:      taskID,
		workerID:    workerID,
		input:       input,
	}
}

// Log will log a message
func (ctx *PTaskContext) Log(message string) {
	_, fn, line, _ := runtime.Caller(1)
	fns := strings.Split(fn, "/")
	fmt.Println("PTASK:", fns[len(fns)-1], line, message)
}

// Param return parameter by name (empty in AsyncTask)
func (ctx *PTaskContext) Param(name string) string {
	return ""
}

// QueryParam return empty in async task
func (ctx *PTaskContext) QueryParam(name string) string {
	return ""
}

// ReadInput return message (return empty in AsyncTask)
func (ctx *PTaskContext) ReadInput() string {
	return ctx.input
}

// ReadInputs return messages in batch (return nil in AsyncTask)
func (ctx *PTaskContext) ReadInputs() []string {
	return nil
}

// Response return response to client
func (ctx *PTaskContext) Response(responseCode int, responseData interface{}) {
	maxLimit := 100
	for true {
		// Just check the limit to prevent infinite loop
		maxLimit--
		if maxLimit < 0 {
			return
		}

		// 1. Get the current task status
		cacher := ctx.Cacher(ctx.cacheServer)
		currentStatusStr, err := cacher.Get(ctx.taskID)
		if err != nil {
			ctx.Log(err.Error())
			return
		}
		currentStatus := map[string]interface{}{}
		if len(currentStatusStr) > 0 {
			err = json.Unmarshal([]byte(currentStatusStr), &currentStatus)
			if err != nil {
				ctx.Log(err.Error())
				return
			}
		}

		// 2. If task is complete, return
		taskStatus, _ := currentStatus["status"].(string)
		if taskStatus == "complete" {
			return
		}
		workers, _ := currentStatus["workers"].([]interface{})
		if len(workers) == 0 {
			ctx.Log("No Workers")
			return
		}

		// 3. Find worker that match ctx, and set the status to complete
		for _, w := range workers {
			worker := w.(map[string]interface{})
			workerID, _ := worker["worker_id"]
			if workerID != ctx.workerID {
				continue
			}

			worker["status"] = "complete"
			worker["response"] = responseData
			worker["code"] = responseCode
			break
		}

		// 4. If all workers has completed, set the task status to complete
		allWorkerComplete := true
		for _, w := range workers {
			worker := w.(map[string]interface{})
			if worker["status"] == "running" {
				allWorkerComplete = false
				break
			}
		}
		if allWorkerComplete {
			currentStatus["status"] = "complete"
		}

		// 5. Before send status update, just check race condition
		if len(currentStatusStr) > 0 {
			hasChanged, err := cacher.HasChanged(ctx.taskID, currentStatusStr)
			if err != nil {
				ctx.Log(err.Error())
				return
			}
			// If race condition happen, just refresh status and try to update again
			if hasChanged {
				continue
			}
		}

		// 6. Save status in cache
		cacher.Set(ctx.taskID, currentStatus, 30*time.Minute)
		break
	}
}

// Now return now
func (ctx *PTaskContext) Now() time.Time {
	return time.Now()
}

// Cacher return cacher
func (ctx *PTaskContext) Cacher(server string) ICacher {
	return ctx.ms.getCacher(server)
}

// Producer return producer
func (ctx *PTaskContext) Producer(servers string) IProducer {
	return ctx.ms.getProducer(servers)
}

// MQ return MQ
func (ctx *PTaskContext) MQ(servers string) IMQ {
	return NewMQ(servers, ctx.ms)
}

// Requester return Requester
func (ctx *PTaskContext) Requester(baseURL string, timeout time.Duration) IRequester {
	return NewRequester(baseURL, timeout, ctx.ms)
}
