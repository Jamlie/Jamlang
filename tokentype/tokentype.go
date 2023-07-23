package tokentype

type TokenType int

const (
    Number TokenType = iota
    Null
    Identifier
    Equals
    OpenParen
    CloseParen
    BinaryOperator
    Whitespace

    Let

    EndOfFile
)

