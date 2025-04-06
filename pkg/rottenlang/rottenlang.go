package rottenlang

import (
	"fmt"
	"strings"

	"github.com/bagaswh/rottenlang/pkg/scanner"
)

type Rottenlang struct {
	Scanner       *scanner.Scanner
	ErrorReporter scanner.ErrorReporter
}

func NewRottenlang(source string, errorReporter scanner.ErrorReporter) *Rottenlang {
	r := strings.NewReader(source)
	scanner := scanner.NewScanner(r, 0)
	return &Rottenlang{
		Scanner:       scanner,
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
		// return
	}
	tokens := d.Scanner.Tokens()
	for _, token := range tokens {
		fmt.Printf("str='%s' tok=%s line=%d literal=%w\n", *token.Str, token.Name(), token.Line, token.Literal)
	}
}
