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
    UnaryOperator
    LogicalOperator
    ComparisonOperator
    Whitespace
    SemiColon

    Comma
    Colon
    Dot
    LSquirly
    RSquirly
    OpenBracket
    CloseBracket

    Function
    Return
    Let
    Constant
    If
    Else
    While
    Loop
    Break
    Continue
    Not
    And
    Or
    Xor

    EndOfFile
)

