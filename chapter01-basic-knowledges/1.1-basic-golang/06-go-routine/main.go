package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	exitChannel := make(chan bool)
	workerCount := 10
	for i := 0; i < workerCount; i++ {
		workerID := fmt.Sprintf("worker-%d", i)
		go worker(workerID, exitChannel)
	}

	time.Sleep(10 * time.Second)
	for i := 0; i < workerCount; i++ {
		exitChannel <- true
	}
	close(exitChannel)
	time.Sleep(2 * time.Second)

	println("Main is exited")
}

func worker(workerID string, exitChannel chan bool) {
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
