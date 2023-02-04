package main

import "encoding/json"

// 基础数据类型，表达式求值后的结果，应该遵循这个接口
type Val interface {
	String() string
}

// String 类型
type StrVal string

func (s *StrVal) String() string {
	return string(*s)
}

// {StrVal:StrVal} 类型
type StrMap map[StrVal]StrVal

func (m *StrMap) String() string {
	buf, _ := json.Marshal(m)
	return string(buf)
}

// [Val] 类型
type ArrayVal []Val

func (a ArrayVal) String() string {
	buf, _ := json.Marshal(a)
	return string(buf)
}

// Expr 表达式，求值后可得 Val

// Statement
// 语法规则的产生式
// 非终结符 ::= 非终结符 终结符
type Statement interface {
	Print()
}

type FunctionStatement struct {
	Name string
	Body BlockStatement
}

func (f FunctionStatement) Print() {
}

type BlockStatement struct {
}

func (b BlockStatement) Print() {

}

type AssignmentStatement struct {
	Name  string
	Value Expr
}
