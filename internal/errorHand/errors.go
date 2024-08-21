package errorHand

import (
	"fmt"
	"os"
)

var HadError = false

func Error(line int, message string) {
	ReportError(line, message)
}

func ParseError(line int, message string) {
	ReportError(line, message)
}

func ReportError(line int, message string) {
	output := fmt.Sprintf("[line %d] Error: %s", line, message)
	fmt.Fprintln(os.Stderr, output)
	HadError = true
}
