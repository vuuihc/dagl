package main

import (
	"reflect"
	"testing"
)

// TestParseConst tests the parser's ability to parse a constant assignment
func TestParseConst(t *testing.T) {
	input := `@foo = "bar";`
	expected := []Statement{AssignStmt{VarName: "foo", Value: StrVal{Type: StrValTypeLiteral, Value: "bar"}}}
	parser := NewParser(input)
	actual := parser.parse()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// TestParseBuiltinFuncCall tests the parser's ability to parse a builtin function call
func TestParseBuiltinFuncCall(t *testing.T) {
	input := `builtin("get_cache", [req], prefix="hello");`
	expected := []Statement{BuiltinFuncCallStmt{FuncName: "get_cache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "prefix", Value: StrVal{Type: StrValTypeLiteral, Value: "hello"}}}}}
	parser := NewParser(input)
	actual := parser.parseBody()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
