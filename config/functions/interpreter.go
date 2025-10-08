package functions

import (
	"fmt"
	"os"
	"strconv"
)

type RuntimeValue interface {
	Type() string
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) Type() string { return "number" }

type StringValue struct {
	Value string
}

func (s *StringValue) Type() string { return "string" }

type BoolValue struct {
	Value bool
}

func (b *BoolValue) Type() string { return "bool" }

type NullValue struct{}

func (n *NullValue) Type() string { return "null" }

type FunctionValue struct {
	Parameters []string
	Body       []Node
	Env        *Environment
}

func (f *FunctionValue) Type() string { return "function" }

type Environment struct {
	parent    *Environment
	variables map[string]RuntimeValue
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		parent:    parent,
		variables: make(map[string]RuntimeValue),
	}
}

func (e *Environment) Define(name string, value RuntimeValue) RuntimeValue {
	e.variables[name] = value
	return value
}

func (e *Environment) Get(name string) (RuntimeValue, error) {
	if val, ok := e.variables[name]; ok {
		return val, nil
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil, fmt.Errorf("undefined variable: %s", name)
}

func (e *Environment) Set(name string, value RuntimeValue) error {
	if _, ok := e.variables[name]; ok {
		e.variables[name] = value
		return nil
	}
	if e.parent != nil {
		return e.parent.Set(name, value)
	}
	return fmt.Errorf("cannot assign to undefined variable: %s", name)
}

type ControlFlow struct {
	Type  string // "break", "continue", "return"
	Value RuntimeValue
}

type Interpreter struct {
	env         *Environment
	controlFlow *ControlFlow
}

func NewInterpreter() *Interpreter {
	env := NewEnvironment(nil)

	// Add built-in functions
	env.Define("bol", &FunctionValue{
		Parameters: []string{"value"},
		Body:       nil, // Native function
	})

	return &Interpreter{
		env:         env,
		controlFlow: nil,
	}
}

func (i *Interpreter) Evaluate(node Node) (RuntimeValue, error) {
	if i.controlFlow != nil {
		return &NullValue{}, nil
	}

	if node == nil {
		return &NullValue{}, nil
	}

	switch n := node.(type) {
	case *Program:
		return i.evalProgram(n)
	case *Declaration:
		return i.evalDeclaration(n)
	case *Assignment:
		return i.evalAssignment(n)
	case *Identifier:
		return i.evalIdentifier(n)
	case *Literal:
		return i.evalLiteral(n)
	case *BinaryExpression:
		return i.evalBinaryExpression(n)
	case *FunctionDeclaration:
		return i.evalFunctionDeclaration(n)
	case *FunctionCall:
		return i.evalFunctionCall(n)
	case *IfStatement:
		return i.evalIfStatement(n)
	case *WhileLoop:
		return i.evalWhileLoop(n)
	case *RepeatLoop:
		return i.evalRepeatLoop(n)
	case *BreakStatement:
		i.controlFlow = &ControlFlow{Type: "break"}
		return &NullValue{}, nil
	case *ContinueStatement:
		i.controlFlow = &ControlFlow{Type: "continue"}
		return &NullValue{}, nil
	case *ReturnStatement:
		val, err := i.Evaluate(n.Value)
		if err != nil {
			return nil, err
		}
		i.controlFlow = &ControlFlow{Type: "return", Value: val}
		return val, nil
	default:
		return nil, fmt.Errorf("unknown node type: %s", node.NodeType())
	}
}

func (i *Interpreter) evalProgram(p *Program) (RuntimeValue, error) {
	var lastValue RuntimeValue = &NullValue{}
	for _, node := range p.Body {
		val, err := i.Evaluate(node)
		if err != nil {
			return nil, err
		}
		lastValue = val
		if i.controlFlow != nil && i.controlFlow.Type == "return" {
			break
		}
	}
	return lastValue, nil
}

func (i *Interpreter) evalDeclaration(d *Declaration) (RuntimeValue, error) {
	value, err := i.Evaluate(d.Value)
	if err != nil {
		return nil, err
	}
	return i.env.Define(d.Name, value), nil
}

func (i *Interpreter) evalAssignment(a *Assignment) (RuntimeValue, error) {
	value, err := i.Evaluate(a.Value)
	if err != nil {
		return nil, err
	}
	err = i.env.Set(a.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (i *Interpreter) evalIdentifier(id *Identifier) (RuntimeValue, error) {
	return i.env.Get(id.Name)
}

func (i *Interpreter) evalLiteral(l *Literal) (RuntimeValue, error) {
	// Try to parse as number
	if num, err := strconv.ParseFloat(l.Value, 64); err == nil {
		return &NumberValue{Value: num}, nil
	}
	// Otherwise it's a string
	return &StringValue{Value: l.Value}, nil
}

func (i *Interpreter) evalBinaryExpression(b *BinaryExpression) (RuntimeValue, error) {
	left, err := i.Evaluate(b.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.Evaluate(b.Right)
	if err != nil {
		return nil, err
	}

	leftNum, leftIsNum := left.(*NumberValue)
	rightNum, rightIsNum := right.(*NumberValue)

	if leftIsNum && rightIsNum {
		switch b.Operator {
		case "+":
			return &NumberValue{Value: leftNum.Value + rightNum.Value}, nil
		case "-":
			return &NumberValue{Value: leftNum.Value - rightNum.Value}, nil
		case "*":
			return &NumberValue{Value: leftNum.Value * rightNum.Value}, nil
		case "/":
			if rightNum.Value == 0 {
				return nil, fmt.Errorf("division by zero")
			}
			return &NumberValue{Value: leftNum.Value / rightNum.Value}, nil
		case "%":
			if rightNum.Value == 0 {
				return nil, fmt.Errorf("modulo by zero")
			}
			leftInt := int64(leftNum.Value)
			rightInt := int64(rightNum.Value)
			return &NumberValue{Value: float64(leftInt % rightInt)}, nil
		case "<":
			return &BoolValue{Value: leftNum.Value < rightNum.Value}, nil
		case ">":
			return &BoolValue{Value: leftNum.Value > rightNum.Value}, nil
		case "<=":
			return &BoolValue{Value: leftNum.Value <= rightNum.Value}, nil
		case ">=":
			return &BoolValue{Value: leftNum.Value >= rightNum.Value}, nil
		case "==":
			return &BoolValue{Value: leftNum.Value == rightNum.Value}, nil
		case "!=":
			return &BoolValue{Value: leftNum.Value != rightNum.Value}, nil
		}
	}

	// String concatenation
	if b.Operator == "+" {
		leftStr := i.toString(left)
		rightStr := i.toString(right)
		return &StringValue{Value: leftStr + rightStr}, nil
	}

	return nil, fmt.Errorf("unsupported operator: %s", b.Operator)
}

func (i *Interpreter) evalFunctionDeclaration(f *FunctionDeclaration) (RuntimeValue, error) {
	fn := &FunctionValue{
		Parameters: f.Parameters,
		Body:       f.Body,
		Env:        i.env,
	}
	return i.env.Define(f.Name, fn), nil
}

func (i *Interpreter) evalFunctionCall(f *FunctionCall) (RuntimeValue, error) {
	// Handle built-in bol() function
	if f.Name == "bol" {
		for _, arg := range f.Arguments {
			val, err := i.Evaluate(arg)
			if err != nil {
				return nil, err
			}
			fmt.Println(i.toString(val))
		}
		return &NullValue{}, nil
	}

	// Get function from environment
	fnVal, err := i.env.Get(f.Name)
	if err != nil {
		return nil, err
	}

	fn, ok := fnVal.(*FunctionValue)
	if !ok {
		return nil, fmt.Errorf("%s is not a function", f.Name)
	}

	// Evaluate arguments
	args := make([]RuntimeValue, len(f.Arguments))
	for idx, arg := range f.Arguments {
		val, err := i.Evaluate(arg)
		if err != nil {
			return nil, err
		}
		args[idx] = val
	}

	// Create new environment for function execution
	funcEnv := NewEnvironment(fn.Env)
	for idx, param := range fn.Parameters {
		if idx < len(args) {
			funcEnv.Define(param, args[idx])
		} else {
			funcEnv.Define(param, &NullValue{})
		}
	}

	// Save current environment and control flow
	prevEnv := i.env
	prevFlow := i.controlFlow
	i.env = funcEnv
	i.controlFlow = nil

	// Execute function body
	var result RuntimeValue = &NullValue{}
	for _, node := range fn.Body {
		val, err := i.Evaluate(node)
		if err != nil {
			i.env = prevEnv
			i.controlFlow = prevFlow
			return nil, err
		}
		result = val
		if i.controlFlow != nil && i.controlFlow.Type == "return" {
			result = i.controlFlow.Value
			break
		}
	}

	// Restore environment and control flow
	i.env = prevEnv
	i.controlFlow = prevFlow

	return result, nil
}

func (i *Interpreter) evalIfStatement(ifStmt *IfStatement) (RuntimeValue, error) {
	condition, err := i.Evaluate(ifStmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(condition) {
		return i.evalBlock(ifStmt.Consequent)
	}

	// Check else-if statements
	for _, elseIf := range ifStmt.ElseIfs {
		cond, err := i.Evaluate(elseIf.Condition)
		if err != nil {
			return nil, err
		}
		if i.isTruthy(cond) {
			return i.evalBlock(elseIf.Consequent)
		}
	}

	// Execute else block
	if len(ifStmt.Alternate) > 0 {
		return i.evalBlock(ifStmt.Alternate)
	}

	return &NullValue{}, nil
}

func (i *Interpreter) evalWhileLoop(w *WhileLoop) (RuntimeValue, error) {
	var lastValue RuntimeValue = &NullValue{}

	for {
		condition, err := i.Evaluate(w.Condition)
		if err != nil {
			return nil, err
		}

		if !i.isTruthy(condition) {
			break
		}

		lastValue, err = i.evalBlock(w.Body)
		if err != nil {
			return nil, err
		}

		if i.controlFlow != nil {
			if i.controlFlow.Type == "break" {
				i.controlFlow = nil
				break
			}
			if i.controlFlow.Type == "continue" {
				i.controlFlow = nil
				continue
			}
			if i.controlFlow.Type == "return" {
				break
			}
		}
	}

	return lastValue, nil
}

func (i *Interpreter) evalRepeatLoop(r *RepeatLoop) (RuntimeValue, error) {
	var lastValue RuntimeValue = &NullValue{}

	for {
		val, err := i.evalBlock(r.Body)
		if err != nil {
			return nil, err
		}
		lastValue = val

		if i.controlFlow != nil {
			if i.controlFlow.Type == "break" {
				i.controlFlow = nil
				break
			}
			if i.controlFlow.Type == "continue" {
				i.controlFlow = nil
				continue
			}
			if i.controlFlow.Type == "return" {
				break
			}
		}
	}

	return lastValue, nil
}

func (i *Interpreter) evalBlock(nodes []Node) (RuntimeValue, error) {
	var lastValue RuntimeValue = &NullValue{}
	for _, node := range nodes {
		val, err := i.Evaluate(node)
		if err != nil {
			return nil, err
		}
		lastValue = val
		if i.controlFlow != nil {
			break
		}
	}
	return lastValue, nil
}

func (i *Interpreter) isTruthy(val RuntimeValue) bool {
	switch v := val.(type) {
	case *BoolValue:
		return v.Value
	case *NumberValue:
		return v.Value != 0
	case *StringValue:
		return v.Value != ""
	case *NullValue:
		return false
	default:
		return true
	}
}

func (i *Interpreter) toString(val RuntimeValue) string {
	switch v := val.(type) {
	case *NumberValue:
		if v.Value == float64(int(v.Value)) {
			return fmt.Sprintf("%d", int(v.Value))
		}
		return fmt.Sprintf("%v", v.Value)
	case *StringValue:
		return v.Value
	case *BoolValue:
		return fmt.Sprintf("%v", v.Value)
	case *NullValue:
		return "null"
	case *FunctionValue:
		return "<function>"
	default:
		return fmt.Sprintf("%v", val)
	}
}

func Run(code string) error {
	tokens := Lexer(code)
	ast := Parser(tokens)
	interpreter := NewInterpreter()
	_, err := interpreter.Evaluate(ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime Error: %v\n", err)
		return err
	}
	return nil
}
