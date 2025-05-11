package parser

import (
	"fmt"

	"github.com/bagaswh/rottenlang/pkg/ast"
)

type GenricParserError struct {
	token          *ast.Token
	where, message string
}

func (err *GenricParserError) Error() string {
	return fmt.Sprintf("Parser error: line=%d col=%d %s: %s", err.token.Line, err.token.Column, err.where, err.message)
}

func NewGenericParserError(token *ast.Token, where, message string) *GenricParserError {
	return &GenricParserError{
		token:   token,
		where:   where,
		message: message,
	}
}
