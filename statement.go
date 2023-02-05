package main

import "fmt"

type StrExpType int

const (
	StrValTypeLiteral StrExpType = iota
	StrValTypeConst
)

type StrVal struct {
	Type  StrExpType
	Value string
}

type NodeExpType int

const (
	NodeExpTypeVar NodeExpType = iota
	NodeExpTypeFunc
)

type NodeExp struct {
	Type  NodeExpType
	Value interface{}
}

// BuiltinFuncArgPair
type ArgPair struct {
	Name  string
	Value StrVal
}

// Statement is an interface for all statements
type Statement interface {
	// String returns the statement in string format
	String() string
}

// AssignStmt is a statement that assigns a value to a constant
type AssignStmt struct {
	VarName string
	Value   StrVal
}

func (a AssignStmt) String() string {
	return fmt.Sprintf("%s = %v", a.VarName, a.Value)
}

// NodeAssignStmt is a statement that assigns a NodeVal to a variable
type NodeAssignStmt struct {
	VarName string
	Value   FuncCallStmtInterface
}

func (a NodeAssignStmt) String() string {
	return fmt.Sprintf("%s = %s", a.VarName, a.Value)
}

// FuncCallStmtInterface is an interface for all function call statements
type FuncCallStmtInterface interface {
	Statement
	FuncCall()
}

// ModelFuncCallStmt is a statement that calls a model function
type ModelFuncCallStmt struct {
	FuncName string
	Args     []ArgPair
	Inputs   []NodeExp
}

func (m ModelFuncCallStmt) String() string {
	return fmt.Sprintf("%s(%v)", m.FuncName, m.Args)
}

func (m ModelFuncCallStmt) FuncCall() {}

// BuiltinFuncCallStmt is a statement that calls a builtin function
type BuiltinFuncCallStmt struct {
	FuncName string
	Args     []ArgPair
	Inputs   []NodeExp
}

func (b BuiltinFuncCallStmt) String() string {
	return fmt.Sprintf("%s(%v)", b.FuncName, b.Args)
}

func (b BuiltinFuncCallStmt) FuncCall() {}

// InlineFuncCallStmt is a statement that calls an inline function
type InlineFuncCallStmt struct {
	FuncName string
	Args     []ArgPair
	Inputs   []NodeExp
}

func (i InlineFuncCallStmt) String() string {
	return fmt.Sprintf("%s(%v)", i.FuncName, i.Args)
}

func (i InlineFuncCallStmt) FuncCall() {}

// IfStmt is a statement that executes a block of statements if a condition is true
type IfStmt struct {
	Cond  NodeExp
	True  []Statement
	False []Statement
}

func (i IfStmt) String() string {
	return fmt.Sprintf("if %v", i.Cond)
}

// FuncStmt is a statement that defines a function
type FuncStmt struct {
	Name   string
	Inputs []string
	Body   []Statement
}

func (f FuncStmt) String() string {
	return fmt.Sprintf("func %s(%s)", f.Name, f.Inputs)
}
