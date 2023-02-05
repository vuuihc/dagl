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

// TestParseInlineFuncCall tests the parser's ability to parse an inline function call
func TestParseInlineFuncCall(t *testing.T) {
	input := `@call(setCache, [req,output]);`
	expected := []Statement{InlineFuncCallStmt{FuncName: "setCache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}, {Type: NodeExpTypeVar, Value: "output"}}}}
	parser := NewParser(input)
	actual := parser.parseBody()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// TestParseNodeAssignment tests the parser's ability to parse a node assignment
func TestParseNodeAssignment(t *testing.T) {
	input := `node=builtin("jq",[input],filter="");`
	expected := []Statement{NodeAssignStmt{VarName: "node", Value: BuiltinFuncCallStmt{FuncName: "jq", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "input"}}, Args: []ArgPair{{Name: "filter", Value: StrVal{Type: StrValTypeLiteral, Value: ""}}}}}}
	parser := NewParser(input)
	actual := parser.parseBody()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// TestParseInlineFunc tests the parser's ability to parse an inline function
func TestParseInlineFunc(t *testing.T) {
	input := `inline func setCache(req, output) {
		builtin("set_cache", [req], output="output");
	}`
	expected := []Statement{FuncStmt{Name: "setCache", Inputs: []string{"req", "output"}, Body: []Statement{BuiltinFuncCallStmt{FuncName: "set_cache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "output", Value: StrVal{Type: StrValTypeLiteral, Value: "output"}}}}}}}
	parser := NewParser(input)
	actual := parser.parse()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// TestParseIfStmt tests the parser's ability to parse an if statement
func TestParseIfStmt(t *testing.T) {
	input := `if (cacheHit) {
		builtin("set_cache", [req], output="output");
	} else {
		model("finder", [req], output="output");
	}}`
	expected := []Statement{IfStmt{Cond: NodeExp{Type: NodeExpTypeVar, Value: "cacheHit"}, True: []Statement{BuiltinFuncCallStmt{FuncName: "set_cache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "output", Value: StrVal{Type: StrValTypeLiteral, Value: "output"}}}}}, False: []Statement{ModelFuncCallStmt{FuncName: "finder", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "output", Value: StrVal{Type: StrValTypeLiteral, Value: "output"}}}}}}}
	parser := NewParser(input)
	actual := parser.parseBody()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}
