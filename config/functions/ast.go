package functions

type Node interface {
	NodeType() string
}

type Program struct {
	Body []Node
}

func (p *Program) NodeType() string { return "Program" }

type Declaration struct {
	Name  string
	Value Node
}

func (d *Declaration) NodeType() string { return "Declaration" }

type Assignment struct {
	Name  string
	Value Node
}

func (a *Assignment) NodeType() string { return "Assignment" }

type Identifier struct {
	Name string
}

func (i *Identifier) NodeType() string { return "Identifier" }

type Literal struct {
	Value string
}

func (l *Literal) NodeType() string { return "Literal" }

type BinaryExpression struct {
	Operator string
	Left     Node
	Right    Node
}

func (b *BinaryExpression) NodeType() string { return "BinaryExpression" }

type FunctionDeclaration struct {
	Name       string
	Parameters []string
	ReturnType string
	Body       []Node
}

func (f *FunctionDeclaration) NodeType() string { return "FunctionDeclaration" }

type FunctionCall struct {
	Name      string
	Arguments []Node
}

func (f *FunctionCall) NodeType() string { return "FunctionCall" }

type IfStatement struct {
	Condition  Node
	Consequent []Node
	ElseIfs    []ElseIfStatement
	Alternate  []Node
}

func (i *IfStatement) NodeType() string { return "IfStatement" }

type ElseIfStatement struct {
	Condition  Node
	Consequent []Node
}

func (e *ElseIfStatement) NodeType() string { return "ElseIfStatement" }

type WhileLoop struct {
	Condition Node
	Body      []Node
}

func (w *WhileLoop) NodeType() string { return "WhileLoop" }

type RepeatLoop struct {
	Body []Node
}

func (r *RepeatLoop) NodeType() string { return "RepeatLoop" }

type BreakStatement struct{}

func (b *BreakStatement) NodeType() string { return "BreakStatement" }

type ContinueStatement struct{}

func (c *ContinueStatement) NodeType() string { return "ContinueStatement" }

type ReturnStatement struct {
	Value Node
}

func (r *ReturnStatement) NodeType() string { return "ReturnStatement" }

type ParserState struct {
	tokens   []map[string]string
	position int
	length   int
}

// NewParser creates a new parser instance(same as the LexerState in lexer.go)
func NewParser(tokens []map[string]string) *ParserState {
	return &ParserState{
		tokens:   tokens,
		position: 0,
		length:   len(tokens),
	}
}

func (p *ParserState) current() map[string]string {
	if p.position >= p.length {
		return nil
	}
	return p.tokens[p.position]
}

func (p *ParserState) peek(offset int) map[string]string {
	pos := p.position + offset
	if pos >= p.length {
		return nil
	}
	return p.tokens[pos]
}

func (p *ParserState) advance() {
	p.position++
}

func (p *ParserState) expect(tokenType, tokenValue string) bool {
	token := p.current()
	if token == nil {
		return false
	}
	if tokenType != "" && token["Type"] != tokenType {
		return false
	}
	if tokenValue != "" && token["Value"] != tokenValue {
		return false
	}
	return true
}

func isOperator(value string) bool {
	operators := map[string]bool{
		"+": true, "-": true, "*": true, "/": true, "%": true,
		"==": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true,
		"&&": true, "||": true, "=": true,
	}
	return operators[value]
}

func (p *ParserState) parseExpression() Node {
	if p.position >= p.length {
		return nil
	}

	token := p.current()

	if token["Type"] == "NUMBER" || token["Type"] == "STRING" {
		p.advance()
		return &Literal{Value: token["Value"]}
	}

	if token["Type"] == "IDENTIFIER" {
		name := token["Value"]
		p.advance()
		if p.expect("PAREN", "(") {
			p.advance()
			args := p.parseArguments()
			if p.expect("PAREN", ")") {
				p.advance()
			}
			return &FunctionCall{Name: name, Arguments: args}
		}

		return &Identifier{Name: name}
	}

	if token["Type"] == "OPERATOR" && isOperator(token["Value"]) {
		op := token["Value"]
		p.advance()
		left := p.parseExpression()
		right := p.parseExpression()
		return &BinaryExpression{
			Operator: op,
			Left:     left,
			Right:    right,
		}
	}
	p.advance()
	return nil
}

func (p *ParserState) parseArguments() []Node {
	args := []Node{}

	for !p.expect("PAREN", ")") && p.position < p.length {
		arg := p.parseExpression()
		if arg != nil {
			args = append(args, arg)
		}

		if p.expect("COMMA", ",") {
			p.advance()
		}
	}

	return args
}

func (p *ParserState) parseBlock() []Node {
	if !p.expect("BRACE", "{") {
		return []Node{}
	}

	p.advance()

	blockTokens := []map[string]string{}
	depth := 1

	for p.position < p.length && depth > 0 {
		token := p.current()

		if token["Type"] == "BRACE" {
			if token["Value"] == "{" {
				depth++
			} else if token["Value"] == "}" {
				depth--
				if depth == 0 {
					break
				}
			}
		}

		if depth > 0 {
			blockTokens = append(blockTokens, token)
		}
		p.advance()
	}

	if p.expect("BRACE", "}") {
		p.advance()
	}
	return Parser(blockTokens).Body
}

func (p *ParserState) parseFunctionDeclaration() Node {
	p.advance()

	if !p.expect("IDENTIFIER", "") {
		return nil
	}

	name := p.current()["Value"]
	p.advance()

	params := []string{}
	if p.expect("PAREN", "(") {
		p.advance()

		for !p.expect("PAREN", ")") && p.position < p.length {
			if p.expect("IDENTIFIER", "") {
				params = append(params, p.current()["Value"])
				p.advance()
			}

			if p.expect("COMMA", ",") {
				p.advance()
			}
		}

		if p.expect("PAREN", ")") {
			p.advance()
		}
	}

	returnType := ""
	if p.expect("COLON", ":") {
		p.advance()

		if p.position < p.length {
			token := p.current()
			if token["Type"] == "IDENTIFIER" || token["Type"] == "KEYWORD" {
				returnType = token["Value"]
				p.advance()
			}
		}
	}

	body := p.parseBlock()

	return &FunctionDeclaration{
		Name:       name,
		Parameters: params,
		ReturnType: returnType,
		Body:       body,
	}
}

func (p *ParserState) parseDeclaration() Node {
	p.advance()

	if !p.expect("IDENTIFIER", "") {
		return nil
	}

	name := p.current()["Value"]
	p.advance()
	if !p.expect("OPERATOR", "=") {
		return nil
	}
	p.advance()

	value := p.parseExpression()

	return &Declaration{Name: name, Value: value}
}

func (p *ParserState) parseAssignment() Node {
	name := p.current()["Value"]
	p.advance()
	if !p.expect("OPERATOR", "=") {
		return nil
	}
	p.advance()

	value := p.parseExpression()

	return &Assignment{Name: name, Value: value}
}

func (p *ParserState) parseIfStatement() Node {
	p.advance()

	condition := p.parseExpression()
	consequent := p.parseBlock()

	elseIfs := []ElseIfStatement{}
	alternate := []Node{}
	for p.position < p.length {
		token := p.current()

		if token != nil && token["Type"] == "KEYWORD" && token["Value"] == "ya fir" {
			p.advance()
			elseIfCond := p.parseExpression()
			elseIfConsequent := p.parseBlock()

			elseIfs = append(elseIfs, ElseIfStatement{
				Condition:  elseIfCond,
				Consequent: elseIfConsequent,
			})
		} else if token != nil && token["Type"] == "KEYWORD" && token["Value"] == "ya" {
			p.advance()
			alternate = p.parseBlock()
			break
		} else {
			break
		}
	}

	return &IfStatement{
		Condition:  condition,
		Consequent: consequent,
		ElseIfs:    elseIfs,
		Alternate:  alternate,
	}
}

func (p *ParserState) parseWhileLoop() Node {
	p.advance()

	condition := p.parseExpression()
	body := p.parseBlock()

	return &WhileLoop{
		Condition: condition,
		Body:      body,
	}
}

func (p *ParserState) parseRepeatLoop() Node {
	p.advance()
	body := p.parseBlock()

	return &RepeatLoop{
		Body: body,
	}
}

func (p *ParserState) parseBreakStatement() Node {
	p.advance()
	return &BreakStatement{}
}

func (p *ParserState) parseContinueStatement() Node {
	p.advance()
	return &ContinueStatement{}
}

func (p *ParserState) parseReturnStatement() Node {
	p.advance()

	var value Node
	token := p.current()

	if token != nil && token["Type"] != "SEMICOLON" && !(token["Type"] == "BRACE" && token["Value"] == "}") {
		value = p.parseExpression()
	}

	return &ReturnStatement{
		Value: value,
	}
}

//main fn
func (p *ParserState) parse() *Program {
	program := &Program{Body: []Node{}}

	for p.position < p.length {
		token := p.current()

		if token == nil {
			break
		}

		var node Node
		if token["Type"] == "KEYWORD" {
			switch token["Value"] {
			case "firseKaro":
				node = p.parseFunctionDeclaration()
			case "ye":
				node = p.parseDeclaration()
			case "agar":
				node = p.parseIfStatement()
			case "jabtak":
				node = p.parseWhileLoop()
			case "dohraye":
				node = p.parseRepeatLoop()
			case "roko":
				node = p.parseBreakStatement()
			case "aage badho":
				node = p.parseContinueStatement()
			case "wapas bhejo":
				node = p.parseReturnStatement()
			default:
				p.advance()
				continue
			}
		} else if token["Type"] == "IDENTIFIER" && p.peek(1) != nil && p.peek(1)["Type"] == "OPERATOR" && p.peek(1)["Value"] == "=" {
			node = p.parseAssignment()
		} else if token["Type"] == "IDENTIFIER" {
			node = p.parseExpression()
		} else {
			p.advance()
			continue
		}

		if node != nil {
			program.Body = append(program.Body, node)
		}
	}

	return program
}
func Parser(tokens []map[string]string) *Program {
	parser := NewParser(tokens)
	return parser.parse()
}
