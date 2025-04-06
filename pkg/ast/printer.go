package ast

import "fmt"

type ASTPrinter struct{}

func NewASTPrinter() *ASTPrinter {
	return &ASTPrinter{}
}

func (p *ASTPrinter) Print(expr Expr) string {
	return expr.Accept(p).(string)
}

func (p *ASTPrinter) VisitBinaryExpr(expr *BinaryExpr) any {
	left := expr.Left().Accept(p).(string)
	right := expr.Right().Accept(p).(string)
	return left + " " + *expr.operator.Str + " " + right
}

func (p *ASTPrinter) VisitLiteralExpr(expr *LiteralExpr) any {
	return fmt.Sprintf("%v", expr.Value())
}

func (p *ASTPrinter) VisitGroupingExpr(expr *GroupingExpr) any {
	return fmt.Sprintf("(%v)", expr.expr.Accept(p).(string))
}
