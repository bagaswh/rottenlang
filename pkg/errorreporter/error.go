package errorreporter

import (
	"fmt"
	"os"
)

type ErrorReporter interface {
	ReportScannerError(line, column int, lexeme, message string)
	ReportParserError(line, column int, where, message string)
}

type StderrErrorReporter struct{}

func (e *StderrErrorReporter) ReportScannerError(line, column int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line=%d col=%d] Error: %s. Snippet: '%s'", line, column, message, where)
}

func (e *StderrErrorReporter) ReportParserError(line, column int, where, message string) {
	fmt.Fprintf(os.Stderr, "[line=%d col=%d] Error: %s. Snippet: '%s'", line, column, message, where)
}
