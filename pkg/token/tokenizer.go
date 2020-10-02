package token

import runtime "github.com/nu11ptr/parsegen/runtime/go"

// *** Potentially Generated ***

type Mode int

const (
	REGULAR Mode = iota
	CHAR_CLASS
)

const (
	// Mode: Regular

	// Char set
	RULE_NAME runtime.TokenType = iota + runtime.EOF + 1
	TOKEN_NAME

	// Sequences
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
	keywords = map[string]runtime.TokenType{
		"fragment": FRAGMENT,
		"skip":     SKIP_ACTION,
		"pushMode": PUSH_ACTION,
		"popMode":  POP_ACTION,
	}
)

type Tokenizer struct {
	lex *runtime.Lexer

	mode Mode
}

func New(lex *runtime.Lexer) *Tokenizer {
	return &Tokenizer{lex: lex, mode: REGULAR}
}

func (t *Tokenizer) processRuleName(tok *runtime.Token) bool {
	// [a-z]
	if !t.lex.MatchCharInRange('a', 'z') {
		return false
	}

	// [A-Za-z0-9_]*
	for t.lex.MatchCharInRange('A', 'Z') || t.lex.MatchCharInRange('a', 'z') ||
		t.lex.MatchCharInRange('0', '9') || t.lex.MatchChar('_') {
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

func (t *Tokenizer) processTokenName(tok *runtime.Token) bool {
	// [A-Z]
	if !t.lex.MatchCharInRange('A', 'Z') {
		return false
	}

	// [A-Za-z0-9_]*
	for t.lex.MatchCharInRange('A', 'Z') || t.lex.MatchCharInRange('a', 'z') ||
		t.lex.MatchCharInRange('0', '9') || t.lex.MatchChar('_') {
	}

	t.lex.BuildTokenData(TOKEN_NAME, tok)

	return true
}

func (t *Tokenizer) charClassNextToken(ch rune, tok *runtime.Token) {
	switch ch {
	case '\\':
		switch t.lex.NextChar() {
		case 'u':
			t.lex.NextChar()

			// HEX_DIGIT+
			t.lex.MarkPos()
			matched := false
			for t.lex.MatchCharInRange('A', 'F') || t.lex.MatchCharInRange('a', 'f') ||
				t.lex.MatchCharInRange('0', '9') {
				matched = true
			}
			if matched {
				t.lex.BuildTokenData(UNICODE_ESCAPE_CHAR, tok)
				break
			}

			// '{'
			t.lex.ResetPos()
			if !t.lex.MatchChar('{') {
				t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
				break
			}

			// HEX_DIGIT+
			matched = false
			for t.lex.MatchCharInRange('A', 'F') || t.lex.MatchCharInRange('a', 'f') ||
				t.lex.MatchCharInRange('0', '9') {
				matched = true
			}
			if !matched {
				t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
				break
			}

			if !t.lex.MatchChar('}') {
				t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
				break
			}

			t.lex.BuildTokenData(UNICODE_ESCAPE_CHAR, tok)
		default:
			t.lex.BuildTokenDataNext(ESCAPE_CHAR, tok)
		}
	case '-':
		t.lex.BuildTokenNext(DASH, tok)
	case ']':
		t.lex.BuildTokenNext(RBRACK, tok)
		t.mode = REGULAR
	default:
		t.lex.BuildTokenDataNext(BASIC_CHAR, tok)
	}
}

func (t *Tokenizer) NextToken(tok *runtime.Token) {
	if t.mode == CHAR_CLASS {
		t.charClassNextToken(t.lex.CurrChar(), tok)
		return
	}

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
	if t.processTokenName(tok) {
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

		t.lex.BuildTokenData(TOKEN_LIT, tok)
	case '-':
		t.lex.NextChar()

		// '>'
		if !t.lex.MatchChar('>') {
			t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
			return
		}

		t.lex.BuildToken(RARROW, tok)
	case '.':
		t.lex.BuildTokenNext(DOT, tok)
	case ':':
		t.lex.BuildTokenNext(COLON, tok)
	case ';':
		t.lex.BuildTokenNext(SEMI, tok)
	case '|':
		t.lex.BuildTokenNext(PIPE, tok)
	case '(':
		t.lex.BuildTokenNext(LPAREN, tok)
	case ')':
		t.lex.BuildTokenNext(RPAREN, tok)
	case '+':
		t.lex.BuildTokenNext(PLUS, tok)
	case '*':
		t.lex.BuildTokenNext(STAR, tok)
	case '?':
		t.lex.BuildTokenNext(QUEST_MARK, tok)
	case '~':
		t.lex.BuildTokenNext(TILDE, tok)
	case ',':
		t.lex.BuildTokenNext(COMMA, tok)
	case '[':
		t.lex.BuildTokenNext(LBRACK, tok)
		t.mode = CHAR_CLASS
	case runtime.EOFChar:
		t.lex.BuildToken(runtime.EOF, tok)
	default:
		t.lex.BuildTokenDataNext(runtime.ILLEGAL, tok)
	}
}
