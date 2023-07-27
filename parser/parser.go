package parser

import (
    "strconv"
    "fmt"
    "os"

    "github.com/Jamlee977/CustomLanguage/ast"
    "github.com/Jamlee977/CustomLanguage/lexer"
    "github.com/Jamlee977/CustomLanguage/tokentype"
)

type Parser struct {
    tokens []lexer.Token
    isFunction bool
    isLoop bool
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
    case tokentype.Let:
        return p.parseVariableDeclaration()
    case tokentype.Constant:
        return p.parseVariableDeclaration()
    case tokentype.Function:
        return p.parseFunctionDeclaration()
    case tokentype.Return:
        if !p.isFunction {
            fmt.Println("Error: Return statement outside of function")
            os.Exit(0)
        }
        return p.parseReturnStatement()
    case tokentype.Break:
        if !p.isLoop {
            fmt.Println("Error: Break statement outside of loop")
            os.Exit(0)
        }
        return p.parseBreakStatement()
    case tokentype.If:
        return p.parseIfStatement()
    case tokentype.Else:
        return p.parseIfStatement()
    case tokentype.While:
        return p.parseWhileStatement()
    case tokentype.Loop:
        return p.parseLoopStatement()
    case tokentype.Import:
        return p.parseImportStatement()
    case tokentype.SemiColon:
        p.eat()
        return &ast.NullLiteral{}
    default:
        return p.parseExpression()
    }
}

func (p *Parser) parseImportStatement() ast.Statement {
    p.eat()
    path := p.expect(tokentype.String, "Error: Expected string after import statement").Value
    return &ast.ImportStatement{Path: path}
}

func (p *Parser) parseLoopStatement() ast.Statement {
    p.eat()
    p.expect(tokentype.LSquirly, "Error: Expected { after loop statement")

    p.isLoop = true
    defer func() { p.isLoop = false }()
    var body []ast.Statement
    for p.at().Type != tokentype.RSquirly {
        body = append(body, p.parseStatement())
    }

    p.expect(tokentype.RSquirly, "Error: Expected } after loop statement")

    return &ast.LoopStatement{Body: body}
}

func (p *Parser) parseWhileStatement() ast.Statement {
    p.eat()
    condition := p.parseExpression()
    p.expect(tokentype.LSquirly, "Error: Expected { after while statement")

    p.isLoop = true
    defer func() { p.isLoop = false }()
    var body []ast.Statement
    for p.at().Type != tokentype.RSquirly {
        body = append(body, p.parseStatement())
    }

    p.expect(tokentype.RSquirly, "Error: Expected } after while statement")

    return &ast.WhileStatement{
        Condition: condition,
        Body: body,
    }
}

func (p *Parser) parseIfStatement() ast.Statement {
    p.eat()
    condition := p.parseExpression()
    p.expect(tokentype.LSquirly, "Error: Expected { after if statement")

    var body []ast.Statement
    for p.at().Type != tokentype.RSquirly {
        body = append(body, p.parseStatement())
    }

    p.expect(tokentype.RSquirly, "Error: Expected } after if statement")

    if p.at().Type == tokentype.Else {
        p.eat()
        p.expect(tokentype.LSquirly, "Error: Expected { after else statement")

        var elseBody []ast.Statement
        for p.at().Type != tokentype.RSquirly {
            elseBody = append(elseBody, p.parseStatement())
        }

        p.expect(tokentype.RSquirly, "Error: Expected } after else statement")

        return &ast.ConditionalStatement{
            Condition: condition,
            Body: body,
            Alternate: elseBody,
        }
    }

    return &ast.ConditionalStatement{
        Condition: condition,
        Body:      body,
    }
}

func (p *Parser) parseFunctionDeclaration() ast.Statement {
    p.eat()
    name := p.expect(tokentype.Identifier, "Error: Expected function name after fn keyword").Value

    args := p.parseArgs()
    var params []string

    for _, arg := range args {
        if arg.Kind() != ast.IdentifierType {
            fmt.Printf("Error: Expected function parameter to be of type string, got %s\n", arg.Kind())
            os.Exit(0)
            return nil
        }

        params = append(params, arg.(*ast.Identifier).Symbol)
    }

    p.expect(tokentype.LSquirly, "Error: Expected '{' after function declaration")

    p.isFunction = true
    defer func() { p.isFunction = false }()
    body := []ast.Statement{}
    for (p.at().Type != tokentype.EndOfFile && p.at().Type != tokentype.RSquirly) {
        body = append(body, p.parseStatement())
    }

    p.expect(tokentype.RSquirly, "Error: Expected '}' after function declaration")
    
    return &ast.FunctionDeclaration{
        Name: name,
        Parameters: params,
        Body: body,
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


func (p *Parser) parseVariableDeclaration() ast.Statement {
    isConstant := p.eat().Type == tokentype.Constant
    identifier := p.expect(tokentype.Identifier, "Error: Expected identifier name after let/const keyword").Value

    if p.at().Type == tokentype.SemiColon {
        p.eat()
        if isConstant {
            fmt.Println("Error: Constant declaration without assignment is not allowed")
            os.Exit(0)
            return nil
        }

        return &ast.VariableDeclaration{
            Identifier: identifier,
            Constant: isConstant,
            Value: &ast.NullLiteral{},
        }
    }

    p.expect(tokentype.Equals, "Expected = after identifier name")
    declaration := &ast.VariableDeclaration{
        Identifier: identifier,
        Constant: isConstant,
        Value: p.parseExpression(),
    }

    if p.at().Type == tokentype.SemiColon {
        p.eat()
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
            Left: left,
            Right: right,
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
            Left: left,
            Right: right,
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
            Left: left,
            Right: right,
        }
    }
    return left
}

func (p *Parser) parseNotExpression() ast.Expression {
    if p.at().Value == "not" {
        p.eat()
        return &ast.LogicalExpression{
            Operator: "not",
            Right: p.parseNotExpression(),
        }
    }
    return p.parseComparisonExpression()
}

func (p *Parser) parseAssignmentExpression() ast.Expression {
    left := p.parseOrExpression()

    if p.at().Type == tokentype.Equals {
        p.eat()
        value := p.parseAssignmentExpression()
        return &ast.AssignmentExpression{
            Assigne: left,
            Value: value,
        }
    }

    return left
}

func (p *Parser) parseObjectExpression() ast.Expression {
    if p.at().Type != tokentype.LSquirly {
        return p.parseAdditiveExpression()
    }

    p.eat()
    properties := []ast.Property{}

    for p.notEndOfFile() && p.at().Type != tokentype.RSquirly {
        key := p.expect(tokentype.Identifier, "Expected identifier as object key").Value

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

        p.expect(tokentype.Colon, "Expected : after object key")
        value := p.parseExpression()
        properties = append(properties, ast.Property{
            Key: key,
            Value: value,
        })

        if p.at().Type != tokentype.RSquirly {
            p.expect(tokentype.Comma, "Expected , after object property")
        }
    }

    p.expect(tokentype.RSquirly, "Object literal must end with a }")
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
            Left: left,
            Operator: operator,
            Right: right,
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
            Left: left,
            Operator: operator,
            Right: right,
        }
    }

    return left
}

func (p *Parser) parseMultiplicativeExpression() ast.Expression {
    left := p.parseCallMemberExpression()

    for p.at().Value == "*" || p.at().Value == "/" || p.at().Value == "%" {
        operator := p.eat().Value

        right := p.parseCallMemberExpression()
        left = &ast.BinaryExpression{
            Left: left,
            Operator: operator,
            Right: right,
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
        Args: p.parseArgs(),
    }

    if p.at().Type == tokentype.OpenParen {
        callExpression = p.parseCallExpression(callExpression)
    }

    return callExpression
}

func (p *Parser) parseArgs() []ast.Expression {
    p.expect(tokentype.OpenParen, "Expected '(' after function name")

    var args []ast.Expression
    if p.at().Type == tokentype.CloseParen {
        args = []ast.Expression{}
    } else {
        args = p.parseArgumentsList()
    }

    p.expect(tokentype.CloseParen, "Expected ')' after function arguments")

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
                fmt.Println("Expected identifier after '.'")
                os.Exit(0)
                return nil
            }
        } else {
            computed = true
            property = p.parseExpression()
            p.expect(tokentype.CloseBracket, "Expected ']' after computed property")
        }

        object = &ast.MemberExpression{
            Object: object,
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
    case tokentype.Number:
        value, err := strconv.ParseFloat(p.eat().Value, 64)
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(0)
            return nil
        }
        return &ast.NumericLiteral{
            Value: value,
        }
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
        p.expect(tokentype.CloseParen, "Expected closing parenthesis")
        return value
    case tokentype.UnaryOperator:
        operator := p.eat().Value
        value := p.parsePrimaryExpression()
        return &ast.UnaryExpression{
            Operator: operator,
            Value: value,
        }
    default:
        fmt.Println("Unexpected token found: ", p.at())
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
        fmt.Println(message)
        os.Exit(0)
    }
    return p.eat()
}

func (p *Parser) notEndOfFile() bool {
    return p.tokens[0].Type != tokentype.EndOfFile
}
