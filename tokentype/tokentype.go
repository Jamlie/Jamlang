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

    Comma
    Colon
    LSquirly
    RSquirly

    Let
    Constant
    If
    Then
    Else
    End

    EndOfFile
)

