package ast

type LexerRule struct {
	Fragment bool
	Name     string
	Rules    LexerAlternatives
}

type LexerNode interface {
	LexerNode()
}

type LexerAlternatives struct {
	Rules []LexerNode
}

func (l *LexerAlternatives) LexerNode() {}

type LexerNot struct {
	Node LexerNode
}

func (l *LexerNot) LexerNode() {}

type LexerZeroOrMore struct {
	Node LexerNode
}

func (l *LexerZeroOrMore) LexerNode() {}

type LexerOneOrMore struct {
	Node LexerNode
}

func (l *LexerOneOrMore) LexerNode() {}

type LexerZeroOrOne struct {
	Node LexerNode
}

func (l *LexerZeroOrOne) LexerNode() {}

type LexerRuleRef struct {
	Name string
}

func (l *LexerRuleRef) LexerNode() {}

type LexerToken struct {
	Token Token
}

func (l *LexerToken) LexerNode() {}

type LexerAnyChar struct{}

func (l *LexerAnyChar) LexerNode() {}

type LexerCharClass struct {
	CharData string // TODO
}

func (l *LexerCharClass) LexerNode() {}
