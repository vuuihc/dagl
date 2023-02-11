package parser

import (
	"reflect"
	"testing"
)

// TestParseConst tests the parser's ability to parse a constant assignment
func TestParseConst(t *testing.T) {
	input := `@foo = "bar";`
	expected := []Statement{AssignStmt{VarName: "foo", Value: StrVal{Type: StrValTypeLiteral, Value: "bar"}}}
	parser := NewParser(input)
	parser.lexer.Next()
	actual := parser.parseConst()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// TestParseBuiltinFuncCall tests the parser's ability to parse a builtin function call
func TestParseBuiltinFuncCall(t *testing.T) {
	input := `builtin("get_cache", [req], prefix="hello");`
	expected := []Statement{FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "get_cache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "prefix", Value: StrVal{Type: StrValTypeLiteral, Value: "hello"}}}}}
	parser := NewParser(input)
	parser.lexer.Next()
	actual := parser.parseFuncCall(FuncCallTypeBuiltin)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// TestParseInlineFuncCall tests the parser's ability to parse an inline function call
func TestParseInlineFuncCall(t *testing.T) {
	input := `@call(setCache, [req,output]);`
	expected := []Statement{FuncCallStmt{Type: FuncCallTypeInline, FuncName: "setCache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}, {Type: NodeExpTypeVar, Value: "output"}}}}
	parser := NewParser(input)
	parser.lexer.Next()
	actual := parser.parseInlineFuncCall()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// TestParseModelFuncCall tests the parser's ability to parse a model function call
func TestParseModelFuncCall(t *testing.T) {
	input := `model("finder", [req], output="output");`
	expected := []Statement{FuncCallStmt{Type: FuncCallTypeModel, FuncName: "finder", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "output", Value: StrVal{Type: StrValTypeLiteral, Value: "output"}}}}}
	parser := NewParser(input)
	parser.lexer.Next()
	actual := parser.parseFuncCall(FuncCallTypeModel)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// TestParseNodeAssignment tests the parser's ability to parse a node assignment
func TestParseNodeAssignment(t *testing.T) {
	input := `node=builtin("jq",[input],filter="");`
	expected := []Statement{NodeAssignStmt{VarName: "node", Value: FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "jq", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "input"}}, Args: []ArgPair{{Name: "filter", Value: StrVal{Type: StrValTypeLiteral, Value: ""}}}}}}
	parser := NewParser(input)
	actual := parser.parseNodeAssign()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// TestParseInlineFunc tests the parser's ability to parse an inline function
func TestParseInlineFunc(t *testing.T) {
	input := `inline func setCache(key, result) {
		cacheReq=builtin("jq",[key,result],filter='{"key": .[0], "payload": .[1], "ttl": 259200000}');
		builtin("set_cache", cacheReq, prefix='ime_rec_bert_ner_v1');
	}
	// test comment`
	expected := []Statement{
		FuncStmt{
			Name:   "setCache",
			Inputs: []string{"key", "result"},
			Body: []Statement{
				NodeAssignStmt{
					VarName: "cacheReq",
					Value:   FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "jq", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "key"}, {Type: NodeExpTypeVar, Value: "result"}}, Args: []ArgPair{{Name: "filter", Value: StrVal{Type: StrValTypeLiteral, Value: `{"key": .[0], "payload": .[1], "ttl": 259200000}`}}}},
				},
				FuncCallStmt{
					Type:     FuncCallTypeBuiltin,
					FuncName: "set_cache",
					Inputs:   []NodeExp{{Type: NodeExpTypeVar, Value: "cacheReq"}},
					Args:     []ArgPair{{Name: "prefix", Value: StrVal{Type: StrValTypeLiteral, Value: "ime_rec_bert_ner_v1"}}}},
			},
		},
		CommentStmt{Comment: "// test comment"},
	}
	parser := NewParser(input)
	parser.lexer.Next()
	actual := parser.Parse()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %#v\n got %#v\n", expected, actual)
	}
}

// TestParseIfStmt tests the parser's ability to parse an if statement
func TestParseIfStmt(t *testing.T) {
	input := `if (cacheHit) {
		builtin("set_cache", [req], output="output");
	} else {
		// test comment in body
		model("finder", [req], output="output");
	}}`
	expected := []Statement{IfStmt{
		Cond:  NodeExp{Type: NodeExpTypeVar, Value: "cacheHit"},
		True:  []Statement{FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "set_cache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "output", Value: StrVal{Type: StrValTypeLiteral, Value: "output"}}}}},
		False: []Statement{CommentStmt{Comment: "// test comment in body"}, FuncCallStmt{Type: FuncCallTypeModel, FuncName: "finder", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "req"}}, Args: []ArgPair{{Name: "output", Value: StrVal{Type: StrValTypeLiteral, Value: "output"}}}}}}}
	parser := NewParser(input)
	parser.lexer.Next()
	actual := parser.parseIfStmt()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// TestParseComment tests the parser's ability to parse a comment
func TestParseComment(t *testing.T) {
	input := `// this is a comment`
	expected := []Statement{CommentStmt{Comment: "// this is a comment"}}
	parser := NewParser(input)
	actual := parser.Parse()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

// TestParseFullCode tests the parser's ability to parse a full code
func TestParseFullCode(t *testing.T) {
	input := `
	@cacheKey='.suggestion_type+"##"+(.filter_retrievers//[]|join("#"))+"##"+(.context//[]|join("#"))+"##"+.query';
	inline func getCacheKey(input) {
		builtin("jq",input,filter=@cacheKey);
	  }
	  
	  inline func setCache(key, result) {
		  cacheReq=builtin("jq",[key,result],filter='{"key": .[0], "payload": .[1], "ttl": 259200000}');
		  builtin("set_cache", cacheReq, prefix='ime_rec_bert_ner_v1');
	  }
	  
	  inline func lookupCache(key){
		builtin("lookup_cache", key, prefix='ime_rec_bert_ner_v1');
	  }
	  
	  func main(input) {
		  input = builtin("jq", input, filter='.payload | fromjson');
		  key=@call(getCacheKey, [input]);
		  cacheRes=@call(lookupCache,[key]);
		  result=builtin("http", input, endpoint='http://192002625-146479.Production/suggestion/',
			  method="post", max_retry_times="3", default_value='{"actions":[]}', timeout="800ms");
		  @call(setCache, [key, result]);
		  cacheMiss=builtin("jq", cacheRes, filter='.found | not');
		  if(cacheMiss){
			builtin("identity",result);
		  }else{
			builtin("jq", cacheRes, filter='.payload');
		  }
	  }
	  
	  // {"payload": "{\"request_id\":\"1674\",\"request_type\":7,\"context\":[],\"context_interval\":[],\"query\":\"红楼梦小姐姐\",\"uid\":\"1674\",\"api_level\":0}"}
	`
	expected := []Statement{
		AssignStmt{VarName: "cacheKey", Value: StrVal{Type: StrValTypeLiteral, Value: ".suggestion_type+\"##\"+(.filter_retrievers//[]|join(\"#\"))+\"##\"+(.context//[]|join(\"#\"))+\"##\"+.query"}},
		FuncStmt{
			Name: "getCacheKey", Inputs: []string{"input"},
			Body: []Statement{
				FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "jq", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "input"}}, Args: []ArgPair{{Name: "filter", Value: StrVal{Type: StrValTypeConst, Value: "cacheKey"}}}},
			},
		},
		FuncStmt{
			Name:   "setCache",
			Inputs: []string{"key", "result"},
			Body: []Statement{
				NodeAssignStmt{VarName: "cacheReq", Value: FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "jq", Inputs: []NodeExp{
					{Type: NodeExpTypeVar, Value: "key"}, {Type: NodeExpTypeVar, Value: "result"},
				}, Args: []ArgPair{
					{Name: "filter", Value: StrVal{Type: StrValTypeLiteral, Value: "{\"key\": .[0], \"payload\": .[1], \"ttl\": 259200000}"}},
				},
				}},
				FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "set_cache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "cacheReq"}}, Args: []ArgPair{{Name: "prefix", Value: StrVal{Type: StrValTypeLiteral, Value: "ime_rec_bert_ner_v1"}}}}},
		}, FuncStmt{
			Name:   "lookupCache",
			Inputs: []string{"key"},
			Body: []Statement{
				FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "lookup_cache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "key"}}, Args: []ArgPair{{Name: "prefix", Value: StrVal{Type: StrValTypeLiteral, Value: "ime_rec_bert_ner_v1"}}}},
			}},
		FuncStmt{
			Name:   "main",
			Inputs: []string{"input"},
			Body: []Statement{
				NodeAssignStmt{VarName: "input", Value: FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "jq", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "input"}}, Args: []ArgPair{{Name: "filter", Value: StrVal{Type: StrValTypeLiteral, Value: ".payload | fromjson"}}}}},
				NodeAssignStmt{VarName: "key", Value: FuncCallStmt{Type: FuncCallTypeInline, FuncName: "getCacheKey", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "input"}}}},
				NodeAssignStmt{VarName: "cacheRes", Value: FuncCallStmt{Type: FuncCallTypeInline, FuncName: "lookupCache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "key"}}}},
				NodeAssignStmt{VarName: "result", Value: FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "http", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "input"}}, Args: []ArgPair{
					{Name: "endpoint", Value: StrVal{Type: StrValTypeLiteral, Value: "http://192002625-146479.Production/suggestion/"}},
					{Name: "method", Value: StrVal{Type: StrValTypeLiteral, Value: "post"}},
					{Name: "max_retry_times", Value: StrVal{Type: StrValTypeLiteral, Value: "3"}},
					{Name: "default_value", Value: StrVal{Type: StrValTypeLiteral, Value: "{\"actions\":[]}"}},
					{Name: "timeout", Value: StrVal{Type: StrValTypeLiteral, Value: "800ms"}},
				}}},
				FuncCallStmt{Type: FuncCallTypeInline, FuncName: "setCache", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "key"}, {Type: NodeExpTypeVar, Value: "result"}}},
				NodeAssignStmt{VarName: "cacheMiss", Value: FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "jq", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "cacheRes"}}, Args: []ArgPair{{Name: "filter", Value: StrVal{Type: StrValTypeLiteral, Value: ".found | not"}}}}},
				IfStmt{
					Cond: NodeExp{Type: NodeExpTypeVar, Value: "cacheMiss"},
					True: []Statement{
						FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "identity", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "result"}}},
					},
					False: []Statement{
						FuncCallStmt{Type: FuncCallTypeBuiltin, FuncName: "jq", Inputs: []NodeExp{{Type: NodeExpTypeVar, Value: "cacheRes"}}, Args: []ArgPair{{Name: "filter", Value: StrVal{Type: StrValTypeLiteral, Value: ".payload"}}}},
					},
				},
			},
		},
		CommentStmt{Comment: `// {"payload": "{\"request_id\":\"1674\",\"request_type\":7,\"context\":[],\"context_interval\":[],\"query\":\"红楼梦小姐姐\",\"uid\":\"1674\",\"api_level\":0}"}`},
	}
	parser := NewParser(input)
	actual := parser.Parse()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %s\n got %s\n", expected, actual)
	}
}
