package ast

import (
	"fmt"
	"log"
	"strings"

	"github.com/nu11ptr/parsegen/pkg/token"
	runtime "github.com/nu11ptr/parsegen/runtime/go"
)

const spaces = 3

type TopLevel struct {
	ParserRulesMap map[string]*ParserRule
	LexerRulesMap  map[string]*LexerRule

	ParserRules []*ParserRule
	LexerRules  []*LexerRule
}

func NewTopLevel(rules []*ParserRule) *TopLevel {
	topLevel := &TopLevel{
		ParserRulesMap: make(map[string]*ParserRule, 16),
		LexerRulesMap:  make(map[string]*LexerRule, 16),
	}
	for _, rule := range rules {
		topLevel.ParserRulesMap[rule.Name] = rule
		topLevel.ParserRules = append(topLevel.ParserRules, rule)
	}
	return topLevel
}

func (t *TopLevel) String() string {
	buff := strings.Builder{}
	buff.WriteString("TopLevel:\n")
	for _, rule := range t.ParserRules {
		buff.WriteString(rule.String(1))
	}
	return buff.String()
}

type ParserRule struct {
	Name  string
	Rules *ParserAlternatives
}

func (p *ParserRule) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString(fmt.Sprintf("└──ParserRule: %s\n", p.Name))
	buff.WriteString(p.Rules.String(indent + 1))
	return buff.String()
}

type ParserNode interface {
	ParserNode()
	String(int) string
}

func NewNestedNode(node ParserNode, suffix *runtime.Token) ParserNode {
	if suffix == nil {
		return node
	}
	switch suffix.Type {
	case token.PLUS:
		return &ParserOneOrMore{Node: node}
	case token.STAR:
		return &ParserZeroOrMore{Node: node}
	case token.QUEST_MARK:
		return &ParserZeroOrOne{Node: node}
	default:
		log.Panicf("Unknown token type: %d", suffix.Type)
		return nil
	}
}

type ParserAlternatives struct {
	Rules [][]ParserNode
}

func (p *ParserAlternatives) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString("└──Alternatives:\n")

	for i, alt := range p.Rules {
		buff.WriteString(strings.Repeat(" ", (indent+1)*spaces))
		buff.WriteString(fmt.Sprintf("└──Alternative %d:\n", i))
		for _, rule := range alt {
			buff.WriteString(rule.String(indent + 2))
		}
	}
	return buff.String()
}

func (p *ParserAlternatives) ParserNode() {}

type ParserZeroOrMore struct {
	Node ParserNode
}

func (p *ParserZeroOrMore) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString("└──ZeroOrMore:\n")
	buff.WriteString(p.Node.String(indent + 1))
	return buff.String()
}

func (p *ParserZeroOrMore) ParserNode() {}

type ParserOneOrMore struct {
	Node ParserNode
}

func (p *ParserOneOrMore) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString("└──OneOrMore:\n")
	buff.WriteString(p.Node.String(indent + 1))
	return buff.String()
}

func (p *ParserOneOrMore) ParserNode() {}

type ParserZeroOrOne struct {
	Node ParserNode
}

func (p *ParserZeroOrOne) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString("└──ZeroOrOne:\n")
	buff.WriteString(p.Node.String(indent + 1))
	return buff.String()
}

func (p *ParserZeroOrOne) ParserNode() {}

type ParserRuleRef struct {
	Name string
}

func (p *ParserRuleRef) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString(fmt.Sprintf("└──ParserRuleRef: %s\n", p.Name))
	return buff.String()
}

func (p *ParserRuleRef) ParserNode() {}

type ParserLexerRuleRef struct {
	Name string
}

func (p *ParserLexerRuleRef) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString(fmt.Sprintf("└──LexerRuleRef: %s\n", p.Name))
	return buff.String()
}

func (p *ParserLexerRuleRef) ParserNode() {}

type ParserToken struct {
	Token *runtime.Token
}

func (p *ParserToken) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString("└──Token Literal:\n")
	if p.Token.Data != "" {
		buff.WriteString(strings.Repeat(" ", (indent+1)*spaces))
		buff.WriteString(fmt.Sprintf("└──Data: %s\n", p.Token.Data))
	}
	return buff.String()
}

func (p *ParserToken) ParserNode() {}
