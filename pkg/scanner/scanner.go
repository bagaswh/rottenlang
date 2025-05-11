package scanner

import (
	"io"
	"strconv"

	"github.com/bagaswh/rottenlang/pkg/ast"
)

func NewScanner(r io.Reader, readBuffer int) *Scanner {
	scanner := &Scanner{
		r:             r,
		buf:           make([]byte, 0),
		scannerErrors: make(map[int][]*GenericScanError),

		line: 1,
	}
	scanner.read()
	return scanner
}

type Scanner struct {
	r       io.Reader
	buf     []byte
	current int

	// start index of lexeme
	start int
	line  int
	// col is the buffer index of the first character in a line
	col int

	tokens []*ast.Token

	scannerErrors map[int][]*GenericScanError
}

func (s *Scanner) read() error {
	return s.readAll()
}

func (s *Scanner) readAll() error {
	bytes, err := io.ReadAll(s.r)
	if err != nil {
		return err
	}
	s.buf = bytes
	return nil
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= (len(s.buf))
}

func (s *Scanner) linecol() int {
	return s.current - s.col
}

func (s *Scanner) advance() string {
	if s.isAtEnd() {
		return ""
	}

	c := s.buf[s.current]
	s.current++
	return string(c)
}

func (s *Scanner) addToken(tokenType ast.TokenType, literal any) {
	tokenStr := string(s.buf[s.start:s.current])
	// TODO: remove hardcode and make it more elegant
	// if tokenType == TokenString {
	// tokenStr += "\""
	// }
	s.tokens = append(s.tokens, ast.NewToken(tokenType, &tokenStr, literal, s.line, s.linecol()))
}

func (s *Scanner) scanToken() error {
	c := s.advance()
	switch c {
	case "(":
		s.addToken(ast.TokenLeftParen, nil)
	case ")":
		s.addToken(ast.TokenRightParen, nil)
	case "{":
		s.addToken(ast.TokenLeftBrace, nil)
	case "}":
		s.addToken(ast.TokenRightBrace, nil)
	case "[":
		s.addToken(ast.TokenLeftBracket, nil)
	case "]":
		s.addToken(ast.TokenRightBracket, nil)
	case ",":
		s.addToken(ast.TokenComma, nil)
	case ".":
		s.addToken(ast.TokenDot, nil)
	case "-":
		s.addToken(ast.TokenMinus, nil)
	case "+":
		s.addToken(ast.TokenPlus, nil)
	case ";":
		s.addToken(ast.TokenSemicolon, nil)
	case "*":
		s.addToken(ast.TokenStar, nil)
	case "!":
		token := ast.TokenBang
		if s.match("=") {
			token = ast.TokenBangEqual
		}
		s.addToken(token, nil)
	case "=":
		token := ast.TokenEqual
		if s.match("=") {
			token = ast.TokenEqualEqual
		}
		s.addToken(token, nil)
	case "<":
		token := ast.TokenLess
		if s.match("=") {
			token = ast.TokenLessEqual
		}
		s.addToken(token, nil)
	case ">":
		token := ast.TokenGreater
		if s.match("=") {
			token = ast.TokenGreaterEqual
		}
		s.addToken(token, nil)
	case "/":
		if s.match("/") {
			s.comment()
		} else if s.match("*") {
			s.cStyleComment()
		} else {
			s.addToken(ast.TokenSlash, nil)
		}
	case "\"":
		s.string()
	case "\n":
		s.newline()
	case CharEOF:
		return nil
	default:
		if s.isDigit(c) {
			s.number()
		} else {
			s.scanError(unexpectedCharacterError(c))
		}
	}
	return nil
}

func (s *Scanner) scanError(errDesc *ScanErrorDescription) {
	_, ok := s.scannerErrors[s.line]
	if !ok {
		s.scannerErrors[s.line] = make([]*GenericScanError, 0)
	}
	if len(s.scannerErrors[s.line]) > 0 {
		return
	}
	theError := NewGenericScanError(errDesc.Message, errDesc.Class, s.line, s.linecol())
	s.scannerErrors[s.line] = append(s.scannerErrors[s.line], theError)
}

func (s *Scanner) newline() {
	s.line++
	s.col = s.current
}

func (s *Scanner) comment() {
	// comment goes until newline
	for s.peek() != "\n" && !s.isAtEnd() {
		s.advance()
	}

	s.addToken(ast.TokenComment, string(s.buf[s.start:s.current]))
}

func (s *Scanner) cStyleComment() {
	// nested comment is possible
	level := 1
	for {
		c := s.peek()
		if c == "/" && s.ahead() == "*" {
			// nested comment
			level++
			s.advance()
		} else if c == "*" && s.ahead() == "/" {
			level--
			s.advance()
		} else if c == "\n" {
			s.newline()
		}

		s.advance()

		if level <= 0 {
			break
		}
	}

	s.addToken(ast.TokenCStyleComment, string(s.buf[s.start:s.current]))

}

func (s *Scanner) string() {
	strValue := ""
	for {
		if s.isAtEnd() {
			break
		}

		// escape char
		if s.peek() == "\\" {
			strValue += s.peek()
			s.advance()
			strValue += s.peek()
			s.advance()
		}

		if s.peek() == "\"" {
			break
		}

		if s.peek() == "\n" {
			s.newline()
		}

		strValue += s.peek()

		s.advance()
	}

	if s.isAtEnd() {
		s.scanError(ErrUnterminatedString)
	}

	// the closing ""
	s.advance()

	s.addToken(ast.TokenString, string(strValue))
}

func (s *Scanner) isDigit(ch string) bool {
	if len(ch) == 0 {
		return false
	}
	return ch[0] >= 48 && ch[0] <= 57
}

func (s *Scanner) number() {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == "." && s.isDigit(s.ahead()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	num := string(s.buf[s.start:s.current])
	numVal, err := strconv.ParseFloat(num, 64)
	if err != nil {
		s.scanError(invalidNumberLiteralError(num))
	}

	s.addToken(ast.TokenNumber, numVal)
}

func (s *Scanner) peek() string {
	if s.isAtEnd() {
		return CharEOF
	}
	return string(s.buf[s.current])
}

func (s *Scanner) ahead() string {
	if s.isAtEnd() || s.current+1 >= len(s.buf) {
		return CharEOF
	}
	return string(s.buf[s.current+1])
}

func (s *Scanner) match(ch string) bool {
	if s.isAtEnd() {
		return false
	}
	if string(s.buf[s.current]) != ch {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) ScanTokens() ([]*ast.Token, error) {
	tokens := make([]*ast.Token, 0)
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	tokens = append(tokens, ast.NewToken(ast.TokenEOF, strPtr(""), nil, s.line, s.linecol()))

	if len(s.scannerErrors) > 0 {
		return nil, ErrScanner
	}

	return tokens, nil
}

func (s *Scanner) Tokens() []*ast.Token {
	return s.tokens
}

func (s *Scanner) HadError() bool {
	return len(s.ScannerErrors()) > 0
}

func (s *Scanner) ScannerErrors() map[int][]*GenericScanError {
	return s.scannerErrors
}

func strPtr(str string) *string {
	return &str
}
