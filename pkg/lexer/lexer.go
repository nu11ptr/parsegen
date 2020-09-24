package lexer

const (
	eof rune = -1
)

type Token struct {
	Type               TokenType
	Data               string
	StartRow, StartCol int
	EndRow, EndCol     int
}

type Lexer struct {
	pos, tokenMark, mark int
	mode                 Mode
}

func (l *Lexer) currChar() rune {
	return eof
}

func (l *Lexer) nextChar() rune {
	return eof
}

func (l *Lexer) buildToken(tt TokenType, t *Token) {
	t.Type = tt
	l.tokenMark = l.pos
}

func (l *Lexer) buildTokenData(tt TokenType, t *Token) {
	t.Type = tt
	t.Data = "" // TODO
	l.tokenMark = l.pos
}

func (l *Lexer) matchSeq(seq string) bool {
	return true
}

func (l *Lexer) matchChar(ch, char rune) bool {
	return ch == char
}

func (l *Lexer) notMatchChar(ch, char rune) bool {
	return ch != char
}

func (l *Lexer) matchCharRange(ch, start, end rune) bool {
	return ch >= start && ch <= end
}

func (l *Lexer) matchCharInSeq(ch rune, seq string) bool {
	for _, c := range seq {
		if ch == c {
			return true
		}
	}
	return false
}

func (l *Lexer) notMatchCharInSeq(ch rune, seq string) bool {
	return !l.matchCharInSeq(ch, seq)
}

// *** Potentially Generated ***

type Mode int

const (
	REGULAR Mode = iota
	CHAR_CLASS
)

type TokenType int

const (
	ILLEGAL TokenType = iota
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
	keywords = map[string]TokenType{
		"fragment": FRAGMENT,
		"skip":     SKIP_ACTION,
		"pushMode": PUSH_ACTION,
		"popMode":  POP_ACTION,
	}
)

func (l *Lexer) processRuleName(ch rune, t *Token) bool {
	// [a-z]
	if !l.matchCharRange(ch, 'a', 'z') {
		return false
	}
	ch = l.nextChar()

	// [A-Za-z0-9_]*
	for l.matchCharRange(ch, 'A', 'Z') || l.matchCharRange(ch, 'a', 'z') ||
		l.matchCharRange(ch, '0', '9') || l.matchChar(ch, '_') {
	}

	l.buildTokenData(RULE_NAME, t)

	// Possible conflicting keyword
	tt, ok := keywords[t.Data]
	if ok {
		t.Type = tt
	}

	return true
}

func (l *Lexer) processTokenName(ch rune, t *Token) bool {
	// [A-Z]
	if !l.matchCharRange(ch, 'A', 'Z') {
		return false
	}

	// [A-Za-z0-9_]*
	for l.matchCharRange(ch, 'A', 'Z') || l.matchCharRange(ch, 'a', 'z') ||
		l.matchCharRange(ch, '0', '9') || l.matchChar(ch, '_') {
	}

	l.buildTokenData(TOKEN_NAME, t)

	return true
}

func (l *Lexer) charClassNextToken(ch rune, t *Token) {
	// ~[\]\\\-]
	if l.matchCharInSeq(ch, "]\\-") {
		l.buildTokenData(BASIC_CHAR, t)
		return
	}

	switch ch {
	case '\\':
	case '-':
	case ']':
	}
}

func (l *Lexer) NextToken(t *Token) {
	ch := l.currChar()

	if l.mode == CHAR_CLASS {
		l.charClassNextToken(ch, t)
		return
	}

	if l.processRuleName(ch, t) {
		return
	}
	if l.processTokenName(ch, t) {
		return
	}

	switch ch {
	case '\'':
		ch = l.nextChar()

		// ('\\\'' | ~'\'')+
		matched := false
		for l.matchSeq("\\'") || l.notMatchChar(ch, '\\') {
			matched = true
		}
		if !matched {
			l.buildTokenData(ILLEGAL, t)
			break
		}

		// '\''
		if ch != '\'' {
			l.buildTokenData(ILLEGAL, t)
			break
		}

		l.buildToken(TOKEN_LIT, t)
	case '-':
		ch = l.nextChar()

		// '>'
		if ch != '>' {
			l.buildTokenData(ILLEGAL, t)
			break
		}

		l.buildToken(RARROW, t)
	case '.':
		l.buildToken(DOT, t)
	case ':':
		l.buildToken(COLON, t)
	case ';':
		l.buildToken(SEMI, t)
	case '|':
		l.buildToken(PIPE, t)
	case '(':
		l.buildToken(LPAREN, t)
	case ')':
		l.buildToken(RPAREN, t)
	case '+':
		l.buildToken(PLUS, t)
	case '*':
		l.buildToken(STAR, t)
	case '?':
		l.buildToken(QUEST_MARK, t)
	case '~':
		l.buildToken(TILDE, t)
	case ',':
		l.buildToken(COMMA, t)
	case '[':
		l.buildToken(LBRACK, t)
	default:
		l.buildTokenData(ILLEGAL, t)
		l.mode = CHAR_CLASS
	}
}
