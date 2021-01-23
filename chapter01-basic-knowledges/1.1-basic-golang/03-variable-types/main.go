package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// 1. Literal types
	// Note: Literal type copy by value
	theString := "This variable type string"
	theInt := 305
	theBool := true
	printString(theString)
	printInt(theInt)
	printBool(theBool)
	printUnderline()
	// 2. Type in Golang are static type
	// printString(theInt) // Compile error

	// 3. interface{} type can accept all types
	printInterface(theString)
	printInterface(theInt)
	printInterface(theBool)
	printUnderline()

	// 4. map type
	// Note: Map variable is the pointer of map, pass by pointer of map
	theMap := map[string]interface{}{}
	theMap["firstname"] = "Chaiyapong"
	theMap["lastname"] = "Lapliengtrakul"
	theMap["citizen_id"] = "1234"
	printMap(theMap)
	// updateMap(theMap) // This will update map that pass in to function
	printMapAsJSON(theMap)
	// Note: We can use 2 variables to get map value
	gender, ok := theMap["gender"]
	if ok {
		fmt.Println("gender is ", gender.(string)) // Casting
	} else {
		fmt.Println("No Gender specify")
	}

	printUnderline()

	// 5. slice type (dynamic array)
	// Note: We should name slice phurally
	// Note: Slices pass by pointer of slice (just like map)
	theSlices := []string{}
	theSlices = append(theSlices, "item 1")
	theSlices = append(theSlices, "item 2")
	theSlices = append(theSlices, "item 3")
	printSlice(theSlices)
	// updateSliceIndex0(theSlices) // This will update value of index 0 of slice
	printSliceAsJSON(theSlices)
	printUnderline()

	// 6. Struct type
	// Note: It is my practice to always create struct variable as pointer (*Type)
	//   to avoid confusion in team member, so will can always assume that strut is a pointer
	theCitizen := &Citizen{
		Firstname: "Chaiyapong",
		Lastname:  "Lapliengtrakul",
		CitizenID: "1234",
	}
	printCitizen(theCitizen)
	// Note: JSON marshal will use struct tag
	printCitizenAsJSON(theCitizen)
	printUnderline()

	// 7. nil value
	// Note: Literal cannot be nil (string, int, bool, ..)
	// theString = nil // Error
	// theInt = nil    // Error
	// theBool = nil   // Error

	// Note: map, slice and struct and be nil
	theMap = nil
	theSlices = nil
	theCitizen = nil

	// Note: assign value to nil map and struct will cause runtime error
	// theMap["key"] = "value" // Runtime Error
	// theCitizen.Firstname = "Chaiyapong" // Runtime Error

	// Note: append value to nil slice is OK
	theSlices = append(theSlices, "value") // OK

	// 8. Check nil or zero using len()
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

	// 9. Enum is just constant in golang
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

func updateMap(input map[string]interface{}) {
	input["updated"] = true
}

func printMapAsJSON(input map[string]interface{}) {
	js, _ := json.Marshal(input)
	fmt.Println("map JSON = ", string(js))
}

func printSlice(input []string) {
	fmt.Println("slice = ", input)
}

func updateSliceIndex0(input []string) {
	input[0] = "Item 0"
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
