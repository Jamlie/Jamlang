package runtimelang

import (
    "fmt"
    "os"

    "github.com/Jamlee977/CustomLanguage/ast"
)

func EvaluateProgram(program ast.Program, env Environment) RuntimeValue {
    var lastEvaluated RuntimeValue = &InitialValue{}
    for _, statement := range program.Body {
        lastEvaluated, _ = Evaluate(statement, env)
    }
    return lastEvaluated
}

func EvaluateBreakStatement(statement ast.BreakStatement, env Environment) RuntimeValue {
    return &BreakType{}
}

func EvaluateReturnStatement(statement ast.ReturnStatement, env Environment) RuntimeValue {
    returnValue, _ := Evaluate(statement.Value, env)
    return returnValue
}

func EvaluateFunctionDeclaration(declaration ast.FunctionDeclaration, env Environment) RuntimeValue {
    fn := FunctionValue{
        Name: declaration.Name,
        Parameters: declaration.Parameters,
        DeclarationEnvironment: env,
        Body: declaration.Body,
    }
    return env.DeclareVariable(declaration.Name, &fn, true)
}

func EvaluateVariableDeclaration(declaration ast.VariableDeclaration, env Environment) RuntimeValue {
    value, _ := Evaluate(declaration.Value, env)
    // if value.Type() == "object" {
    //     value = value.(ObjectValue).Clone()
    // }
    return env.DeclareVariable(declaration.Identifier, value, declaration.Constant)
}

func EvaluateIdentifier(identifier *ast.Identifier, env Environment) RuntimeValue {
    if identifier == nil {
        return MakeNullValue()
    }
    value := env.LookupVariable(identifier.Symbol)
    if value == nil {
        fmt.Printf("Undefined variable %s\n", identifier.Symbol)
        os.Exit(0)
        return nil
    }

    return value
}
