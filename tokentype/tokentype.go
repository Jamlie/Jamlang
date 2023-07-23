package tokentype

type TokenType int

const (
    Number TokenType = iota
    Identifier
    Equals
    OpenParen
    CloseParen
    BinaryOperator
    Whitespace

    Let

    EndOfFile
)

