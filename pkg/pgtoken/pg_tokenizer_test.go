package pgtoken_test

import (
	"testing"

	"github.com/nu11ptr/parsegen/pkg/pgtoken"
	runtime "github.com/nu11ptr/parsegen/runtime/go"
	"github.com/stretchr/testify/assert"
)

const (
	code = `
parser = 'parse.g4'

code('go') {
	top_level -> *ast.TopLevel {{
		return ast.NewTopLevel(parseRules)
	}}

	parse_rule -> *ast.ParserRule {{
		return &ast.ParserRule{ruleNameTok.Data, ruleBody}
	}}

	rule_body.sub1 {{
		return &ruleBodySub1{pipeTok: pipeTok, ruleSects: ruleSects}
	}}
}
`
)

var (
	tokens = []runtime.Token{
		// Parser stmt
		{Type: pgtoken.PARSER},
		{Type: pgtoken.EQUALS},
		{Type: pgtoken.STRING, Data: "'parse.g4'"},

		// Code entry
		{Type: pgtoken.CODE},
		{Type: pgtoken.LPAREN},
		{Type: pgtoken.STRING, Data: "'go'"},
		{Type: pgtoken.RPAREN},
		{Type: pgtoken.LBRACE},

		// top_level
		{Type: pgtoken.RULE_NAME, Data: "top_level"},
		{Type: pgtoken.TYPE, Data: "-> *ast.TopLevel "},
		{Type: pgtoken.CODE_BLOCK, Data: `{{
		return ast.NewTopLevel(parseRules)
	}}`},

		// parse_rule
		{Type: pgtoken.RULE_NAME, Data: "parse_rule"},
		{Type: pgtoken.TYPE, Data: "-> *ast.ParserRule "},
		{Type: pgtoken.CODE_BLOCK, Data: `{{
		return &ast.ParserRule{ruleNameTok.Data, ruleBody}
	}}`},

		// rule_body.sub1
		{Type: pgtoken.RULE_NAME, Data: "rule_body.sub1"},
		{Type: pgtoken.CODE_BLOCK, Data: `{{
		return &ruleBodySub1{pipeTok: pipeTok, ruleSects: ruleSects}
	}}`},

		// Code exit/EOF
		{Type: pgtoken.RBRACE},
		{Type: runtime.EOF},
	}
)

func TestTokenizer(t *testing.T) {
	lex := runtime.NewLexerFromString(code)
	tokenizer := pgtoken.New(lex)

	for _, tok2 := range tokens {
		var tok runtime.Token
		tokenizer.NextToken(&tok)
		assert.Equal(t, tok2.Type, tok.Type)
		assert.Equal(t, tok2.Data, tok.Data)
	}
}
