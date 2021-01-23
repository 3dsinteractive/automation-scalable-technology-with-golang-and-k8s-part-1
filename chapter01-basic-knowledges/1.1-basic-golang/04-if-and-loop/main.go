package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	citizenID := "1234"
	password := "Helloworld"
	gender := Female

	// 1. If
	if citizenID == "1234" && password == "Helloworld" {
		println("CitizenID 1234, you are logged in")
	} else if citizenID == "5678" && password == "Helloworld" {
		println("CitizenID 5678, you are logged in")
	} else {
		println("You are not logged in")
	}
	printUnderline()

	// 2. Switch
	switch gender {
	case Female:
		println("You are Female")
	case Male:
		println("You are Male")
	default:
		println("Gender is not specify")
	}
	printUnderline()

	// 3. For i
	for i := 0; i < 10; i++ {
		println(fmt.Sprintf("Loop i=%d", i))
	}
	printUnderline()

	// 4. For range
	countries := []string{
		"Thailand",
		"Japan",
		"China",
		"Korea",
		"Vietnam",
	}
	for i, country := range countries {
		println(fmt.Sprintf("Country %d=%s", i, country))
	}
	// Note: for range without index use underscore at first variable
	// for _, country := range countries {
	// 	println(fmt.Sprintf("Country %s", country))
	// }
	printUnderline()

	// 5. For condition
	i := 0
	for i < 10 {
		i++
		println(fmt.Sprintf("For i=%d", i))
	}
	printUnderline()

	// 6. For Infinite
	n := 0
	for {
		n++
		println(fmt.Sprintf("For n=%d", n))
		if n > 10 {
			break
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
