grammar lex_parser
	;

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
	| CHAR_CLASS
	;

char_class: '[' (INDIV_CHAR | CHAR_RANGE)+ ']';

// *** Lexer ***

fragment ALPHA_NUM: [A-Za-z0-9_];

RULE_NAME: [a-z] ALPHA_NUM*;

TOKEN_NAME: [A-Z] ALPHA_NUM*;

TOKEN_LIT: '\'' ('\\\'' | ~'\'')+ '\'';

CHAR_CLASS: '[' ~']' ']'; // TODO: Expand greatly
