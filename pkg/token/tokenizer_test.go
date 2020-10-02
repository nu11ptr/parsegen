package token_test

import (
	"testing"

	"github.com/nu11ptr/parsegen/pkg/token"
	runtime "github.com/nu11ptr/parsegen/runtime/go"
	"github.com/stretchr/testify/assert"
)

const (
	parserCode = `
top_level: (parse_rule | lex_rule)* EOF;

// *** Parser parser ***

parse_rule: RULE_NAME ':' rule_body ';';
`

	lexerCode = `
fragment WS: [ \t\r\n\f]+ -> skip;

/* a comment
   same comment */

RULE_NAME: [a-z] NAME;

BOGUS: [\uffff\u{abcd}];
`
)

var (
	parserTokens = []runtime.Token{
		// top_level
		{Type: token.RULE_NAME, Data: "top_level"},
		{Type: token.COLON},
		{Type: token.LPAREN},
		{Type: token.RULE_NAME, Data: "parse_rule"},
		{Type: token.PIPE},
		{Type: token.RULE_NAME, Data: "lex_rule"},
		{Type: token.RPAREN},
		{Type: token.STAR},
		{Type: token.TOKEN_NAME, Data: "EOF"},
		{Type: token.SEMI},

		// parse_rule
		{Type: token.RULE_NAME, Data: "parse_rule"},
		{Type: token.COLON},
		{Type: token.TOKEN_NAME, Data: "RULE_NAME"},
		{Type: token.TOKEN_LIT, Data: "':'"},
		{Type: token.RULE_NAME, Data: "rule_body"},
		{Type: token.TOKEN_LIT, Data: "';'"},
		{Type: token.SEMI},

		{Type: runtime.EOF},
	}

	lexerTokens = []runtime.Token{
		// WS
		{Type: token.FRAGMENT},
		{Type: token.TOKEN_NAME, Data: "WS"},
		{Type: token.COLON},
		{Type: token.LBRACK},
		{Type: token.BASIC_CHAR, Data: " "},
		{Type: token.ESCAPE_CHAR, Data: "\\t"},
		{Type: token.ESCAPE_CHAR, Data: "\\r"},
		{Type: token.ESCAPE_CHAR, Data: "\\n"},
		{Type: token.ESCAPE_CHAR, Data: "\\f"},
		{Type: token.RBRACK},
		{Type: token.PLUS},
		{Type: token.RARROW},
		{Type: token.SKIP_ACTION},
		{Type: token.SEMI},

		// RULE_NAME
		{Type: token.TOKEN_NAME, Data: "RULE_NAME"},
		{Type: token.COLON},
		{Type: token.LBRACK},
		{Type: token.BASIC_CHAR, Data: "a"},
		{Type: token.DASH},
		{Type: token.BASIC_CHAR, Data: "z"},
		{Type: token.RBRACK},
		{Type: token.TOKEN_NAME, Data: "NAME"},
		{Type: token.SEMI},

		// BOGUS
		{Type: token.TOKEN_NAME, Data: "BOGUS"},
		{Type: token.COLON},
		{Type: token.LBRACK},
		{Type: token.UNICODE_ESCAPE_CHAR, Data: "\\uffff"},
		{Type: token.UNICODE_ESCAPE_CHAR, Data: "\\u{abcd}"},
		{Type: token.RBRACK},
		{Type: token.SEMI},

		{Type: runtime.EOF},
	}
)

func TestParserTokenizer(t *testing.T) {
	lex := runtime.NewFromString(parserCode)
	tokenizer := token.NewParseGen(lex)

	for _, tok2 := range parserTokens {
		var tok runtime.Token
		tokenizer.NextToken(&tok)
		assert.Equal(t, tok2.Type, tok.Type)
		assert.Equal(t, tok2.Data, tok.Data)
	}
}

func TestLexerTokenizer(t *testing.T) {
	lex := runtime.NewFromString(lexerCode)
	tokenizer := token.NewParseGen(lex)

	for _, tok2 := range lexerTokens {
		var tok runtime.Token
		tokenizer.NextToken(&tok)
		assert.Equal(t, tok2.Type, tok.Type)
		assert.Equal(t, tok2.Data, tok.Data)
	}
}
