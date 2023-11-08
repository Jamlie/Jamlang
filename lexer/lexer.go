package lexer

import (
    "github.com/Jamlie/Jamlang/tokentype"
    "strings"
    "fmt"
    "os"
    "strconv"
)

type Token struct {
    Type tokentype.TokenType
    Value string
}

var Keywords map[string]tokentype.TokenType = map[string]tokentype.TokenType{
    "var": tokentype.Var,
    "let": tokentype.Let,
    "const": tokentype.Constant,
    "fn": tokentype.Function,
    "return": tokentype.Return,
    "if": tokentype.If,
    "elseif": tokentype.ElseIf,
    "else": tokentype.Else,
    "while": tokentype.While,
    "loop": tokentype.Loop,
    "foreach": tokentype.ForEach,
    "for": tokentype.For,
    "in": tokentype.In,
    "break": tokentype.Break,
    "not": tokentype.LogicalOperator,
    "and": tokentype.LogicalOperator,
    "or": tokentype.LogicalOperator,
    "import": tokentype.Import,
    "class": tokentype.Class,
    "type": tokentype.Type,
}


func createToken(value string, tokenType tokentype.TokenType) Token {
    return Token{
        Type: tokenType,
        Value: value,
    }
}

func isAlpha(src string) bool {
    return strings.ContainsAny(src, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_")
}

func isInt(src string) bool {
    _, err := strconv.Atoi(src)
    return err == nil
}

func isFloat(src string) bool {
    _, err := strconv.ParseFloat(src, 64)
    return err == nil
}

func isWhitespace(src string) bool {
    return src == " " || src == "\t" || src == "\n" || src == "\r"
}

func Tokenize(sourceCode string) []Token {
    tokens := []Token{}
    src := strings.Split(sourceCode, "")

    for len(src) > 0 {
        if src[0] == "(" {
            tokens = append(tokens, createToken(src[0], tokentype.OpenParen))
            src = src[1:] 
        } else if src[0] == ")" {
            tokens = append(tokens, createToken(src[0], tokentype.CloseParen))
            src = src[1:]
        }  else if src[0] == "{" {
            tokens = append(tokens, createToken(src[0], tokentype.LSquirly))
            src = src[1:]
        }  else if src[0] == "}" {
            tokens = append(tokens, createToken(src[0], tokentype.RSquirly))
            src = src[1:]
        } else if src[0] == "[" {
            tokens = append(tokens, createToken(src[0], tokentype.OpenBracket))
            src = src[1:]
        } else if src[0] == "]" {
            tokens = append(tokens, createToken(src[0], tokentype.CloseBracket))
            src = src[1:]
        } else if src[0] == "+" || src[0] == "-" || src[0] == "*" || src[0] == "/" || src[0] == "%" || src[0] == "&" || src[0] == "|" || src[0] == "^" {
            if (src[0] == "-" && isInt(src[1])) || (src[0] == "-" && isFloat(src[1])) || (src[0] == "-" && isAlpha(src[1])) {
                tokens = append(tokens, createToken(src[0], tokentype.UnaryOperator))
                src = src[1:]
                continue
            }
            if src[0] == "+" && src[1] == "+" {
                tokens = append(tokens, createToken("++", tokentype.UnaryOperator))
                src = src[2:]
                continue
            }
            if src[0] == "-" && src[1] == "-" {
                tokens = append(tokens, createToken("--", tokentype.UnaryOperator))
                src = src[2:]
                continue
            }
            if src[0] == "*" && src[1] == "*" {
                tokens = append(tokens, createToken("**", tokentype.BinaryOperator))
                src = src[2:]
                continue
            }
            if src[0] == "/" && src[1] == "*" {
                tokens = append(tokens, createToken("/*", tokentype.OpenComment))
                src = src[2:]
                continue
            }
            if src[0] == "*" && src[1] == "/" {
                tokens = append(tokens, createToken("*/", tokentype.CloseComment))
                src = src[2:]
                continue
            }
            if src[0] == "/" && src[1] == "/" {
                tokens = append(tokens, createToken("//", tokentype.BinaryOperator))
                src = src[2:]
                continue
            }
            tokens = append(tokens, createToken(src[0], tokentype.BinaryOperator))
            src = src[1:]
        } else if src[0] == "=" {
            tokens = append(tokens, createToken(src[0], tokentype.Equals))
            src = src[1:]
        } else if src[0] == ">" {
            if src[1] == ">" {
                tokens = append(tokens, createToken(">>", tokentype.BinaryOperator))
                src = src[2:]
                continue
            }
            tokens = append(tokens, createToken(src[0], tokentype.ComparisonOperator))
            src = src[1:]
        } else if src[0] == "<" {
            if src[1] == "<" {
                tokens = append(tokens, createToken("<<", tokentype.BinaryOperator))
                src = src[2:]
                continue
            }
            tokens = append(tokens, createToken(src[0], tokentype.ComparisonOperator))
            src = src[1:]
        } else if src[0] == ">" && src[1] == "=" {
            tokens = append(tokens, createToken(">=", tokentype.ComparisonOperator))
            src = src[2:]
        } else if src[0] == "<" && src[1] == "=" {
            tokens = append(tokens, createToken("<=", tokentype.ComparisonOperator))
            src = src[2:]
        } else if src[0] == "=" && src[1] == "=" {
            tokens = append(tokens, createToken("==", tokentype.ComparisonOperator))
            src = src[2:]
        } else if src[0] == "!" && src[1] == "=" {
            tokens = append(tokens, createToken("!=", tokentype.ComparisonOperator))
            src = src[2:]
        } else if src[0] == ";" {
            tokens = append(tokens, createToken(src[0], tokentype.SemiColon))
            src = src[1:]
        } else if src[0] == "," {
            tokens = append(tokens, createToken(src[0], tokentype.Comma))
            src = src[1:]
        } else if src[0] == "." {
            tokens = append(tokens, createToken(src[0], tokentype.Dot))
            src = src[1:]
        } else if src[0] == ":" {
            if src[1] == ":" {
                tokens = append(tokens, createToken("::", tokentype.ColonColon))
                src = src[2:]
                continue
            }
            tokens = append(tokens, createToken(src[0], tokentype.Colon))
            src = src[1:]
        } else if src[0] == "\"" || src[0] == "`" || src[0] == "'" {
            quotationOrBacktick := src[0]
            src = src[1:]
            str := ""
            if quotationOrBacktick == "'" {
                for len(src) > 0 && src[0] != "'" {
                    str += src[0]
                    src = src[1:]
                }
            } else if quotationOrBacktick == "`" {
                for len(src) > 0 && src[0] != "`" {
                    str += src[0]
                    src = src[1:]
                }
            } else {
                for len(src) > 0 && src[0] != "\"" {
                    str += src[0]
                    src = src[1:]
                }
            }

            if len(src) == 0 {
                fmt.Fprintln(os.Stderr, "Error: Unterminated string")
                os.Exit(0)
            }

            tokens = append(tokens, createToken(str, tokentype.String))
            src = src[1:]
        } else {
            if isInt(src[0]) || (src[0] == "-" && isInt(src[1])) {
                num := ""
                isFloatNum := false

                for len(src) > 0 && (isInt(src[0]) || (!isFloatNum && src[0] == "." && len(src) > 1 && isInt(src[1]))) {
                    if src[0] == "." {
                        isFloatNum = true
                    }
                    num += src[0]
                    src = src[1:]
                }

                for len(src) > 0 && isInt(src[0]) {
                    num += src[0]
                    src = src[1:]
                }

                if isFloatNum {
                    tokens = append(tokens, createToken(num, tokentype.Float))
                } else {
                    tokens = append(tokens, createToken(num, tokentype.Integer))
                }
            } else if isAlpha(src[0]) {
                identifier := ""
                for len(src) > 0 && isAlpha(src[0]) {
                    identifier += src[0]
                    src = src[1:]
                }

                reserved, ok := Keywords[identifier]
                if ok {
                    tokens = append(tokens, createToken(identifier, reserved))
                } else {
                    tokens = append(tokens, createToken(identifier, tokentype.Identifier))
                }
            } else if isWhitespace(src[0]) {
                src = src[1:]
            } else {
                fmt.Fprintln(os.Stderr, "Error: Invalid character '" + string(src[0]) + "'")
                os.Exit(0)
            }
        }
    }

    tokens = append(tokens, createToken("EndOfFile", tokentype.EndOfFile))
    return tokens
}
