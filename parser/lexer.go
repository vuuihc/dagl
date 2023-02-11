package parser

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Token int

const (
	ILEGAL            Token = iota // 错误，返回的v是错误信息
	IDENTIFIER                     // 关键字，字母下划线开始，由字母数字下划线组成
	STRING                         // 字符串 `` "" '' 所包围的内容
	SEMICOLON                      // 分号
	LEFT_PARENTHESIS               // 左括号
	RIGHT_PARENTHESIS              // 右括号
	LEFT_CURLY_BRACE               // 左花括号
	RIGHT_CURLY_BRACE              // 右花括号
	LEFT_SQUARE_BRACKET
	RIGHT_SQUARE_BRACKET
	ASSIGNMENT // 等号
	COMMA
	AT
	COMMENT
	EOF // eof
)

func (t Token) String() string {
	switch t {
	case ILEGAL:
		return "ilegal"
	case IDENTIFIER:
		return "identifier"
	case STRING:
		return "string"
	case SEMICOLON:
		return ";"
	case LEFT_PARENTHESIS:
		return "("
	case RIGHT_PARENTHESIS:
		return ")"
	case LEFT_CURLY_BRACE:
		return "{"
	case RIGHT_CURLY_BRACE:
		return "}"
	case LEFT_SQUARE_BRACKET:
		return "["
	case RIGHT_SQUARE_BRACKET:
		return "]"
	case ASSIGNMENT:
		return "="
	case COMMA:
		return ","
	case AT:
		return "@"
	case EOF:
		return "EOF"
	case COMMENT:
		return "COMMENT"
	}
	return ""
}

var literalToToken = map[rune]Token{
	';': SEMICOLON,
	'(': LEFT_PARENTHESIS,
	')': RIGHT_PARENTHESIS,
	'{': LEFT_CURLY_BRACE,
	'}': RIGHT_CURLY_BRACE,
	'=': ASSIGNMENT,
	',': COMMA,
	'[': LEFT_SQUARE_BRACKET,
	']': RIGHT_SQUARE_BRACKET,
	'@': AT,
}

type TokenData struct {
	Token Token
	Value interface{}
}

type lexer struct {
	input     string
	pos       int // 在整个input中的位置
	line      int // 行号，用于报错
	lineBegin int
	queue     []TokenData
}

func newLexer(input string) *lexer {
	return &lexer{input: input}
}

// 主流程
// 遍历输入，根据当前（或连续几个）字符，解析下一个token
// 情况1，正常解析，会解析出下一个token 和对应的值
// 情况2，已经遍历完成，返回EOF token
// 情况3，解析错误，返回ilegal 和 报错信息
func (l *lexer) Next() (t Token, v interface{}) {
	if len(l.queue) > 0 {
		tok := l.queue[len(l.queue)-1]
		l.queue = l.queue[:len(l.queue)-1]
		return tok.Token, tok.Value
	}
	item, err := l.nextItem()
	if err == io.EOF {
		return EOF, nil
	}
	switch item {
	case ';', '(', ')', '{', '}', '=', ',', '[', ']', '@':
		return literalToToken[item], nil
	case '\'', '"', '`':
		return l.scanString(item)
	case '/':
		return l.scanComment()
	default:
		if unicode.IsSpace(item) {
			return l.Next()
		}
		return l.scanIdentifier(item)
	}
}

// LookAhead 用于预读下一个token
// 会返回下一个token，但是不会移动指针
func (l *lexer) LookAhead() (t Token, v interface{}) {
	pos := l.pos
	line := l.line
	lineBegin := l.lineBegin
	t, v = l.Next()
	l.pos = pos
	l.line = line
	l.lineBegin = lineBegin
	return
}

// Back 用于回退一个token
// 会将指针回退到上一个token的位置
func (l *lexer) Back(tok Token, val interface{}) {
	l.queue = append(l.queue, TokenData{tok, val})
}

func (l *lexer) scanComment() (Token, interface{}) {
	item, err := l.nextItem()
	if err != nil || item != '/' {
		return ILEGAL, l.errorf("ilegal char, do you mean `//` ?")
	}
	buf := bytes.NewBufferString("//")
	for {
		item, err := l.nextItem()
		if err == io.EOF || item == '\n' {
			return COMMENT, buf.String()
		}
		buf.WriteRune(item)

	}
}

func (l *lexer) scanString(pre rune) (Token, interface{}) {
	buf := bytes.NewBufferString("")
	for {
		item, err := l.nextItem()
		if err != nil || pre != '`' && item == '\n' {
			return ILEGAL, l.errorf("waiting for %b", pre)
		}
		if item == pre {
			return STRING, buf.String()
		}
		buf.WriteRune(item)
	}
}

func (l *lexer) scanIdentifier(begin rune) (Token, interface{}) {
	if !l.isIdentifier(begin) {
		return ILEGAL, l.errorf("ilegal begin of identifier")
	}
	buf := bytes.NewBufferString(string(begin))
	for {
		item, err := l.nextItem()
		if err == io.EOF || !l.isIdentifier(item) {
			if !l.isIdentifier(item) {
				l.backItem(item)
			}
			if buf.Len() > 0 {
				return IDENTIFIER, buf.String()
			}
			return ILEGAL, l.errorf("empty identifier") // should not in here
		}
		buf.WriteRune(item)
	}
}

// 标志符
func (l *lexer) isIdentifier(item rune) bool {
	return item == '_' || unicode.IsLetter(item) || unicode.IsDigit(item)
}

func (l *lexer) nextItem() (rune, error) {
	if l.pos >= len(l.input) {
		return 0, io.EOF
	}
	r, size := utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += size
	if r == '\n' {
		l.line++
		l.lineBegin = l.pos
	}
	return r, nil
}

// backItem is used to back one rune
func (l *lexer) backItem(item rune) {
	if item == '\n' {
		l.line--
	}
	l.pos -= utf8.RuneLen(item)
}

func (l *lexer) errorf(format string, args ...interface{}) string {
	msg := bytes.NewBufferString("")
	msg.WriteString(fmt.Sprintf("line %d:%s\n", l.line, fmt.Sprintf(format, args...)))
	msg.WriteString(strings.Split(l.input, "\n")[l.line])
	msg.WriteByte('\n')
	msg.WriteString(strings.Repeat(" ", l.pos-l.lineBegin-1))
	msg.WriteString("^\n")
	return msg.String()
}
