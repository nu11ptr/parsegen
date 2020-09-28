package token_test

import (
	"testing"

	"github.com/nu11ptr/parsegen/pkg/lexer"
	"github.com/nu11ptr/parsegen/pkg/token"
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
	parserTokens = []lexer.Token{
		// top_level
		lexer.Token{Type: token.RULE_NAME, Data: "top_level"},
		lexer.Token{Type: token.COLON},
		lexer.Token{Type: token.LPAREN},
		lexer.Token{Type: token.RULE_NAME, Data: "parse_rule"},
		lexer.Token{Type: token.PIPE},
		lexer.Token{Type: token.RULE_NAME, Data: "lex_rule"},
		lexer.Token{Type: token.RPAREN},
		lexer.Token{Type: token.STAR},
		lexer.Token{Type: token.TOKEN_NAME, Data: "EOF"},
		lexer.Token{Type: token.SEMI},

		// parse_rule
		lexer.Token{Type: token.RULE_NAME, Data: "parse_rule"},
		lexer.Token{Type: token.COLON},
		lexer.Token{Type: token.TOKEN_NAME, Data: "RULE_NAME"},
		lexer.Token{Type: token.TOKEN_LIT, Data: "':'"},
		lexer.Token{Type: token.RULE_NAME, Data: "rule_body"},
		lexer.Token{Type: token.TOKEN_LIT, Data: "';'"},
		lexer.Token{Type: token.SEMI},

		lexer.Token{Type: token.EOF},
	}

	lexerTokens = []lexer.Token{
		// WS
		lexer.Token{Type: token.FRAGMENT},
		lexer.Token{Type: token.TOKEN_NAME, Data: "WS"},
		lexer.Token{Type: token.COLON},
		lexer.Token{Type: token.LBRACK},
		lexer.Token{Type: token.BASIC_CHAR, Data: " "},
		lexer.Token{Type: token.ESCAPE_CHAR, Data: "\\t"},
		lexer.Token{Type: token.ESCAPE_CHAR, Data: "\\r"},
		lexer.Token{Type: token.ESCAPE_CHAR, Data: "\\n"},
		lexer.Token{Type: token.ESCAPE_CHAR, Data: "\\f"},
		lexer.Token{Type: token.RBRACK},
		lexer.Token{Type: token.PLUS},
		lexer.Token{Type: token.RARROW},
		lexer.Token{Type: token.SKIP_ACTION},
		lexer.Token{Type: token.SEMI},

		// RULE_NAME
		lexer.Token{Type: token.TOKEN_NAME, Data: "RULE_NAME"},
		lexer.Token{Type: token.COLON},
		lexer.Token{Type: token.LBRACK},
		lexer.Token{Type: token.BASIC_CHAR, Data: "a"},
		lexer.Token{Type: token.DASH},
		lexer.Token{Type: token.BASIC_CHAR, Data: "z"},
		lexer.Token{Type: token.RBRACK},
		lexer.Token{Type: token.TOKEN_NAME, Data: "NAME"},
		lexer.Token{Type: token.SEMI},

		// BOGUS
		lexer.Token{Type: token.TOKEN_NAME, Data: "BOGUS"},
		lexer.Token{Type: token.COLON},
		lexer.Token{Type: token.LBRACK},
		lexer.Token{Type: token.UNICODE_ESCAPE_CHAR, Data: "\\uffff"},
		lexer.Token{Type: token.UNICODE_ESCAPE_CHAR, Data: "\\u{abcd}"},
		lexer.Token{Type: token.RBRACK},
		lexer.Token{Type: token.SEMI},

		lexer.Token{Type: token.EOF},
	}
)

func TestParserTokenizer(t *testing.T) {
	lex := lexer.NewFromString(parserCode)
	tokenizer := token.NewParseGen(lex)

	for _, tok2 := range parserTokens {
		var tok lexer.Token
		tokenizer.NextToken(&tok)
		assert.Equal(t, tok2.Type, tok.Type)
		assert.Equal(t, tok2.Data, tok.Data)
	}
}

func TestLexerTokenizer(t *testing.T) {
	lex := lexer.NewFromString(lexerCode)
	tokenizer := token.NewParseGen(lex)

	for _, tok2 := range lexerTokens {
		var tok lexer.Token
		tokenizer.NextToken(&tok)
		assert.Equal(t, tok2.Type, tok.Type)
		assert.Equal(t, tok2.Data, tok.Data)
	}
}
