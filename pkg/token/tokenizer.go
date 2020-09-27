package token

// *** Potentially Generated ***

import "github.com/nu11ptr/parsegen/pkg/lexer"

type Mode int

const (
	REGULAR Mode = iota
	CHAR_CLASS
)

const (
	ILLEGAL lexer.TokenType = iota
	EOF

	// Mode: Regular

	// Char set
	RULE_NAME
	TOKEN_NAME
	TOKEN_LIT

	// Keywords
	FRAGMENT
	SKIP_ACTION
	PUSH_ACTION
	POP_ACTION

	// Basic Sequences
	RARROW
	DOT
	COLON
	SEMI
	PIPE
	LPAREN
	RPAREN
	PLUS
	STAR
	QUEST_MARK
	TILDE
	COMMA
	LBRACK

	// Mode: CHAR_CLASS

	// Char set
	BASIC_CHAR

	// Basic Sequences
	UNICODE_ESCAPE_CHAR
	ESCAPE_CHAR
	DASH
	RBRACK
)

var (
	keywords = map[string]lexer.TokenType{
		"fragment": FRAGMENT,
		"skip":     SKIP_ACTION,
		"pushMode": PUSH_ACTION,
		"popMode":  POP_ACTION,
	}
)

type Tokenizer struct {
	lexer.Lexer

	mode Mode
}

func (t *Tokenizer) processRuleName(tok *lexer.Token) bool {
	// [a-z]
	if !t.MatchCharInRange('a', 'z') {
		return false
	}

	// [A-Za-z0-9_]*
	for t.MatchCharInRange('A', 'Z') || t.MatchCharInRange('a', 'z') ||
		t.MatchCharInRange('0', '9') || t.MatchChar('_') {
	}

	t.BuildTokenData(RULE_NAME, tok)

	// Possible conflicting keyword
	tt, ok := keywords[tok.Data]
	if ok {
		tok.Type = tt
		tok.Data = ""
	}

	return true
}

func (t *Tokenizer) processTokenName(tok *lexer.Token) bool {
	// [A-Z]
	if !t.MatchCharInRange('A', 'Z') {
		return false
	}

	// [A-Za-z0-9_]*
	for t.MatchCharInRange('A', 'Z') || t.MatchCharInRange('a', 'z') ||
		t.MatchCharInRange('0', '9') || t.MatchChar('_') {
	}

	t.BuildTokenData(TOKEN_NAME, tok)

	return true
}

func (t *Tokenizer) charClassNextToken(ch rune, tok *lexer.Token) {
	// Skip
	switch ch {
	case '/':
		ch = t.NextChar()

		switch ch {
		// '//'
		case '/':
			t.NextChar()

			// ~[\r\n]*
			for t.MatchCharExceptInSeq("\r\n") {
			}
			t.DiscardTokenData()
		// '/*'
		case '*':
			t.NextChar()
			t.MatchUntilSeq("*/")
			t.DiscardTokenData()
		default:
			t.BuildTokenDataNext(ILLEGAL, tok)
			return
		}
	// [ \t\r\n\f]+
	case ' ', '\t', '\r', '\n', '\f':
		ch = t.NextChar()

		for t.MatchCharInSeq(" \t\r\n\f") {
		}
		t.DiscardTokenData()
	}

	// ~[\]\\\-]
	if t.MatchCharExceptInSeq("]\\-") {
		t.BuildTokenData(BASIC_CHAR, tok)
		return
	}

	switch ch {
	case '\\':
		ch = t.NextChar()

		switch ch {
		case 'u':
			// HEX_DIGIT+
			t.MarkPos()
			matched := false
			for t.MatchCharInRange('A', 'F') || t.MatchCharInRange('a', 'f') ||
				t.MatchCharInRange('0', '9') {
				matched = true
			}
			if matched {
				t.BuildTokenData(UNICODE_ESCAPE_CHAR, tok)
				break
			}

			// '{'
			t.ResetPos()
			if !t.MatchChar('{') {
				t.BuildTokenDataNext(ILLEGAL, tok)
				break
			}

			// HEX_DIGIT+
			matched = false
			for t.MatchCharInRange('A', 'F') || t.MatchCharInRange('a', 'f') ||
				t.MatchCharInRange('0', '9') {
				matched = true
			}
			if !matched {
				t.BuildTokenDataNext(ILLEGAL, tok)
				break
			}

			if !t.MatchChar('}') {
				t.BuildTokenDataNext(ILLEGAL, tok)
				break
			}

			t.BuildTokenData(UNICODE_ESCAPE_CHAR, tok)
		default:
			t.BuildTokenDataNext(ESCAPE_CHAR, tok)
		}
	case '-':
		t.BuildTokenNext(DASH, tok)
	case ']':
		t.BuildTokenNext(RBRACK, tok)
		t.mode = REGULAR
	default:
		t.BuildTokenDataNext(ILLEGAL, tok)
	}
}

func (t *Tokenizer) NextToken(tok *lexer.Token) {
	ch := t.CurrChar()

	if t.mode == CHAR_CLASS {
		t.charClassNextToken(ch, tok)
		return
	}

	// Skip
	switch ch {
	case '/':
		ch = t.NextChar()

		switch ch {
		// '//'
		case '/':
			t.NextChar()

			// ~[\r\n]*
			for t.MatchCharExceptInSeq("\r\n") {
			}
			t.DiscardTokenData()
		// '/*'
		case '*':
			t.NextChar()
			t.MatchUntilSeq("*/")
			t.DiscardTokenData()
		default:
			t.BuildTokenDataNext(ILLEGAL, tok)
			return
		}
	// [ \t\r\n\f]+
	case ' ', '\t', '\r', '\n', '\f':
		ch = t.NextChar()

		for t.MatchCharInSeq(" \t\r\n\f") {
		}
		t.DiscardTokenData()
	}

	if t.processRuleName(tok) {
		return
	}
	if t.processTokenName(tok) {
		return
	}

	switch ch {
	case '\'':
		t.NextChar()

		// ('\\\'' | ~'\'')+
		matched := false
		for t.MatchSeq("\\'") || t.MatchCharExcept('\\') {
			matched = true
		}
		if !matched {
			t.BuildTokenDataNext(ILLEGAL, tok)
			return
		}

		// '\''
		if !t.MatchChar('\'') {
			t.BuildTokenDataNext(ILLEGAL, tok)
			return
		}

		t.BuildToken(TOKEN_LIT, tok)
	case '-':
		ch = t.NextChar()

		// '>'
		if !t.MatchChar('>') {
			t.BuildTokenDataNext(ILLEGAL, tok)
			return
		}

		t.BuildToken(RARROW, tok)
	case '.':
		t.BuildTokenNext(DOT, tok)
	case ':':
		t.BuildTokenNext(COLON, tok)
	case ';':
		t.BuildTokenNext(SEMI, tok)
	case '|':
		t.BuildTokenNext(PIPE, tok)
	case '(':
		t.BuildTokenNext(LPAREN, tok)
	case ')':
		t.BuildTokenNext(RPAREN, tok)
	case '+':
		t.BuildTokenNext(PLUS, tok)
	case '*':
		t.BuildTokenNext(STAR, tok)
	case '?':
		t.BuildTokenNext(QUEST_MARK, tok)
	case '~':
		t.BuildTokenNext(TILDE, tok)
	case ',':
		t.BuildTokenNext(COMMA, tok)
	case '[':
		t.BuildTokenNext(LBRACK, tok)
	default:
		t.BuildTokenDataNext(ILLEGAL, tok)
		t.mode = CHAR_CLASS
	}
}
