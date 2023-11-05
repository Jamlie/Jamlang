package runtimelang

import (
    "fmt"
    "os"
    "strconv"
    "math"
    "io/ioutil"

    "github.com/Jamlie/Jamlang/ast"
    "github.com/Jamlie/Jamlang/parser"
)

var (
    IsReturnError = fmt.Errorf("return statement error")
    IsBreakError = fmt.Errorf("break statement error")
)

func EvaluateFunctionDeclaration(expr ast.FunctionDeclaration, env *Environment) (RuntimeValue, error) {
    var fn FunctionValue
    if expr.Name == "" {
        expr.IsAnonymous = true
    }
    if expr.IsAnonymous {
        fn = FunctionValue{
            Body: expr.CloneBody(),
            Parameters: expr.CloneParameters(),
            DeclarationEnvironment: *env,
            IsAnonymous: true,
        }
    } else {
        fn = FunctionValue{
            Name: expr.Name,
            Body: expr.CloneBody(),
            Parameters: expr.CloneParameters(),
            DeclarationEnvironment: *env,
        }

        env.DeclareVariable(expr.Name, fn, true, ast.FunctionType)
    }

    return fn, nil
}

func EvaluateClassDeclaration(expr ast.ClassDeclaration, env *Environment) (RuntimeValue, error) {
    class := ObjectValue{
        Properties: map[string]RuntimeValue{},
        IsClass: true,
    }
    actualClass := ClassValue{
        Name: expr.Name,
        Methods: make(map[string]*FunctionValue),
    }

    scope := NewEnvironment(env)
    
    for _, method := range expr.Body {
        if method.Kind() == ast.FunctionDeclarationType {
            if method.(*ast.FunctionDeclaration).Name == "constructor" {
                method := method.(*ast.FunctionDeclaration)
                actualClass.Constructor = &FunctionValue{
                    Name: method.Name,
                    Body: method.CloneBody(),
                    Parameters: method.CloneParameters(),
                    DeclarationEnvironment: *scope,
                }

                env.DeclareVariable(expr.Name, actualClass.Constructor, true, ast.FunctionType)

                continue
            }
            method := method.(*ast.FunctionDeclaration)

            class.Properties[method.Name] = &FunctionValue{
                Name: method.Name,
                Body: method.CloneBody(),
                Parameters: method.CloneParameters(),
                DeclarationEnvironment: *scope,
            }

        } else if method.Kind() == ast.VariableDeclarationType {
            method := method.(*ast.VariableDeclaration)
            _, err := Evaluate(method.Value, *scope)
            if err != nil {
                return MakeNullValue(), err
            }
        } else {
            return MakeNullValue(), fmt.Errorf("unexpected method type")
        }
    }

    scope.DeclareVariable("this", class, false, ast.ObjectType)
    return class, nil
}

func EvaluateImportExpression(expr ast.ImportStatement, env *Environment) (RuntimeValue, error) {
    file, err := os.Open(expr.Path)
    if err != nil {
        return MakeNullValue(), err
    }
    defer file.Close()

    parser := parser.NewParser()
    var fileString string
    if fileBytes, err := ioutil.ReadAll(file); err != nil {
        return MakeNullValue(), err
    } else {
        fileString = string(fileBytes)
    }

    path := expr.Path
    if path[0] == '.' {
        path = path[1:]
    }
    path = path[:len(path)-4]
    if path[0] == '/' {
        path = path[1:]
    }


    program := parser.ProduceAST(fileString)
    for _, statement := range program.Body {
        if statement.Kind() == ast.FunctionDeclarationType {
            function := statement.(*ast.FunctionDeclaration)
            fn, _ := Evaluate(function, *env)
            if _, ok := fn.(FunctionValue); !ok {
                continue
            }

            if function.Name[0] >= 'A' && function.Name[0] <= 'Z' {
                env.variables[function.Name] = fn
            }
        } else if statement.Kind() == ast.VariableDeclarationType {
            variable := statement.(*ast.VariableDeclaration)
            value, _ := Evaluate(variable.Value, *env)
            if _, ok := value.(RuntimeValue); !ok {
                continue
            }

            if variable.Identifier[0] >= 'A' && variable.Identifier[0] <= 'Z' {
                env.variables[variable.Identifier] = value
            }
        }
    }

    return MakeNullValue(), nil
}

func EvaluateForExpression(expr ast.ForStatement, env *Environment) (RuntimeValue, error) {
    scope := NewEnvironment(env)

    initialVariable := expr.Init.(*ast.VariableDeclaration).Identifier
    _, err := Evaluate(expr.Init, *env)
    defer func() {
        env.RemoveVariable(initialVariable)
    }()

    if err != nil {
        return MakeNullValue(), err
    }

    for {
        if expr.Condition != nil {
            condition, err := Evaluate(expr.Condition, *scope)
            if err != nil {
                return MakeNullValue(), err
            }
            if condition.Type() == Bool {
                if !condition.(BoolValue).Value {
                    break
                }
            } else {
                return MakeNullValue(), fmt.Errorf("for loop condition must be a boolean value")
            }
        }

        for _, statement := range expr.Body {
            if statement.Kind() == ast.BreakStatementType {
                return MakeNullValue(), nil
            }

            _, err := Evaluate(statement, *scope)

            if err != nil {
                return MakeNullValue(), err
            }
        }
        if expr.Update != nil {
            _, err := Evaluate(expr.Update, *scope)
            if err != nil {
                return MakeNullValue(), err
            }
        }

        for k, v := range scope.variables {
            if _, ok := env.variables[k]; ok {
                env.variables[k] = v
            }
        }

        scope = NewEnvironment(env)
    }

    return MakeNullValue(), nil
}

func EvaluateForEachExpression(expr ast.ForEachStatement, env *Environment) (RuntimeValue, error) {
    scope := NewEnvironment(env)

    array, err := Evaluate(expr.Collection, *env)
    if err != nil {
        return MakeNullValue(), err
    }

    if array.Type() == Array {
        scope.DeclareVariable(expr.Variable, MakeNullValue(), false, ast.AnyType)
        for _, element := range array.(ArrayValue).Values {
            scope.AssignVariable(expr.Variable, element)

            for _, statement := range expr.Body {
                if statement.Kind() == ast.ReturnStatementType {
                    return Evaluate(statement, *scope)
                }
                _, err := Evaluate(statement, *scope)
                if err != nil {
                    return MakeNullValue(), err
                }
            }
        }
    } else if array.Type() == Tuple {
        scope.DeclareVariable(expr.Variable, MakeNullValue(), false, ast.AnyType)
        for _, element := range array.(TupleValue).Values {
            scope.AssignVariable(expr.Variable, element)

            for _, statement := range expr.Body {
                if statement.Kind() == ast.ReturnStatementType {
                    return Evaluate(statement, *scope)
                }
                _, err := Evaluate(statement, *scope)
                if err != nil {
                    return MakeNullValue(), err
                }
            }
        }
    } else if array.Type() == String {
        scope.DeclareVariable(expr.Variable, MakeNullValue(), false, ast.StringType)
        for _, element := range array.(StringValue).Value {
            scope.AssignVariable(expr.Variable, StringValue{Value: string(element)})

            for _, statement := range expr.Body {
                if statement.Kind() == ast.ReturnStatementType {
                    return Evaluate(statement, *scope)
                }
                _, err := Evaluate(statement, *scope)
                if err != nil {
                    return MakeNullValue(), err
                }
            }
        }
    } else {
        return MakeNullValue(), fmt.Errorf("Cannot iterate over non-iterable type")
    }

    return MakeNullValue(), nil
}

func EvaluateLoopExpression(expr ast.LoopStatement, env *Environment) (RuntimeValue, error) {
    scope := NewEnvironment(env)

    for {
        for _, statement := range expr.Body {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Fprintln(os.Stderr, err)
                    os.Exit(0)
                }
                return result, IsReturnError
            }
            if statement.Kind() == ast.BreakStatementType {
                return MakeNullValue(), nil
            }

            _, err := Evaluate(statement, *scope)
            if err == IsBreakError {
                return MakeNullValue(), nil
            }
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
        }

        // ! why is this here?
        // for k, v := range scope.variables {
        //     if _, ok := env.variables[k]; ok {
        //         env.variables[k] = v
        //     }
        // }

        scope = NewEnvironment(env)
    }
}

func EvaluateWhileExpression(expr ast.WhileStatement, env *Environment) (RuntimeValue, error) {
    condition, err := Evaluate(expr.Condition, *env)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }

    scope := NewEnvironment(env)

    if condition.Type() != Bool {
        return MakeNullValue(), fmt.Errorf("while statement condition must be a boolean")
    }

    for condition.Get() == true {
        for _, statement := range expr.Body {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Fprintln(os.Stderr, err)
                    os.Exit(0)
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
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
        }

        // ! why was this a thing?
        // for k, v := range scope.variables {
        //     if _, ok := env.variables[k]; ok {
        //         env.variables[k] = v
        //     }
        // }

        scope = NewEnvironment(env)
        condition, err = Evaluate(expr.Condition, *env)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(0)
        }
    }

    return MakeNullValue(), nil
}

func EvaluateConditionalExpression(expr ast.ConditionalStatement, env *Environment) (RuntimeValue, error) {
    condition, err := Evaluate(expr.Condition, *env)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }

    scope := NewEnvironment(env)
    // scope.variables = env.variables

    if condition.Type() != Bool {
        return MakeNullValue(), fmt.Errorf("if statement condition must be a boolean")
    }

    if condition.Get() == true {
        for _, statement := range expr.Body {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Fprintln(os.Stderr, err)
                    os.Exit(0)
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
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
        }

        // for k, v := range scope.variables {
        //     if _, ok := env.variables[k]; ok {
        //         env.variables[k] = v
        //     }
        // }

        scope = NewEnvironment(env)
    } else {
        for _, statement := range expr.Alternate {
            if statement.Kind() == ast.ReturnStatementType {
                result, err := Evaluate(statement, *scope)
                if err != nil {
                    fmt.Fprintln(os.Stderr, err)
                    os.Exit(0)
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
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
        }

        // for k, v := range scope.variables {
        //     if _, ok := env.variables[k]; ok {
        //         env.variables[k] = v
        //     }
        // }

        scope = NewEnvironment(env)
    }


    return MakeNullValue(), nil
}

func EvaluateCallExpression(expr ast.CallExpression, env Environment) RuntimeValue {
    var args []RuntimeValue
    for _, arg := range expr.Args {
        value, err := Evaluate(arg, env)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(0)
        }

        args = append(args, value)
    }

    function, err := Evaluate(expr.Caller, env)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }

    if function.Type() == NativeFunction {
        result := function.(NativeFunctionValue).Call(args, env)
        return result
    } else if function.Type() == Function {
        fn := function.(FunctionValue)
        scope := NewEnvironment(&fn.DeclarationEnvironment)

        for i := 0; i < len(fn.Parameters); i++ {
            if i >= len(args) {
                fmt.Fprintln(os.Stderr, "Error: Not enough arguments")
                os.Exit(0)
            }
            varname := fn.Parameters[i]
            scope.DeclareVariable(varname, args[i], false, ast.FunctionType)
        }
        var result RuntimeValue = MakeNullValue()
        for _, stmt := range fn.Body {
            if stmt.Kind() == ast.ReturnStatementType {
                result, err = Evaluate(stmt, *scope)
                if err != nil {
                    fmt.Fprintln(os.Stderr, err)
                    os.Exit(0)
                }
                return result
            }
            result, err = Evaluate(stmt, *scope)
            if err == IsReturnError {
                return result
            }
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
        }

        return result
    }

    fmt.Fprintln(os.Stderr, "Error: Not a function")
    os.Exit(0)
    return nil
}

func EvaluateMemberExpression(expr ast.MemberExpression, env Environment) RuntimeValue {
    obj, err := Evaluate(expr.Object, env)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }

    if expr.Computed {
        property, err := Evaluate(expr.Property, env)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(0)
        }
        if _, ok := obj.(ArrayValue); ok {
            if property.(NumberValue).Value >= float64(len(obj.(ArrayValue).Values)) {
                fmt.Fprintln(os.Stderr, "Error: Index out of bounds")
                os.Exit(0)
            }

            if property.(NumberValue).Value < 0 {
                if -property.(NumberValue).Value > float64(len(obj.(ArrayValue).Values)) {
                    fmt.Fprintln(os.Stderr, "Error: Index out of bounds")
                    os.Exit(0)
                }
                return obj.(ArrayValue).Values[int(property.(NumberValue).Value) + len(obj.(ArrayValue).Values)]
            }

            return obj.(ArrayValue).Values[int(property.(NumberValue).Value)]
        }

        if _, ok := obj.(TupleValue); ok {
            if property.(NumberValue).Value >= float64(len(obj.(TupleValue).Values)) {
                fmt.Fprintln(os.Stderr, "Error: Index out of bounds")
                os.Exit(0)
            }

            if property.(NumberValue).Value < 0 {
                if -property.(NumberValue).Value > float64(len(obj.(TupleValue).Values)) {
                    fmt.Fprintln(os.Stderr, "Error: Index out of bounds")
                    os.Exit(0)
                }
                return obj.(TupleValue).Values[int(property.(NumberValue).Value) + len(obj.(TupleValue).Values)]
            }

            return obj.(TupleValue).Values[int(property.(NumberValue).Value)]
        }

        if _, ok := obj.(StringValue); ok {
            if property.(NumberValue).Value >= float64(len(obj.(StringValue).Value)) {
                fmt.Fprintln(os.Stderr, "Error: Index out of bounds")
                os.Exit(0)
            }

            if property.(NumberValue).Value < 0 {
                if -property.(NumberValue).Value > float64(len(obj.(StringValue).Value)) {
                    fmt.Fprintln(os.Stderr, "Error: Index out of bounds")
                    os.Exit(0)
                }
                return MakeStringValue(string(obj.(StringValue).Value[int(property.(NumberValue).Value) + len(obj.(StringValue).Value)]))
            }

            return MakeStringValue(string(obj.(StringValue).Value[int(property.(NumberValue).Value)]))
        }

        return obj.(ObjectValue).Properties[property.(StringValue).Value]
    } else {
        if _, ok := obj.(NullValue); ok {
            return obj
        }

        if _, ok := obj.(ArrayValue); ok {
            switch expr.Property.(*ast.Identifier).Symbol {
            case "length":
                return MakeNumberValue(float64(len(obj.(ArrayValue).Values)))
            case "push":
                if arr, ok := obj.(ArrayValue); ok {
                    a := &arr.Values
                    return jamlangArrayPush(&a)
                }
            case "pop":
                return jamlangArrayPop(obj.(ArrayValue).Values)
            case "shift":
                return jamlangArrayShift(obj.(ArrayValue).Values)
            case "contains":
                return jamlangArrayContains(obj.(ArrayValue).Values)
            case "insert":
                return jamlangArrayInsertInto(obj.(ArrayValue).Values)
            case "pushAll":
                return jamlangArrayPushAll(obj.(ArrayValue).Values)
            default:
                fmt.Fprintln(os.Stderr, "Error: Array does not have property " + expr.Property.(*ast.Identifier).Symbol)
                os.Exit(0)
            }
        }

        if _, ok := obj.(TupleValue); ok {
            switch expr.Property.(*ast.Identifier).Symbol {
            case "length":
                return MakeNumberValue(float64(len(obj.(TupleValue).Values)))
            default:
                fmt.Fprintln(os.Stderr, "Error: Tuple does not have property " + expr.Property.(*ast.Identifier).Symbol)
                os.Exit(0)
            }
        }

        if _, ok := obj.(StringValue); ok {
            switch expr.Property.(*ast.Identifier).Symbol {
            case "length":
                return MakeNumberValue(float64(len(obj.(StringValue).Value)))
            case "toUpper":
                return jamlangStringToUpper(obj.(StringValue).Value)
            case "toLower":
                return jamlangStringToLower(obj.(StringValue).Value)
            case "contains":
                return jamlangStringContains(obj.(StringValue).Value)
            case "split":
                return jamlangStringSplit(obj.(StringValue).Value)
            case "equalsIgnoreCase":
                return jamlangStringEqualsIgnoreCase(obj.(StringValue).Value)
            case "startsWith":
                return jamlangStringStartsWith(obj.(StringValue).Value)
            case "endsWith":
                return jamlangStringEndsWith(obj.(StringValue).Value)
            case "indexOf":
                return jamlangStringIndexOf(obj.(StringValue).Value)
            case "lastIndexOf":
                return jamlangStringLastIndexOf(obj.(StringValue).Value)
            case "substring":
                return jamlangStringSubstring(obj.(StringValue).Value)
            case "replace":
                return jamlangStringReplace(obj.(StringValue).Value)
            case "trim":
                return jamlangStringTrim(obj.(StringValue).Value)
            case "trimLeft":
                return jamlangStringTrimLeft(obj.(StringValue).Value)
            case "trimRight":
                return jamlangStringTrimRight(obj.(StringValue).Value)
            case "repeat":
                return jamlangStringRepeat(obj.(StringValue).Value)
            case "leftPad":
                return jamlangStringLeftPad(obj.(StringValue).Value)
            case "rightPad":
                return jamlangStringRightPad(obj.(StringValue).Value)
            default:
                fmt.Fprintln(os.Stderr, "Error: String has no property " + expr.Property.(*ast.Identifier).Symbol)
                os.Exit(0)
            }
        }

        if _, ok := obj.(ObjectValue); !ok {
            fmt.Fprintf(os.Stderr, "Error: %s has no property %s\n", obj.Type(), expr.Property.(*ast.Identifier).Symbol)
            os.Exit(0)
        }

        return obj.(ObjectValue).Properties[expr.Property.(*ast.Identifier).Symbol]
    }
}

func EvaluateArrayExpression(expr ast.ArrayLiteral, env Environment) RuntimeValue {
    array := ArrayValue{
        Values: make([]RuntimeValue, len(expr.Elements)),
    }

    for i, element := range expr.Elements {
        array.Values[i], _ = Evaluate(element, env)
    }

    return array
}

func EvaluateObjectExpression(obj ast.ObjectLiteral, env Environment) RuntimeValue {
    object := ObjectValue{
        Properties: make(map[string]RuntimeValue),
    }

    for _, property := range obj.Properties {
        key := property.Key
        var value RuntimeValue
        var err error
        if property.Value != nil {
            value, err = Evaluate(property.Value, env)
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
        } else {
            value = env.LookupVariable(key)
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
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }
    rhs, err := Evaluate(binaryExpression.Right, env)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }

    if lhs.Type() == Number && rhs.Type() == Number {
        return EvaluateNumericBinaryExpression(lhs.(NumberValue), rhs.(NumberValue), binaryExpression.Operator)
    } else if lhs.Type() == String && rhs.Type() == String {
        return EvaluateStringBinaryExpression(lhs.(StringValue), rhs.(StringValue), binaryExpression.Operator)
    } else if lhs.Type() == String && rhs.Type() == Number {
        return EvaluateStringNumericBinaryExpression(lhs.(StringValue), rhs.(NumberValue), binaryExpression.Operator)
    } else if lhs.Type() == Number && rhs.Type() == String {
        return EvaluateNumericStringBinaryExpression(lhs.(NumberValue), rhs.(StringValue), binaryExpression.Operator)
    } else if lhs.Type() == Null || rhs.Type() == Null {
        return EvaluateNullBinaryExpression(lhs, rhs, binaryExpression.Operator)
    } else if lhs.Type() == Object && rhs.Type() == Object {
        return EvaluateObjectBinaryExpression(lhs.(ObjectValue), rhs.(ObjectValue), binaryExpression.Operator)
    }


    return MakeNullValue()
}


func EvaluateObjectBinaryExpression(lhs ObjectValue, rhs ObjectValue, op string) RuntimeValue {
    if op == "==" {
        for key, value := range lhs.Properties {
            if !value.Equals(rhs.Properties[key]) {
                return BoolValue{false}
            }
        }
        return BoolValue{true}
    } else if op == "!=" {
        for key, value := range lhs.Properties {
            if !value.Equals(rhs.Properties[key]) {
                return BoolValue{true}
            }
        }
        return BoolValue{false}
    }

    return MakeNullValue()
}

func EvaluateNullBinaryExpression(lhs RuntimeValue, rhs RuntimeValue, op string) RuntimeValue {
    if lhs.Type() == Null && rhs.Type() == Null {
        if op == "==" {
            return BoolValue{true}
        }
        return BoolValue{false}
    } else if lhs.Type() == Bool && rhs.Type() == Null {
        if op == "==" {
            return BoolValue{false}
        }
        return BoolValue{true}
    } else if lhs.Type() == Null && rhs.Type() == Bool {
        if op == "==" {
            return BoolValue{false}
        }
        return BoolValue{true}
    } else if lhs.Type() == Object && rhs.Type() == Null {
        if op == "==" {
            return BoolValue{false}
        }
        return BoolValue{true}
    } else if lhs.Type() == Null && rhs.Type() == Object {
        if op == "==" {
            return BoolValue{false}
        }
        return BoolValue{true}
    }

    return MakeBoolValue(false)
}

func EvaluateStringNumericBinaryExpression(lhs StringValue, rhs NumberValue, op string) RuntimeValue {
    if op == "+" {
        rhsAsString := strconv.FormatFloat(rhs.Value, 'f', -1, 64)
        return StringValue{lhs.Value + rhsAsString}
    }

    fmt.Fprintf(os.Stderr, "Unknown operator %s for string\n", op)
    os.Exit(0)
    return nil
}

func EvaluateNumericStringBinaryExpression(lhs NumberValue, rhs StringValue, op string) RuntimeValue {
    if op == "+" {
        lhsAsString := strconv.FormatFloat(lhs.Value, 'f', -1, 64)

        return StringValue{lhsAsString + rhs.Value}
    }

    fmt.Fprintf(os.Stderr, "Unknown operator %s for string\n", op)
    os.Exit(0)
    return nil
}

func EvaluateStringBinaryExpression(lhs, rhs StringValue, op string) RuntimeValue {
    if op == "+" {
        return StringValue{lhs.Value + rhs.Value}
    } else if op == "==" {
        return BoolValue{lhs.Value == rhs.Value}
    } else if op == "!=" {
        return BoolValue{lhs.Value != rhs.Value}
    }

    fmt.Fprintf(os.Stderr, "Unknown operator %s for string\n", op)
    os.Exit(0)
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
    } else if op == "**" {
        result = math.Pow(lhs.Value, rhs.Value)
    } else if op == "/" {
        if rhs.Value == 0 {
            fmt.Fprintf(os.Stderr, "Division by zero\n")
            os.Exit(0)
            return nil
        }
        result = lhs.Value / rhs.Value
    } else if op == "//" {
        if rhs.Value == 0 {
            fmt.Fprintf(os.Stderr, "Division by zero\n")
            os.Exit(0)
            return nil
        }
        result = math.Floor(lhs.Value / rhs.Value)
    } else if op == "%" {
        result = math.Mod(lhs.Value, rhs.Value)
    } else if op == "&" {
        result = float64(int64(lhs.Value) & int64(rhs.Value))
    } else if op == "|" {
        result = float64(int64(lhs.Value) | int64(rhs.Value))
    } else if op == "^" {
        result = float64(int64(lhs.Value) ^ int64(rhs.Value))
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
    } else if op == "<<" {
        result = float64(int64(lhs.Value) << int64(rhs.Value))
    } else if op == ">>" {
        result = float64(int64(lhs.Value) >> int64(rhs.Value))
    } else {
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", op)
        os.Exit(0)
        return nil

    }

    return NumberValue{result}
}

func EvaluateUnaryExpression(node ast.UnaryExpression, env Environment) RuntimeValue {
    value, err := Evaluate(node.Value, env)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
    }

    switch node.Operator {
    case "!":
        if value.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: ! operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }
        return BoolValue{!value.(BoolValue).Value}
    case "++":
        if value.Type() != Number {
            fmt.Fprintln(os.Stderr, "Error: ++ operator can only be applied to number values")
            os.Exit(0)
            return nil
        }
        if node.Value.Kind() == ast.MemberExpressionType {
            val, err := Evaluate(node.Value.(*ast.MemberExpression).Property, env)
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
            arr, err := Evaluate(node.Value.(*ast.MemberExpression).Object, env)
            if err != nil {
                fmt.Fprintln(os.Stderr, err)
                os.Exit(0)
            }
            arr.(ArrayValue).Values[int(ToGoNumberValue(val.(NumberValue)))] = NumberValue{val.(NumberValue).Value + 1}
            return NumberValue{val.(NumberValue).Value + 1}
        }
        env.AssignVariable(node.Value.(*ast.Identifier).Symbol, NumberValue{value.(NumberValue).Value + 1})
        return NumberValue{value.(NumberValue).Value + 1}
    case "--":
        if value.Type() != Number {
            fmt.Fprintln(os.Stderr, "Error: -- operator can only be applied to number values")
            os.Exit(0)
            return nil
        }
        env.AssignVariable(node.Value.(*ast.Identifier).Symbol, NumberValue{value.(NumberValue).Value - 1})
        return NumberValue{value.(NumberValue).Value - 1}
    case "-":
        if value.Type() != Number {
            fmt.Fprintln(os.Stderr, "Error: - operator can only be applied to number values")
            os.Exit(0)
            return nil
        }
        return NumberValue{-value.(NumberValue).Value}
    default:
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", node.Operator)
        os.Exit(0)
        return nil
    }
}

func EvaluateLogicalExpression(node ast.LogicalExpression, env Environment) RuntimeValue {
    switch node.Operator {
    case "and":
        left, err := Evaluate(node.Left, env)
        if err != nil {
            return nil
        }
        if left.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: and operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }
        if left.(BoolValue).Value == false {
            return BoolValue{false}
        }

        right, err := Evaluate(node.Right, env)
        if err != nil {
            return nil
        }
        if right.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: and operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }
        return BoolValue{right.(BoolValue).Value}
    case "or":
        left, err := Evaluate(node.Left, env)
        if err != nil {
            return nil
        }
        if left.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: or operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }
        if left.(BoolValue).Value == true {
            return BoolValue{true}
        }

        right, err := Evaluate(node.Right, env)
        if err != nil {
            return nil
        }
        if right.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: or operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }
        return BoolValue{right.(BoolValue).Value}
    case "xor":
        left, err := Evaluate(node.Left, env)
        if err != nil {
            return nil
        }
        if left.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: xor operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }

        right, err := Evaluate(node.Right, env)
        if err != nil {
            return nil
        }
        if right.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: xor operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }
        return BoolValue{left.(BoolValue).Value != right.(BoolValue).Value}
    case "not":
        operand, err := Evaluate(node.Right, env)
        if err != nil {
            return nil
        }
        if operand.Type() != Bool {
            fmt.Fprintln(os.Stderr, "Error: not operator can only be applied to boolean values")
            os.Exit(0)
            return nil
        }
        return BoolValue{!operand.(BoolValue).Value}
    default:
        fmt.Fprintln(os.Stderr, "Error: unknown operator")
        os.Exit(0)
        return nil
    }
}

func EvaluateAssignment(node ast.AssignmentExpression, env Environment) RuntimeValue {
    if node.Assignee.Kind() == ast.MemberExpressionType {
        objectLiteral := node.Assignee.(*ast.MemberExpression).Object
        objectValue, _ := Evaluate(objectLiteral, env)
        if objectValue.Type() == Array {
            index, _ := Evaluate(node.Assignee.(*ast.MemberExpression).Property, env)
            if index.Type() != Number {
                fmt.Fprintln(os.Stderr, "Error: array index must be a number")
                os.Exit(0)
                return nil
            }
            if index.(NumberValue).Value >= float64(len(objectValue.(ArrayValue).Values)) {
                fmt.Fprintln(os.Stderr, "Error: array index out of bounds")
                os.Exit(0)
                return nil
            }
            if index.(NumberValue).Value < 0 {
                if -index.(NumberValue).Value > float64(len(objectValue.(ArrayValue).Values)) {
                    fmt.Fprintln(os.Stderr, "Error: array index out of bounds")
                    os.Exit(0)
                    return nil
                }
                objectValue.(ArrayValue).Values[len(objectValue.(ArrayValue).Values) + int(index.(NumberValue).Value)], _ = Evaluate(node.Value, env)
                return objectValue
            }

            objectValue.(ArrayValue).Values[int(index.(NumberValue).Value)], _ = Evaluate(node.Value, env)
            return objectValue
        }
        if objectValue.Type() == Null {
            return objectValue.(NullValue)
        }
        if objectValue.Type() == String {
            fmt.Fprintln(os.Stderr, "Error: string does not support assignment")
            os.Exit(0)
        }
        objectValue.(ObjectValue).Properties[node.Assignee.(*ast.MemberExpression).Property.(*ast.Identifier).Symbol], _ = Evaluate(node.Value, env)
        return objectValue
    }

    if node.Assignee.Kind() != ast.IdentifierType {
        fmt.Fprintln(os.Stderr, "Error: Left side of assignment must be a variable")
        os.Exit(0)
        return nil
    }

    variableName := node.Assignee.(*ast.Identifier).Symbol
    environment, err := Evaluate(node.Value, env)
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(0)
        return nil
    }
    return env.AssignVariable(variableName, environment)
}
