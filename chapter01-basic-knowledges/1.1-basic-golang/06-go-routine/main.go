package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {

	workerCount := 10

	// 1. Use channel to receive data from go routine
	//    Send data from go routine back to outside
	responseChannel := make(chan string)
	for i := 0; i < workerCount; i++ {
		workerID := fmt.Sprintf("worker-%d", i)
		go worker1(workerID, responseChannel)
	}

	for i := 0; i < workerCount; i++ {
		res := <-responseChannel
		println(res)
	}

	close(responseChannel)
	println("All response returned")

	// 2. Use channel to signal go module to exit
	//    Send data from outside go model into module
	exitChannel := make(chan bool)
	for i := 0; i < workerCount; i++ {
		workerID := fmt.Sprintf("worker-%d", i)
		go worker3(workerID, exitChannel)
	}

	time.Sleep(10 * time.Second)
	for i := 0; i < workerCount; i++ {
		exitChannel <- true
	}
	close(exitChannel)
	time.Sleep(2 * time.Second)
	println("Main is exited")
}

func worker1(workerID string, responseChannel chan string) {
	// Simulate request latency
	time.Sleep(1 * time.Second)
	responseChannel <- (workerID + " Response")
}

func worker3(workerID string, exitChannel chan bool) {
	i := 0
	for true {
		i++

		select {
		case <-exitChannel:
			println(fmt.Sprintf("Worker=%s has exited", workerID))
			return
		default:
			println(fmt.Sprintf("Worker=%s, Counter=%d", workerID, i))
			time.Sleep(1 * time.Second)
		}
	}
}

func printUnderline() {
	fmt.Println("---")
}

func printString(input string) {
	fmt.Println("string = ", input)
}

func printInt(input int) {
	fmt.Println("int = ", input)
}

func printBool(input bool) {
	fmt.Println("bool = ", input)
}

func printInterface(input interface{}) {
	fmt.Println("interface = ", input)
}

func printMap(input map[string]interface{}) {
	fmt.Println("map = ", input)
}

func printMapAsJSON(input map[string]interface{}) {
	js, _ := json.Marshal(input)
	fmt.Println("map JSON = ", string(js))
}

func printSlice(input []string) {
	fmt.Println("slice = ", input)
}

func printSliceAsJSON(input []string) {
	js, _ := json.Marshal(input)
	fmt.Println("slice JSON = ", string(js))
}

// Citizen is type represent person in country
type Citizen struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	CitizenID string `json:"citizen_id"`
}

func printCitizen(input *Citizen) {
	fmt.Println("Citizen = ", input)
}

func printCitizenAsJSON(input *Citizen) {
	js, _ := json.Marshal(input)
	fmt.Println("Citizen JSON = ", string(js))
}

// GenderType is enum for Gender
type GenderType string

const (
	// Unspecify is gender type for Unspecify
	Unspecify GenderType = "UNSPECIFY"
	// Male is gender type for Male
	Male GenderType = "MALE"
	// Female is gender type for Female
	Female GenderType = "FEMALE"
)

func printGender(input GenderType) {
	fmt.Println("Gender = ", input)
}
