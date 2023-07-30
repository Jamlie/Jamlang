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
    ColonColon
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
    ForEach
    For
    In
    Break
    Not
    And
    Or
    Xor
    Import
    Class

    EndOfFile
)

