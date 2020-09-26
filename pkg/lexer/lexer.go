package lexer

import (
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

const (
	Eof rune = '\uffff'
	Err rune = '\ufffe'
	bom rune = '\ufeff'
)

type TokenType int

type Token struct {
	Type               TokenType
	Data               string
	StartRow, StartCol int32
	EndRow, EndCol     int32
}

type Lexer struct {
	pos, nextPos, tokenStart, mark, markNext int
	row, col, markRow, markCol               int32
	startRow, startCol, endRow, endCol       int32
	currCh, markCh                           rune
	input                                    []byte
}

func NewFromString(input string) *Lexer {
	l := &Lexer{
		input: []byte(input), row: 1, col: 0, // inc'd first time by NextChar
		startCol: 1, startRow: 1, endRow: 1, endCol: 1,
	}
	l.NextChar()
	return l
}

func NewFromReader(r io.Reader) (*Lexer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	l := &Lexer{input: b, row: 1, col: 0, // inc'd first time by NextChar
		startCol: 1, startRow: 1, endRow: 1, endCol: 1,
	}
	l.NextChar()
	return l, nil
}

func NewFromFile(filename string) (*Lexer, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewFromReader(f)
}

func (l *Lexer) readChar() (ch rune, size int) {
	ch, size = rune(l.input[l.pos]), 1

	if ch >= utf8.RuneSelf {
		// Is not a single byte wide, so fallback to full UTF8 decode
		ch, size = utf8.DecodeRune(l.input[l.pos:])
		if ch == utf8.RuneError {
			if size > 0 {
				// TODO: record illegal encoding error
				ch = Err
			} else {
				// TODO: record error
				ch = Err
			}
			return
		} else if ch == bom && l.pos > 0 {
			// TODO: record illegal byte order mark
			ch = Err
			return
		}
	}

	return
}

func (l *Lexer) CurrChar() rune {
	return l.currCh
}

func (l *Lexer) NextChar() rune {
	// Did we reach end of line on prev char?
	l.endRow, l.endCol = l.row, l.col
	if l.currCh == '\n' {
		l.row++
		l.col = 1
	} else if l.currCh != Eof {
		l.col++
	}

	l.pos = l.nextPos

	// Are we done?
	if l.nextPos >= len(l.input) {
		l.currCh = Eof
		return l.currCh
	}

	ch, size := l.readChar()
	l.currCh = ch

	l.nextPos += size
	return ch
}

func (l *Lexer) MarkPos() {
	l.mark, l.markNext, l.markCol, l.markRow = l.pos, l.nextPos, l.col, l.row
	l.markCh = l.currCh
}

func (l *Lexer) ResetPos() {
	l.pos, l.nextPos, l.col, l.row = l.mark, l.markNext, l.markCol, l.markRow
	l.currCh = l.markCh
}

// *** Build/Discard token ***

func (l *Lexer) BuildToken(tt TokenType, t *Token) {
	t.Type = tt
	t.StartRow, t.EndRow = l.startRow, l.endRow
	t.StartCol, t.EndCol = l.startCol, l.endCol

	l.DiscardTokenData()
}

func (l *Lexer) BuildTokenNext(tt TokenType, t *Token) {
	l.NextChar()
	l.BuildToken(tt, t)
}

func (l *Lexer) BuildTokenData(tt TokenType, t *Token) {
	t.Data = string(l.input[l.tokenStart:l.pos])
	l.BuildToken(tt, t)
}

func (l *Lexer) BuildTokenDataNext(tt TokenType, t *Token) {
	l.NextChar()
	l.BuildTokenData(tt, t)
}

func (l *Lexer) DiscardTokenData() {
	l.tokenStart, l.startRow, l.startCol = l.pos, l.row, l.col
}

func (l *Lexer) DiscardTokenDataNext() {
	l.NextChar()
	l.DiscardTokenData()
}

// *** Matchers ***

func (l *Lexer) MatchChar(char rune) bool {
	if l.CurrChar() != char {
		return false
	}

	l.NextChar()
	return true
}

func (l *Lexer) MatchCharExcept(char rune) bool {
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

func (l *Lexer) MatchCharExceptInRange(start, end rune) bool {
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

func (l *Lexer) MatchCharExceptInSeq(seq string) bool {
	ch := l.CurrChar()

	for _, c := range seq {
		if ch == c {
			return false
		}
	}

	l.NextChar()
	return true
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

func (l *Lexer) MatchUntilSeq(seq string) {
outer:
	for ch := l.CurrChar(); ch != Eof; ch = l.NextChar() {
		l.MarkPos()
		for _, c := range seq {
			if c != ch {
				continue outer
			}
			ch = l.NextChar()
		}
		l.ResetPos() // We don't match the chars in the seq itself
		break
	}
}
