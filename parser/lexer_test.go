package parser

import (
	"testing"
)

// TestLexerBuiltinFuncCall tests the lexer's ability to parse a builtin function call
func TestLexerBuiltinFuncCall(t *testing.T) {
	input := `builtin("get_cache", [req], prefix="hello");`
	expected := []Token{
		IDENTIFIER, LEFT_PARENTHESIS, STRING, COMMA, LEFT_SQUARE_BRACKET, IDENTIFIER, RIGHT_SQUARE_BRACKET, COMMA, IDENTIFIER, ASSIGNMENT, STRING, RIGHT_PARENTHESIS, SEMICOLON, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, _ := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v", e, token)
		}
	}
}

// TestLexerInlineFuncCall tests the lexer's ability to parse an inline function call
func TestLexerInlineFuncCall(t *testing.T) {
	input := `@call(setCache, [req,output]);`
	expected := []Token{
		AT, IDENTIFIER, LEFT_PARENTHESIS, IDENTIFIER, COMMA, LEFT_SQUARE_BRACKET, IDENTIFIER, COMMA, IDENTIFIER, RIGHT_SQUARE_BRACKET, RIGHT_PARENTHESIS, SEMICOLON, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, _ := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v", e, token)
		}
	}
}

// TestLexerConstStr test the lexer's ability to parse a const string
func TestLexerConstStr(t *testing.T) {
	input := `@s = "hello";`
	expected := []Token{
		AT, IDENTIFIER, ASSIGNMENT, STRING, SEMICOLON, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, _ := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v", e, token)
		}
	}
}

// TestLexerNodeAssignment tests the lexer's ability to parse a node assignment
func TestLexerNodeAssignment(t *testing.T) {
	input := `node=builtin("jq",[input],filter="");`
	expected := []Token{
		IDENTIFIER, ASSIGNMENT, IDENTIFIER, LEFT_PARENTHESIS, STRING, COMMA, LEFT_SQUARE_BRACKET, IDENTIFIER, RIGHT_SQUARE_BRACKET, COMMA, IDENTIFIER, ASSIGNMENT, STRING, RIGHT_PARENTHESIS, SEMICOLON, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, _ := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v", e, token)
		}
	}
}

// TestlexerIfStmt tests the lexer's ability to parse an if statement
func TestLexerIfStmt(t *testing.T) {
	input := `if (cacheHit) { @call(setCache, [req, output]); } else { model("model_name", [req]); }`
	expected := []Token{
		IDENTIFIER, LEFT_PARENTHESIS, IDENTIFIER, RIGHT_PARENTHESIS, LEFT_CURLY_BRACE, AT, IDENTIFIER, LEFT_PARENTHESIS, IDENTIFIER, COMMA, LEFT_SQUARE_BRACKET, IDENTIFIER, COMMA, IDENTIFIER, RIGHT_SQUARE_BRACKET, RIGHT_PARENTHESIS, SEMICOLON, RIGHT_CURLY_BRACE, IDENTIFIER, LEFT_CURLY_BRACE, IDENTIFIER, LEFT_PARENTHESIS, STRING, COMMA, LEFT_SQUARE_BRACKET, IDENTIFIER, RIGHT_SQUARE_BRACKET, RIGHT_PARENTHESIS, SEMICOLON, RIGHT_CURLY_BRACE, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, v := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v, %v", e, token, v)
		}
	}
}

// TestLexerFuncDecl tests the lexer's ability to parse a function declaration
func TestLexerFuncDecl(t *testing.T) {
	input := `func setCache(req, output) { builtin("set_cache", [req, output], prefix="hello"); }`
	expected := []Token{
		IDENTIFIER, IDENTIFIER, LEFT_PARENTHESIS, IDENTIFIER, COMMA, IDENTIFIER, RIGHT_PARENTHESIS, LEFT_CURLY_BRACE, IDENTIFIER, LEFT_PARENTHESIS, STRING, COMMA, LEFT_SQUARE_BRACKET, IDENTIFIER, COMMA, IDENTIFIER, RIGHT_SQUARE_BRACKET, COMMA, IDENTIFIER, ASSIGNMENT, STRING, RIGHT_PARENTHESIS, SEMICOLON, RIGHT_CURLY_BRACE, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, v := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v, %v", e, token, v)
		}
	}
}

// TestLexerInlineFuncDecl tests the lexer's ability to parse an inline function declaration
func TestLexerInlineFuncDecl(t *testing.T) {
	input := `inline func(req, output) { builtin("set_cache", [req, output], prefix="hello"); };`
	expected := []Token{
		IDENTIFIER, IDENTIFIER, LEFT_PARENTHESIS, IDENTIFIER, COMMA, IDENTIFIER, RIGHT_PARENTHESIS, LEFT_CURLY_BRACE, IDENTIFIER, LEFT_PARENTHESIS, STRING, COMMA, LEFT_SQUARE_BRACKET, IDENTIFIER, COMMA, IDENTIFIER, RIGHT_SQUARE_BRACKET, COMMA, IDENTIFIER, ASSIGNMENT, STRING, RIGHT_PARENTHESIS, SEMICOLON, RIGHT_CURLY_BRACE, SEMICOLON, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, v := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v, %v", e, token, v)
		}
	}
}

// TestLexerComment tests the lexer's ability to parse a comment
func TestLexerComment(t *testing.T) {
	input := `// this is a comment`
	expected := []Token{
		COMMENT, EOF,
	}
	l := newLexer(input)
	for _, e := range expected {
		token, _ := l.Next()
		if token != e {
			t.Fatalf("expected %v, got %v", e, token)
		}
	}
}
