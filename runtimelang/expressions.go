package runtimelang

import (
    "fmt"
    "os"
    "strconv"
    "math"

    "github.com/Jamlee977/CustomLanguage/ast"
)

func EvaluateObjectExpression(obj ast.ObjectLiteral, env Environment) RuntimeValue {
    object := ObjectValue{
        Properties: make(map[string]RuntimeValue),
    }

    for _, property := range obj.Properties {
        key := property.Key
        value, err := Evaluate(property.Value, env)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        var runtimeValue RuntimeValue
        if value == nil {
            runtimeValue = env.LookupVariable(key)
        } else {
            runtimeValue = value
        }

        object.Properties[key] = runtimeValue
    }

    return object
}

func EvaluateConditionalExpression(expr ast.ConditionalExpression, env Environment) RuntimeValue {
    condition, err := Evaluate(expr.Condition, env)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    conditionValue, ok := condition.(BoolValue)
    if !ok {
        fmt.Println("Condition must be a boolean value")
        os.Exit(1)
    }

    if conditionValue.Value {
        return EvaluateStatement(expr.Consequent, env)
    } else {
        if expr.Alternate == nil {
            return NullValue{}
        }
        return EvaluateStatement(expr.Alternate, env)
    }
}

func EvaluateStatement(statement ast.Statement, env Environment) RuntimeValue {
    value, _ := Evaluate(statement, env)
    return value
}

func EvaluateBinaryExpression(binaryExpression ast.BinaryExpression, env Environment) RuntimeValue {
    lhs, err := Evaluate(binaryExpression.Left, env)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    rhs, err := Evaluate(binaryExpression.Right, env)
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

    return MakeNullValue()
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

func EvaluateAssignment(node ast.AssignmentExpression, env Environment) RuntimeValue {
    if node.Assigne.Kind() != ast.IdentifierType {
        fmt.Println("Left hand side of assignment must be an identifier")
        os.Exit(1)
        return nil
    }

    variableName := node.Assigne.(*ast.Identifier).Symbol
    environment, err := Evaluate(node.Value, env)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
        return nil
    }
    return env.AssignVariable(variableName, environment)
}
