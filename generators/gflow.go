package generators

import (
	"encoding/json"
	"fmt"
	"log"

	"emperror.dev/emperror"
	"github.com/vuuihc/gfc/parser"
)

type Node struct {
	// Type 是节点要执行的任务名字。
	//      约定 builtin.XXX 为系统内建任务
	//           op.XXX     为注册到gflow的Op名字
	Type string `msg:"type" json:"type"`

	// Args 节点执行时候的参数。节点执行时也会传递给节点执行期
	Args map[string][]string `msg:"args,omitempty" json:"args,omitempty"`

	// InDegree 节点的入度。当运行时入度为0时，则该节点可以被调度执行。
	InDegree int `msg:"in_degree" json:"in_degree"`

	// Inputs 这个节点依赖的其他节点的输出。int为其他节点在Graph中的Offset。
	// 需要注意的是，输入数量不一定和入度一样。如果一个节点只需要部分输入，则InDegree小于inputs的数量。
	// 这个输入会归并成 json array传递给节点执行器
	Inputs []int `msg:"inputs,omitempty" json:"inputs,omitempty"`

	// Dependencies 表示这个节点的依赖节点。依赖节点与Inputs的区别是，依赖节点的结果并不会传递给Op执行期。
	// 但是依赖节点必须先于该节点执行。
	// 在分支逻辑中需要用到依赖节点
	Dependencies []int `msg:"dependencies,omitempty" json:"dependencies,omitempty"`

	// IsResponse 如果为True，则当这个节点执行完之后，这个图的Response就是这个节点的输出。
	// 需要注意的是，图的IsResponse节点不一定是图的最后一个节点。在返回Response之后，系统
	// 还可以继续执行一些操作。
	IsResponse bool `msg:"is_response,omitempty" json:"is_response,omitempty"`
}

type Graph struct {
	// 图只由节点构成。每个节点只有一个输出。每个节点可以是其他节点的输入。
	Nodes []Node `msg:"nodes" json:"nodes"`
}

// NewNode creates a new node
func (g *Graph) AddNode(node Node) int {
	if node.InDegree == 0 {
		if len(node.Dependencies) == 0 {
			node.InDegree = len(node.Inputs)
		} else {
			node.InDegree = len(node.Dependencies)
		}
	}
	g.Nodes = append(g.Nodes, node)
	return len(g.Nodes) - 1
}

// MarshalToJson marshals a graph to json
func (g *Graph) MarshalToJson() []byte {
	js, err := json.Marshal(g)
	emperror.Panic(err)
	return js
}

// Stack is a map of variables and their values
type Stack map[string]interface{}

// Copy copies a stack
func (s Stack) Copy() Stack {
	newStack := make(Stack)
	for k, v := range s {
		newStack[k] = v
	}
	return newStack
}

// GFGenerator is a gflow generator
// Inputs: statements []parser.Statement
// Outputs: gflow Graph
type GFGenerator struct {
	statements []parser.Statement
	graph      *Graph
}

// NewGFGenerator creates a new gflow generator
func NewGFGenerator(statements []parser.Statement) *GFGenerator {
	return &GFGenerator{
		statements: statements,
		graph: &Graph{Nodes: []Node{{
			Type: "builtin.start",
		}}},
	}
}

// reportError reports an error
func (g *GFGenerator) reportErrorf(stmt parser.Statement, format string, args ...interface{}) {
	ctx := fmt.Sprintf("current statement: %s\n", stmt)
	msg := fmt.Sprintf(format, args...)
	log.Fatalf("generator error: %s\n%s", msg, ctx)
}

// GenerateGraph generates a gflow graph
func (g *GFGenerator) GenerateGraph() *Graph {
	stack := Stack{}
	for _, statement := range g.statements {
		switch v := statement.(type) {
		case parser.AssignStmt: // const definition
			stack[v.VarName] = v.Value
			break
		case parser.FuncStmt: // func definition
			stack[v.Name] = v
			break
		case parser.FuncCallStmt: // func call
			g.newFuncCallNode(&v, stack, nil)
			break
		case parser.NodeAssignStmt:
			g.newNodeAssignNode(&v, stack, nil)
			break
		case parser.IfStmt:
			g.newIfNode(&v, stack, nil)
			break
		case parser.CommentStmt:
			continue
		default:
			g.reportErrorf(statement, "unknown statement type")
		}
	}
	mainFunc, ok := stack["main"].(parser.FuncStmt)
	if !ok {
		log.Fatalf("main function not found")
	}
	if len(mainFunc.Inputs) != 1 {
		log.Fatalf("main function should have only one input")
	}
	stack[mainFunc.Inputs[0]] = 0
	g.newInlineFuncCallNode(&parser.FuncCallStmt{
		FuncName: mainFunc.Name,
		Inputs: []parser.NodeExp{{
			Type:  parser.NodeExpTypeVar,
			Value: mainFunc.Inputs[0],
		}},
	}, stack, nil)
	return g.graph
}

// generateWithDependency generates a node with dependency
func (g *GFGenerator) generateWithDependency(stmt parser.Statement, stack Stack, dependencies []int) int {
	switch v := stmt.(type) {
	case parser.NodeValStmt:
		return g.newFuncCallNode(&parser.FuncCallStmt{
			Type:     parser.FuncCallTypeBuiltin,
			FuncName: "identity",
			Inputs:   []parser.NodeExp{{Type: parser.NodeExpTypeVar, Value: v.Name}},
		}, stack, dependencies)
	case parser.FuncCallStmt:
		return g.newFuncCallNode(&v, stack, dependencies)
	case parser.NodeAssignStmt:
		return g.newNodeAssignNode(&v, stack, dependencies)
	case parser.IfStmt:
		return g.newIfNode(&v, stack, dependencies)
	case parser.CommentStmt:
		return -2
	default:
		g.reportErrorf(stmt, "unknown statement type")
		return -1
	}
}

func (gf *GFGenerator) newFuncCallNode(stmt *parser.FuncCallStmt, stack Stack, dependencies []int) int {
	var nodeType string
	switch stmt.Type {
	case parser.FuncCallTypeModel:
		nodeType = "model." + stmt.FuncName
		break
	case parser.FuncCallTypeBuiltin:
		nodeType = "builtin." + stmt.FuncName
		break
	default:
		return gf.newInlineFuncCallNode(stmt, stack, dependencies)
	}
	node := Node{
		Type:         nodeType,
		Args:         make(map[string][]string),
		Dependencies: dependencies,
	}
	// fill inputs
	for _, input := range stmt.Inputs {
		switch input.Type {
		case parser.NodeExpTypeVar:
			nodeVar, ok := input.Value.(string)
			if !ok {
				gf.reportErrorf(stmt, "invalid input value %v", input.Value)
			}
			inputNode, ok := stack[nodeVar].(int)
			if !ok {
				gf.reportErrorf(stmt, "invalid input node %v", input.Value)
			}
			node.Inputs = append(node.Inputs, inputNode)
			break
		default:
			gf.reportErrorf(stmt, "unknown input type %v", input.Type)
		}
	}

	// fill args
	for _, arg := range stmt.Args {
		switch arg.Value.Type {
		case parser.StrValTypeConst:
			v, ok := stack[arg.Value.Value].(string)
			if !ok {
				gf.reportErrorf(stmt, "const string not found: %v", arg.Value.Value)
			}
			node.Args[arg.Name] = append(node.Args[arg.Name], v)
			break
		case parser.StrValTypeLiteral:
			node.Args[arg.Name] = append(node.Args[arg.Name], arg.Value.Value)
			break
		default:
			gf.reportErrorf(stmt, "unknown arg type %v", arg.Value.Type)
		}
	}
	return gf.graph.AddNode(node)
}

// newNodeAssignNode creates a new node assign node
func (gf *GFGenerator) newNodeAssignNode(stmt *parser.NodeAssignStmt, stack Stack, dependencies []int) int {
	nodeID := gf.newFuncCallNode(&stmt.Value, stack, dependencies)
	stack[stmt.VarName] = nodeID
	return nodeID
}

// newIfNode creates a new if node
func (gf *GFGenerator) newIfNode(stmt *parser.IfStmt, stack Stack, dependencies []int) int {
	var condNodeID int
	switch stmt.Cond.Type {
	case parser.NodeExpTypeVar:
		var ok bool
		condNodeID, ok = stack[stmt.Cond.Value.(string)].(int)
		if !ok {
			gf.reportErrorf(stmt, "invalid cond value of Node Type Var %v", stmt.Cond.Value)
		}
		break
	case parser.NodeExpTypeFuncCall:
		condNodeID = gf.newFuncCallNode(stmt.Cond.Value.(*parser.FuncCallStmt), stack, dependencies)
		break
	default:
		gf.reportErrorf(stmt, "unknown cond type %v", stmt.Cond.Type)
		return -1
	}
	if len(stmt.True) == 0 {
		gf.reportErrorf(stmt, "if statement should have at least one true statement")
		return -1
	}
	// create when_true node
	trueNode := Node{
		Type: "builtin.when_true",
	}
	trueNode.Inputs = append(trueNode.Inputs, condNodeID)
	trueNodeID := gf.graph.AddNode(trueNode)

	trueStack := stack.Copy()
	var trueEndID int
	for _, statement := range stmt.True {
		if id := gf.generateWithDependency(statement, trueStack, []int{trueNodeID}); id >= 0 {
			trueEndID = id
		}
	}
	if stmt.False != nil {
		falseNode := Node{
			Type: "builtin.when_false",
		}
		falseNode.Inputs = append(falseNode.Inputs, condNodeID)
		falseNodeID := gf.graph.AddNode(falseNode)
		var falseEndID int
		for _, statement := range stmt.False {
			if id := gf.generateWithDependency(statement, stack, []int{falseNodeID}); id >= 0 {
				falseEndID = id
			}
		}
		anyNode := Node{
			Type: "builtin.when_any",
		}
		anyNode.Inputs = append(anyNode.Inputs, trueEndID, falseEndID)
		return gf.graph.AddNode(anyNode)
	} else {
		return trueEndID
	}
}

func (gf *GFGenerator) newInlineFuncCallNode(stmt *parser.FuncCallStmt, stack Stack, dependencies []int) int {
	funcStmt, ok := stack[stmt.FuncName].(parser.FuncStmt)
	if !ok {
		gf.reportErrorf(stmt, "inline function not found: %v", stmt.FuncName)
		return -1
	}
	// create a new stack
	newStack := stack.Copy()
	// fill Inputs to newStack
	if len(funcStmt.Inputs) != len(stmt.Inputs) {
		gf.reportErrorf(stmt, "input length mismatch: %v", stmt.FuncName)
		return -1
	}
	for i, input := range stmt.Inputs {
		switch input.Type {
		case parser.NodeExpTypeVar:
			nodeVar, ok := input.Value.(string)
			v, ok := stack[nodeVar].(int)
			if !ok {
				gf.reportErrorf(stmt, "const string not found: %v", nodeVar)
			}
			newStack[funcStmt.Inputs[i]] = v
			break
		default:
			gf.reportErrorf(stmt, "unknown input type %v", input.Type)
		}
	}
	if len(funcStmt.Body) == 0 {
		gf.reportErrorf(stmt, "empty function body: %v", stmt.FuncName)
		return -1
	}
	var lastNodeID int
	for _, stmt := range funcStmt.Body {
		if id := gf.generateWithDependency(stmt, newStack, dependencies); id >= 0 {
			lastNodeID = id
		}
	}
	return lastNodeID
}
