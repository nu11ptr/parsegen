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
	buff := strings.Builder{}

	buff.WriteString("Body:\n")
	buff.WriteString(strings.Repeat(" ", 1*spaces))
	buff.WriteString(fmt.Sprintf("└──Parser: %s\n", b.Parser))
	buff.WriteString(b.CodeBlocks.String(1))

	return buff.String()
}

type CodeBlocks struct {
	Language string
	Blocks   []*CodeBlock
}

func NewCodeBlocks(lang string, blocks []*CodeBlock) *CodeBlocks {
	return &CodeBlocks{Language: parseString(lang), Blocks: blocks}
}

func (c *CodeBlocks) String(indent int) string {
	buff := strings.Builder{}

	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString("└──Code Blocks:\n")

	buff.WriteString(strings.Repeat(" ", (indent+1)*spaces))
	buff.WriteString(fmt.Sprintf("└──Language: %s\n", c.Language))

	for _, block := range c.Blocks {
		buff.WriteString(block.String(indent + 1))
	}

	return buff.String()
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

func (c *CodeBlock) String(indent int) string {
	buff := strings.Builder{}
	buff.WriteString(strings.Repeat(" ", indent*spaces))
	buff.WriteString("└──Code Block:\n")

	buff.WriteString(strings.Repeat(" ", (indent+1)*spaces))
	buff.WriteString(fmt.Sprintf("└──Rule: %s\n", c.Rule))

	if c.Type != "" {
		buff.WriteString(strings.Repeat(" ", (indent+1)*spaces))
		buff.WriteString(fmt.Sprintf("└──Type: %s\n", c.Type))
	}

	buff.WriteString(strings.Repeat(" ", (indent+1)*spaces))
	buff.WriteString(fmt.Sprintf("└──Code: {{ %s }}\n", c.Code))
	return buff.String()
}
