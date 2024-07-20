package main

import "fmt"

func ReportError(line int, message string) {
	output := fmt.Sprintf("[Line %d] Error: %s", line, message)
	fmt.Println(output)
}
