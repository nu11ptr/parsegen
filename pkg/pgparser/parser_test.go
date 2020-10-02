package pgparser_test

import (
	"testing"

	"github.com/nu11ptr/parsegen/pkg/pgparser"
	"github.com/nu11ptr/parsegen/pkg/pgtoken"
	runtime "github.com/nu11ptr/parsegen/runtime/go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	grammar = `parser = 'parse.g4'

   code {
       top_level -> *ast.TopLevel {{ 
           return ast.NewTopLevel(parseRules) 
       }}
   
       parse_rule -> *ast.ParserRule {{
           return &ast.ParserRule{ruleNameTok.Data, ruleBody}
       }}

       rule_body -> *ast.ParserAlternatives {{
         parserNodes := [][]ast.ParserNode{ruleSects}
         for _, node := range ruleBodySub1s {
             parserNodes = append(parserNodes, node.ruleSects)
         }
         return &ast.ParserAlternatives{Rules: parserNodes}
       }}

       rule_body.sub1 {{
         return &ruleBodySub1{pipeTok: pipeTok, ruleSects: ruleSects}
       }}
   }
`

	expected = `Body:
   └──Parser: parse.g4
   └──Code Blocks:
      └──Code Blocks:
         └──Rule: top_level
         └──Type: *ast.TopLevel
         └──Code: {{ return ast.NewTopLevel(parseRules) }}
      └──Code Blocks:
         └──Rule: parse_rule
         └──Type: *ast.ParserRule
         └──Code: {{ return &ast.ParserRule{ruleNameTok.Data, ruleBody} }}
      └──Code Blocks:
         └──Rule: rule_body
         └──Type: *ast.ParserAlternatives
         └──Code: {{ parserNodes := [][]ast.ParserNode{ruleSects}
         for _, node := range ruleBodySub1s {
             parserNodes = append(parserNodes, node.ruleSects)
         }
         return &ast.ParserAlternatives{Rules: parserNodes} }}
      └──Code Blocks:
         └──Rule: rule_body.sub1
         └──Code: {{ return &ruleBodySub1{pipeTok: pipeTok, ruleSects: ruleSects} }}
`
)

func TestParser(t *testing.T) {
	lex := runtime.NewFromString(grammar)
	tokenizer := pgtoken.NewTokenizer(lex)
	parse := runtime.NewParser(tokenizer)
	parsegen := pgparser.New(parse)

	ast := parsegen.ParseBody()
	require.NotNil(t, ast)
	assert.Equal(t, expected, ast.String())
}
