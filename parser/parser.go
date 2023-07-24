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
    case tokentype.SemiColon:
        p.eat()
        return &ast.NullLiteral{}
    default:
        return p.parseExpression()
    }
}

func (p *Parser) parseVariableDeclaration() ast.Statement {
    isConstant := p.eat().Type == tokentype.Constant
    identifier := p.expect(tokentype.Identifier, "Expected identifier name after let/const keyword").Value

    if p.at().Type == tokentype.SemiColon {
        p.eat()
        if isConstant {
            fmt.Println("Constant declaration without assignment is not allowed")
            os.Exit(1)
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

func (p *Parser) parseAssignmentExpression() ast.Expression {
    left := p.parseAdditiveExpression()

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
    left := p.parsePrimaryExpression()

    for p.at().Value == "*" || p.at().Value == "/" || p.at().Value == "%" {
        operator := p.eat().Value

        right := p.parsePrimaryExpression()
        left = &ast.BinaryExpression{
            Left: left,
            Operator: operator,
            Right: right,
        }
    }

    return left
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
            panic(err)
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
    default:
        fmt.Println("Unexpected token found: ", p.at())
        panic("Unexpected token found")
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

func (p *Parser) expect(token tokentype.TokenType, message string) lexer.Token {
    if p.at().Type != token {
        fmt.Println(message)
        os.Exit(1)
    }
    return p.eat()
}

func (p *Parser) notEndOfFile() bool {
    return p.tokens[0].Type != tokentype.EndOfFile
}
