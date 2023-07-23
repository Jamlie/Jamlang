package runtimelang

import (
    "fmt"
    "math"
    "os"
    "strconv"

    "github.com/Jamlee977/CustomLanguage/ast"
)

func Evaluate(astNode ast.Statement) (RuntimeValue, error) {
    switch astNode.Kind() {
    case ast.NumericLiteralType:
        return NumberValue{astNode.(*ast.NumericLiteral).Value}, nil
    case ast.StringLiteralType:
        return StringValue{astNode.(*ast.StringLiteral).Value}, nil
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
    } else if lhs.Type() == "string" && rhs.Type() == "string" {
        return EvaluateStringBinaryExpression(lhs.(StringValue), rhs.(StringValue), binaryExpression.Operator)
    } else if lhs.Type() == "string" && rhs.Type() == "number" {
        return EvaluateStringNumericBinaryExpression(lhs.(StringValue), rhs.(NumberValue), binaryExpression.Operator)
    } else if lhs.Type() == "number" && rhs.Type() == "string" {
        return EvaluateNumericStringBinaryExpression(lhs.(NumberValue), rhs.(StringValue), binaryExpression.Operator)
    }

    return NullValue{"null"}
}

func EvaluateStringNumericBinaryExpression(lhs StringValue, rhs NumberValue, op string) RuntimeValue {
    if op == "+" {
        rhsAsString := strconv.FormatFloat(rhs.Value, 'f', -1, 64)

        return StringValue{lhs.Value + rhsAsString}
    }

    err := fmt.Errorf("Unknown operator %s for string", op)
    panic(err)
}

func EvaluateNumericStringBinaryExpression(lhs NumberValue, rhs StringValue, op string) RuntimeValue {
    if op == "+" {
        lhsAsString := strconv.FormatFloat(lhs.Value, 'f', -1, 64)

        return StringValue{lhsAsString + rhs.Value}
    }

    err := fmt.Errorf("Unknown operator %s for string", op)
    panic(err)
}

func EvaluateStringBinaryExpression(lhs, rhs StringValue, op string) RuntimeValue {
    if op == "+" {
        return StringValue{lhs.Value + rhs.Value}
    }

    err := fmt.Errorf("Unknown operator %s for string", op)
    panic(err)
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
    var lastEvaluated RuntimeValue = &InitialValue{}
    for _, statement := range program.Body {
        lastEvaluated, _ = Evaluate(statement)
    }
    return lastEvaluated
}
