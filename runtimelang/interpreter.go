package runtimelang

import (
    "fmt"
    "math"
    "os"

    "github.com/Jamlee977/CustomLanguage/ast"
)

func Evaluate(astNode ast.Statement) (RuntimeValue, error) {
    switch astNode.Kind() {
    case ast.NumericLiteralType:
        return NumberValue{astNode.(*ast.NumericLiteral).Value}, nil
    case ast.NullLiteralType:
        return NullValue{"null"}, nil
    case ast.BinaryExpressionType:
        binaryExpression, ok := astNode.(*ast.BinaryExpression)
        if !ok {
            err := fmt.Errorf("Expected BinaryExpression, got %T", astNode)
            panic(err)
        }
        return EvaluateBinaryExpression(*binaryExpression), nil
    case ast.ProgramType:
        program, ok := astNode.(*ast.Program)
        if !ok {
            err := fmt.Errorf("Expected Program, got %T", astNode)
            panic(err)
        }
        return EvaluateProgram(*program), nil
    default:
        err := fmt.Errorf("Unknown AST node type %T", astNode)
        panic(err)
    }
}

func EvaluateBinaryExpression(binaryExpression ast.BinaryExpression) RuntimeValue {
    lhs, err := Evaluate(binaryExpression.Left)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    rhs, err := Evaluate(binaryExpression.Right)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if lhs.Type() == "number" && rhs.Type() == "number" {
        return EvaluateNumericBinaryExpression(lhs.(NumberValue), rhs.(NumberValue), binaryExpression.Operator)
    }

    return NullValue{"null"}
}

func EvaluateNumericBinaryExpression(lhs, rhs NumberValue, op string) RuntimeValue {
    var result float64 = 0
    if op == "+" {
        result = lhs.Value + rhs.Value
    } else if op == "-" {
        result = lhs.Value - rhs.Value
    } else if op == "*" {
        result = lhs.Value * rhs.Value
    } else if op == "/" {
        if rhs.Value == 0 {
            err := fmt.Errorf("Division by zero")
            panic(err)
        }
        result = lhs.Value / rhs.Value
    } else if op == "%" {
        result = math.Mod(lhs.Value, rhs.Value)
    } else {
        err := fmt.Errorf("Unknown operator: %s", op)
        panic(err)
    }

    return NumberValue{result}
}

func EvaluateProgram(program ast.Program) RuntimeValue {
    var lastEvaluated RuntimeValue = NullValue{"null"}
    for _, statement := range program.Body {
        lastEvaluated, _ = Evaluate(statement)
    }
    return lastEvaluated
}
