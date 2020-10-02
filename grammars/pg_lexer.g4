lexer grammar pg_lexer
	;

RULE_NAME: [a-z] [A-Za-z0-9_.]*;

STRING: '\'' ('\\\'' | ~'\'')+ '\'';

TYPE: '->' ~'{{'+;

CODE_BLOCK: '{{' ~'}}'+ '}}';

// *** Skip ***

COMMENT: '//' ~[\r\n]* -> skip;

ML_COMMENT: '/*' ~'*/'* '*/' -> skip;

WS: [ \t\r\n\f]+ -> skip;

// *** Keywords ***

PARSER: 'parser';

CODE: 'code';

// *** Basic Sequences ****

EQUALS: '=';

LBRACE: '{';

RBRACE: '}';
