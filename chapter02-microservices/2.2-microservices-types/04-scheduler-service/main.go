// Create and maintain by Chaiyapong Lapliengtrakul (chaiyapong@3dsinteractive.com), All right reserved (2021 - Present)
package main

import (
	"fmt"
	"time"
)

func main() {
	ms := NewMicroservice()

	timer := 1 * time.Second
	exitScheduler := ms.Schedule(timer, func(ctx IContext) error {
		now := ctx.Now()
		ctx.Log(fmt.Sprintf("Tick at %s", now.Format("15:04:05")))
		return nil
	})

	defer func() { exitScheduler <- true }()
	defer ms.Cleanup()

	ms.Start()
}
