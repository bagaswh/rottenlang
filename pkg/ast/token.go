package ast

type TokenType byte

const (
	TokenLeftParen TokenType = iota
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenLeftBracket
	TokenRightBracket
	TokenComma
	TokenDot
	TokenMinus
	TokenPlus
	TokenSemicolon
	TokenSlash
	TokenStar

	TokenEqual
	TokenEqualEqual
	TokenBang
	TokenBangEqual
	TokenGreater
	TokenGreaterEqual
	TokenLess
	TokenLessEqual

	TokenIdentifier
	TokenString
	TokenNumber

	TokenAnd
	TokenOr
	TokenIf
	TokenElse
	TokenFunc
	TokenFor
	TokenWhile
	TokenReturn
	TokenTrue
	TokenFalse
	TokenNil
	TokenVar
	TokenConst

	TokenComment
	TokenCStyleComment

	TokenEOF
)

var keywords = map[string]TokenType{
	"and":               TokenAnd,
	"or":                TokenOr,
	"chat is this real": TokenIf,
	"else":              TokenElse,
	"cap":               TokenFalse,
	"nocap":             TokenTrue,
	"for":               TokenFor,
	"func":              TokenFunc,
	"nil":               TokenNil,

	// Additional Gen Alpha keywords
	"vibes":           TokenVar,          // For variable declaration
	"purrr":           TokenReturn,       // Return keyword
	"slay":            TokenConst,        // Constant declaration
	"skibidi":         TokenWhile,        // While loop
	"rizz":            TokenAnd,          // Alternative for "and"
	"mid":             TokenLess,         // Less than
	"bussin":          TokenGreater,      // Greater than
	"no_tea_no_shade": TokenEqualEqual,   // Equality check
	"fr_fr":           TokenBangEqual,    // Not equal
	"based":           TokenPlus,         // Addition
	"cringe":          TokenMinus,        // Subtraction
	"yeet":            TokenSlash,        // Division
	"ong":             TokenStar,         // Multiplication
	"lowkey":          TokenLessEqual,    // Less than or equal
	"highkey":         TokenGreaterEqual, // Greater than or equal
	"deadass":         TokenBang,         // Logical NOT
	"iykyk":           TokenLeftBrace,    // Start block
	"periodt":         TokenRightBrace,   // End block
}

type Token struct {
	Type         TokenType
	Str          *string
	Literal      any
	Line, Column int
}

func NewToken(tokenType TokenType, str *string, literal any, line, column int) *Token {
	return &Token{
		Type:    tokenType,
		Str:     str,
		Literal: literal,
		Line:    line,
		Column:  column,
	}
}

func (t *Token) Name() string {
	switch t.Type {
	case TokenLeftParen:
		return "LEFT_PAREN"
	case TokenRightParen:
		return "RIGHT_PAREN"
	case TokenLeftBrace:
		return "LEFT_BRACE"
	case TokenRightBrace:
		return "RIGHT_BRACE"
	case TokenLeftBracket:
		return "LEFT_BRACKET"
	case TokenRightBracket:
		return "RIGHT_BRACKET"
	case TokenComma:
		return "COMMA"
	case TokenDot:
		return "DOT"
	case TokenMinus:
		return "MINUS"
	case TokenPlus:
		return "PLUS"
	case TokenSemicolon:
		return "SEMICOLON"
	case TokenSlash:
		return "SLASH"
	case TokenStar:
		return "STAR"
	case TokenEqual:
		return "EQUAL"
	case TokenEqualEqual:
		return "EQUAL_EQUAL"
	case TokenBang:
		return "BANG"
	case TokenBangEqual:
		return "BANG_EQUAL"
	case TokenGreater:
		return "GREATER"
	case TokenGreaterEqual:
		return "GREATER_EQUAL"
	case TokenLess:
		return "LESS"
	case TokenLessEqual:
		return "LESS_EQUAL"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenString:
		return "STRING"
	case TokenNumber:
		return "NUMBER"
	case TokenAnd:
		return "AND"
	case TokenOr:
		return "OR"
	case TokenIf:
		return "IF"
	case TokenElse:
		return "ELSE"
	case TokenFunc:
		return "FUNC"
	case TokenFor:
		return "FOR"
	case TokenWhile:
		return "WHILE"
	case TokenReturn:
		return "RETURN"
	case TokenTrue:
		return "TRUE"
	case TokenFalse:
		return "FALSE"
	case TokenNil:
		return "NIL"
	case TokenVar:
		return "VAR"
	case TokenConst:
		return "CONST"
	case TokenComment:
		return "COMMENT"
	case TokenCStyleComment:
		return "C_STYLE_COMMENT"
	case TokenEOF:
		return "EOF"
	default:
		return "UNKNOWN"
	}
}
