package lexer

import "unicode/utf8"

const (
	eof rune = '\uffff'
	err rune = '\ufffe'
	bom rune = '\ufeff'
)

type Token struct {
	Type               TokenType
	Data               string
	StartRow, StartCol int
	EndRow, EndCol     int
}

type Lexer struct {
	pos, tokenMark, mark                                 int
	mode                                                 Mode
	currCh                                               rune
	input                                                []byte
	startCol, endCol, startRow, endRow, markCol, markRow int
}

func (l *Lexer) readChar() (ch rune) {
	size := 1
	ch = rune(l.input[l.pos])

	if ch >= utf8.RuneSelf {
		// Is not a single byte wide, so fallback to full UTF8 decode
		ch, size = utf8.DecodeRune(l.input[l.pos:])
		if ch == utf8.RuneError {
			if size > 0 {
				// TODO: record illegal encoding error
				ch = err
			} else {
				// TODO: record error and return EOF???
				ch = err
			}
			return
		} else if ch == bom && l.pos > 0 {
			// TODO: record illegal byte order mark
			ch = err
			return
		}
	}

	l.pos += size
	return
}

func (l *Lexer) currChar() rune {
	return l.currCh
}

func (l *Lexer) nextChar() rune {
	return eof
}

func (l *Lexer) markPos() {
	l.mark = l.pos
}

func (l *Lexer) resetPos() {
	l.pos = l.mark
}

func (l *Lexer) buildToken(tt TokenType, t *Token) {
	t.Type = tt
	l.tokenMark = l.pos
}

func (l *Lexer) buildTokenNext(tt TokenType, t *Token) {
	l.nextChar()
	l.buildToken(tt, t)
}

func (l *Lexer) buildTokenData(tt TokenType, t *Token) {
	t.Type = tt
	t.Data = "" // TODO
	l.tokenMark = l.pos
}

func (l *Lexer) buildTokenDataNext(tt TokenType, t *Token) {
	l.nextChar()
	l.buildTokenData(tt, t)
}

func (l *Lexer) discardTokenData() {
	l.tokenMark = l.pos
}

func (l *Lexer) discardTokenDataNext() {
	l.nextChar()
	l.discardTokenData()
}

func (l *Lexer) matchSeq(seq string) bool {
	l.markPos()
	ch := l.currChar()

	for _, c := range seq {
		if ch != c {
			l.resetPos()
			return false
		}
		ch = l.nextChar()
	}
	return true
}

func (l *Lexer) matchChar(char rune) bool {
	if l.currChar() != char {
		return false
	}

	l.nextChar()
	return true
}

func (l *Lexer) notMatchChar(char rune) bool {
	if l.currChar() == char {
		return false
	}

	l.nextChar()
	return true
}

func (l *Lexer) matchCharInRange(start, end rune) bool {
	ch := l.currChar()
	if ch < start || ch > end {
		return false
	}

	l.nextChar()
	return true
}

func (l *Lexer) notMatchCharInRange(start, end rune) bool {
	ch := l.currChar()
	if ch >= start && ch <= end {
		return false
	}

	l.nextChar()
	return true
}

func (l *Lexer) matchCharInSeq(seq string) bool {
	ch := l.currChar()

	for _, c := range seq {
		if ch == c {
			l.nextChar()
			return true
		}
	}
	return false
}

func (l *Lexer) notMatchCharInSeq(seq string) bool {
	ch := l.currChar()

	for _, c := range seq {
		if ch == c {
			return false
		}
	}

	l.nextChar()
	return true
}

func (l *Lexer) matchUntilSeq(seq string) {
outer:
	for ch := l.currChar(); ch != eof; ch = l.nextChar() {
		for _, c := range seq {
			if c != ch {
				continue outer
			}
			ch = l.nextChar()
		}
		break
	}
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

func (l *Lexer) processRuleName(t *Token) bool {
	// [a-z]
	if !l.matchCharInRange('a', 'z') {
		return false
	}

	// [A-Za-z0-9_]*
	for l.matchCharInRange('A', 'Z') || l.matchCharInRange('a', 'z') ||
		l.matchCharInRange('0', '9') || l.matchChar('_') {
	}

	l.buildTokenData(RULE_NAME, t)

	// Possible conflicting keyword
	tt, ok := keywords[t.Data]
	if ok {
		t.Type = tt
		t.Data = ""
	}

	return true
}

func (l *Lexer) processTokenName(t *Token) bool {
	// [A-Z]
	if !l.matchCharInRange('A', 'Z') {
		return false
	}

	// [A-Za-z0-9_]*
	for l.matchCharInRange('A', 'Z') || l.matchCharInRange('a', 'z') ||
		l.matchCharInRange('0', '9') || l.matchChar('_') {
	}

	l.buildTokenData(TOKEN_NAME, t)

	return true
}

func (l *Lexer) charClassNextToken(ch rune, t *Token) {
	// Skip
	switch ch {
	case '/':
		ch = l.nextChar()

		switch ch {
		// '//'
		case '/':
			l.nextChar()

			// ~[\r\n]*
			for l.notMatchCharInSeq("\r\n") {
			}
			l.discardTokenData()
		// '/*'
		case '*':
			l.nextChar()
			l.matchUntilSeq("*/")
			l.discardTokenData()
		default:
			l.buildTokenDataNext(ILLEGAL, t)
			return
		}
	// [ \t\r\n\f]+
	case ' ', '\t', '\r', '\n', '\f':
		ch = l.nextChar()

		for l.matchCharInSeq(" \t\r\n\f") {
		}
		l.discardTokenData()
	}

	// ~[\]\\\-]
	if l.notMatchCharInSeq("]\\-") {
		l.buildTokenData(BASIC_CHAR, t)
		return
	}

	switch ch {
	case '\\':
		ch = l.nextChar()

		switch ch {
		case 'u':
			// HEX_DIGIT+
			l.markPos()
			matched := false
			for l.matchCharInRange('A', 'F') || l.matchCharInRange('a', 'f') ||
				l.matchCharInRange('0', '9') {
				matched = true
			}
			if matched {
				l.buildTokenData(UNICODE_ESCAPE_CHAR, t)
				break
			}

			// '{'
			l.resetPos()
			if !l.matchChar('{') {
				l.buildTokenDataNext(ILLEGAL, t)
				break
			}

			// HEX_DIGIT+
			matched = false
			for l.matchCharInRange('A', 'F') || l.matchCharInRange('a', 'f') ||
				l.matchCharInRange('0', '9') {
				matched = true
			}
			if !matched {
				l.buildTokenDataNext(ILLEGAL, t)
				break
			}

			if !l.matchChar('}') {
				l.buildTokenDataNext(ILLEGAL, t)
				break
			}

			l.buildTokenData(UNICODE_ESCAPE_CHAR, t)
		default:
			l.buildTokenDataNext(ESCAPE_CHAR, t)
		}
	case '-':
		l.buildTokenNext(DASH, t)
	case ']':
		l.buildTokenNext(RBRACK, t)
		l.mode = REGULAR
	default:
		l.buildTokenDataNext(ILLEGAL, t)
	}
}

func (l *Lexer) NextToken(t *Token) {
	ch := l.currChar()

	if l.mode == CHAR_CLASS {
		l.charClassNextToken(ch, t)
		return
	}

	// Skip
	switch ch {
	case '/':
		ch = l.nextChar()

		switch ch {
		// '//'
		case '/':
			l.nextChar()

			// ~[\r\n]*
			for l.notMatchCharInSeq("\r\n") {
			}
			l.discardTokenData()
		// '/*'
		case '*':
			l.nextChar()
			l.matchUntilSeq("*/")
			l.discardTokenData()
		default:
			l.buildTokenDataNext(ILLEGAL, t)
			return
		}
	// [ \t\r\n\f]+
	case ' ', '\t', '\r', '\n', '\f':
		ch = l.nextChar()

		for l.matchCharInSeq(" \t\r\n\f") {
		}
		l.discardTokenData()
	}

	if l.processRuleName(t) {
		return
	}
	if l.processTokenName(t) {
		return
	}

	switch ch {
	case '\'':
		l.nextChar()

		// ('\\\'' | ~'\'')+
		matched := false
		for l.matchSeq("\\'") || l.notMatchChar('\\') {
			matched = true
		}
		if !matched {
			l.buildTokenDataNext(ILLEGAL, t)
			return
		}

		// '\''
		if !l.matchChar('\'') {
			l.buildTokenDataNext(ILLEGAL, t)
			return
		}

		l.buildToken(TOKEN_LIT, t)
	case '-':
		ch = l.nextChar()

		// '>'
		if !l.matchChar('>') {
			l.buildTokenDataNext(ILLEGAL, t)
			return
		}

		l.buildToken(RARROW, t)
	case '.':
		l.buildTokenNext(DOT, t)
	case ':':
		l.buildTokenNext(COLON, t)
	case ';':
		l.buildTokenNext(SEMI, t)
	case '|':
		l.buildTokenNext(PIPE, t)
	case '(':
		l.buildTokenNext(LPAREN, t)
	case ')':
		l.buildTokenNext(RPAREN, t)
	case '+':
		l.buildTokenNext(PLUS, t)
	case '*':
		l.buildTokenNext(STAR, t)
	case '?':
		l.buildTokenNext(QUEST_MARK, t)
	case '~':
		l.buildTokenNext(TILDE, t)
	case ',':
		l.buildTokenNext(COMMA, t)
	case '[':
		l.buildTokenNext(LBRACK, t)
	default:
		l.buildTokenDataNext(ILLEGAL, t)
		l.mode = CHAR_CLASS
	}
}
