package runtimelang

import (
    "strconv"
    
    "github.com/Jamlee977/CustomLanguage/ast"
)

type ValueType string

const (
    Number ValueType = "number"
    Null ValueType = "null"
    String ValueType = "string"
    Bool ValueType = "bool"
    Object ValueType = "object"
    NativeFunction ValueType = "native_function"
    Function ValueType = "function"
    Break ValueType = "break"
)

type RuntimeValue interface {
    Get() any
    Type() ValueType
    ToString() string
}

type InitialValue struct {}

func (v InitialValue) Type() ValueType {
    return Null
}

func (v InitialValue) Get() any {
    return ""
}

func (v InitialValue) ToString() string {
    return ""
}

type NullValue struct {
    Value string
}

func (v NullValue) Type() ValueType {
    return Null
}

func (v NullValue) Get() any {
    return v.Value
}

func (v NullValue) ToString() string {
    return "null"
}

func MakeNullValue() NullValue {
    return NullValue{Value: "null"}
}

type NumberValue struct {
    Value float64
}

func (v NumberValue) Type() ValueType {
    return Number
}

func (v NumberValue) Get() any {
    return v.Value
}

func (v NumberValue) ToString() string {
    return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

func MakeNumberValue(value float64) NumberValue {
    return NumberValue{Value: value}
} 

type StringValue struct {
    Value string
}

func (v StringValue) Type() ValueType {
    return String
}

func (v StringValue) Get() any {
    return v.Value
}

func (v StringValue) ToString() string {
    return v.Value
}

func MakeStringValue(value string) StringValue {
    return StringValue{Value: value}
}

type BoolValue struct {
    Value bool
}

func (v BoolValue) Type() ValueType {
    return Bool
}

func (v BoolValue) Get() any {
    return v.Value
}

func (v BoolValue) ToString() string {
    return strconv.FormatBool(v.Value)
}

func MakeBoolValue(value bool) BoolValue {
    return BoolValue{Value: value}
}

type ObjectValue struct {
    Properties map[string]RuntimeValue
}

func (v ObjectValue) Type() ValueType {
    return Object
}

func (v ObjectValue) Get() any {
    str := "{ "
    counter := 0
    for key, value := range v.Properties {
        counter++
        str += key + ": "

        switch value.Type() {
        case Null:
            str += "null"
        case Number:
            str += value.ToString()
        case String:
            str += value.ToString()
        case Bool:
            str += value.ToString()
        case Object:
            str += value.ToString()            
        default:
            str += "unknown"
        }

        if counter < len(v.Properties) {
            str += ", "
        }
    }
    str += " }"
    return str
}

func (v ObjectValue) ToString() string {
    return v.Get().(string)
}

type FunctionCall func(args []RuntimeValue, env Environment) RuntimeValue

type NativeFunctionValue struct {
    Call FunctionCall
}

func (v NativeFunctionValue) Type() ValueType {
    return NativeFunction
}

func (v NativeFunctionValue) Get() any {
    return v.Call
}

func (v NativeFunctionValue) ToString() string {
    return "native function"
}

func MakeNativeFunction(call FunctionCall) NativeFunctionValue {
    return NativeFunctionValue{Call: call}
}

type FunctionValue struct {
    Name string
    Parameters []string
    DeclarationEnvironment Environment
    Body []ast.Statement
}

func (v FunctionValue) Type() ValueType {
    return Function
}

func (v FunctionValue) Get() any {
    str := "fn " + v.Name + "("
    for i, param := range v.Parameters {
        str += param
        if i < len(v.Parameters) - 1 {
            str += ", "
        }
    }
    str += ") { ... }"
    return str
}

func (v FunctionValue) ToString() string {
    return "function"
}

type ReturnValue struct {
    Value RuntimeValue
}

type BreakType struct {}

func (v BreakType) Type() ValueType {
    return Break
}

func (v BreakType) Get() any {
    return nil
}

func (v BreakType) ToString() string {
    return "break"
}
