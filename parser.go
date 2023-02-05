package main

import "fmt"

// NewParser returns a new parser
func NewParser(input string) *parser {
	return &parser{lexer: newLexer(input)}
}

type parser struct {
	lexer *lexer
}

func (p *parser) parse() (statements []Statement) {
	for {
		tok, v := p.lexer.Next()
		switch tok {
		case EOF:
			return
		case ILEGAL:
			p.reportErrorf("ilegal token: %s", v)
			return
		case AT:
			statements = p.parseConst()
			break
		case IDENTIFIER:
			switch v {
			case "inline":
				statements = p.parseInlineFunc()
				break
			case "func":
				statements = p.parseFunc()
				break
			default:
				p.reportErrorf("unexpected token: %s", tok)
			}
			return
		default:
			p.reportErrorf("unexpected token: %s", tok)
			return
		}
	}
}

func (p *parser) parseConst() (statements []Statement) {
	tok, v := p.checkTokenType(IDENTIFIER)
	constName := v.(string)
	p.checkTokenType(ASSIGNMENT)
	tok, v = p.lexer.Next()
	switch tok {
	case AT:
		_, v = p.checkTokenType(IDENTIFIER)
		statements = append(statements, AssignStmt{VarName: constName, Value: StrVal{Type: StrValTypeConst, Value: v.(string)}})
		break
	case STRING:
		statements = []Statement{AssignStmt{VarName: constName, Value: StrVal{Type: StrValTypeLiteral, Value: v.(string)}}}
		break
	default:
		p.reportErrorf("expect string, got %s", tok)
		return
	}
	p.checkTokenType(SEMICOLON)
	return
}

func (p *parser) parseInlineFunc() (statements []Statement) {
	p.checkTokenAndValue(IDENTIFIER, "func")
	statements = p.parseFunc()
	return
}

func (p *parser) parseFunc() (statements []Statement) {
	tok, v := p.checkTokenType(IDENTIFIER)
	funcName := v.(string)
	p.checkTokenType(LEFT_PARENTHESIS)
	var inputs []string
	for {
		tok, v = p.lexer.Next()
		if tok == RIGHT_PARENTHESIS {
			break
		}
		if len(inputs) > 0 {
			if tok != COMMA {
				p.reportErrorf("expect comma, got %s", tok)
				return
			}
			tok, v = p.lexer.Next()
		}
		if tok != IDENTIFIER {
			p.reportErrorf("expect identifier, got %s", tok)
			return
		}
		inputs = append(inputs, v.(string))
	}
	p.checkTokenType(LEFT_CURLY_BRACE)
	statements = p.parseBody()
	statements = []Statement{FuncStmt{Name: funcName, Inputs: inputs, Body: statements}}
	return
}

func (p *parser) parseBody() (statements []Statement) {
	for {
		var stmts []Statement
		tok, v := p.lexer.Next()
		switch tok {
		case EOF:
			p.reportErrorf("unexpected EOF")
			return
		case ILEGAL:
			p.reportErrorf("ilegal token: %s", v)
			return
		case AT:
			stmts = p.parseInlineFuncCall()
			break
		case IDENTIFIER:
			switch v {
			case "model", "builtin":
				stmts = p.parseBuiltinFuncCall(v)
				break
			case "if":
				stmts = p.parseIfStmt()
				break
			default:
				t1, _ := p.lexer.LookAhead()
				if t1 == ASSIGNMENT {
					p.lexer.Back(tok, v)
					stmts = p.parseNodeAssign()
				}
				p.reportErrorf("unexpected token: %s, %s\n", tok, v)
			}
			break
		case RIGHT_CURLY_BRACE:
			return
		default:
			p.reportErrorf("unexpected token: %s", tok)
			return
		}
		statements = append(statements, stmts...)
	}
}

// parseInlineFuncCall parses inline function call.
// @call(funcName, [input1, input2, ...])
func (p *parser) parseInlineFuncCall() (statements []Statement) {
	p.checkTokenAndValue(IDENTIFIER, "call")
	p.checkTokenType(LEFT_PARENTHESIS)
	_, v := p.checkTokenType(IDENTIFIER)
	funcName := v.(string)
	p.checkTokenType(COMMA)
	inputs := p.parseInputs()
	p.checkTokenType(RIGHT_PARENTHESIS)
	p.checkTokenType(SEMICOLON)
	statements = []Statement{InlineFuncCallStmt{FuncName: funcName, Inputs: inputs}}
	return
}

func (p *parser) parseBuiltinFuncCall(_type interface{}) (statements []Statement) {
	p.checkTokenType(LEFT_PARENTHESIS)
	_, v := p.checkTokenType(STRING)
	funcName := v.(string)
	p.checkTokenType(COMMA)
	inputs := p.parseInputs()
	p.checkTokenType(COMMA)
	argPairs := p.parseArgPairs()
	p.checkTokenType(SEMICOLON)
	if _type == "model" {
		statements = []Statement{ModelFuncCallStmt{FuncName: funcName, Inputs: inputs, Args: argPairs}}
		return
	}
	statements = []Statement{BuiltinFuncCallStmt{FuncName: funcName, Inputs: inputs, Args: argPairs}}
	return
}

// parseInputs parses inputs of a function.
func (p *parser) parseInputs() (inputs []NodeExp) {
	p.checkTokenType(LEFT_SQUARE_BRACKET)
	for {
		tok, v := p.lexer.Next()
		if tok == RIGHT_SQUARE_BRACKET {
			return
		}
		if len(inputs) > 0 {
			if tok != COMMA {
				p.reportErrorf("expect comma, got %s", tok)
				return
			}
			tok, v = p.lexer.Next()
		}
		if tok != IDENTIFIER {
			p.reportErrorf("expect identifier, got %s", tok)
			return
		}
		inputs = append(inputs, NodeExp{Type: NodeExpTypeVar, Value: v.(string)})
	}
}

// parseArgPairs parses argument pairs of a function.
func (p *parser) parseArgPairs() (argPairs []ArgPair) {
	for {
		tok, v := p.lexer.Next()
		if tok == RIGHT_PARENTHESIS {
			break
		}
		if len(argPairs) > 0 {
			if tok != COMMA {
				p.reportErrorf("expect comma, got %s", tok)
				return
			}
			tok, v = p.lexer.Next()
		}
		if tok != IDENTIFIER {
			p.reportErrorf("expect identifier, got %s", tok)
			return
		}
		argName := v.(string)
		p.checkTokenType(ASSIGNMENT)
		var argValue StrVal
		tok, v = p.lexer.Next()
		if tok == AT {
			_, v = p.checkTokenType(IDENTIFIER)
			argValue = StrVal{Type: StrValTypeConst, Value: v.(string)}
		} else if tok == STRING {
			argValue = StrVal{Type: StrValTypeLiteral, Value: v.(string)}
		} else {
			p.reportErrorf("expect string, got %s", tok)
			return
		}
		argPairs = append(argPairs, ArgPair{Name: argName, Value: argValue})
	}
	return
}

func (p *parser) parseNodeAssign() (statements []Statement) {
	tok, v := p.checkTokenType(IDENTIFIER)
	nodeName := v.(string)
	p.checkTokenType(ASSIGNMENT)
	tok, v = p.lexer.Next()
	switch tok {
	case IDENTIFIER:
		switch v.(string) {
		case "model", "builtin":
			stmts := p.parseBuiltinFuncCall(v)
			statements = append(statements, NodeAssignStmt{VarName: nodeName, Value: stmts[0].(FuncCallStmtInterface)})
			break
		default:
			p.reportErrorf("expect builtin, got %s", v.(string))
		}
	case AT:
		stmts := p.parseInlineFuncCall()
		statements = append(statements, NodeAssignStmt{VarName: nodeName, Value: stmts[0].(FuncCallStmtInterface)})
	default:
		p.reportErrorf("expect @call or identifier, got %s", tok)
	}
	return
}

// parserIfStmt parses if statement
func (p *parser) parseIfStmt() (statements []Statement) {
	p.checkTokenType(LEFT_PARENTHESIS)
	var cond NodeExp
	// parse condition
	// condition can be nodeVar or builtinFuncCall or inlineFuncCall
	tok, v := p.lexer.Next()
	switch tok {
	case IDENTIFIER:
		switch v.(string) {
		case "model", "builtin":
			stmts := p.parseBuiltinFuncCall(v)
			cond = NodeExp{Type: NodeExpTypeFunc, Value: stmts[0].(FuncCallStmtInterface)}
			break
		default:
			cond = NodeExp{Type: NodeExpTypeVar, Value: v.(string)}
		}
	case AT:
		stmts := p.parseInlineFuncCall()
		cond = NodeExp{Type: NodeExpTypeFunc, Value: stmts[0].(FuncCallStmtInterface)}
	default:
		p.reportErrorf("expect builtin or @call or identifier, got %s", tok)
	}
	p.checkTokenType(RIGHT_PARENTHESIS)
	// parse true body
	p.checkTokenType(LEFT_CURLY_BRACE)
	trueStmts := p.parseBody()
	var falseStmts []Statement
	tok, v = p.lexer.LookAhead()
	if tok == IDENTIFIER && v.(string) == "else" {
		p.lexer.Next()
		tok, v = p.lexer.Next()
		if tok != LEFT_CURLY_BRACE {
			p.reportErrorf("expect {, got %s", tok)
		}
		falseStmts = p.parseBody()
	}
	statements = []Statement{IfStmt{Cond: cond, True: trueStmts, False: falseStmts}}
	return
}

// checkTokenAndValue checks if the next token is the expected one
func (p *parser) checkTokenAndValue(tok Token, v interface{}) (Token, interface{}) {
	tok2, v2 := p.lexer.Next()
	if tok2 != tok || v2 != v {
		p.reportErrorf("expect %s %s, got %s %s", tok, v, tok2, v2)
	}
	return tok2, v2
}

// checkTokenType checks if the next token is the expected one
func (p *parser) checkTokenType(tok Token) (Token, interface{}) {
	tok2, v := p.lexer.Next()
	if tok2 != tok {
		p.reportErrorf("line %d: %s:\nexpect %s, got %s:%s\n", p.lexer.line, p.lexer.input[p.lexer.lineBegin:p.lexer.pos], tok, tok2, v)
	}
	return tok2, v
}

// reportErrorf reports error
func (p *parser) reportErrorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
