package lexer

import (
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

const (
	eof rune = '\uffff'
	err rune = '\ufffe'
	bom rune = '\ufeff'
)

type TokenType int

type Token struct {
	Type               TokenType
	Data               string
	StartCol, StartRow int32
	EndCol, EndRow     int32
}

type Lexer struct {
	pos, nextPos, tokenStart, mark, markNext       int
	col, row, startCol, startRow, markCol, markRow int32
	currCh                                         rune
	input                                          []byte
}

func NewFromString(input string) *Lexer {
	l := &Lexer{input: []byte(input), col: 0, row: 1, startCol: 1, startRow: 1}
	l.NextChar()
	return l
}

func NewFromReader(r io.Reader) (*Lexer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	l := &Lexer{input: b, col: 0, row: 1, startCol: 1, startRow: 1}
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
				ch = err
			} else {
				// TODO: record error
				ch = err
			}
			return
		} else if ch == bom && l.pos > 0 {
			// TODO: record illegal byte order mark
			ch = err
			return
		}
	}

	return
}

func (l *Lexer) CurrChar() rune {
	return l.currCh
}

func (l *Lexer) NextChar() rune {
	// Are we done?
	if l.nextPos >= len(l.input) {
		l.currCh = eof
		return l.currCh
	}

	// Did we reach end of line on prev char?
	if l.currCh == '\n' {
		l.row++
		l.col = 1
	} else {
		l.col++
	}

	l.pos = l.nextPos

	ch, size := l.readChar()
	l.currCh = ch

	l.nextPos += size
	return ch
}

func (l *Lexer) MarkPos() {
	l.mark, l.markNext, l.markCol, l.markRow = l.pos, l.nextPos, l.col, l.row
}

func (l *Lexer) ResetPos() {
	l.pos, l.nextPos, l.col, l.row = l.mark, l.markNext, l.markCol, l.markRow
}

func (l *Lexer) BuildToken(tt TokenType, t *Token) {
	t.Type = tt
	t.StartRow, t.EndRow = l.startRow, l.row
	t.StartCol, t.EndCol = l.startCol, l.col

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
