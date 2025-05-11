package printer

import (
	"fmt"
	"strconv"

	"github.com/bagaswh/rottenlang/pkg/ast"
)

type ASTPrinter struct{}

func NewASTPrinter() *ASTPrinter {
	return &ASTPrinter{}
}

func (p *ASTPrinter) Print(expr ast.Expr) string {
	return expr.Accept(p).(string)
}

func (p *ASTPrinter) VisitBinaryExpr(expr *ast.BinaryExpr) any {
	left := expr.Left().Accept(p).(string)
	right := expr.Right().Accept(p).(string)
	return left + " " + *expr.Operator().Lexeme + " " + right
}

func (p *ASTPrinter) VisitLiteralExpr(expr *ast.LiteralExpr) any {
	s := ""
	v := expr.Value()
	switch theV := v.(type) {
	case float64:
		s = strconv.FormatFloat(theV, 'f', 6, 64)
	case string:
		s = theV
	}
	return fmt.Sprintf("%v", s)
}

func (p *ASTPrinter) VisitGroupingExpr(expr *ast.GroupingExpr) any {
	return fmt.Sprintf("(%v)", expr.Expr().Accept(p).(string))
}

func (p *ASTPrinter) VisitUnaryExpr(expr *ast.UnaryExpr) any {
	return fmt.Sprintf("%s%v", *expr.Operator().Lexeme, expr.Operator().Lexeme)
}
