package ast

type Token int // TODO: Placeholder

type TopLevel struct {
	ParserRules map[string]*ParserRule
	LexerRules  map[string]*LexerRule
}

type ParserRule struct {
	Name  string
	Rules ParserAlternatives
}

type ParserNode interface {
	ParserNode()
}

type ParserAlternatives struct {
	Rules []ParserNode
}

func (p *ParserAlternatives) ParserNode() {}

type ParserZeroOrMore struct {
	Node ParserNode
}

func (p *ParserZeroOrMore) ParserNode() {}

type ParserOneOrMore struct {
	Node ParserNode
}

func (p *ParserOneOrMore) ParserNode() {}

type ParserZeroOrOne struct {
	Node ParserNode
}

func (p *ParserZeroOrOne) ParserNode() {}

type ParserRuleRef struct {
	Name string
}

func (p *ParserRuleRef) ParserNode() {}

type ParserLexerRuleRef struct {
	Name string
}

func (p *ParserLexerRuleRef) ParserNode() {}

type ParserToken struct {
	Token Token
}

func (p *ParserToken) ParserNode() {}
