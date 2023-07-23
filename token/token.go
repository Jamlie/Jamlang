package token

import (
    "github.com/Jamlee977/CustomLanguage/tokentype"
    "strings"
    "fmt"
    "os"
)

type Token struct {
    Type tokentype.TokenType
    Value string
}

var Keywords map[string]tokentype.TokenType = map[string]tokentype.TokenType{
    "let": tokentype.Let,
}


func createToken(value string, tokenType tokentype.TokenType) Token {
    return Token{
        Type: tokenType,
        Value: value,
    }
}

func isAlpha(src string) bool {
    return strings.ToUpper(src) != strings.ToLower(src)
}

func isInt(src string) bool {
    return src[0] >= 48 && src[0] <= 57
}

func isWhitespace(src string) bool {
    return src == " " || src == "\t" || src == "\n"
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
        } else if src[0] == "+" || src[0] == "-" || src[0] == "*" || src[0] == "/" || src[0] == "%" {
            tokens = append(tokens, createToken(src[0], tokentype.BinaryOperator))
            src = src[1:]
        } else if src[0] == "=" {
            tokens = append(tokens, createToken(src[0], tokentype.Equals))
            src = src[1:]
        } else {
            if isInt(src[0]) {
                num := ""
                for len(src) > 0 && isInt(src[0]) {
                    num += src[0]
                    src = src[1:]
                }

                tokens = append(tokens, createToken(num, tokentype.Number))
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
                tokens = append(tokens, createToken(src[0], tokentype.Whitespace))
                src = src[1:]
            } else {
                fmt.Println("Unknown token: " + src[0])
                os.Exit(1)
            }
        }
    }

    tokens = append(tokens, createToken("EndOfFile", tokentype.EndOfFile))
    return tokens
}
