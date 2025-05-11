package scanner

import (
	"errors"
	"fmt"
)

var (
	ErrScanner = errors.New("scan error")
)

var (
	ErrUnterminatedString        = &ScanErrorDescription{Message: "unterminated string", Class: ErrClassUnterminatedString}
	ErrUnterminatedNumberLiteral = &ScanErrorDescription{Message: "unterminated number literal", Class: ErrClassUnterminatedNumberLiteral}
)

var (
	ErrClassUnexpectedCharacter       = "UnexpectedCharacter"
	ErrClassInvalidNumberLiteral      = "InvalidNumberLiteral"
	ErrClassUnterminatedString        = "UnterminatedString"
	ErrClassUnterminatedNumberLiteral = "UnterminatedNumberLiteral"
)

type ScanError interface {
	Error() string
	Class() string
}

const (
	CharEOF = "\\0"
)

type GenericScanError struct {
	message   string
	class     string
	line, col int
}

func NewGenericScanError(message, class string, line, col int) *GenericScanError {
	return &GenericScanError{
		message: message,
		class:   class,
		line:    line,
		col:     col,
	}
}

func (err GenericScanError) Error() string {
	return fmt.Sprintf("error at line %d, column %d: %s", err.line, err.col, err.message)
}

func (err GenericScanError) Class() string {
	return err.class
}

type ScanErrorDescription struct {
	Message string
	Class   string
}

func unexpectedCharacterError(ch string) *ScanErrorDescription {
	return &ScanErrorDescription{
		Message: fmt.Sprintf("unexpected character '%s'", ch),
		Class:   ErrClassUnexpectedCharacter,
	}
}

func invalidNumberLiteralError(num string) *ScanErrorDescription {
	return &ScanErrorDescription{
		Message: fmt.Sprintf("invalid number literal '%s'", num),
		Class:   ErrClassInvalidNumberLiteral,
	}
}
