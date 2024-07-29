package main

import (
	"fmt"
	"os"
)

func Error(line int, message string) {
	ReportError(line, message)
}

func ReportError(line int, message string) {
	output := fmt.Sprintf("[line %d] Error: %s", line, message)
	fmt.Fprintln(os.Stderr, output)
	hadError = true
}
