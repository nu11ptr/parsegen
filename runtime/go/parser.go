package runtime

type Parser struct {
	t      Tokenizer
	tokens []Token
	pos    int
}

// NewParser creates a new parser with a given tokenizer
func NewParser(t Tokenizer) *Parser {
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

func (p *Parser) CurrToken() *Token {
	return &p.tokens[p.pos]
}

func (p *Parser) NextToken() *Token {
	p.pos++

	if p.pos < len(p.tokens) {
		return &p.tokens[p.pos]
	}

	// Get a new token from the tokenizer and append it to our token history before returning it
	var tok Token
	p.t.NextToken(&tok)
	p.tokens = append(p.tokens, tok)
	return &p.tokens[0]
}

func (p *Parser) MatchTokenOrRollback(tt TokenType, oldPos int) *Token {
	tok := p.CurrToken()
	if tok.Type != tt {
		// Failed - rollback
		p.SetPos(oldPos)
		return nil
	}
	p.NextToken()
	return tok
}

func (p *Parser) TryMatchToken(tt TokenType) *Token {
	tok := p.CurrToken()
	if tok.Type != tt {
		return nil
	}
	p.NextToken()
	return tok
}
