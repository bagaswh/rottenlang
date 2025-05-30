package ast

type Visitor interface {
	VisitBinaryExpr(expr *BinaryExpr) any
	VisitUnaryExpr(expr *UnaryExpr) any
	VisitLiteralExpr(expr *LiteralExpr) any
	VisitGroupingExpr(expr *GroupingExpr) any
}

type Expr interface {
	Accept(visitor Visitor) any
}

// BinaryExpr

type BinaryExpr struct {
	left     Expr
	operator *Token
	right    Expr
}

func (e *BinaryExpr) Accept(visitor Visitor) any {
	return visitor.VisitBinaryExpr(e)
}

func (e *BinaryExpr) Left() Expr {
	return e.left
}

func (e *BinaryExpr) Operator() *Token {
	return e.operator
}

func (e *BinaryExpr) Right() Expr {
	return e.right
}

func NewBinaryExpr(left Expr, operator *Token, right Expr) *BinaryExpr {
	return &BinaryExpr{
		left:     left,
		operator: operator,
		right:    right,
	}
}

// UnaryExpr

type UnaryExpr struct {
	operator *Token
	right    Expr
}

func (e *UnaryExpr) Accept(visitor Visitor) any {
	return visitor.VisitUnaryExpr(e)
}

func (e *UnaryExpr) Operator() *Token {
	return e.operator
}

func (e *UnaryExpr) Right() Expr {
	return e.right
}

func NewUnaryExpr(operator *Token, right Expr) *UnaryExpr {
	return &UnaryExpr{
		operator: operator,
		right:    right,
	}
}

// LiteralExpr

type LiteralExpr struct {
	value any
}

func (e *LiteralExpr) Accept(visitor Visitor) any {
	return visitor.VisitLiteralExpr(e)
}

func (e *LiteralExpr) Value() any {
	return e.value
}

func NewLiteralExpr(value any) *LiteralExpr {
	return &LiteralExpr{
		value: value,
	}
}

// GroupingExpr
type GroupingExpr struct {
	expr Expr
}

func (e *GroupingExpr) Expr() Expr {
	return e.expr
}

func (e *GroupingExpr) Accept(visitor Visitor) any {
	return visitor.VisitGroupingExpr(e)
}

func NewGroupingExpr(expr Expr) *GroupingExpr {
	return &GroupingExpr{
		expr: expr,
	}
}
