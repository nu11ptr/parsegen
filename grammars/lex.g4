lexer grammar lex
	;

fragment HEX_DIGIT: [A-Fa-f0-9];

fragment NAME: [A-Za-z0-9_]*;

fragment COMMENT: '//' ~[\r\n]* -> skip;

fragment ML_COMMENT
	: '/*' .*? '*/' -> skip
	; // TODO: Will we support reluctant matchers?

fragment WS: [ \t\r\n\f]+ -> skip;

RULE_NAME: [a-z] NAME;

TOKEN_NAME: [A-Z] NAME;

TOKEN_LIT
	: '\'' ('\\\'' | ~'\'')+ '\''
	; // TODO: Handle escape chars as fragment

// *** Skip ***

REG_COMMENT: COMMENT;

REG_ML_COMMENT: ML_COMMENT;

REG_WS: WS;

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
