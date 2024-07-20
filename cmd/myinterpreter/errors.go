package main

import (
	"fmt"
	"os"
)

func ReportError(line int, message string) {
	output := fmt.Sprintf("[line %d] Error: %s", line, message)
	fmt.Fprintln(os.Stderr, output)
}
