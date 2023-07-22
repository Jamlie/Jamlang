package main

import (
    "fmt"
    "github.com/Jamlee977/CustomLanguage/token"
)

func main() {
    source := "let x = 5 + (foo * bar)"
    tokens := token.Tokenize(source)
    for _, token := range tokens {
        fmt.Println(token)
    }
}
