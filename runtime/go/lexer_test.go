package runtime_test

import (
	"testing"

	runtime "github.com/nu11ptr/parsegen/runtime/go"
	"github.com/stretchr/testify/assert"
)

const (
	input = `abc
	de/* blah blah
	  */fghi😊`

	bogus runtime.TokenType = iota
)

func assertToken(t *testing.T, tok *runtime.Token, tt runtime.TokenType, data string,
	sr, sc, er, ec int) {

	assert.Equal(t, bogus, tok.Type, "Bad token type")
	assert.Equal(t, data, tok.Data, "Bad data")
	assert.Equal(t, int32(sr), tok.StartRow, "Bad start row")
	assert.Equal(t, int32(sc), tok.StartCol, "Bad start col")
	assert.Equal(t, int32(er), tok.EndRow, "Bad end row")
	assert.Equal(t, int32(ec), tok.EndCol, "Bad end col")
}

func TestLexer(t *testing.T) {
	lex := runtime.NewFromString(input)
	var tok runtime.Token

	t.Run("MatchChar", func(t *testing.T) {
		assert.False(t, lex.MatchChar('b'))
		assert.True(t, lex.MatchChar('a'))

		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "a", 1, 1, 1, 1)
	})

	t.Run("MatchCharExcept", func(t *testing.T) {
		assert.False(t, lex.MatchCharExcept('b'))
		assert.True(t, lex.MatchCharExcept('a'))

		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "b", 1, 2, 1, 2)
	})

	t.Run("MatchSeq", func(t *testing.T) {
		assert.False(t, lex.MatchSeq("c\n\tdf"))
		assert.True(t, lex.MatchSeq("c\n\tde"))

		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "c\n\tde", 1, 3, 2, 3)
	})

	t.Run("MatchUntilSeq/DiscardTokenData", func(t *testing.T) {
		lex.MatchUntilSeq("*/")
		lex.DiscardTokenData()

		assert.True(t, lex.MatchSeq("*/"))
		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "*/", 3, 4, 3, 5)
	})

	t.Run("MatchCharInRange", func(t *testing.T) {
		assert.False(t, lex.MatchCharInRange('m', 'z'))
		assert.True(t, lex.MatchCharInRange('a', 'f'))

		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "f", 3, 6, 3, 6)
	})

	t.Run("MatchCharExceptInRange", func(t *testing.T) {
		assert.False(t, lex.MatchCharExceptInRange('a', 'z'))
		assert.True(t, lex.MatchCharExceptInRange('a', 'f'))

		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "g", 3, 7, 3, 7)
	})

	t.Run("MatchCharInSeq", func(t *testing.T) {
		assert.False(t, lex.MatchCharInSeq("abcdefg"))
		assert.True(t, lex.MatchCharInSeq("hijklmnop"))

		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "h", 3, 8, 3, 8)
	})

	t.Run("MatchCharExceptInSeq", func(t *testing.T) {
		assert.False(t, lex.MatchCharExceptInSeq("abcdefghi"))
		assert.True(t, lex.MatchCharExceptInSeq("jklmnop"))

		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "i", 3, 9, 3, 9)
	})

	t.Run("Match Unicode", func(t *testing.T) {
		assert.Equal(t, lex.CurrChar(), '😊')

		lex.NextChar()
		lex.BuildTokenData(bogus, &tok)
		assertToken(t, &tok, bogus, "😊", 3, 10, 3, 10)
	})

	t.Run("Match EOF", func(t *testing.T) {
		assert.Equal(t, lex.CurrChar(), runtime.EOF)
	})
}
