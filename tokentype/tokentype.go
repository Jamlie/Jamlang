package tokentype

type TokenType int

const (
    Number TokenType = iota
    String
    Identifier
    Equals
    OpenParen
    CloseParen
    BinaryOperator
    Whitespace
    SemiColon

    Let
    Constant
    If
    Then
    End

    EndOfFile
)

