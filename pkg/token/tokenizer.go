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

func (l *Tokenizer) processRuleName(t *lexer.Token) bool {
	// [a-z]
	if !l.MatchCharInRange('a', 'z') {
		return false
	}

	// [A-Za-z0-9_]*
	for l.MatchCharInRange('A', 'Z') || l.MatchCharInRange('a', 'z') ||
		l.MatchCharInRange('0', '9') || l.MatchChar('_') {
	}

	l.BuildTokenData(RULE_NAME, t)

	// Possible conflicting keyword
	tt, ok := keywords[t.Data]
	if ok {
		t.Type = tt
		t.Data = ""
	}

	return true
}

func (l *Tokenizer) processTokenName(t *lexer.Token) bool {
	// [A-Z]
	if !l.MatchCharInRange('A', 'Z') {
		return false
	}

	// [A-Za-z0-9_]*
	for l.MatchCharInRange('A', 'Z') || l.MatchCharInRange('a', 'z') ||
		l.MatchCharInRange('0', '9') || l.MatchChar('_') {
	}

	l.BuildTokenData(TOKEN_NAME, t)

	return true
}

func (l *Tokenizer) charClassNextToken(ch rune, t *lexer.Token) {
	// Skip
	switch ch {
	case '/':
		ch = l.NextChar()

		switch ch {
		// '//'
		case '/':
			l.NextChar()

			// ~[\r\n]*
			for l.NotMatchCharInSeq("\r\n") {
			}
			l.DiscardTokenData()
		// '/*'
		case '*':
			l.NextChar()
			l.MatchUntilSeq("*/")
			l.DiscardTokenData()
		default:
			l.BuildTokenDataNext(ILLEGAL, t)
			return
		}
	// [ \t\r\n\f]+
	case ' ', '\t', '\r', '\n', '\f':
		ch = l.NextChar()

		for l.MatchCharInSeq(" \t\r\n\f") {
		}
		l.DiscardTokenData()
	}

	// ~[\]\\\-]
	if l.NotMatchCharInSeq("]\\-") {
		l.BuildTokenData(BASIC_CHAR, t)
		return
	}

	switch ch {
	case '\\':
		ch = l.NextChar()

		switch ch {
		case 'u':
			// HEX_DIGIT+
			l.MarkPos()
			matched := false
			for l.MatchCharInRange('A', 'F') || l.MatchCharInRange('a', 'f') ||
				l.MatchCharInRange('0', '9') {
				matched = true
			}
			if matched {
				l.BuildTokenData(UNICODE_ESCAPE_CHAR, t)
				break
			}

			// '{'
			l.ResetPos()
			if !l.MatchChar('{') {
				l.BuildTokenDataNext(ILLEGAL, t)
				break
			}

			// HEX_DIGIT+
			matched = false
			for l.MatchCharInRange('A', 'F') || l.MatchCharInRange('a', 'f') ||
				l.MatchCharInRange('0', '9') {
				matched = true
			}
			if !matched {
				l.BuildTokenDataNext(ILLEGAL, t)
				break
			}

			if !l.MatchChar('}') {
				l.BuildTokenDataNext(ILLEGAL, t)
				break
			}

			l.BuildTokenData(UNICODE_ESCAPE_CHAR, t)
		default:
			l.BuildTokenDataNext(ESCAPE_CHAR, t)
		}
	case '-':
		l.BuildTokenNext(DASH, t)
	case ']':
		l.BuildTokenNext(RBRACK, t)
		l.mode = REGULAR
	default:
		l.BuildTokenDataNext(ILLEGAL, t)
	}
}

func (l *Tokenizer) NextToken(t *lexer.Token) {
	ch := l.CurrChar()

	if l.mode == CHAR_CLASS {
		l.charClassNextToken(ch, t)
		return
	}

	// Skip
	switch ch {
	case '/':
		ch = l.NextChar()

		switch ch {
		// '//'
		case '/':
			l.NextChar()

			// ~[\r\n]*
			for l.NotMatchCharInSeq("\r\n") {
			}
			l.DiscardTokenData()
		// '/*'
		case '*':
			l.NextChar()
			l.MatchUntilSeq("*/")
			l.DiscardTokenData()
		default:
			l.BuildTokenDataNext(ILLEGAL, t)
			return
		}
	// [ \t\r\n\f]+
	case ' ', '\t', '\r', '\n', '\f':
		ch = l.NextChar()

		for l.MatchCharInSeq(" \t\r\n\f") {
		}
		l.DiscardTokenData()
	}

	if l.processRuleName(t) {
		return
	}
	if l.processTokenName(t) {
		return
	}

	switch ch {
	case '\'':
		l.NextChar()

		// ('\\\'' | ~'\'')+
		matched := false
		for l.MatchSeq("\\'") || l.NotMatchChar('\\') {
			matched = true
		}
		if !matched {
			l.BuildTokenDataNext(ILLEGAL, t)
			return
		}

		// '\''
		if !l.MatchChar('\'') {
			l.BuildTokenDataNext(ILLEGAL, t)
			return
		}

		l.BuildToken(TOKEN_LIT, t)
	case '-':
		ch = l.NextChar()

		// '>'
		if !l.MatchChar('>') {
			l.BuildTokenDataNext(ILLEGAL, t)
			return
		}

		l.BuildToken(RARROW, t)
	case '.':
		l.BuildTokenNext(DOT, t)
	case ':':
		l.BuildTokenNext(COLON, t)
	case ';':
		l.BuildTokenNext(SEMI, t)
	case '|':
		l.BuildTokenNext(PIPE, t)
	case '(':
		l.BuildTokenNext(LPAREN, t)
	case ')':
		l.BuildTokenNext(RPAREN, t)
	case '+':
		l.BuildTokenNext(PLUS, t)
	case '*':
		l.BuildTokenNext(STAR, t)
	case '?':
		l.BuildTokenNext(QUEST_MARK, t)
	case '~':
		l.BuildTokenNext(TILDE, t)
	case ',':
		l.BuildTokenNext(COMMA, t)
	case '[':
		l.BuildTokenNext(LBRACK, t)
	default:
		l.BuildTokenDataNext(ILLEGAL, t)
		l.mode = CHAR_CLASS
	}
}
