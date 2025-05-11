package printer

import (
	"testing"

	"github.com/bagaswh/rottenlang/pkg/ast"
	"github.com/bagaswh/rottenlang/pkg/types"
)

func TestPrint_BinaryExpr(t *testing.T) {
	expr := ast.NewBinaryExpr(ast.NewLiteralExpr(1), ast.NewToken(ast.TokenPlus, types.StrPtr("+"), nil, 0, 0), ast.NewLiteralExpr(2))
	astPrinter := NewASTPrinter()
	result := astPrinter.VisitBinaryExpr(expr).(string)
	if result != "1 + 2" {
		t.Errorf("got %s, want %s", result, "1 + 2")
	}
}

// func TestPrint_GroupingExpr(t *testing.T) {
// 	expr := GroupingExpr{
// 		expr: NewBinaryExpr(
// 			NewBinaryExpr(
// 				NewLiteralExpr(9), NewToken(TokenStar, types.StrPtr("*"), nil, 0, 0), NewLiteralExpr(3),
// 			),
// 			NewToken(TokenStar, types.StrPtr("/"), nil, 0, 0),
// 			NewGroupingExpr(
// 				NewBinaryExpr(
// 					NewLiteralExpr(2), NewToken(TokenStar, types.StrPtr("-"), nil, 0, 0), NewLiteralExpr(3),
// 				),
// 			),
// 		),
// 	}
// 	astPrinter := NewASTPrinter()
// 	result := astPrinter.VisitGroupingExpr(&expr).(string)
// 	expected := "(9 * 3 / (2 - 3))"
// 	if result != expected {
// 		t.Errorf("got %s, want %s", result, expected)
// 	}
// }
