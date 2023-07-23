package ast

import "strconv"

type NodeType string

const (
    ProgramType NodeType = "Program"
    NumericLiteralType NodeType = "NumericLiteral"
    IdentifierType NodeType = "Identifier"
    BinaryExpressionType NodeType = "BinaryExpression"
)


type Statement interface {
    Kind() NodeType
    ToString() string
}

type Program struct {
    Body []Statement
}

func (p *Program) Kind() NodeType {
    return ProgramType
}

func (p *Program) ToString() string {
    s := ""
    for _, statement := range p.Body {
        s += statement.ToString()
    }

    return s
}

type Expression interface {
    Statement
}

type BinaryExpression struct {
    Left Expression
    Right Expression
    Operator string
}

func (b *BinaryExpression) Kind() NodeType {
    return BinaryExpressionType
}

func (b *BinaryExpression) ToString() string {
    return "(" + b.Left.ToString() + " " + b.Operator + " " + b.Right.ToString() + ")"
}

type Identifier struct {
    Symbol string
}

func (i *Identifier) Kind() NodeType {
    return IdentifierType
}

func (i *Identifier) ToString() string {
    return i.Symbol
}

type NumericLiteral struct {
    Value float64
}

func (n *NumericLiteral) Kind() NodeType {
    return NumericLiteralType
}

func (n *NumericLiteral) ToString() string {
    return strconv.FormatFloat(n.Value, 'f', -1, 64)
}
