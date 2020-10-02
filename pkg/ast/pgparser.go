package ast

import (
	"fmt"
	"strings"

	runtime "github.com/nu11ptr/parsegen/runtime/go"
)

type Body struct {
	Parser     string
	CodeBlocks []*CodeBlock
}

func NewBody(parser string, codeBlocks []*CodeBlock) *Body {
	// Slicing is safe because we know first and last char are single quote (aka ASCII)
	return &Body{Parser: parser[1 : len(parser)-1], CodeBlocks: codeBlocks}
}

func (b *Body) String() string {
	buff := strings.Builder{}
	buff.WriteString("Body:\n")

	buff.WriteString(strings.Repeat(" ", 1*spaces))
	buff.WriteString(fmt.Sprintf("└──Parser: %s\n", b.Parser))

	buff.WriteString(strings.Repeat(" ", 1*spaces))
	buff.WriteString("└──Code Blocks:\n")

	for _, block := range b.CodeBlocks {
		buff.WriteString(block.String(2))
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
	buff.WriteString("└──Code Blocks:\n")

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
