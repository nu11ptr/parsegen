package ast

import "strings"

type Body struct {
	Parser     string
	CodeBlocks []*CodeBlock
}

func NewBody(parser string, codeBlocks []*CodeBlock) *Body {
	// Slicing is safe because we know first and last char are single quote (aka ASCII)
	return &Body{Parser: parser[1 : len(parser)-1], CodeBlocks: codeBlocks}
}

type CodeBlock struct {
	Type string
	Code string
}

func NewCodeBlock(type_, code string) *CodeBlock {
	// It is safe to slice because we know prefixes and suffixes and that they are ASCII
	t := ""
	if type_ != "" {
		t = strings.TrimSpace(type_[2:])
	}
	c := strings.TrimSpace(code[2 : len(code)-2])
	return &CodeBlock{Type: t, Code: c}
}
