package pgparser

import (
	"github.com/nu11ptr/parsegen/pkg/ast"
	"github.com/nu11ptr/parsegen/pkg/pgtoken"
	runtime "github.com/nu11ptr/parsegen/runtime/go"
)

// *** Below will be automatically generated

type Parser struct {
	p *runtime.Parser

	bodyMap       map[int]*ast.Body
	parserDecl    map[int]*string
	codeBlocksMap map[int][]*ast.CodeBlock
	codeBlockMap  map[int]*ast.CodeBlock
}

func New(p *runtime.Parser) *Parser {
	return &Parser{
		p:             p,
		bodyMap:       make(map[int]*ast.Body, 8),
		parserDecl:    make(map[int]*string, 8),
		codeBlocksMap: make(map[int][]*ast.CodeBlock, 8),
		codeBlockMap:  make(map[int]*ast.CodeBlock, 8),
	}
}

// *** body ***

func (p *Parser) memoParseBody() *ast.Body {
	pos := p.p.Pos()
	body, ok := p.bodyMap[pos]
	if ok {
		return body
	}
	body = p.ParseBody()
	// Memoize what we did here in case this exact rule/position is needed again
	p.bodyMap[pos] = body
	return body
}

// ParseBody parses the "body" parser rule
func (p *Parser) ParseBody() *ast.Body {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### parser_decl ###
	parserDecl := p.memoParseParserDecl()
	if parserDecl == nil {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}

	// ### code_blocks ###
	codeBlocks := p.memoParseCodeBlocks()
	if codeBlocks == nil {
		// Rule failed - rollback
		p.p.SetPos(oldPos)
		return nil
	}

	// ### EOF ###
	eofTok := p.p.MatchTokenOrRollback(runtime.EOF, oldPos)
	if eofTok == nil {
		return nil
	}

	return ast.NewBody(*parserDecl, codeBlocks)
}

// *** parser_decl ***

func (p *Parser) memoParseParserDecl() *string {
	pos := p.p.Pos()
	parserDecl, ok := p.parserDecl[pos]
	if ok {
		return parserDecl
	}
	parserDecl = p.ParseParserDecl()
	// Memoize what we did here in case this exact rule/position is needed again
	p.parserDecl[pos] = parserDecl
	return parserDecl
}

// ParseParserDecl parses the "parser_decl" parser rule
func (p *Parser) ParseParserDecl() *string {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### 'parser' ###
	parserTok := p.p.MatchTokenOrRollback(pgtoken.PARSER, oldPos)
	if parserTok == nil {
		return nil
	}

	// ### '=' ###
	equalsTok := p.p.MatchTokenOrRollback(pgtoken.EQUALS, oldPos)
	if equalsTok == nil {
		return nil
	}

	// ### STRING ###
	stringTok := p.p.MatchTokenOrRollback(pgtoken.STRING, oldPos)
	if stringTok == nil {
		return nil
	}

	return &stringTok.Data
}

// *** code_blocks ***

func (p *Parser) memoParseCodeBlocks() []*ast.CodeBlock {
	pos := p.p.Pos()
	codeBlocks, ok := p.codeBlocksMap[pos]
	if ok {
		return codeBlocks
	}
	codeBlocks = p.ParseCodeBlocks()
	// Memoize what we did here in case this exact rule/position is needed again
	p.codeBlocksMap[pos] = codeBlocks
	return codeBlocks
}

// ParseCodeBlocks parses the "code_blocks" parser rule
func (p *Parser) ParseCodeBlocks() []*ast.CodeBlock {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### 'code' ###
	codeTok := p.p.MatchTokenOrRollback(pgtoken.CODE, oldPos)
	if codeTok == nil {
		return nil
	}

	// ### '{' ###
	lbraceTok := p.p.MatchTokenOrRollback(pgtoken.LBRACE, oldPos)
	if lbraceTok == nil {
		return nil
	}

	// ### code_block* ###
	codeBlocks := []*ast.CodeBlock{}
	for {
		codeBlock := p.memoParseCodeBlock()
		if codeBlock == nil {
			break
		}
		codeBlocks = append(codeBlocks, codeBlock)
	}

	// ### '}' ###
	rbraceTok := p.p.MatchTokenOrRollback(pgtoken.RBRACE, oldPos)
	if rbraceTok == nil {
		return nil
	}

	return codeBlocks
}

// *** code_block ***

func (p *Parser) memoParseCodeBlock() *ast.CodeBlock {
	pos := p.p.Pos()
	codeBlock, ok := p.codeBlockMap[pos]
	if ok {
		return codeBlock
	}
	codeBlock = p.ParseCodeBlock()
	// Memoize what we did here in case this exact rule/position is needed again
	p.codeBlockMap[pos] = codeBlock
	return codeBlock
}

// ParseCodeBlock parses the "code_block" parser rule
func (p *Parser) ParseCodeBlock() *ast.CodeBlock {
	// Rule can fail - might need to rollback
	oldPos := p.p.Pos()

	// ### RULE_NAME ###
	ruleNameTok := p.p.MatchTokenOrRollback(pgtoken.RULE_NAME, oldPos)
	if ruleNameTok == nil {
		return nil
	}

	// ### TYPE? ###
	typeTok := p.p.TryMatchToken(pgtoken.TYPE)

	// ### CODE_BLOCK ###
	codeBlockTok := p.p.MatchTokenOrRollback(pgtoken.CODE_BLOCK, oldPos)
	if codeBlockTok == nil {
		return nil
	}

	return ast.NewCodeBlock(ruleNameTok.Data, typeTok, codeBlockTok.Data)
}
