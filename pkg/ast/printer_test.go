package ast

import (
	"testing"

	"github.com/bagaswh/rottenlang/pkg/types"
)

func TestPrint_BinaryExpr(t *testing.T) {
	expr := BinaryExpr{
		left:     NewLiteralExpr(1),
		right:    NewLiteralExpr(2),
		operator: NewToken(TokenPlus, types.StrPtr("+"), nil, 0, 0),
	}
	astPrinter := NewASTPrinter()
	result := astPrinter.VisitBinaryExpr(&expr).(string)
	if result != "1 + 2" {
		t.Errorf("got %s, want %s", result, "1 + 2")
	}
}
