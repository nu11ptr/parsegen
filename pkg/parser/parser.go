package parser

import (
	"github.com/nu11ptr/parsegen/pkg/ast"
	"github.com/nu11ptr/parsegen/pkg/lexer"
	"github.com/nu11ptr/parsegen/pkg/token"
)

type Parser struct {
	t      lexer.Tokenizer
	tokens []lexer.Token
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
}

func (p *Parser) CurrToken() *lexer.Token {
	return &p.tokens[p.pos]
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
	return &p.tokens[0]
}

func (p *Parser) MatchTokenOrRollback(tt lexer.TokenType, oldPos int) *lexer.Token {
	tok := p.CurrToken()
	if tok.Type != tt {
		// Failed - rollback
		p.SetPos(oldPos)
		return nil
	}
	p.NextToken()
	return tok
}

// *** Below will be automatically generated

type NewParseGenParser struct {
	p *Parser

	topLevelMap     map[int]*ast.TopLevel
	parseRuleMap    map[int]*ast.ParserRule
	ruleBodyMap     map[int]*ast.ParserAlternatives
	ruleBodySub1Map map[int]*ruleBodySub1
}

func NewParseGen(p *Parser) *NewParseGenParser {
	return &NewParseGenParser{
		p:               p,
		topLevelMap:     make(map[int]*ast.TopLevel, 8),
		parseRuleMap:    make(map[int]*ast.ParserRule, 8),
		ruleBodyMap:     make(map[int]*ast.ParserAlternatives, 8),
		ruleBodySub1Map: make(map[int]*ruleBodySub1, 8),
	}
}

// *** top_level ***

func (p *NewParseGenParser) memoParseTopLevel() *ast.TopLevel {
	pos := p.p.Pos()
	topLevel, ok := p.topLevelMap[pos]
	if ok {
		return topLevel
	}
	topLevel = p.ParseTopLevel()
	// Memoize what we did here in case this exact rule/position is needed again
	p.topLevelMap[pos] = topLevel
	return topLevel
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
	eofTok := p.p.MatchTokenOrRollback(token.EOF, oldPos)
	if eofTok == nil {
		return nil
	}

	return ast.NewTopLevel(parseRules)
}

// *** parse_rule ***

func (p *NewParseGenParser) memoParseParseRule() *ast.ParserRule {
	pos := p.p.Pos()
	parseRule, ok := p.parseRuleMap[pos]
	if ok {
		return parseRule
	}
	parseRule = p.ParseParseRule()
	// Memoize what we did here in case this exact rule/position is needed again
	p.parseRuleMap[pos] = parseRule
	return parseRule
}

// ParseParseRule parses the "parse_rule" parser rule
func (p *NewParseGenParser) ParseParseRule() *ast.ParserRule {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### RULE_NAME ###
	ruleNameTok := p.p.MatchTokenOrRollback(token.RULE_NAME, oldPos)
	if ruleNameTok == nil {
		return nil
	}

	// ### ':' ###
	colonTok := p.p.MatchTokenOrRollback(token.COLON, oldPos)
	if colonTok == nil {
		return nil
	}

	// ### rule_body ###
	ruleBody := p.memoParseRuleBody()
	if ruleBody == nil {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}

	// ### ';' ###
	semiTok := p.p.MatchTokenOrRollback(token.SEMI, oldPos)
	if semiTok == nil {
		return nil
	}

	return &ast.ParserRule{ruleNameTok.Data, ruleBody}
}

// *** rule_body ***

func (p *NewParseGenParser) memoParseRuleBody() *ast.ParserAlternatives {
	pos := p.p.Pos()
	ruleBody, ok := p.ruleBodyMap[pos]
	if ok {
		return ruleBody
	}
	ruleBody = p.ParseRuleBody()
	// Memoize what we did here in case this exact rule/position is needed again
	p.ruleBodyMap[pos] = ruleBody
	return ruleBody
}

// ParseRuleBody parses the "rule_body" parser rule
func (p *NewParseGenParser) ParseRuleBody() *ast.ParserAlternatives {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### rule_sect+ ###
	ruleSects := []ast.ParserNode{}
	matched := false
	for {
		ruleSect := p.memoParseRuleSect()
		if ruleSect == nil {
			break
		}
		matched = true
		ruleSects = append(ruleSects, ruleSect)
	}
	if !matched {
		// Failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}

	// ### ('|' rule_sect+)* ###
	ruleBodySub1s := []*ruleBodySub1{}
	for {
		ruleBodySub1 := p.memoParseRuleBodySub1()
		if ruleBodySub1 == nil {
			break
		}
		ruleBodySub1s = append(ruleBodySub1s, ruleBodySub1)
	}

	parserNodes := [][]ast.ParserNode{ruleSects}
	for _, node := range ruleBodySub1s {
		parserNodes = append(parserNodes, node.ruleSects)
	}
	return &ast.ParserAlternatives{Rules: parserNodes}
}

// *** rule_body (sub 1) ***

type ruleBodySub1 struct {
	pipeTok   *lexer.Token
	ruleSects []ast.ParserNode
}

func (p *NewParseGenParser) memoParseRuleBodySub1() *ruleBodySub1 {
	pos := p.p.Pos()
	ruleBodySub1, ok := p.ruleBodySub1Map[pos]
	if ok {
		return ruleBodySub1
	}
	ruleBodySub1 = p.ParseRuleBodySub1()
	// Memoize what we did here in case this exact rule/position is needed again
	p.ruleBodySub1Map[pos] = ruleBodySub1
	return ruleBodySub1
}

func (p *NewParseGenParser) ParseRuleBodySub1() *ruleBodySub1 {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### '|' ###
	pipeTok := p.p.MatchTokenOrRollback(token.PIPE, oldPos)
	if pipeTok == nil {
		return nil
	}

	// ### rule_sect+ ###
	ruleSects := []ast.ParserNode{}
	matched := false
	for {
		ruleSect := p.memoParseRuleSect()
		if ruleSect == nil {
			break
		}
		matched = true
		ruleSects = append(ruleSects, ruleSect)
	}
	if !matched {
		// Failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}

	return &ruleBodySub1{pipeTok: pipeTok, ruleSects: ruleSects}
}

// *** rule_select ***

func (p *NewParseGenParser) memoParseRuleSect() ast.ParserNode {
	return nil
}

func (p *NewParseGenParser) ParseRuleSect() ast.ParserNode {
	return nil
}
