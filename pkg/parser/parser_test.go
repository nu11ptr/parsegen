package parser_test

import (
	"testing"

	"github.com/nu11ptr/parsegen/pkg/parser"
	"github.com/nu11ptr/parsegen/pkg/token"
	runtime "github.com/nu11ptr/parsegen/runtime/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	grammar = `top_level: (parse_rule | lex_rule)* EOF;

// *** Parser parser ***

parse_rule: RULE_NAME ':' rule_body ';';

rule_body: rule_sect+ ('|' rule_sect+)*;

rule_sect: rule_part suffix?;

rule_part
	: '(' rule_body ')'
	| RULE_NAME
	| TOKEN_NAME
	| TOKEN_LIT
	;

suffix: '+' | '*' | '?';
`

	expected = `TopLevel:
   └──ParserRule: top_level
      └──Alternatives:
         └──Alternative 0:
            └──ZeroOrMore:
               └──Alternatives:
                  └──Alternative 0:
                     └──ParserRuleRef: parse_rule
                  └──Alternative 1:
                     └──ParserRuleRef: lex_rule
            └──LexerRuleRef: EOF
   └──ParserRule: parse_rule
      └──Alternatives:
         └──Alternative 0:
            └──LexerRuleRef: RULE_NAME
            └──Token Literal:
               └──Data: ':'
            └──ParserRuleRef: rule_body
            └──Token Literal:
               └──Data: ';'
   └──ParserRule: rule_body
      └──Alternatives:
         └──Alternative 0:
            └──OneOrMore:
               └──ParserRuleRef: rule_sect
            └──ZeroOrMore:
               └──Alternatives:
                  └──Alternative 0:
                     └──Token Literal:
                        └──Data: '|'
                     └──OneOrMore:
                        └──ParserRuleRef: rule_sect
   └──ParserRule: rule_sect
      └──Alternatives:
         └──Alternative 0:
            └──ParserRuleRef: rule_part
            └──ZeroOrOne:
               └──ParserRuleRef: suffix
   └──ParserRule: rule_part
      └──Alternatives:
         └──Alternative 0:
            └──Token Literal:
               └──Data: '('
            └──ParserRuleRef: rule_body
            └──Token Literal:
               └──Data: ')'
         └──Alternative 1:
            └──LexerRuleRef: RULE_NAME
         └──Alternative 2:
            └──LexerRuleRef: TOKEN_NAME
         └──Alternative 3:
            └──LexerRuleRef: TOKEN_LIT
   └──ParserRule: suffix
      └──Alternatives:
         └──Alternative 0:
            └──Token Literal:
               └──Data: '+'
         └──Alternative 1:
            └──Token Literal:
               └──Data: '*'
         └──Alternative 2:
            └──Token Literal:
               └──Data: '?'
`
)

func TestParser(t *testing.T) {
	lex := runtime.NewLexerFromString(grammar)
	tokenizer := token.New(lex)
	parse := runtime.NewParser(tokenizer)
	parsegen := parser.New(parse)

	ast := parsegen.ParseTopLevel()
	require.NotNil(t, ast)
	assert.Equal(t, expected, ast.String())
}
