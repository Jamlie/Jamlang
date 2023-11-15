package parser

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Jamlie/Jamlang/ast"
	"github.com/Jamlie/Jamlang/internal"
	"github.com/Jamlie/Jamlang/lexer"
	"github.com/Jamlie/Jamlang/tokentype"
)

type Parser struct {
	tokens     []lexer.Token
	isFunction bool
	isLoop     bool
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ProduceAST(sourceCode string) ast.Program {
	p.tokens = lexer.Tokenize(sourceCode)
	program := ast.Program{
		Body: []ast.Statement{},
	}

	for p.notEndOfFile() {
		if p.at().Type == tokentype.Whitespace {
			p.eat()
			continue
		}
		program.Body = append(program.Body, p.parseStatement())
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.at().Type {
	case tokentype.OpenComment:
		return p.parseComment()
	case tokentype.Let, tokentype.Constant:
		return p.parseVariableDeclaration()
	case tokentype.Function:
		return p.parseFunctionDeclaration()
	case tokentype.Return:
		if !p.isFunction {
			fmt.Fprintf(os.Stderr, "Error on line %d: Return statement outside of function\n", internal.Line())
			os.Exit(0)
		}
		return p.parseReturnStatement()
	// case tokentype.Class:
	// 	return p.parseClassDeclaration()
	case tokentype.Break:
		if !p.isLoop {
			fmt.Fprintf(os.Stderr, "Error on line %d: Break statement outside of loop\n", internal.Line())
			os.Exit(0)
		}
		return p.parseBreakStatement()
	case tokentype.If:
		return p.parseIfStatement()
	case tokentype.ElseIf:
		fmt.Fprintf(os.Stderr, "Error on line %d: Else if statement outside of if statement", internal.Line())
		os.Exit(0)
		return nil
	case tokentype.Else:
		fmt.Fprintf(os.Stderr, "Error on line %d: Else statement outside of if statement", internal.Line())
		os.Exit(0)
		return nil
	case tokentype.While:
		return p.parseWhileStatement()
	case tokentype.Loop:
		return p.parseLoopStatement()
	case tokentype.ForEach:
		return p.parseForEachStatement()
	case tokentype.For:
		return p.parseForStatement()
	case tokentype.Import:
		return p.parseImportStatement()
	case tokentype.SemiColon:
		p.eat()
		return &ast.NullLiteral{}
	default:
		return p.parseExpression()
}
}

func (p *Parser) parseComment() ast.Statement {
	for p.at().Type != tokentype.CloseComment {
		p.eat()
	}
	p.expect(tokentype.CloseComment, fmt.Sprintf("Error on line %d: Expected close comment", internal.Line()))
	return &ast.NullLiteral{}
}

func (p *Parser) parseImportStatement() ast.Statement {
	p.eat()
	path := p.expect(tokentype.String, fmt.Sprintf("Error on line %d: Expected string after import statement", internal.Line())).Value
	p.expect(tokentype.SemiColon, fmt.Sprintf("Error on line %d: Expected ';' after import statement", internal.Line()))
	return &ast.ImportStatement{Path: path}
}

func (p *Parser) parseForStatement() ast.Statement {
	p.eat()
	p.isLoop = true
	defer func() { p.isLoop = false }()
	init := p.parseStatement()
	p.expect(tokentype.SemiColon, fmt.Sprintf("Error on line %d: Expected ';' after for statement", internal.Line()))
	condition := p.parseExpression()
	p.expect(tokentype.SemiColon, fmt.Sprintf("Error on line %d: Expected ';' after for statement", internal.Line()))
	increment := p.parseExpression()
	p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected '{' after for statement", internal.Line()))
	var body []ast.Statement
	for p.at().Type != tokentype.RSquirly {
		body = append(body, p.parseStatement())
	}
	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected '}' after for statement", internal.Line()))
	return &ast.ForStatement{Init: init, Condition: condition, Update: increment, Body: body}
}

func (p *Parser) parseForEachStatement() ast.Statement {
	p.eat()

	p.isLoop = true
	defer func() { p.isLoop = false }()

	value := p.expect(tokentype.Identifier, fmt.Sprintf("Error on line %d: Expected identifier in for each statement", internal.Line())).Value
	if p.at().Type == tokentype.Comma {
		key := value
		p.eat()
		val := p.expect(tokentype.Identifier, fmt.Sprintf("Error on line %d: Expected identifier in for each statement", internal.Line())).Value
		p.expect(tokentype.In, fmt.Sprintf("Error on line %d: Expected in after identifier in for each statement", internal.Line()))
		obj := p.parseExpression()

		p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after for each statement", internal.Line()))

		var body []ast.Statement
		for p.at().Type != tokentype.RSquirly {
			body = append(body, p.parseStatement())
		}

		p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after for each statement", internal.Line()))

		return &ast.ForEachStatement{Key: key, Value: val, Variable: "", Collection: obj, Body: body}
	}
	p.expect(tokentype.In, fmt.Sprintf("Error on line %d: Expected in after identifier in for each statement", internal.Line()))
	array := p.parseExpression()

	p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after for each statement", internal.Line()))

	var body []ast.Statement
	for p.at().Type != tokentype.RSquirly {
		body = append(body, p.parseStatement())
	}

	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after for each statement", internal.Line()))

	return &ast.ForEachStatement{Variable: value, Key: "", Value: "", Collection: array, Body: body}
}

func (p *Parser) parseLoopStatement() ast.Statement {
	p.eat()
	p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after loop statement", internal.Line()))

	p.isLoop = true
	defer func() { p.isLoop = false }()
	var body []ast.Statement
	for p.at().Type != tokentype.RSquirly {
		body = append(body, p.parseStatement())
	}

	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after loop statement", internal.Line()))

	return &ast.LoopStatement{Body: body}
}

func (p *Parser) parseWhileStatement() ast.Statement {
	p.eat()
	condition := p.parseExpression()
	p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after while statement", internal.Line()))

	p.isLoop = true
	defer func() { p.isLoop = false }()
	var body []ast.Statement
	for p.at().Type != tokentype.RSquirly {
		body = append(body, p.parseStatement())
	}

	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after while statement", internal.Line()))

	return &ast.WhileStatement{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseIfStatement() ast.Statement {
	p.eat()
	condition := p.parseExpression()
	p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after if statement", internal.Line()))

	var body []ast.Statement
	for p.at().Type != tokentype.RSquirly {
		body = append(body, p.parseStatement())
	}

	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after if statement", internal.Line()))

	var elseifCondition []ast.Expression
	var elseifBody [][]ast.Statement
	if p.at().Type == tokentype.ElseIf {
		p.eat()
		elseifCondition = append(elseifCondition, p.parseExpression())
		p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after else if statement", internal.Line()))

		var elseifBodyTemp []ast.Statement
		for p.at().Type != tokentype.RSquirly {
			elseifBodyTemp = append(elseifBodyTemp, p.parseStatement())
		}
		elseifBody = append(elseifBody, elseifBodyTemp)
		p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after else if statement", internal.Line()))

		for p.at().Type == tokentype.ElseIf {
			p.eat()
			elseifCondition = append(elseifCondition, p.parseExpression())
			p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after else if statement", internal.Line()))

			var elseifBodyTemp []ast.Statement
			for p.at().Type != tokentype.RSquirly {
				elseifBodyTemp = append(elseifBodyTemp, p.parseStatement())
			}
			elseifBody = append(elseifBody, elseifBodyTemp)
			p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after else if statement", internal.Line()))
		}
		// p.eat()
		// elseifCondition = p.parseExpression()
		// p.expect(tokentype.LSquirly, "Error: Expected { after else if statement")
		//
		// for p.at().Type != tokentype.RSquirly {
		//     var elseifBodyTemp []ast.Statement
		//     for p.at().Type != tokentype.RSquirly {
		//         elseifBodyTemp = append(elseifBodyTemp, p.parseStatement())
		//     }
		//     elseifBody = append(elseifBody, elseifBodyTemp)
		//     p.expect(tokentype.RSquirly, "Error: Expected } after else if statement")
		//     if p.at().Type == tokentype.ElseIf {
		//         p.eat()
		//         elseifCondition = p.parseExpression()
		//         p.expect(tokentype.LSquirly, "Error: Expected { after else if statement")
		//     }
		// }
		//
		// p.expect(tokentype.RSquirly, "Error: Expected } after else if statement")
	}

	if p.at().Type == tokentype.Else {
		p.eat()
		p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after else statement", internal.Line()))

		var elseBody []ast.Statement
		for p.at().Type != tokentype.RSquirly {
			elseBody = append(elseBody, p.parseStatement())
		}

		p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after else statement", internal.Line()))

		return &ast.ConditionalStatement{
			Condition: condition,
			Body:      body,
			Alternate: elseBody,
			ElseIfBodies:    elseifBody,
			ElseIfConditions: elseifCondition,
		}
	}

	return &ast.ConditionalStatement{
		Condition: condition,
		Body:      body,
		ElseIfConditions: elseifCondition,
		ElseIfBodies:    elseifBody,
	}
}

// func (p *Parser) parseClassDeclaration() ast.Statement {
// 	p.eat()
// 	name := p.expect(tokentype.Identifier, fmt.Sprintf("Error on line %d: Expected class name after class keyword", internal.Line())).Value
// 	p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected { after class name", internal.Line()))
//
// 	var body []ast.Statement
// 	for p.at().Type != tokentype.RSquirly {
// 		body = append(body, p.parseStatement())
// 	}
//
// 	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected } after class declaration", internal.Line()))
//
// 	return &ast.ClassDeclaration{
// 		Name: name,
// 		Body: body,
// 	}
// }

func (p *Parser) parseFunctionDeclaration() ast.Statement {
	p.eat()
	var name string

	if p.at().Type != tokentype.OpenParen {
		name = p.expect(tokentype.Identifier, fmt.Sprintf("Error on line %d: Expected function name after fn keyword", internal.Line())).Value
	}

	args := p.parseArgs()
	var params []string

	for _, arg := range args {
		if arg.Kind() != ast.IdentifierType {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected function parameter to be of type string, got %s\n", internal.Line(), arg.Kind())
			os.Exit(0)
			return nil
		}

		params = append(params, arg.(*ast.Identifier).Symbol)
	}

	returnType := ast.AnyType

	if p.at().Type == tokentype.Colon {
		p.eat()
		var err error
		returnType, err = p.parseType()
		if err != nil {
			fmt.Printf("Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
		}
	}

	p.expect(tokentype.LSquirly, fmt.Sprintf("Error on line %d: Expected '{' after function declaration", internal.Line()))

	if !p.isFunction {
		p.isFunction = true
		defer func() { p.isFunction = false }()
	}

	body := []ast.Statement{}
	for p.at().Type != tokentype.EndOfFile && p.at().Type != tokentype.RSquirly {
		body = append(body, p.parseStatement())
	}

	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Expected '}' after function declaration", internal.Line()))

	return &ast.FunctionDeclaration{
		Name:       name,
		Parameters: params,
		Body:       body,
		ReturnType: returnType,
	}
}

func (p *Parser) parseBreakStatement() ast.Statement {
	p.eat()
	return &ast.BreakStatement{}
}

func (p *Parser) parseReturnStatement() ast.Statement {
	p.eat()
	return &ast.ReturnStatement{
		Value: p.parseExpression(),
	}
}

var Types = map[string]ast.VariableType{
	"str":    ast.StringType,
	"i8":     ast.Int8Type,
	"i16":    ast.Int16Type,
	"i32":    ast.Int32Type,
	"i64":    ast.Int64Type,
	"f32":    ast.Float32Type,
	"f64":    ast.Float64Type,
	"bool":   ast.BoolType,
	"list":   ast.ArrayType,
	"tuple":  ast.TupleType,
	"fn":     ast.FunctionType,
	"object": ast.ObjectType,
	"any":    ast.AnyType,
}

func (p *Parser) parseType() (ast.VariableType, error) {
	if p.at().Type == tokentype.Identifier {
		if t, ok := Types[p.at().Value]; ok {
			p.eat()
			return t, nil
		}
	}

	return ast.VariableType(""), fmt.Errorf("Error on line %d: Expected type", internal.Line())
}

func (p *Parser) parseVariableDeclaration() ast.Statement {
	isConstant := p.eat().Type == tokentype.Constant
	identifier := p.expect(tokentype.Identifier, fmt.Sprintf("Error on line %d: Expected identifier name after let/const keyword", internal.Line())).Value


	if p.at().Type == tokentype.SemiColon {
		p.eat()
		if isConstant {
			fmt.Fprintf(os.Stderr, "Error on line %d: Constant declaration without assignment is not allowed\n", internal.Line())
			os.Exit(0)
			return nil
		}

		return &ast.VariableDeclaration{
			Identifier: identifier,
			Constant:   isConstant,
			Value:      &ast.NullLiteral{},
			Type:       ast.AnyType,
		}
	}

	var varType = ast.AnyType

	if p.at().Type == tokentype.Colon {
		var err error
		p.eat()
		varType, err = p.parseType()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: Unknown type\n", internal.Line())
			os.Exit(0)
		}

		if isConstant && p.at().Type == tokentype.SemiColon {
			fmt.Fprintf(os.Stderr, "Error on line %d: Constant declaration without assignment is not allowed", internal.Line())
			os.Exit(0)
			return nil
		}
		if p.at().Type == tokentype.SemiColon {
			p.eat()
			return &ast.VariableDeclaration{
				Identifier:        identifier,
				Constant:          isConstant,
				Value:             &ast.NullLiteral{},
				Type:              varType,
				IsUserDefinedType: false,
			}
		}
	}

	p.expect(tokentype.Equals, fmt.Sprintf("Error on line %d: Expected = after identifier name", internal.Line()))

	var declaration ast.Statement
	declaration = &ast.VariableDeclaration{
		Identifier:        identifier,
		Constant:          isConstant,
		Value:             p.parseExpression(),
		Type:              varType,
		IsUserDefinedType: false,
	}

	if !p.isLoop {
		if p.at().Type == tokentype.SemiColon {
			p.eat()
		}
	}

	return declaration
}

func (p *Parser) parseExpression() ast.Expression {
	return p.parseAssignmentExpression()
}

func (p *Parser) parseLogicalExpression() ast.Expression {
	return p.parseOrExpression()
}

func (p *Parser) parseOrExpression() ast.Expression {
	left := p.parseAndExpression()

	for p.at().Value == "or" {
		p.eat()
		right := p.parseAndExpression()
		left = &ast.LogicalExpression{
			Operator: "or",
			Left:     left,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseAndExpression() ast.Expression {
	left := p.parseXorExpression()

	for p.at().Value == "and" {
		p.eat()
		right := p.parseXorExpression()
		left = &ast.LogicalExpression{
			Operator: "and",
			Left:     left,
			Right:    right,
		}
	}
	return left
}

func (p *Parser) parseXorExpression() ast.Expression {
	left := p.parseNotExpression()

	for p.at().Value == "xor" {
		p.eat()
		right := p.parseNotExpression()
		left = &ast.LogicalExpression{
			Operator: "xor",
			Left:     left,
			Right:    right,
		}
	}
	return left
}

func (p *Parser) parseNotExpression() ast.Expression {
	if p.at().Value == "not" {
		p.eat()
		return &ast.LogicalExpression{
			Operator: "not",
			Right:    p.parseNotExpression(),
		}
	}
	return p.parseComparisonExpression()
}

func (p *Parser) parseArrayExpression() ast.Expression {
	if p.at().Type != tokentype.OpenBracket {
		return p.parseBitwise()
	}
	p.eat()
	elements := []ast.Expression{}
	for p.at().Type != tokentype.CloseBracket {
		elements = append(elements, p.parseExpression())
		if p.at().Type == tokentype.Comma {
			p.eat()
		}
	}
	p.expect(tokentype.CloseBracket, fmt.Sprintf("Error on line %d: Expected closing bracket after array expression", internal.Line()))
	return &ast.ArrayLiteral{Elements: elements}
}

func (p *Parser) parseAssignmentExpression() ast.Expression {
	left := p.parseOrExpression()

	if p.at().Type == tokentype.Equals {
		p.eat()
		value := p.parseAssignmentExpression()
		return &ast.AssignmentExpression{
			Assignee: left,
			Value:   value,
		}
	}

	return left
}

func (p *Parser) parseObjectExpression() ast.Expression {
	if p.at().Type != tokentype.LSquirly {
		return p.parseArrayExpression()
	}

	p.eat()
	properties := []ast.Property{}

	for p.notEndOfFile() && p.at().Type != tokentype.RSquirly {
		if p.at().Type == tokentype.Integer {
			v := p.parsePrimaryExpression()
			key := v.(*ast.NumericIntegerLiteral)
			p.expect(tokentype.Colon, fmt.Sprintf("Error on line %d: Expected : after object key", internal.Line()))
			value := p.parseExpression()
			properties = append(properties, ast.Property{
				Key:   strconv.Itoa(int(key.Value)),
				Value: value,
			})
		} else if p.at().Type == tokentype.Float {
			fmt.Fprintf(os.Stderr, "Error on line %d: Floats are not allowed as object keys", internal.Line())
			os.Exit(0)
		} else if p.at().Type == tokentype.String {
			v := p.parsePrimaryExpression()
			key := v.(*ast.StringLiteral)
			p.expect(tokentype.Colon, fmt.Sprintf("Error on line %d: Expected : after object key", internal.Line()))
			value := p.parseExpression()
			properties = append(properties, ast.Property{
				Key:   key.Value,
				Value: value,
			})
		} else if p.at().Type == tokentype.Identifier {
			key := p.expect(tokentype.Identifier, fmt.Sprintf("Error on line %d: Expected identifier as object key", internal.Line())).Value
			if p.at().Type == tokentype.Comma {
				p.eat()
				properties = append(properties, ast.Property{
					Key: key,
				})
				continue
			} else if p.at().Type == tokentype.RSquirly {
				properties = append(properties, ast.Property{
					Key: key,
				})
				continue
			}

			p.expect(tokentype.Colon, fmt.Sprintf("Error on line %d: Expected : after object key", internal.Line()))
			value := p.parseExpression()
			properties = append(properties, ast.Property{
				Key:   key,
				Value: value,
			})
		}

		if p.at().Type != tokentype.RSquirly {
			p.expect(tokentype.Comma, fmt.Sprintf("Error on line %d: Expected , after object property", internal.Line()))
		}
	}

	p.expect(tokentype.RSquirly, fmt.Sprintf("Error on line %d: Object literal must end with a }", internal.Line()))
	return &ast.ObjectLiteral{
		Properties: properties,
	}
}

func (p *Parser) parseComparisonExpression() ast.Expression {
	left := p.parseObjectExpression()

	for p.at().Value == ">" || p.at().Value == "<" || (p.at().Value == "=" && p.peek().Value == "=") || p.at().Value == "!=" {
		operator := p.eat().Value
		if p.at().Value == "=" {
			operator += p.eat().Value
		}

		right := p.parseObjectExpression()

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseBitwise() ast.Expression {
	left := p.parseBitwiseShiftBit()

	for p.at().Value == "&" || p.at().Value == "|" || p.at().Value == "^" {
		operator := p.eat().Value

		right := p.parseBitwiseShiftBit()

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseBitwiseShiftBit() ast.Expression {
	left := p.parseAdditiveExpression()

	for p.at().Value == "<<" || p.at().Value == ">>" || p.at().Value == ">>>" {
		operator := p.eat().Value

		right := p.parseAdditiveExpression()

		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseAdditiveExpression() ast.Expression {
	left := p.parseMultiplicativeExpression()

	for p.at().Value == "+" || p.at().Value == "-" {
		operator := p.eat().Value

		right := p.parseMultiplicativeExpression()
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseMultiplicativeExpression() ast.Expression {
	left := p.parseCallMemberExpression()

	for p.at().Value == "*" || p.at().Value == "/" || p.at().Value == "%" || p.at().Value == "**" || p.at().Value == "//" {
		operator := p.eat().Value

		right := p.parseCallMemberExpression()
		left = &ast.BinaryExpression{
			Left:     left,
			Operator: operator,
			Right:    right,
		}
	}

	return left
}

func (p *Parser) parseCallMemberExpression() ast.Expression {
	member := p.parseMemberExpression()

	if p.at().Type == tokentype.OpenParen {
		return p.parseCallExpression(member)
	}

	return member
}

func (p *Parser) parseCallExpression(caller ast.Expression) ast.Expression {
	var callExpression ast.Expression = &ast.CallExpression{
		Caller: caller,
		Args:   p.parseArgs(),
	}

	if p.at().Type == tokentype.OpenParen {
		callExpression = p.parseCallExpression(callExpression)
	}

	return callExpression
}

func (p *Parser) parseArgs() []ast.Expression {
	p.expect(tokentype.OpenParen, fmt.Sprintf("Error on line %d: Expected '(' after function name", internal.Line()))

	var args []ast.Expression
	if p.at().Type == tokentype.CloseParen {
		args = []ast.Expression{}
	} else {
		args = p.parseArgumentsList()
	}

	p.expect(tokentype.CloseParen, fmt.Sprintf("Error on line %d: Expected ')' after function arguments", internal.Line()))

	return args
}

func (p *Parser) parseArgumentsList() []ast.Expression {
	args := []ast.Expression{p.parseAssignmentExpression()}

	for p.at().Type == tokentype.Comma {
		p.eat()
		args = append(args, p.parseAssignmentExpression())
	}

	return args
}

func (p *Parser) parseMemberExpression() ast.Expression {
	object := p.parsePrimaryExpression()

	for p.at().Type == tokentype.Dot || p.at().Type == tokentype.OpenBracket {
		operator := p.eat()
		var property ast.Expression
		var computed bool

		if operator.Type == tokentype.Dot {
			computed = false
			property = p.parsePrimaryExpression()

			if property.Kind() != ast.IdentifierType {
				fmt.Fprintln(os.Stderr, fmt.Sprintf("Error on line %d: Expected identifier after '.'", internal.Line()))
				os.Exit(0)
				return nil
			}
		} else {
			computed = true
			property = p.parseExpression()
			p.expect(tokentype.CloseBracket, fmt.Sprintf("Error on line %d: Expected ']' after computed property", internal.Line()))
		}

		object = &ast.MemberExpression{
			Object:   object,
			Property: property,
			Computed: computed,
		}
	}

	return object
}

func (p *Parser) parsePrimaryExpression() ast.Expression {
	token := p.at().Type

	switch token {
	case tokentype.Identifier:
		return &ast.Identifier{
			Symbol: p.eat().Value,
		}
	case tokentype.Integer:
		value, err := strconv.ParseInt(p.eat().Value, 10, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(0)
			return nil
		}
		return &ast.NumericIntegerLiteral{Value: value}
	case tokentype.Float:
		value, err := strconv.ParseFloat(p.eat().Value, 64)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(0)
			return nil
		}
		return &ast.NumericFloatLiteral{Value: value}
	case tokentype.String:
		return &ast.StringLiteral{
			Value: p.eat().Value,
		}
	case tokentype.Whitespace:
		p.eat()
		return p.parsePrimaryExpression()
	case tokentype.OpenParen:
		p.eat()
		value := p.parseExpression()
		p.expect(tokentype.CloseParen, fmt.Sprintf("Error on line %d: Expected closing parenthesis", internal.Line()))
		return value
	case tokentype.UnaryOperator:
		operator := p.eat().Value
		value := p.parsePrimaryExpression()
		return &ast.UnaryExpression{
			Operator: operator,
			Value:    value,
		}
	case tokentype.Function:
		return p.parseFunctionDeclaration()
	default:
		fmt.Fprintf(os.Stderr, "Error on line %d: Unexpected token found: %s", internal.Line(), p.at().Value)
		os.Exit(0)
		return nil
}
}

func (p *Parser) at() lexer.Token {
	return p.tokens[0]
}

func (p *Parser) eat() lexer.Token {
	prev := p.tokens[0]
	p.tokens = p.tokens[1:]
	return prev
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[1]
}

func (p *Parser) expect(token tokentype.TokenType, message string) lexer.Token {
	if p.at().Type != token {
		fmt.Fprintln(os.Stderr, message)
		os.Exit(0)
	}
	return p.eat()
}

func (p *Parser) notEndOfFile() bool {
	return p.tokens[0].Type != tokentype.EndOfFile
}
