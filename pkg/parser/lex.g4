lexer grammar lex
	;

fragment HEX_DIGIT: [A-Fa-f0-9];

fragment NAME: [A-Za-z0-9_]*;

RULE_NAME: [a-z] NAME;

TOKEN_NAME: [A-Z] NAME;

TOKEN_LIT
	: '\'' ('\\\'' | ~'\'')+ '\''
	; // TODO: Handle escape chars as fragment

// *** Skip ***

COMMENT: '//' ~[\r\n]* -> skip;

ML_COMMENT
	: '/*' .*? '*/' -> skip
	; // TODO: Will we support reluctant matchers?

WS: [ \t\r\n\f]+ -> skip;

// *** Keywords ***

FRAGMENT: 'fragment';

SKIP_ACTION: 'skip';

PUSH_ACTION: 'pushMode';

POP_ACTION: 'popMode';

// *** Basic Sequences ****

RARROW: '->';

DOT: '.';

COLON: ':';

SEMI: ';';

PIPE: '|';

LPAREN: '(';

RPAREN: ')';

PLUS: '+';

STAR: '*';

QUEST_MARK: '?';

TILDE: '~';

COMMA: ',';

LBRACK: '[' -> pushMode(CHAR_CLASS);

// *** Lexer: CHAR_CLASS ***

mode CHAR_CLASS
	;

UNICODE_ESCAPE_CHAR: '\\u' (HEX_DIGIT+ | '{' HEX_DIGIT+ '}');

ESCAPE_CHAR: '\\' .;

BASIC_CHAR: ~[\]\\\-];

DASH: '-';

RBRACK: ']' -> popMode;