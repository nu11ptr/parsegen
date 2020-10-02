parser grammar pg_parser
	;

options {
	tokenVocab = pg_lexer;
}

body: parser_decl code_blocks EOF;

parser_decl: 'parser' '=' STRING;

code_blocks: 'code' '{' code_block* '}';

code_block: RULE_NAME TYPE? CODE_BLOCK;