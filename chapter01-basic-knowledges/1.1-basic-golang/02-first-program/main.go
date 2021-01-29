// 1. Package name must be main for main package
package main

// 2. Import section is the list of package dependency
import (
	"automationworkshop/main/formatter"
	"fmt"
)

// 3. main() is the function where application start
func main() {
	fmt.Println("My First Program")

	formatter.MyPrintln("This is from MyPrintln")
	// This function is not export
	// formatter.myPrintln("This is from MyPrintln")
}
