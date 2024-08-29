package errorHand

import (
	"fmt"
	"os"
)

var HadError = false

func Error(line int, message string) {
	message = "Error: " + message
	ReportError(line, message)
}

func ParseError(token string, line int, message string) {
	message = "Error at " + "'" + token + "': " + message
	ReportError(line, message)
}

func ReportError(line int, message string) {
	output := fmt.Sprintf("[line %d] %s", line, message)
	fmt.Fprintln(os.Stderr, output)
	HadError = true
}
