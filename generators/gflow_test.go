package generators

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vuuihc/gfc/parser"
)

// getStatementWithParser parses a statement and returns the statement
func testWithCodeAndGraph(t *testing.T, code string, expected *Graph) {
	p := parser.NewParser(code)
	statments := p.Parse()
	generator := NewGFGenerator(statments)
	graph := generator.GenerateGraph()
	require.Equal(t, expected.MarshalToJson(), graph.MarshalToJson(), "expected %s \nactual %s \n", expected.MarshalToJson(), graph.MarshalToJson())
}

// TestGenerateBuiltinIdentity tests the generator's ability to generate a graph for a builtin identity function
func TestGenerateBuiltinIdentity(t *testing.T) {
	code := `func main(input) {builtin("identity", [input]);}`
	expected := &Graph{
		Nodes: []Node{
			{Type: "builtin.start"},
			{Type: "builtin.identity", Inputs: []int{0}, InDegree: 1},
		},
	}
	testWithCodeAndGraph(t, code, expected)
}

// TestGenerateInlineFunc tests the generator's ability to generate a graph for a inline function
func TestGenerateInlineFunc(t *testing.T) {
	code := `
	inline func setCache(key, result) {
		cacheReq=builtin("jq",[key,result],filter='{"key": .[0], "payload": .[1], "ttl": 259200000}');
		builtin("set_cache", cacheReq, prefix='ime_rec_bert_ner_v1');
	}
	func main(input) {
		result=builtin("jq",input,filter='{"key": .key, "payload": .payload}');
		@call(setCache, [input, result]);
	}`
	expected := &Graph{
		Nodes: []Node{
			{Type: "builtin.start"},
			{Type: "builtin.jq", Inputs: []int{0}, Args: map[string][]string{"filter": {"{\"key\": .key, \"payload\": .payload}"}}, InDegree: 1},
			{Type: "builtin.jq", Inputs: []int{0, 1}, Args: map[string][]string{"filter": {"{\"key\": .[0], \"payload\": .[1], \"ttl\": 259200000}"}}, InDegree: 2},
			{Type: "builtin.set_cache", Inputs: []int{2}, Args: map[string][]string{"prefix": {"ime_rec_bert_ner_v1"}}, InDegree: 1},
		},
	}
	testWithCodeAndGraph(t, code, expected)
}

// TestGenerateFullGraph tests the generator's ability to generate a graph for a full graph
func TestGenerateFullGraph(t *testing.T) {
	code := `
	inline func getCacheKey(input) {
		builtin("jq",input,filter='.suggestion_type+"##"+(.filter_retrievers//[]|join("#"))+"##"+(.context//[]|join("#"))+"##"+.query');
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
			  method='post', max_retry_times="3", default_value='{"actions":[]}', timeout="800ms");
		  @call(setCache, [key, result]);
		  cacheMiss=builtin("jq", cacheRes, filter='.found | not');
		  if(cacheMiss){
			result;
		  }else{
			builtin("jq", cacheRes, filter='.payload');
		  }
	  }
	  
	  // {"payload": "{\"request_id\":\"1674\",\"request_type\":7,\"context\":[],\"context_interval\":[],\"query\":\"红楼梦小姐姐\",\"uid\":\"1674\",\"api_level\":0}"}
`
	expected := &Graph{
		Nodes: []Node{
			{Type: "builtin.start"},
			{Type: "builtin.jq", Inputs: []int{0}, Args: map[string][]string{"filter": {".payload | fromjson"}}, InDegree: 1},
			{Type: "builtin.jq", Inputs: []int{1}, Args: map[string][]string{"filter": {".suggestion_type+\"##\"+(.filter_retrievers//[]|join(\"#\"))+\"##\"+(.context//[]|join(\"#\"))+\"##\"+.query"}}, InDegree: 1},
			{Type: "builtin.lookup_cache", Inputs: []int{2}, Args: map[string][]string{"prefix": {"ime_rec_bert_ner_v1"}}, InDegree: 1},
			{Type: "builtin.http", Inputs: []int{1}, Args: map[string][]string{"endpoint": {"http://xxxxx.Production/suggestion/"}, "method": {"post"}, "max_retry_times": {"3"}, "default_value": {"{\"actions\":[]}"}, "timeout": {"800ms"}}, InDegree: 1},
			{Type: "builtin.jq", Inputs: []int{2, 4}, Args: map[string][]string{"filter": {"{\"key\": .[0], \"payload\": .[1], \"ttl\": 259200000}"}}, InDegree: 2},
			{Type: "builtin.set_cache", Inputs: []int{5}, Args: map[string][]string{"prefix": {"ime_rec_bert_ner_v1"}}, InDegree: 1},
			{Type: "builtin.jq", Inputs: []int{3}, Args: map[string][]string{"filter": {".found | not"}}, InDegree: 1},
			{Type: "builtin.when_true", Inputs: []int{7}, InDegree: 1},
			{Type: "builtin.identity", Inputs: []int{4}, Dependencies: []int{8}, InDegree: 1},
			{Type: "builtin.when_false", Inputs: []int{7}, InDegree: 1},
			{Type: "builtin.jq", Inputs: []int{3}, Args: map[string][]string{"filter": {".payload"}}, Dependencies: []int{10}, InDegree: 1},
			{Type: "builtin.when_any", Inputs: []int{9, 11}, InDegree: 2},
		},
	}
	testWithCodeAndGraph(t, code, expected)
}
