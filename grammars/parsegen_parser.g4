parser grammar parsegen_parser
	;

options {
	tokenVocab = parsegen_lexer;
}

body: parser_decl code_blocks;

parser_decl: 'parser' '=' STRING;

code_blocks: 'code' '{' code_block* '}';

code_block: RULE_NAME TYPE? CODE_BLOCK;