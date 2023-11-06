package runtimelang

import (
    "fmt"
    "os"
)

func EvaluateI8BinaryExpression(lhs Int8Value, rhs RuntimeValue, op string) RuntimeValue {
    var result int8 = 0
    if op == "+" {
        return internal_I8Plus(lhs, rhs)
    } else if op == "-" {
        return internal_I8Minus(lhs, rhs)
    } else if op == "*" {
        return internal_I8Mult(lhs, rhs)
    } else if op == "**" {
        return internal_I8Pow(lhs, rhs)
    } else if op == "/" {
        return internal_I8Div(lhs, rhs)
    } else if op == "//" {
        return internal_I8IntDiv(lhs, rhs)
    } else if op == "%" {
        return internal_I8Mod(lhs, rhs)
    } else if op == "&" {
        return internal_I8BitwiseAnd(lhs, rhs)
    } else if op == "|" {
        return internal_I8BitwiseOr(lhs, rhs)
    } else if op == "^" {
        return internal_I8BitwiseXor(lhs, rhs)
    } else if op == ">" {
        return internal_I8GreaterThan(lhs, rhs)
    } else if op == "<" {
        return internal_I8LessThan(lhs, rhs)
    } else if op == ">=" {
        return internal_I8GreaterThanEqual(lhs, rhs)
    } else if op == "<=" {
        return internal_I8LessThanEqual(lhs, rhs)
    } else if op == "==" {
        return internal_I8Equal(lhs, rhs)
    } else if op == "!=" {
        return internal_I8NotEqual(lhs, rhs)
    } else if op == "<<" {
        return internal_I8LeftShift(lhs, rhs)
    } else if op == ">>" {
        return internal_I8RightShift(lhs, rhs)
    } else {
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", op)
        os.Exit(0)
    }

    return Int8Value{result}
}

func EvaluateI16BinaryExpression(lhs Int16Value, rhs RuntimeValue, op string) RuntimeValue {
    var result int16 = 0
    if op == "+" {
        return internal_I16Plus(lhs, rhs)
    } else if op == "-" {
        return internal_I16Minus(lhs, rhs)
    } else if op == "*" {
        return internal_I16Mult(lhs, rhs)
    } else if op == "**" {
        return internal_I16Pow(lhs, rhs)
    } else if op == "/" {
        return internal_I16Div(lhs, rhs)
    } else if op == "//" {
        return internal_I16IntDiv(lhs, rhs)
    } else if op == "%" {
        return internal_I16Mod(lhs, rhs)
    } else if op == "&" {
        return internal_I16BitwiseAnd(lhs, rhs)
    } else if op == "|" {
        return internal_I16BitwiseOr(lhs, rhs)
    } else if op == "^" {
        return internal_I16BitwiseXor(lhs, rhs)
    } else if op == ">" {
        return internal_I16GreaterThan(lhs, rhs)
    } else if op == "<" {
        return internal_I16LessThan(lhs, rhs)
    } else if op == ">=" {
        return internal_I16GreaterThanEqual(lhs, rhs)
    } else if op == "<=" {
        return internal_I16LessThanEqual(lhs, rhs)
    } else if op == "==" {
        return internal_I16Equal(lhs, rhs)
    } else if op == "!=" {
        return internal_I16NotEqual(lhs, rhs)
    } else if op == "<<" {
        return internal_I16LeftShift(lhs, rhs)
    } else if op == ">>" {
        return internal_I16RightShift(lhs, rhs)
    } else {
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", op)
        os.Exit(0)
    }

    return Int16Value{result}
}

func EvaluateI32BinaryExpression(lhs Int32Value, rhs RuntimeValue, op string) RuntimeValue {
    var result int32 = 0
    if op == "+" {
        return internal_I32Plus(lhs, rhs)
    } else if op == "-" {
        return internal_I32Minus(lhs, rhs)
    } else if op == "*" {
        return internal_I32Mult(lhs, rhs)
    } else if op == "**" {
        return internal_I32Pow(lhs, rhs)
    } else if op == "/" {
        return internal_I32Div(lhs, rhs)
    } else if op == "//" {
        return internal_I32IntDiv(lhs, rhs)
    } else if op == "%" {
        return internal_I32Mod(lhs, rhs)
    } else if op == "&" {
        return internal_I32BitwiseAnd(lhs, rhs)
    } else if op == "|" {
        return internal_I32BitwiseOr(lhs, rhs)
    } else if op == "^" {
        return internal_I32BitwiseXor(lhs, rhs)
    } else if op == ">" {
        return internal_I32GreaterThan(lhs, rhs)
    } else if op == "<" {
        return internal_I32LessThan(lhs, rhs)
    } else if op == ">=" {
        return internal_I32GreaterThanEqual(lhs, rhs)
    } else if op == "<=" {
        return internal_I32LessThanEqual(lhs, rhs)
    } else if op == "==" {
        return internal_I32Equal(lhs, rhs)
    } else if op == "!=" {
        return internal_I32NotEqual(lhs, rhs)
    } else if op == "<<" {
        return internal_I32LeftShift(lhs, rhs)
    } else if op == ">>" {
        return internal_I32RightShift(lhs, rhs)
    } else {
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", op)
        os.Exit(0)
    }

    return Int32Value{result}
}

func EvaluateI64BinaryExpression(lhs Int64Value, rhs RuntimeValue, op string) RuntimeValue {
    var result int64 = 0
    if op == "+" {
        return internal_I64Plus(lhs, rhs)
    } else if op == "-" {
        return internal_I64Minus(lhs, rhs)
    } else if op == "*" {
        return internal_I64Mult(lhs, rhs)
    } else if op == "**" {
        return internal_I64Pow(lhs, rhs)
    } else if op == "/" {
        return internal_I64Div(lhs, rhs)
    } else if op == "//" {
        return internal_I64IntDiv(lhs, rhs)
    } else if op == "%" {
        return internal_I64Mod(lhs, rhs)
    } else if op == "&" {
        return internal_I64BitwiseAnd(lhs, rhs)
    } else if op == "|" {
        return internal_I64BitwiseOr(lhs, rhs)
    } else if op == "^" {
        return internal_I64BitwiseXor(lhs, rhs)
    } else if op == ">" {
        return internal_I64GreaterThan(lhs, rhs)
    } else if op == "<" {
        return internal_I64LessThan(lhs, rhs)
    } else if op == ">=" {
        return internal_I64GreaterThanEqual(lhs, rhs)
    } else if op == "<=" {
        return internal_I64LessThanEqual(lhs, rhs)
    } else if op == "==" {
        return internal_I64Equal(lhs, rhs)
    } else if op == "!=" {
        return internal_I64NotEqual(lhs, rhs)
    } else if op == "<<" {
        return internal_I64LeftShift(lhs, rhs)
    } else if op == ">>" {
        return internal_I64RightShift(lhs, rhs)
    } else {
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", op)
        os.Exit(0)
    }

    return Int64Value{result}
}

func EvaluateF32BinaryExpression(lhs Float32Value, rhs RuntimeValue, op string) RuntimeValue {
    var result float32 = 0
    if op == "+" {
        return internal_F32Plus(lhs, rhs)
    } else if op == "-" {
        return internal_F32Minus(lhs, rhs)
    } else if op == "*" {
        return internal_F32Mult(lhs, rhs)
    } else if op == "**" {
        return internal_F32Pow(lhs, rhs)
    } else if op == "/" {
        return internal_F32Div(lhs, rhs)
    } else if op == "//" {
        return internal_F32IntDiv(lhs, rhs)
    } else if op == "%" {
        return internal_F32Mod(lhs, rhs)
    } else if op == ">" {
        return internal_F32GreaterThan(lhs, rhs)
    } else if op == "<" {
        return internal_F32LessThan(lhs, rhs)
    } else if op == ">=" {
        return internal_F32GreaterThanEqual(lhs, rhs)
    } else if op == "<=" {
        return internal_F32LessThanEqual(lhs, rhs)
    } else if op == "==" {
        return internal_F32Equal(lhs, rhs)
    } else if op == "!=" {
        return internal_F32NotEqual(lhs, rhs)
    } else {
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", op)
        os.Exit(0)
    }

    return Float32Value{result}
}

func EvaluateF64BinaryExpression(lhs Float64Value, rhs RuntimeValue, op string) RuntimeValue {
    var result float64 = 0
    if op == "+" {
        return internal_F64Plus(lhs, rhs)
    } else if op == "-" {
        return internal_F64Minus(lhs, rhs)
    } else if op == "*" {
        return internal_F64Mult(lhs, rhs)
    } else if op == "**" {
        return internal_F64Pow(lhs, rhs)
    } else if op == "/" {
        return internal_F64Div(lhs, rhs)
    } else if op == "//" {
        return internal_F64IntDiv(lhs, rhs)
    } else if op == "%" {
        return internal_F64Mod(lhs, rhs)
    } else if op == ">" {
        return internal_F64GreaterThan(lhs, rhs)
    } else if op == "<" {
        return internal_F64LessThan(lhs, rhs)
    } else if op == ">=" {
        return internal_F64GreaterThanEqual(lhs, rhs)
    } else if op == "<=" {
        return internal_F64LessThanEqual(lhs, rhs)
    } else if op == "==" {
        return internal_F64Equal(lhs, rhs)
    } else if op == "!=" {
        return internal_F64NotEqual(lhs, rhs)
    } else {
        fmt.Fprintf(os.Stderr, "Error: Unknown operator: %s\n", op)
        os.Exit(0)
    }

    return Float64Value{result}
}

