package runtimelang

import (
    "fmt"
    "os"
    "strconv"
    "math"

    "github.com/Jamlee977/CustomLanguage/ast"
)

var (
    IsReturnError = fmt.Errorf("return statement error")
    IsBreakError = fmt.Errorf("break statement error")
)

func EvaluateLoopExpression(expr ast.LoopStatement, env Environment) (RuntimeValue, error) {
    scope := NewEnvironment(&env)

    for {
        for _, statement := range expr.Body {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
                return result, IsReturnError
            }
            if statement.Kind() == ast.BreakStatementType {
                return MakeNullValue(), IsBreakError
            }
            _, err := Evaluate(statement, *scope)
            if err == IsBreakError {
                return MakeNullValue(), nil
            }
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
    }
}

func EvaluateWhileExpression(expr ast.WhileStatement, env Environment) (RuntimeValue, error) {
    condition, err := Evaluate(expr.Condition, env)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    scope := NewEnvironment(&env)

    for condition.Get() == true {
        for _, statement := range expr.Body {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
                return result, IsReturnError
            }
            if statement.Kind() == ast.BreakStatementType {
                return MakeNullValue(), IsBreakError
            }
            _, err := Evaluate(statement, *scope)
            if err == IsBreakError {
                return MakeNullValue(), nil
            }
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
        condition, err = Evaluate(expr.Condition, env)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
    }

    return MakeNullValue(), nil
}

func EvaluateConditionalExpression(expr ast.ConditionalStatement, env Environment) (RuntimeValue, error) {
    condition, err := Evaluate(expr.Condition, env)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    scope := NewEnvironment(&env)

    if condition.Get() == true {
        for _, statement := range expr.Body {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }

                // IsReturnError is a special error that is used to indicate that a return statement has been reached
                return result, IsReturnError
            } else if statement.Kind() == ast.BreakStatementType {
                return MakeNullValue(), IsBreakError
            }
            _, err := Evaluate(statement, *scope)
            if err == IsBreakError {
                return MakeNullValue(), nil
            }
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
    } else {
        for _, statement := range expr.Alternate {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }

                return result, IsReturnError
            } else if statement.Kind() == ast.BreakStatementType {
                return MakeNullValue(), IsBreakError
            }
            _, err := Evaluate(statement, *scope)
            if err == IsBreakError {
                return MakeNullValue(), nil
            }
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }
    }

    return MakeNullValue(), nil
}

func EvaluateCallExpression(expr ast.CallExpression, env Environment) RuntimeValue {
    var args []RuntimeValue
    for _, arg := range expr.Args {
        value, err := Evaluate(arg, env)
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        args = append(args, value)
    }

    function, err := Evaluate(expr.Caller, env)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    if function.Type() == NativeFunction {
        result := function.(NativeFunctionValue).Call(args, env)
        return result
    } else if function.Type() == Function {
        fn := function.(*FunctionValue)
        scope := NewEnvironment(&fn.DeclarationEnvironment)

        for i := 0; i < len(fn.Parameters); i++ {
            if i >= len(args) {
                fmt.Println("Error: Not enough arguments")
                os.Exit(1)
            }
            varname := fn.Parameters[i]
            scope.DeclareVariable(varname, args[i], false)
        }

        var result RuntimeValue = MakeNullValue()
        for _, stmt := range fn.Body {
            if stmt.Kind() == ast.ReturnStatementType {
                result, err = Evaluate(stmt, *scope)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }
                return result
            }
            result, err = Evaluate(stmt, *scope)
            if err == IsReturnError {
                return result
            }
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }
        }

        return result
    }

    fmt.Println("Error: Not a function")
    os.Exit(1)
    return nil
}

func EvaluateMemberExpression(expr ast.MemberExpression, env Environment) RuntimeValue {
    panic("Not implemented")
}

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
        if _, ok := value.(RuntimeValue); ok {
            runtimeValue = value
        } else {
            runtimeValue = value
        }

        object.Properties[key] = runtimeValue
    }

    return object
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

    fmt.Printf("Unknown operator %s for string\n", op)
    os.Exit(1)
    return nil
}

func EvaluateNumericStringBinaryExpression(lhs NumberValue, rhs StringValue, op string) RuntimeValue {
    if op == "+" {
        lhsAsString := strconv.FormatFloat(lhs.Value, 'f', -1, 64)

        return StringValue{lhsAsString + rhs.Value}
    }

    fmt.Printf("Unknown operator %s for string\n", op)
    os.Exit(1)
    return nil

}

func EvaluateStringBinaryExpression(lhs, rhs StringValue, op string) RuntimeValue {
    if op == "+" {
        return StringValue{lhs.Value + rhs.Value}
    }

    fmt.Printf("Unknown operator %s for string\n", op)
    os.Exit(1)
    return nil

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
            fmt.Printf("Division by zero\n")
            os.Exit(1)
            return nil
        }
        result = lhs.Value / rhs.Value
    } else if op == "%" {
        result = math.Mod(lhs.Value, rhs.Value)
    } else if op == ">" {
        isGreaterThan := lhs.Value > rhs.Value
        if isGreaterThan {
            return BoolValue{true}
        }
        return BoolValue{false}
    } else if op == "<" {
        isLessThan := lhs.Value < rhs.Value
        if isLessThan {
            return BoolValue{true}
        }
        return BoolValue{false}
    } else if op == ">=" {
        isGreaterThanOrEqual := lhs.Value >= rhs.Value
        if isGreaterThanOrEqual {
            return BoolValue{true}
        }
        return BoolValue{false}
    } else if op == "<=" {
        isLessThanOrEqual := lhs.Value <= rhs.Value
        if isLessThanOrEqual {
            return BoolValue{true}
        }
        return BoolValue{false}
    } else if op == "==" {
        isEqual := lhs.Value == rhs.Value
        if isEqual {
            return BoolValue{true}
        }
        return BoolValue{false}
    } else if op == "!=" {
        isNotEqual := lhs.Value != rhs.Value
        if isNotEqual {
            return BoolValue{true}
        }
        return BoolValue{false}
    } else {
        fmt.Printf("Error: Unknown operator: %s\n", op)
        os.Exit(1)
        return nil

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
