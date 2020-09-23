parser grammar parse
	;

options {
	tokenVocab = lex;
}

top_level: (rule | lex_rule)* EOF;

// *** Parser parser ***

rule: RULE_NAME ':' rule_body ';';

rule_body: rule_sect+ ('|' rule_sect+)*;

rule_sect
	: ('(' rule_body ')' | RULE_NAME | TOKEN_NAME | TOKEN_LIT) suffix?
	;

suffix: '+' | '*' | '?';

// *** Lexer parser ***

lex_rule: 'fragment'? TOKEN_NAME ':' lex_rule_body ';';

lex_rule_body: lex_rule_sect+ ('|' lex_rule_sect+)*;

lex_rule_sect: lex_rule_part suffix? | '~' lex_rule_part;

lex_rule_part
	: '(' lex_rule_body ')'
	| TOKEN_NAME
	| TOKEN_LIT
	| '.'
	| char_set
	;

char_set: '[' (char | char_range)+ ']';

char: UNICODE_ESCAPE_CHAR | ESCAPE_CHAR | BASIC_CHAR;

char_range: char '-' char;
