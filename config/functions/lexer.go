package functions

type Token struct {
	Type  string
	Value string
}

const (
	TokenKeyword    = "KEYWORD"
	TokenIdentifier = "IDENTIFIER"
	TokenNumber     = "NUMBER"
	TokenString     = "STRING"
	TokenOperator   = "OPERATOR"
	TokenParen      = "PAREN"
	TokenBrace      = "BRACE"
	TokenComma      = "COMMA"
	TokenColon      = "COLON"
	TokenSemicolon  = "SEMICOLON"
)

// lexerState helps to encapsulate the state of the lexical analyzer
type LexerState struct {
	input    []rune
	position int
	length   int
	tokens   []map[string]string
}

// new Lexer instance
func NewLexer(input string) *LexerState {
	runes := []rune(input) //type rune = int32
	return &LexerState{
		input:    runes,
		position: 0,
		length:   len(runes),
		tokens:   make([]map[string]string, 0),
	}
}

//main fn
func (l *LexerState) Tokenize() []map[string]string {
	for l.position < l.length {
		if l.skipWhitespace() {
			continue
		}
		if l.skipComment() {
			continue
		}
		if l.scanString() {
			continue
		}
		if l.scanNumber() {
			continue
		}
		if l.scanIdentifierOrKeyword() {
			continue
		}
		if l.scanMultiCharOperator() {
			continue
		}
		if l.scanSingleCharToken() {
			continue
		}
		l.position++
	}
	return l.tokens
}

func (l *LexerState) current() rune {
	if l.position >= l.length {
		return 0
	}
	return l.input[l.position]
}

func (l *LexerState) peek(offset int) rune {
	pos := l.position + offset
	if pos >= l.length {
		return 0
	}
	return l.input[pos]
}

func (l *LexerState) advance() {
	l.position++
}

func (l *LexerState) addToken(tokenType, value string) {
	l.tokens = append(l.tokens, map[string]string{"Type": tokenType, "Value": value})
}

func (l *LexerState) skipWhitespace() bool {
	r := l.current()
	if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
		l.advance()
		return true
	}
	return false
}

func (l *LexerState) skipComment() bool {
	if l.current() == '/' && l.peek(1) == '/' {
		l.position += 2
		for l.position < l.length && l.current() != '\n' {
			l.advance()
		}
		return true
	}
	return false
}
func (l *LexerState) scanString() bool {
	quote := l.current()
	if quote != '"' && quote != '\'' {
		return false
	}

	l.advance()
	value := []rune{}

	for l.position < l.length {
		r := l.current()

		if r == '\\' && l.position+1 < l.length {
			value = append(value, l.peek(1))
			l.position += 2
			continue
		}

		if r == quote {
			l.advance()
			l.addToken(TokenString, string(value))
			return true
		}

		value = append(value, r)
		l.advance()
	}
	l.addToken(TokenString, string(value))
	return true
}

func (l *LexerState) scanNumber() bool {
	if !isDigit(l.current()) {
		return false
	}

	start := l.position
	for l.position < l.length && (isDigit(l.current()) || l.current() == '.') {
		l.advance()
	}

	l.addToken(TokenNumber, string(l.input[start:l.position]))
	return true
}

func (l *LexerState) scanIdentifierOrKeyword() bool {
	if !isLetter(l.current()) {
		return false
	}

	start := l.position
	l.advance()

	for l.position < l.length && (isLetter(l.current()) || isDigit(l.current())) {
		l.advance()
	}

	word := string(l.input[start:l.position])

	switch word {
	case "ya":
		word = l.tryConsumeMultiWordKeyword(word, "fir")
	case "aage":
		word = l.tryConsumeMultiWordKeyword(word, "badho")
	case "wapas":
		word = l.tryConsumeMultiWordKeyword(word, "bhejo")
	}

	if isKeyword(word) {
		l.addToken(TokenKeyword, word)
	} else {
		l.addToken(TokenIdentifier, word)
	}

	return true
}

func (l *LexerState) tryConsumeMultiWordKeyword(firstWord, secondWord string) string {
	savedPos := l.position

	for l.position < l.length && isWhitespace(l.current()) {
		l.advance()
	}

	if l.position < l.length && isLetter(l.current()) {
		start := l.position
		l.advance()

		for l.position < l.length && (isLetter(l.current()) || isDigit(l.current())) {
			l.advance()
		}

		nextWord := string(l.input[start:l.position])
		if nextWord == secondWord {
			return firstWord + " " + secondWord
		}

		l.position = savedPos
	} else {
		l.position = savedPos
	}

	return firstWord
}

func (l *LexerState) scanMultiCharOperator() bool {
	if l.position+1 >= l.length {
		return false
	}

	twoChar := string([]rune{l.current(), l.peek(1)})

	switch twoChar {
	case "==", "!=", "<=", ">=", "&&", "||":
		l.addToken(TokenOperator, twoChar)
		l.position += 2
		return true
	}

	return false
}

func (l *LexerState) scanSingleCharToken() bool {
	r := l.current()

	switch r {
	case '=', '+', '-', '*', '/', '%', '<', '>':
		l.addToken(TokenOperator, string(r))
		l.advance()
		return true
	case '(':
		l.addToken(TokenParen, "(")
		l.advance()
		return true
	case ')':
		l.addToken(TokenParen, ")")
		l.advance()
		return true
	case '{':
		l.addToken(TokenBrace, "{")
		l.advance()
		return true
	case '}':
		l.addToken(TokenBrace, "}")
		l.advance()
		return true
	case ',':
		l.addToken(TokenComma, ",")
		l.advance()
		return true
	case ':':
		l.addToken(TokenColon, ":")
		l.advance()
		return true
	case ';':
		l.addToken(TokenSemicolon, ";")
		l.advance()
		return true
	}

	return false
}
func isLetter(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_'
}
func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}

func isKeyword(word string) bool {
	keywords := map[string]bool{
		"ye":          true, // var/const
		"bol":         true, // tell(println)
		"agar":        true, // if
		"ya":          true,
		"fir":         true,
		"ya fir":      true, // else if
		"firseKaro":   true, //func
		"jabtak":      true, // while
		"dohraye":     true, // repeat
		"roko":        true, // break
		"aage badho":  true, // continue
		"aage":        true,
		"badho":       true,
		"wapas bhejo": true, // return
		"wapas":       true,
		"bhejo":       true,
	}
	return keywords[word]
}

func Lexer(input string) []map[string]string {
	lexer := NewLexer(input)
	return lexer.Tokenize()
}