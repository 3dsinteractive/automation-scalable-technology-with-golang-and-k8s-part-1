// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// IMicroservice is interface for centralized service management
type IMicroservice interface {
	Start() error
	Stop()
	Cleanup() error
	Log(message string)

	// Scheduler Services
	Schedule(timer time.Duration, h ServiceHandleFunc) chan bool /*exit channel*/
}

// Microservice is the centralized service management
type Microservice struct {
	exitChannel chan bool
}

// ServiceHandleFunc is the handler for each Microservice
type ServiceHandleFunc func(ctx IContext) error

// NewMicroservice is the constructor function of Microservice
func NewMicroservice() *Microservice {
	return &Microservice{}
}

// Start start all registered services
func (ms *Microservice) Start() error {
	// There are 2 ways to exit from Microservices
	// 1. The SigTerm can be send from outside program such as from k8s
	// 2. Send true to ms.exitChannel
	osQuit := make(chan os.Signal, 1)
	ms.exitChannel = make(chan bool, 1)
	signal.Notify(osQuit, syscall.SIGTERM, syscall.SIGINT)
	exit := false
	for {
		if exit {
			break
		}
		select {
		case <-osQuit:
			// if exitHTTP != nil {
			// 	exitHTTP <- true
			// }
			exit = true
		case <-ms.exitChannel:
			// if exitHTTP != nil {
			// 	exitHTTP <- true
			// }
			exit = true
		}
	}

	return nil
}

// Stop stop the services
func (ms *Microservice) Stop() {
	if ms.exitChannel == nil {
		return
	}
	ms.exitChannel <- true
}

// Cleanup clean resources up from every registered services before exit
func (ms *Microservice) Cleanup() error {
	return nil
}

// Log log message to console
func (ms *Microservice) Log(tag string, message string) {
	fmt.Println(tag+": ", message)
}

// Schedule will run handler at timer period
func (ms *Microservice) Schedule(timer time.Duration, h ServiceHandleFunc) chan bool /*exit channel*/ {

	// exitChan must be call exitChan <- true from caller to exit scheduler
	exitChan := make(chan bool, 1)
	go func() {
		t := time.NewTicker(timer)
		done := make(chan bool, 1)
		isExit := false
		isExitMutex := sync.Mutex{}

		go func() {
			<-exitChan
			isExitMutex.Lock()
			isExit = true
			isExitMutex.Unlock()
			// Stop Tick() and send done message to exit for loop below
			// Ref: From the documentation http://golang.org/pkg/time/#Ticker.Stop
			// Stop turns off a ticker. After Stop, no more ticks will be sent.
			// Stop does not close the channel, to prevent a read from the channel succeeding incorrectly.
			t.Stop()
			done <- true
		}()

		for {
			select {
			case execTime := <-t.C:
				isExitMutex.Lock()
				if isExit {
					isExitMutex.Unlock()
					// Done in the next round
					continue
				}
				isExitMutex.Unlock()

				now := time.Now()
				// The schedule that older than 10s, will be skip, because t.C is buffer at size 1
				diff := now.Sub(execTime).Seconds()
				if diff > 10 {
					continue
				}
				h(NewSchedulerContext(ms))
			case <-done:
				return
			}
		}
	}()

	return exitChan
}
