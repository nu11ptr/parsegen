package pgtoken

import runtime "github.com/nu11ptr/parsegen/runtime/go"

// *** Potentially Generated ***

const (
	// Char set
	RULE_NAME runtime.TokenType = iota + runtime.EOF + 1

	// Sequences
	STRING
	TYPE
	CODE_BLOCK

	// Keywords
	PARSER
	CODE

	// Basic Sequences
	EQUALS
	LBRACE
	RBRACE
	LPAREN
	RPAREN
)

var (
	keywords = map[string]runtime.TokenType{
		"parser": PARSER,
		"code":   CODE,
	}
)

type Tokenizer struct {
	lex *runtime.Lexer
}

func New(lex *runtime.Lexer) *Tokenizer {
	return &Tokenizer{lex: lex}
}

func (t *Tokenizer) processRuleName(tok *runtime.Token) bool {
	// [a-z]
	if !t.lex.MatchCharInRange('a', 'z') {
		return false
	}

	// [A-Za-z0-9_]*
	for t.lex.MatchCharInRange('A', 'Z') || t.lex.MatchCharInRange('a', 'z') ||
		t.lex.MatchCharInRange('0', '9') || t.lex.MatchChar('_') || t.lex.MatchChar('.') {
	}

	t.lex.BuildTokenData(RULE_NAME, tok)

	// Possible conflicting keyword
	tt, ok := keywords[tok.Data]
	if ok {
		tok.Type = tt
		tok.Data = ""
	}

	return true
}

func (t *Tokenizer) NextToken(tok *runtime.Token) {
	// Skip
	skipping := true
	for skipping {
		switch t.lex.CurrChar() {
		case '/':
			switch t.lex.NextChar() {
			// '//'
			case '/':
				t.lex.NextChar()

				// ~[\r\n]*
				for t.lex.MatchCharExceptInSeq("\r\n") {
				}
				t.lex.DiscardTokenData()
			// '/*'
			case '*':
				t.lex.NextChar()
				t.lex.MatchUntilSeq("*/")
				t.lex.MatchSeq("*/")
				t.lex.DiscardTokenData()
			default:
				t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
				return
			}
		// [ \t\r\n\f]+
		case ' ', '\t', '\r', '\n', '\f':
			t.lex.NextChar()

			for t.lex.MatchCharInSeq(" \t\r\n\f") {
			}
			t.lex.DiscardTokenData()
		default:
			skipping = false
		}
	}

	if t.processRuleName(tok) {
		return
	}

	switch t.lex.CurrChar() {
	case '\'':
		t.lex.NextChar()

		// ('\\\'' | ~'\'')+
		matched := false
		for t.lex.MatchSeq("\\'") || t.lex.MatchCharExcept('\'') {
			matched = true
		}
		if !matched {
			t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
			return
		}

		// '\''
		if !t.lex.MatchChar('\'') {
			t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
			return
		}

		t.lex.BuildTokenData(STRING, tok)
	case '-':
		t.lex.NextChar()

		// '->'
		if !t.lex.MatchChar('>') {
			t.lex.BuildTokenData(runtime.ILLEGAL, tok)
			return
		}

		// ~'{{'+
		if !t.lex.MatchUntilSeq("{{") {
			t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
			return
		}

		t.lex.BuildTokenData(TYPE, tok)
	case '=':
		t.lex.BuildTokenNext(EQUALS, tok)
	case '{':
		t.lex.NextChar()

		if !t.lex.MatchChar('{') {
			t.lex.BuildToken(LBRACE, tok)
			return
		}

		// ~'}}'*
		if !t.lex.MatchUntilSeq("}}") {
			t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
			return
		}

		t.lex.MatchSeq("}}")
		t.lex.BuildTokenData(CODE_BLOCK, tok)
	case '}':
		t.lex.BuildTokenNext(RBRACE, tok)
	case '(':
		t.lex.BuildTokenNext(LPAREN, tok)
	case ')':
		t.lex.BuildTokenNext(RPAREN, tok)
	case runtime.EOFChar:
		t.lex.BuildToken(runtime.EOF, tok)
	default:
		t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
	}
}
