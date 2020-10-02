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
		runtime.Token{Type: token.RULE_NAME, Data: "top_level"},
		runtime.Token{Type: token.COLON},
		runtime.Token{Type: token.LPAREN},
		runtime.Token{Type: token.RULE_NAME, Data: "parse_rule"},
		runtime.Token{Type: token.PIPE},
		runtime.Token{Type: token.RULE_NAME, Data: "lex_rule"},
		runtime.Token{Type: token.RPAREN},
		runtime.Token{Type: token.STAR},
		runtime.Token{Type: token.TOKEN_NAME, Data: "EOF"},
		runtime.Token{Type: token.SEMI},

		// parse_rule
		runtime.Token{Type: token.RULE_NAME, Data: "parse_rule"},
		runtime.Token{Type: token.COLON},
		runtime.Token{Type: token.TOKEN_NAME, Data: "RULE_NAME"},
		runtime.Token{Type: token.TOKEN_LIT, Data: "':'"},
		runtime.Token{Type: token.RULE_NAME, Data: "rule_body"},
		runtime.Token{Type: token.TOKEN_LIT, Data: "';'"},
		runtime.Token{Type: token.SEMI},

		runtime.Token{Type: runtime.EOF},
	}

	lexerTokens = []runtime.Token{
		// WS
		runtime.Token{Type: token.FRAGMENT},
		runtime.Token{Type: token.TOKEN_NAME, Data: "WS"},
		runtime.Token{Type: token.COLON},
		runtime.Token{Type: token.LBRACK},
		runtime.Token{Type: token.BASIC_CHAR, Data: " "},
		runtime.Token{Type: token.ESCAPE_CHAR, Data: "\\t"},
		runtime.Token{Type: token.ESCAPE_CHAR, Data: "\\r"},
		runtime.Token{Type: token.ESCAPE_CHAR, Data: "\\n"},
		runtime.Token{Type: token.ESCAPE_CHAR, Data: "\\f"},
		runtime.Token{Type: token.RBRACK},
		runtime.Token{Type: token.PLUS},
		runtime.Token{Type: token.RARROW},
		runtime.Token{Type: token.SKIP_ACTION},
		runtime.Token{Type: token.SEMI},

		// RULE_NAME
		runtime.Token{Type: token.TOKEN_NAME, Data: "RULE_NAME"},
		runtime.Token{Type: token.COLON},
		runtime.Token{Type: token.LBRACK},
		runtime.Token{Type: token.BASIC_CHAR, Data: "a"},
		runtime.Token{Type: token.DASH},
		runtime.Token{Type: token.BASIC_CHAR, Data: "z"},
		runtime.Token{Type: token.RBRACK},
		runtime.Token{Type: token.TOKEN_NAME, Data: "NAME"},
		runtime.Token{Type: token.SEMI},

		// BOGUS
		runtime.Token{Type: token.TOKEN_NAME, Data: "BOGUS"},
		runtime.Token{Type: token.COLON},
		runtime.Token{Type: token.LBRACK},
		runtime.Token{Type: token.UNICODE_ESCAPE_CHAR, Data: "\\uffff"},
		runtime.Token{Type: token.UNICODE_ESCAPE_CHAR, Data: "\\u{abcd}"},
		runtime.Token{Type: token.RBRACK},
		runtime.Token{Type: token.SEMI},

		runtime.Token{Type: runtime.EOF},
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
