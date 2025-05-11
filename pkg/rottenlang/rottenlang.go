package rottenlang

import (
	"fmt"
	"strings"

	"github.com/bagaswh/rottenlang/pkg/errorreporter"
	"github.com/bagaswh/rottenlang/pkg/parser"
	"github.com/bagaswh/rottenlang/pkg/printer"
	"github.com/bagaswh/rottenlang/pkg/scanner"
)

type Rottenlang struct {
	Scanner       *scanner.Scanner
	Parser        *parser.Parser
	ErrorReporter errorreporter.ErrorReporter
}

func NewRottenlang(source string, errorReporter errorreporter.ErrorReporter) *Rottenlang {
	r := strings.NewReader(source)
	scanner := scanner.NewScanner(r, 0)
	parser := parser.NewParser(errorReporter)
	return &Rottenlang{
		Scanner:       scanner,
		Parser:        parser,
		ErrorReporter: errorReporter,
	}
}

func (d *Rottenlang) Run(line string) {}

func (d *Rottenlang) Scan() {
	_, err := d.Scanner.ScanTokens()
	if err != nil {
		if err == scanner.ErrScanner {
			for _, lineErrs := range d.Scanner.ScannerErrors() {
				firstErr := lineErrs[0]
				fmt.Println(firstErr.Error())
			}
		}
	}

	if d.Scanner.HadError() {
		return
	}

	tokens := d.Scanner.Tokens()
	currentLine := -1
	start := -1
	astPrinter := printer.NewASTPrinter()
	for i, token := range tokens {
		if start == -1 {
			start = i
		}
		if currentLine != token.Line {
			if currentLine != -1 {
				fmt.Println("---------------")
				d.Parser.SetTokens(tokens[start:i])
				expr := d.Parser.Parse()
				fmt.Printf("%d-%d; line %d: %#s\n", start, i, token.Line, astPrinter.Print(expr))
				start = -1
			}
			currentLine = token.Line
		}
		// fmt.Printf("str='%s' tok=%s line=%d literal=%w\n", *token.Lexeme, token.Name(), token.Line, token.Literal)
	}

}
