lexer grammar lex
	;

fragment HEX_DIGIT: [A-Fa-f0-9];

fragment NAME: [A-Za-z0-9_]*;

RULE_NAME: [a-z] NAME;

TOKEN_NAME: [A-Z] NAME;

TOKEN_LIT
	: '\'' ('\\\'' | ~'\'')+ '\''
	; // TODO: Handle escape chars as frag

LBRACK: '[' -> pushMode(CHAR_CLASS);

FRAGMENT: 'fragment';

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

// *** Lexer: CHAR_CLASS ***

mode CHAR_CLASS
	;

UNICODE_ESCAPE_CHAR: '\\u' (HEX_DIGIT+ | '{' HEX_DIGIT+ '}');

ESCAPE_CHAR: '\\' .;

BASIC_CHAR: ~[\]\\\-];

DASH: '-';

RBRACK: ']' -> popMode;