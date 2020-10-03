package ast

import (
	"fmt"
	"strings"
)

const spaces = 3

type print struct {
	buff   strings.Builder
	indent int
}

func (p *print) PushIndent() { p.indent++ }

func (p *print) PopIndent() { p.indent-- }

func (p *print) WriteString(str string) {
	if p.indent > 0 {
		p.buff.WriteString(strings.Repeat(" ", p.indent*spaces))
		p.buff.WriteString("└──")
	}
	p.buff.WriteString(fmt.Sprintf("%s:\n", str))
}

func (p *print) WriteStringPair(str, str2 string) {
	if p.indent > 0 {
		p.buff.WriteString(strings.Repeat(" ", p.indent*spaces))
		p.buff.WriteString("└──")
	}
	p.buff.WriteString(fmt.Sprintf("%s: %s\n", str, str2))
}

func (p *print) String() string {
	return p.buff.String()
}
