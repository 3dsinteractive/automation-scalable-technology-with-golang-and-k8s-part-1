package main

import "fmt"

func main() {

	ms := NewMicroservice()
	ms.POST("/member", func(ctx IContext) error {
		fmt.Println("POST /member")
		return nil
	})

	defer ms.Cleanup()
	ms.Start()
}
