package runtimelang

import (
    "fmt"

    "github.com/Jamlee977/CustomLanguage/ast"
)

func EvaluateProgram(program ast.Program, env Environment) RuntimeValue {
    var lastEvaluated RuntimeValue = &InitialValue{}
    for _, statement := range program.Body {
        lastEvaluated, _ = Evaluate(statement, env)
    }
    return lastEvaluated
}

func EvaluateVariableDeclaration(declaration ast.VariableDeclaration, env Environment) RuntimeValue {
    value, _ := Evaluate(declaration.Value, env)
    return env.DeclareVariable(declaration.Identifier, value, declaration.Constant)
}

func EvaluateIdentifier(identifier *ast.Identifier, env Environment) RuntimeValue {
    if identifier == nil {
        return MakeNullValue()
    }
    value := env.LookupVariable(identifier.Symbol)
    if value == nil {
        err := fmt.Errorf("Undefined variable %s", identifier.Symbol)
        panic(err)
    }

    return value
}
