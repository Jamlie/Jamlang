package ast

import "strconv"
import "bytes"

type NodeType string

const (
    ProgramType NodeType = "Program"
    VariableDeclarationType NodeType = "VariableDeclaration"
    AssignmentExpressionType NodeType = "AssignmentExpression"
    MemberExpressionType NodeType = "MemberExpression"
    CallExpressionType NodeType = "CallExpression"
    ConditionalStatementType NodeType = "ConditionalStatement"
    WhileStatementType NodeType = "WhileStatement"
    LoopStatementType NodeType = "LoopStatement"
    ForEachStatementType NodeType = "ForEachStatement"
    ForStatementType NodeType = "ForStatement"
    FunctionDeclarationType NodeType = "FunctionDeclaration"
    ReturnStatementType NodeType = "ReturnStatement"
    BreakStatementType NodeType = "BreakStatement"
    ImportStatementType NodeType = "ImportStatement"
    ClassDeclarationType NodeType = "ClassDeclaration"
    CommentType NodeType = "Comment"


    PropertyType NodeType = "Property"
    ObjectLiteralType NodeType = "ObjectLiteral"
    ArrayLiteralType NodeType = "ArrayLiteral"
    NumericLiteralType NodeType = "NumericLiteral"
    IdentifierType NodeType = "Identifier"
    BinaryExpressionType NodeType = "BinaryExpression"
    UnaryExpressionType NodeType = "UnaryExpression"
    LogicalExpressionType NodeType = "LogicalExpression"
    StringLiteralType NodeType = "StringLiteral"
    NullLiteralType NodeType = "NullLiteral"
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

type VariableDeclaration struct{
    Constant bool
    Identifier string
    Value Expression
}

func (v *VariableDeclaration) Kind() NodeType {
    return VariableDeclarationType
}

func (v *VariableDeclaration) ToString() string {
    s := ""
    if v.Constant {
        s += "const "
    } else {
        s += "let "
    }

    s += v.Identifier + " = " + v.Value.ToString() + ";\n"

    return s
}

type FunctionDeclaration struct {
    Parameters []string
    Name string
    Body []Statement
    IsAnonymous bool
}

func (f *FunctionDeclaration) Kind() NodeType {
    return FunctionDeclarationType
}

func (f *FunctionDeclaration) ToString() string {
    var s string
    if f.IsAnonymous {
        s = "function("
    } else {
        s = "function " + f.Name + "("
    }
    for i, param := range f.Parameters {
        if i > 0 {
            s += ", "
        }
        s += param
    }
    s += ") {\n"

    for _, statement := range f.Body {
        s += statement.ToString()
    }

    s += "}\n"

    return s
}

func (f *FunctionDeclaration) CloneBody() []Statement {
    body := make([]Statement, len(f.Body))
    copy(body, f.Body)
    return body
}

func (f *FunctionDeclaration) CloneParameters() []string {
    params := make([]string, len(f.Parameters))
    copy(params, f.Parameters)
    return params
}

type ReturnStatement struct {
    Value Expression
}

func (r *ReturnStatement) Kind() NodeType {
    return ReturnStatementType
}

func (r *ReturnStatement) ToString() string {
    return "return " + r.Value.ToString()
}

type BreakStatement struct {}

func (b *BreakStatement) Kind() NodeType {
    return BreakStatementType
}

func (b *BreakStatement) ToString() string {
    return "break"
}

type ImportStatement struct {
    Path string
}

func (i *ImportStatement) Kind() NodeType {
    return ImportStatementType
}

func (i *ImportStatement) ToString() string {
    return "import " + i.Path
}

type ClassDeclaration struct {
    Name string
    Body []Statement
}

func (c *ClassDeclaration) Kind() NodeType {
    return ClassDeclarationType
}

func (c *ClassDeclaration) ToString() string {
    s := "class " + c.Name + " {\n"

    for _, statement := range c.Body {
        s += statement.ToString()
    }

    s += "}\n"

    return s
}

type Comment struct {
    Text string
}

func (c *Comment) Kind() NodeType {
    return CommentType
}

func (c *Comment) ToString() string {
    return "/*" + c.Text + "*/"
}

type Expression interface {
    Statement
}

type ConditionalStatement struct {
    Condition Expression
    Body []Statement
    Alternate []Statement
}

func (c *ConditionalStatement) Kind() NodeType {
    return ConditionalStatementType
}

func (c *ConditionalStatement) ToString() string {
    s := "if (" + c.Condition.ToString() + ") {\n"
    for _, statement := range c.Body {
        s += statement.ToString()
    }
    s += "} else {\n"
    for _, statement := range c.Alternate {
        s += statement.ToString()
    }
    s += "}\n"

    return s
}

type WhileStatement struct {
    Condition Expression
    Body []Statement
}

func (w *WhileStatement) Kind() NodeType {
    return WhileStatementType
}

func (w *WhileStatement) ToString() string {
    s := "while (" + w.Condition.ToString() + ") {\n"
    for _, statement := range w.Body {
        s += statement.ToString()
    }
    s += "}\n"

    return s
}

type LoopStatement struct {
    Body []Statement
}

func (l *LoopStatement) Kind() NodeType {
    return LoopStatementType
}

func (l *LoopStatement) ToString() string {
    s := "loop {\n"
    for _, statement := range l.Body {
        s += statement.ToString()
    }
    s += "}\n"

    return s
}

type ForEachStatement struct {
    Variable string
    Collection Expression
    Body []Statement
}

func (f *ForEachStatement) Kind() NodeType {
    return ForEachStatementType
}

func (f *ForEachStatement) ToString() string {
    s := "foreach (" + f.Variable + " in " + f.Collection.ToString() + ") {\n"
    for _, statement := range f.Body {
        s += statement.ToString()
    }
    s += "}\n"

    return s
}

type ForStatement struct {
    Init Statement
    Condition Expression
    Update Expression
    Body []Statement
}

func (f *ForStatement) Kind() NodeType {
    return ForStatementType
}

func (f *ForStatement) ToString() string {
    s := "for (" + f.Init.ToString() + "; " + f.Condition.ToString() + "; " + f.Update.ToString() + ") {\n"
    for _, statement := range f.Body {
        s += statement.ToString()
    }
    s += "}\n"

    return s
}

type AssignmentExpression struct {
    Assigne Expression
    Value Expression
}

func (a *AssignmentExpression) Kind() NodeType {
    return AssignmentExpressionType
}

func (a *AssignmentExpression) ToString() string {
    return a.Assigne.ToString() + " = " + a.Value.ToString()
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

type UnaryExpression struct {
    Value Expression
    Operator string
}

func (u *UnaryExpression) Kind() NodeType {
    return UnaryExpressionType
}

func (u *UnaryExpression) ToString() string {
    return u.Operator + u.Value.ToString()
}

type LogicalExpression struct {
    Left Expression
    Right Expression
    Operator string
}

func (l *LogicalExpression) Kind() NodeType {
    return LogicalExpressionType
}

func (l *LogicalExpression) ToString() string {
    if l.Operator == "not" {
        return l.Operator + l.Right.ToString()
    } else {
        return "(" + l.Left.ToString() + " " + l.Operator + " " + l.Right.ToString() + ")"
    }
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

type StringLiteral struct {
    Value string
}

func (s *StringLiteral) Kind() NodeType {
    return StringLiteralType
}

func (s *StringLiteral) ToString() string {
    return "\"" + s.Value + "\""
}

type NullLiteral struct {}

func (n *NullLiteral) Kind() NodeType {
    return NullLiteralType
}

func (n *NullLiteral) ToString() string {
    return "null"
}

type Property struct {
    Key string
    Value Expression
}

func (p *Property) Kind() NodeType {
    return PropertyType
}

func (p *Property) ToString() string {
    return p.Key + ": " + p.Value.ToString()
}

type ObjectLiteral struct {
    Properties []Property
}

func (o *ObjectLiteral) Kind() NodeType {
    return ObjectLiteralType
}

func (o *ObjectLiteral) ToString() string {
    var buffer bytes.Buffer
    buffer.WriteString("{")
    for i, p := range o.Properties {
        buffer.WriteString(p.ToString())
        if i < len(o.Properties) - 1 {
            buffer.WriteString(", ")
        }
    }
    buffer.WriteString("}")
    return buffer.String()
}

type ArrayLiteral struct {
    Elements []Expression
}

func (a *ArrayLiteral) Kind() NodeType {
    return ArrayLiteralType
}

func (a *ArrayLiteral) ToString() string {
    var buffer bytes.Buffer
    buffer.WriteString("[")
    for i, e := range a.Elements {
        buffer.WriteString(e.ToString())
        if i < len(a.Elements) - 1 {
            buffer.WriteString(", ")
        }
    }
    buffer.WriteString("]")
    return buffer.String()
}

type CallExpression struct {
    Args []Expression
    Caller Expression
}

func (c CallExpression) Kind() NodeType {
    return CallExpressionType
}

func (c CallExpression) ToString() string {
    var buffer bytes.Buffer
    buffer.WriteString(c.Caller.ToString())
    buffer.WriteString("(")
    for i, arg := range c.Args {
        buffer.WriteString(arg.ToString())
        if i < len(c.Args) - 1 {
            buffer.WriteString(", ")
        }
    }
    buffer.WriteString(")")
    return buffer.String()
}

type MemberExpression struct {
    Object Expression
    Property Expression
    Computed bool
}

func (m *MemberExpression) Kind() NodeType {
    return MemberExpressionType
}

func (m *MemberExpression) ToString() string {
    if m.Computed {
        return m.Object.ToString() + "[" + m.Property.ToString() + "]"
    } else {
        return m.Object.ToString() + "." + m.Property.ToString()
    }
}
