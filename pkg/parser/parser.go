package parser

import (
	"errors"
	"fmt"

	"github.com/bagaswh/rottenlang/pkg/ast"
	"github.com/bagaswh/rottenlang/pkg/errorreporter"
)

type Parser struct {
	tokens        []*ast.Token
	current       int
	errorReporter errorreporter.ErrorReporter
}

func NewParser(errorReporter errorreporter.ErrorReporter) *Parser {
	return &Parser{
		errorReporter: errorReporter,
	}
}

func (p *Parser) SetTokens(tokens []*ast.Token) {
	p.tokens = tokens
	p.Reset()
}

func (p *Parser) Reset() {
	p.current = 0
}

func (p *Parser) Parse() ast.Expr {
	if p.tokens == nil {
		panic(errors.New("tokens is nil"))
	}

	// defer func() {
	// 	r := recover()
	// 	log.Fatalln(r)
	// 	switch r.(type) {
	// 	case GenricParserError:
	// 		p.synchronize()
	// 		return
	// 		// default:
	// 		// 	panic(r)
	// 	}
	// }()
	return p.expression()
}

func (p *Parser) synchronize() {}

func (p *Parser) expression() ast.Expr {
	return p.equality()
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()
	for p.match(ast.TokenBangEqual, ast.TokenEqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}
	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(ast.TokenGreater, ast.TokenGreaterEqual, ast.TokenLess, ast.TokenLessEqual) {
		operator := p.previous()
		right := p.term()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(ast.TokenMinus, ast.TokenPlus) {
		operator := p.previous()
		right := p.factor()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(ast.TokenSlash, ast.TokenStar) {
		operator := p.previous()
		right := p.unary()
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	fmt.Println("unary(): token type:", p.peek().Name(), "current", p.current)
	for p.match(ast.TokenBang, ast.TokenMinus, ast.TokenPlus) {
		operator := p.previous()
		fmt.Println("unary():  operator:", operator)
		right := p.unary()
		return ast.NewUnaryExpr(operator, right)
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(ast.TokenFalse) {
		return ast.NewLiteralExpr(false)
	}
	if p.match(ast.TokenTrue) {
		return ast.NewLiteralExpr(true)
	}
	if p.match(ast.TokenNil) {
		return ast.NewLiteralExpr(nil)
	}

	if p.match(ast.TokenNumber) {
		fmt.Println("match token number at line", p.peek().Line)
		return ast.NewLiteralExpr(p.previous().Literal)
	}

	if p.match(ast.TokenLeftParen) {
		expr := p.expression()
		p.consume(ast.TokenRightParen, "Expect ')' after expression")
		return ast.NewGroupingExpr(expr)
	}

	return ast.NewLiteralExpr(p.peek().Literal)
}

func (p *Parser) consume(tokenType ast.TokenType, errorMessageWhenNotMatched string) *ast.Token {
	if p.check(tokenType) {
		return p.advance()
	}
	err := p.error(p.peek(), errorMessageWhenNotMatched)
	fmt.Println(err)
	panic(err)
}

func (p *Parser) error(token *ast.Token, message string) *GenricParserError {
	var err *GenricParserError
	if token.Type == ast.TokenEOF {
		err = NewGenericParserError(token, " at end", message)
	} else {
		err = NewGenericParserError(token, fmt.Sprintf("at '%s'", *token.Lexeme), message)
	}
	p.errorReporter.ReportParserError(err.token.Line, err.token.Column, err.where, err.message)
	return err
}

// func (p *Parser) report(tokenType ast.TokenType, )

func (p *Parser) peek() *ast.Token {
	if p.current < len(p.tokens) {
		return p.tokens[p.current]
	}
	return ast.EOF
}

func (p *Parser) match(tokens ...ast.TokenType) bool {
	for _, token := range tokens {
		if p.check(token) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(tokenType ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	curr := p.peek()
	return curr.Type == tokenType
}

func (p *Parser) previous() *ast.Token {
	if p.current <= 0 {
		return p.tokens[0]
	}

	return p.tokens[p.current-1]
}

func (p *Parser) advance() *ast.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == ast.TokenEOF
}
