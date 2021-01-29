package formatter

import "fmt"

// MyPrintln will print input string
func MyPrintln(input string) {
	myPrintln(input)
}

// myPrintln will not export
func myPrintln(input string) {
	fmt.Println(input)
}
