// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"sync"
	"time"
)

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
