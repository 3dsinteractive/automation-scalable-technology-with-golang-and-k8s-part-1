package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// 1. Literal types
	// Literal type copy by value
	theString := "This variable type string"
	theInt := 305
	theBool := true
	printString(theString)
	// printString(theInt) // Compile error
	printInt(theInt)
	printBool(theBool)
	printUnderline()

	// 2. interface{} type
	printInterface(theString)
	printInterface(theInt)
	printInterface(theBool)
	printUnderline()

	// 3. map type
	// - Map pass by pointer of map
	theMap := map[string]interface{}{}
	theMap["firstname"] = "Chaiyapong"
	theMap["lastname"] = "Lapliengtrakul"
	theMap["citizen_id"] = "1234"
	printMap(theMap)
	printMapAsJSON(theMap)
	printUnderline()

	// 4. slice type (array)
	// - We should name slice phurally
	// - Slices pass by pointer of slice
	theSlices := []string{}
	theSlices = append(theSlices, "item 1")
	theSlices = append(theSlices, "item 2")
	theSlices = append(theSlices, "item 3")
	printSlice(theSlices)
	printSliceAsJSON(theSlices)
	printUnderline()

	// 5. Struct type
	// - It is my practice to always create struct variable as pointer
	//   to avoid confusion in team member
	// - JSON marshal will use struct tag
	theCitizen := &Citizen{
		Firstname: "Chaiyapong",
		Lastname:  "Lapliengtrakul",
		CitizenID: "1234",
	}
	printCitizen(theCitizen)
	printCitizenAsJSON(theCitizen)
	printUnderline()

	// 6. nil value
	// - Literal cannot be nil
	// map, slice and struct and be nil
	// theString = nil // Error
	// theInt = nil    // Error
	// theBool = nil   // Error
	theMap = nil
	theSlices = nil
	theCitizen = nil
	// theMap["key"] = "value" // Runtime Error
	// theCitizen.Firstname = "Chaiyapong" // Runtime Error
	theSlices = append(theSlices, "value") // OK

	// 7. Check nil or zero using len()
	if len(theMap) == 0 {
		fmt.Println("theMap is nil")
	}
	if len(theSlices) == 0 {
		fmt.Println("theSlices is nil")
	}
	if theCitizen == nil {
		fmt.Println("theCitizen is nil")
	}
	printUnderline()

	// 8. Enum
	theGender := Male
	theGender = Female
	theGender = Unspecify
	switch theGender {
	case Male:
		fmt.Println("Gender is Male")
	case Female:
		fmt.Println("Gender is Female")
	case Unspecify:
		fmt.Println("Gender is Unspecify")
	}
	printGender(theGender)
	printUnderline()
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
