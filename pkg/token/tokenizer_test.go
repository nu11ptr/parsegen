package token_test

import (
	"testing"

	"github.com/nu11ptr/parsegen/pkg/lexer"
	"github.com/nu11ptr/parsegen/pkg/token"
	"github.com/stretchr/testify/assert"
)

const (
	code = `
top_level: (parse_rule | lex_rule)* EOF;

// *** Parser parser ***

parse_rule: RULE_NAME ':' rule_body ';';
`
)

var (
	tokens = []lexer.Token{
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
)

func TestTokenizer(t *testing.T) {
	lex := lexer.NewFromString(code)
	tokenizer := token.NewParseGen(lex)

	for _, tok2 := range tokens {
		var tok lexer.Token
		tokenizer.NextToken(&tok)
		assert.Equal(t, tok2.Type, tok.Type)
		assert.Equal(t, tok2.Data, tok.Data)
	}
}
