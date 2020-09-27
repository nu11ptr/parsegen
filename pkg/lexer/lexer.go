package lexer

import (
	"io"
	"io/ioutil"
	"os"
	"unicode/utf8"
)

const (
	// EOF represents the end of the file
	EOF rune = '\uffff'
	// Err represents an error that occurred during lexing
	Err rune = '\ufffe'
	bom rune = '\ufeff'
)

// TokenType is an enum type for different token type values
type TokenType int

// Token represents a single token output by the lexer and contains the type of
// token, the start/end coordinates, and optionally the string data
type Token struct {
	Type               TokenType
	Data               string
	StartRow, StartCol int32
	EndRow, EndCol     int32
}

// A Tokenizer tokenizes an input stream. It typically will use a Lexer as the
// helper library, but it is not required to do so
type Tokenizer interface {
	NextToken(*Token)
}

// Lexer represents all the internal state needed to perform lexing
type Lexer struct {
	pos, nextPos, tokenStart, mark, markNext int
	row, col, markRow, markCol               int32
	startRow, startCol, endRow, endCol       int32
	currCh, markCh                           rune
	input                                    []byte
}

// NewFromBytes creates a new lexer from a byte array. The byte array should
// be backed by UTF-8 data
func NewFromBytes(input []byte) *Lexer {
	l := &Lexer{
		input: input, row: 1, col: 0, // inc'd first time by NextChar
		startCol: 1, startRow: 1, endRow: 1, endCol: 1,
	}
	l.NextChar()
	return l
}

// NewFromString creates a new lexer from string input data
func NewFromString(input string) *Lexer {
	return NewFromBytes([]byte(input))
}

// NewFromReader creates a new lexer from a reader. Since this there is no
// requirement for a "ReadCloser", it is the responsibility of the caller to
// close the reader if that is required
func NewFromReader(r io.Reader) (*Lexer, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewFromBytes(b), nil
}

// NewFromFile creates a new lexer from an input file
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

// CurrChar returns the current character, but in no way consumes it
func (l *Lexer) CurrChar() rune {
	return l.currCh
}

// NextChar returns the next character in the input buffer advancing any required
// position information and updating the current character. If the end of the
// input is reached, the EOF character will be returned and made current.
func (l *Lexer) NextChar() rune {
	// Did we reach end of line on prev char?
	l.endRow, l.endCol = l.row, l.col
	if l.currCh == '\n' {
		l.row++
		l.col = 1
	} else if l.currCh != EOF {
		l.col++
	}

	l.pos = l.nextPos

	// Are we done?
	if l.nextPos >= len(l.input) {
		l.currCh = EOF
		return l.currCh
	}

	ch, size := l.readChar()
	l.currCh = ch

	l.nextPos += size
	return ch
}

// MarkPos saves all position/current char information for possible later
// restoration. Each time it is called, it overwrites the previous mark, so
// it cannot be called recursively.
func (l *Lexer) MarkPos() {
	l.mark, l.markNext, l.markCol, l.markRow = l.pos, l.nextPos, l.col, l.row
	l.markCh = l.currCh
}

// ResetPos restores a previously marked location in the input data. Both MarkPos
// and ResetPos are only safe to be called without any invocations to Build or
// DiscardToken in between.
func (l *Lexer) ResetPos() {
	l.pos, l.nextPos, l.col, l.row = l.mark, l.markNext, l.markCol, l.markRow
	l.currCh = l.markCh
}

// *** Build/Discard token ***

// BuildToken builds a token with the given token type, but no data
func (l *Lexer) BuildToken(tt TokenType, t *Token) {
	t.Type = tt
	t.StartRow, t.EndRow = l.startRow, l.endRow
	t.StartCol, t.EndCol = l.startCol, l.endCol

	l.DiscardTokenData()
}

// BuildTokenNext builds a token with the given token type (but no data) after
// advancing to the next character in the input. This version must be used after
// matching using primitives instead of the match functions
func (l *Lexer) BuildTokenNext(tt TokenType, t *Token) {
	l.NextChar()
	l.BuildToken(tt, t)
}

// BuildTokenData builds a token with the given token type and string data
func (l *Lexer) BuildTokenData(tt TokenType, t *Token) {
	t.Data = string(l.input[l.tokenStart:l.pos])
	l.BuildToken(tt, t)
}

// BuildTokenDataNext builds a token with the given token type and string data after
// advancing to the next character in the input. This version must be used after
// matching using primitives instead of the match functions
func (l *Lexer) BuildTokenDataNext(tt TokenType, t *Token) {
	l.NextChar()
	l.BuildTokenData(tt, t)
}

// DiscardTokenData discards any matched characers and resets the start of the
// next potential token to the current position
func (l *Lexer) DiscardTokenData() {
	l.tokenStart, l.startRow, l.startCol = l.pos, l.row, l.col
}

// DiscardTokenDataNext discards any matched characters and resets the start of the
// next potential token to the current position after advancing to the next
// character in the input. This version must be used after matching using
// primitives instead of the match functions
func (l *Lexer) DiscardTokenDataNext() {
	l.NextChar()
	l.DiscardTokenData()
}

// *** Matchers ***

// MatchChar attempts to match the char given and returns true if it does or false
// otherwise
func (l *Lexer) MatchChar(char rune) bool {
	if l.CurrChar() != char {
		return false
	}

	l.NextChar()
	return true
}

// MatchCharExcept attempts to match any char except the one given and returns true
// if it does or false otherwise
func (l *Lexer) MatchCharExcept(char rune) bool {
	if l.CurrChar() == char {
		return false
	}

	l.NextChar()
	return true
}

// MatchCharInRange attempts to match any char between and inclusive of the start
// and end characters and returns true if it does or false otherwise
func (l *Lexer) MatchCharInRange(start, end rune) bool {
	ch := l.CurrChar()
	if ch < start || ch > end {
		return false
	}

	l.NextChar()
	return true
}

// MatchCharExceptInRange attempts to match any char except those between and
// inclusive of the start and end characters. It returns true if it does or
// false otherwise
func (l *Lexer) MatchCharExceptInRange(start, end rune) bool {
	ch := l.CurrChar()
	if ch >= start && ch <= end {
		return false
	}

	l.NextChar()
	return true
}

// MatchCharInSeq attempts to match a character if it is in the given
// sequence and returns true if it does or false otherwise
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

// MatchCharExceptInSeq attempts to match a character that is not in the given
// sequence and returns true if it does or false otherwise
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

// MatchSeq attempts to match the exact sequence of characers given and returns
// true if it does or false otherwise
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

// MatchUntilSeq attempts to match any character sequence except the one given.
// If it matches, it does not include the given sequence in the match. If it does
// not match, it will continue to try until it reach end of file.
func (l *Lexer) MatchUntilSeq(seq string) {
outer:
	for ch := l.CurrChar(); ch != EOF; ch = l.NextChar() {
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
