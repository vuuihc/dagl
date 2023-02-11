package parser

import "fmt"

type StrExpType int

const (
	StrValTypeLiteral StrExpType = iota
	StrValTypeConst
)

func (s StrExpType) String() string {
	switch s {
	case StrValTypeLiteral:
		return "literal"
	case StrValTypeConst:
		return "const"
	default:
		return "unknown"
	}
}

type StrVal struct {
	Type  StrExpType
	Value string
}

type NodeExpType int

const (
	NodeExpTypeVar      NodeExpType = iota // node
	NodeExpTypeFuncCall                    // @call(xxx)
)

type NodeExp struct {
	Type  NodeExpType
	Value interface{}
}

func (n NodeExp) String() string {
	switch n.Type {
	case NodeExpTypeVar:
		return fmt.Sprintf("%s\n", n.Value)
	case NodeExpTypeFuncCall:
		return fmt.Sprintf("%s\n", n.Value)
	default:
		return "unknown"
	}
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
	return fmt.Sprintf("%s = %v\n", a.VarName, a.Value)
}

// NodeAssignStmt is a statement that assigns a NodeVal to a variable
type NodeAssignStmt struct {
	VarName string
	Value   FuncCallStmt
}

func (a NodeAssignStmt) String() string {
	return fmt.Sprintf("%s = %s\n", a.VarName, a.Value)
}

type FuncCallType int

func (f FuncCallType) String() string {
	switch f {
	case FuncCallTypeModel:
		return "model"
	case FuncCallTypeBuiltin:
		return "builtin"
	case FuncCallTypeInline:
		return "inline"
	default:
		return "unknown"
	}
}

const (
	FuncCallTypeModel FuncCallType = iota
	FuncCallTypeBuiltin
	FuncCallTypeInline
)

// FuncCallStmt is a statement that calls a model function
type FuncCallStmt struct {
	Type     FuncCallType
	FuncName string
	Args     []ArgPair
	Inputs   []NodeExp
}

func (m FuncCallStmt) String() string {
	return fmt.Sprintf("[%s]%s(%v)(%v)\n", m.Type, m.FuncName, m.Inputs, m.Args)
}

// IfStmt is a statement that executes a block of statements if a condition is true
type IfStmt struct {
	Cond  NodeExp
	True  []Statement
	False []Statement
}

func (i IfStmt) String() string {
	return fmt.Sprintf(`if (%s){
		%v
	} else {
		%v
	}\n`, i.Cond, i.True, i.False)
}

// FuncStmt is a statement that defines a function
type FuncStmt struct {
	Name   string
	Inputs []string
	Body   []Statement
}

func (f FuncStmt) String() string {
	return fmt.Sprintf(`func %s(%s){
	%v
}\n`, f.Name, f.Inputs, f.Body)
}

// NodeValStmt is a statement that return a node
type NodeValStmt struct {
	Name string
}

func (n NodeValStmt) String() string {
	return fmt.Sprintf("%s\n", n.Name)
}

// CommentStmt is a statement that is a comment
type CommentStmt struct {
	Comment string
}

func (c CommentStmt) String() string {
	return fmt.Sprintf("%s\n", c.Comment)
}
