// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// ptaskWorker register worker node for ParallelTask
func (ms *Microservice) ptaskWorkerNode(path string, cacheServer string, mqServers string, h ServiceHandleFunc) {
	topic := escapeName("ptask", path)
	mq := NewMQ(mqServers, ms)
	err := mq.CreateTopicR(topic, 5, 1, time.Hour*24*30)
	if err != nil {
		ms.Log("PTASK", err.Error())
		return
	}

	ms.Consume(mqServers, topic, "ptask", -1, func(ctx IContext) error {
		message := map[string]interface{}{}
		err := json.Unmarshal([]byte(ctx.ReadInput()), &message)
		if err != nil {
			ms.Log("PTASK", err.Error())
			return err
		}
		taskID, _ := message["task_id"].(string)
		workerID, _ := message["worker_id"].(string)
		input, _ := message["input"].(string)
		return h(NewPTaskContext(ms, cacheServer, taskID, workerID, input))
	})
}

// PTaskWorkerNode register worker node for ParallelTask
func (ms *Microservice) PTaskWorkerNode(path string, cacheServer string, mqServers string, h ServiceHandleFunc) {
	go ms.ptaskWorkerNode(path, cacheServer, mqServers, h)
}

func (ms *Microservice) handlePTaskPOST(path string, cacheServer string, mqServers string, ctx IContext) error {
	topic := escapeName("ptask", path)

	// 1. Read Input
	input := ctx.ReadInput()
	taskIDParam := ctx.QueryParam("task_id")
	workerCountStr := ctx.QueryParam("worker_count")

	if len(taskIDParam) == 0 {
		return fmt.Errorf("task_id in query param is required")
	}

	// 2. Get status of current task
	// - If it is running, then return
	// - If it is not running, then start task
	taskID := "ptask-" + taskIDParam
	cacher := ctx.Cacher(cacheServer)
	statusStr, err := cacher.Get(taskID)
	if err != nil {
		ms.Log("PTASK", err.Error())
		return err
	}
	status := map[string]interface{}{}
	if len(statusStr) != 0 {
		err = json.Unmarshal([]byte(statusStr), &status)
		if err != nil {
			ms.Log("PTASK", err.Error())
			return err
		}
		taskStatus, _ := status["status"].(string)
		if taskStatus == "running" {
			return nil
		}
	}

	// 3. Create new task status and save in cache
	workerCount, err := strconv.Atoi(workerCountStr)
	if err != nil {
		workerCount = 3 // default workers size
	}
	workers := []map[string]interface{}{}
	messages := []map[string]interface{}{}
	for i := 0; i < workerCount; i++ {
		workerID := taskID + "-" + randString()
		message := map[string]interface{}{
			"task_id":   taskID,
			"worker_id": workerID,
			"input":     input,
		}
		messages = append(messages, message)

		workers = append(workers, map[string]interface{}{
			"status":    "running",
			"worker_id": workerID,
			"code":      "",
			"response":  "",
			"error":     "",
		})
	}
	status = map[string]interface{}{
		"status":  "running",
		"workers": workers,
	}

	// 4. Set Status in Cache
	expire := time.Minute * 30
	cacher.Set(taskID, status, expire)

	// 5. Send message to start ptask
	prod := ctx.Producer(mqServers)
	for _, message := range messages {
		prod.SendMessage(topic, "", message)
	}

	// 6. Response task_id
	res := map[string]string{
		"task_id": taskIDParam,
	}
	ctx.Response(http.StatusOK, res)
	return nil
}

func (ms *Microservice) handlePTaskGET(path string, cacheServer string, mqServers string, ctx IContext) error {

	// 1. Read Input
	taskIDParam := ctx.QueryParam("task_id")

	if len(taskIDParam) == 0 {
		return fmt.Errorf("task_id in query param is required")
	}

	// 2. Get status of current task
	// - If it is running, then return
	// - If it is not running, then start task
	taskID := "ptask-" + taskIDParam
	cacher := ctx.Cacher(cacheServer)
	statusStr, err := cacher.Get(taskID)
	if err != nil {
		ms.Log("PTASK", err.Error())
		return err
	}
	status := map[string]interface{}{}
	if len(statusStr) > 0 {
		err = json.Unmarshal([]byte(statusStr), &status)
		if err != nil {
			ms.Log("PTASK", err.Error())
			return err
		}
	}

	ctx.Response(http.StatusOK, status)
	return nil
}

// PTask register handler to start/stop/status ParallelTask
func (ms *Microservice) PTaskEndpoint(path string, cacheServer string, mqServers string) {
	// Start PTask
	ms.POST(path, func(ctx IContext) error {
		return ms.handlePTaskPOST(path, cacheServer, mqServers, ctx)
	})
	// Get PTask Status
	ms.GET(path, func(ctx IContext) error {
		return ms.handlePTaskGET(path, cacheServer, mqServers, ctx)
	})
}
