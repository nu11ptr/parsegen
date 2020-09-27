package parser

import (
	"github.com/nu11ptr/parsegen/pkg/ast"
	"github.com/nu11ptr/parsegen/pkg/lexer"
	"github.com/nu11ptr/parsegen/pkg/token"
)

type Parser struct {
	t      lexer.Tokenizer
	tokens []lexer.Token
	curr   *lexer.Token
	pos    int
}

// NewParser creates a new parser with a given tokenizer
func NewParser(t lexer.Tokenizer) *Parser {
	p := &Parser{t: t, pos: -1}
	p.NextToken()
	return p
}

func (p *Parser) Pos() int {
	return p.pos
}

func (p *Parser) SetPos(pos int) {
	p.pos = pos
	p.curr = &p.tokens[pos]
}

func (p *Parser) CurrToken() *lexer.Token {
	return p.curr
}

func (p *Parser) NextToken() *lexer.Token {
	p.pos++

	if p.pos < len(p.tokens) {
		return &p.tokens[p.pos]
	}

	// Get a new token from the tokenizer and append it to our token history before returning it
	var tok lexer.Token
	p.t.NextToken(&tok)
	p.tokens = append(p.tokens, tok)
	p.curr = &p.tokens[0]
	return p.curr
}

// *** Below will be automatically generated

type NewParseGenParser struct {
	p *Parser

	topLevelMap  map[int]*ast.TopLevel
	parseRuleMap map[int]*ast.ParserRule
	ruleBodyMap  map[int]*ast.ParserAlternatives
}

func NewParseGen(p *Parser) *NewParseGenParser {
	return &NewParseGenParser{
		p:            p,
		topLevelMap:  make(map[int]*ast.TopLevel, 8),
		parseRuleMap: make(map[int]*ast.ParserRule, 8),
		ruleBodyMap:  make(map[int]*ast.ParserAlternatives, 8),
	}
}

func (p *NewParseGenParser) memoParseTopLevel() *ast.TopLevel {
	topLevel, ok := p.topLevelMap[p.p.Pos()]
	if ok {
		return topLevel
	}
	return p.ParseTopLevel()
}

// ParseTopLevel parses the "top_level" parser rule
func (p *NewParseGenParser) ParseTopLevel() *ast.TopLevel {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### parse_rule* ###
	parseRules := []*ast.ParserRule{}
	for {
		parseRule := p.memoParseParseRule()
		if parseRule == nil {
			break
		}
		parseRules = append(parseRules, parseRule)
	}

	// ### EOF ###
	eofTok := p.p.CurrToken()
	if eofTok.Type != token.EOF {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}
	p.p.NextToken()

	topLevel := ast.NewTopLevel(parseRules)
	// Memoize what we did here in case this exact rule/position is needed again
	p.topLevelMap[oldPos] = topLevel
	return topLevel
}

func (p *NewParseGenParser) memoParseParseRule() *ast.ParserRule {
	parseRule, ok := p.parseRuleMap[p.p.Pos()]
	if ok {
		return parseRule
	}
	return p.ParseParseRule()
}

// ParseParseRule parses the "parse_rule" parser rule
func (p *NewParseGenParser) ParseParseRule() *ast.ParserRule {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### RULE_NAME ###
	ruleNameTok := p.p.CurrToken()
	if ruleNameTok.Type != token.RULE_NAME {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}
	p.p.NextToken()

	// ### ':' ###
	colonTok := p.p.CurrToken()
	if colonTok.Type != token.COLON {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}
	p.p.NextToken()

	// ### rule_body ###
	ruleBody := p.memoParseRuleBody()
	if ruleBody == nil {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}

	// ### ';' ###
	semiTok := p.p.CurrToken()
	if semiTok.Type != token.SEMI {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}
	p.p.NextToken()

	parseRule := &ast.ParserRule{ruleNameTok.Data, ruleBody}
	p.parseRuleMap[oldPos] = parseRule
	return parseRule
}

func (p *NewParseGenParser) memoParseRuleBody() *ast.ParserAlternatives {
	ruleBody, ok := p.ruleBodyMap[p.p.Pos()]
	if ok {
		return ruleBody
	}
	return p.ParseRuleBody()
}

// ParseRuleBody parses the "rule_body" parser rule
func (p *NewParseGenParser) ParseRuleBody() *ast.ParserAlternatives {
	return nil
}
