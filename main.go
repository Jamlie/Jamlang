package main

import (
    "github.com/Jamlie/Jamlang/jamlang"
    "github.com/Jamlie/Jamlang/runtimelang"
)

var (
    env = runtimelang.CreateGlobalEnvironment()
)

func main() {
    jamlang.CallMain(env)
}

/*
{
    class := ClassValue{
        Name: expr.Name,
        Methods: make(map[string]FunctionValue),
    }
    for _, method := range expr.Body {
        if method.Kind() == ast.FunctionDeclarationType {
            method := method.(*ast.FunctionDeclaration)
            if method.Name == "constructor" {
                class.Constructor = FunctionValue{
                    Name: method.Name,
                    Parameters: method.Parameters,
                    Body: method.Body,
                    DeclarationEnvironment: *env,
                }

                continue
            }
            class.Methods[method.Name] = FunctionValue{
                Name: method.Name,
                Parameters: method.Parameters,
                Body: method.Body,
                DeclarationEnvironment: *env,
            }

            env.DeclareVariable(method.Name, class.Methods[method.Name], true)
        } else if method.Kind() == ast.VariableDeclarationType {
            method := method.(*ast.VariableDeclaration)
            value, err := Evaluate(method.Value, *env)
            if err != nil {
                return MakeNullValue(), err
            }
            env.DeclareVariable(method.Identifier, value, true)
            continue
        } else {
            return MakeNullValue(), fmt.Errorf("invalid method declaration")
        }
    }
    env.DeclareVariable(expr.Name, class.Constructor, true)
    return MakeNullValue(), nil
}
*/
