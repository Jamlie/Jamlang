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

// func EvaluateCondition(condition RuntimeValue) bool {
//     if condition == nil {
//         return false
//     }
//     switch condition.(type) {
//     case BoolValue:
//         return condition.(BoolValue).Value
//     case NumberValue:
//         return condition.(NumberValue).Value != 0
//     case StringValue:
//         return condition.(StringValue).Value != ""
//     case NullValue:
//         return false
//     default:
//         return true
//     }
// }

// func EvaluateConditionalStatement(statement ast.ConditionalStatement, env Environment) RuntimeValue {
//     condition, _ := Evaluate(statement.Condition, env)
//     if EvaluateCondition(condition) {
//         for _, statement := range statement.Body {
//             Evaluate(statement, env)
//         }
//     } else {
//         for _, statement := range statement.Alternate {
//             Evaluate(statement, env)
//         }
//     }
//     
//     return MakeNullValue()
// }

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
    return env.DeclareVariable(declaration.Identifier, value, declaration.Constant)
}

func EvaluateIdentifier(identifier *ast.Identifier, env Environment) RuntimeValue {
    if identifier == nil {
        return MakeNullValue()
    }
    value := env.LookupVariable(identifier.Symbol)
    if value == nil {
        fmt.Printf("Undefined variable %s\n", identifier.Symbol)
        os.Exit(1)
        return nil
    }

    return value
}
