package ast

import (
	"fmt"
	"strings"

	runtime "github.com/nu11ptr/parsegen/runtime/go"
)

func parseString(str string) string {
	// Slicing is safe because we know first and last char are single quote (aka ASCII)
	return str[1 : len(str)-1]
}

type Body struct {
	Parser     string
	CodeBlocks *CodeBlocks
}

func NewBody(parser string, codeBlocks *CodeBlocks) *Body {
	return &Body{Parser: parseString(parser), CodeBlocks: codeBlocks}
}

func (b *Body) String() string {
	print := new(print)
	b.print(print)
	return print.String()
}

func (b *Body) print(print *print) {
	print.WriteString("Body")
	print.PushIndent()
	print.WriteStringPair("Parser", b.Parser)
	b.CodeBlocks.print(print)
	print.PopIndent()
}

type CodeBlocks struct {
	Language string
	Blocks   []*CodeBlock
}

func NewCodeBlocks(lang string, blocks []*CodeBlock) *CodeBlocks {
	return &CodeBlocks{Language: parseString(lang), Blocks: blocks}
}

func (c *CodeBlocks) String() string {
	print := new(print)
	c.print(print)
	return print.String()
}

func (c *CodeBlocks) print(print *print) {
	print.WriteString("Code Blocks")
	print.PushIndent()
	print.WriteStringPair("Language", c.Language)

	for _, block := range c.Blocks {
		block.print(print)
	}
	print.PopIndent()
}

type CodeBlock struct {
	Rule string
	Type string
	Code string
}

func NewCodeBlock(rule string, type_ *runtime.Token, code string) *CodeBlock {
	// It is safe to slice because we know prefixes and suffixes and that they are ASCII
	t := ""
	if type_ != nil {
		t = strings.TrimSpace(type_.Data[2:])
	}
	c := strings.TrimSpace(code[2 : len(code)-2])
	return &CodeBlock{Rule: rule, Type: t, Code: c}
}

func (c *CodeBlock) String() string {
	print := new(print)
	c.print(print)
	return print.String()
}

func (c *CodeBlock) print(print *print) {
	print.WriteString("Code Block")
	print.PushIndent()
	print.WriteStringPair("Rule", c.Rule)
	if c.Type != "" {
		print.WriteStringPair("Type", c.Type)
	}
	print.WriteStringPair("Code", fmt.Sprintf("{{ %s }}", c.Code))
	print.PopIndent()
}
