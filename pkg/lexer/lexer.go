package lexer

import "unicode/utf8"

const (
	eof rune = '\uffff'
	err rune = '\ufffe'
	bom rune = '\ufeff'
)

type TokenType int

type Token struct {
	Type               TokenType
	Data               string
	StartRow, StartCol int
	EndRow, EndCol     int
}

type Lexer struct {
	pos, tokenMark, mark                                 int
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

func (l *Lexer) CurrChar() rune {
	return l.currCh
}

func (l *Lexer) NextChar() rune {
	return eof
}

func (l *Lexer) MarkPos() {
	l.mark = l.pos
}

func (l *Lexer) ResetPos() {
	l.pos = l.mark
}

func (l *Lexer) BuildToken(tt TokenType, t *Token) {
	t.Type = tt
	l.tokenMark = l.pos
}

func (l *Lexer) BuildTokenNext(tt TokenType, t *Token) {
	l.NextChar()
	l.BuildToken(tt, t)
}

func (l *Lexer) BuildTokenData(tt TokenType, t *Token) {
	t.Type = tt
	t.Data = "" // TODO
	l.tokenMark = l.pos
}

func (l *Lexer) BuildTokenDataNext(tt TokenType, t *Token) {
	l.NextChar()
	l.BuildTokenData(tt, t)
}

func (l *Lexer) DiscardTokenData() {
	l.tokenMark = l.pos
}

func (l *Lexer) DiscardTokenDataNext() {
	l.NextChar()
	l.DiscardTokenData()
}

func (l *Lexer) MatchSeq(seq string) bool {
	l.MarkPos()
	ch := l.CurrChar()

	for _, c := range seq {
		if ch != c {
			l.ResetPos()
			return false
		}
		ch = l.NextChar()
	}
	return true
}

func (l *Lexer) MatchChar(char rune) bool {
	if l.CurrChar() != char {
		return false
	}

	l.NextChar()
	return true
}

func (l *Lexer) NotMatchChar(char rune) bool {
	if l.CurrChar() == char {
		return false
	}

	l.NextChar()
	return true
}

func (l *Lexer) MatchCharInRange(start, end rune) bool {
	ch := l.CurrChar()
	if ch < start || ch > end {
		return false
	}

	l.NextChar()
	return true
}

func (l *Lexer) NotMatchCharInRange(start, end rune) bool {
	ch := l.CurrChar()
	if ch >= start && ch <= end {
		return false
	}

	l.NextChar()
	return true
}

func (l *Lexer) MatchCharInSeq(seq string) bool {
	ch := l.CurrChar()

	for _, c := range seq {
		if ch == c {
			l.NextChar()
			return true
		}
	}
	return false
}

func (l *Lexer) NotMatchCharInSeq(seq string) bool {
	ch := l.CurrChar()

	for _, c := range seq {
		if ch == c {
			return false
		}
	}

	l.NextChar()
	return true
}

func (l *Lexer) MatchUntilSeq(seq string) {
outer:
	for ch := l.CurrChar(); ch != eof; ch = l.NextChar() {
		for _, c := range seq {
			if c != ch {
				continue outer
			}
			ch = l.NextChar()
		}
		break
	}
}
