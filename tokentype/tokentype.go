package tokentype

type TokenType int

const (
    Number TokenType = iota
    Integer
    Float
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
    OpenComment
    CloseComment

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
    Var
    Let
    Constant
    If
    ElseIf
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

