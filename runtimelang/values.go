package runtimelang

import (
    "os"
    "strconv"
    // "strings"

    "github.com/Jamlie/Jamlang/ast"
)

type ValueType string

const (
    Number         ValueType = "number"
    Null           ValueType = "null"
    String         ValueType = "string"
    Bool           ValueType = "bool"
    Object         ValueType = "object"
    Array          ValueType = "array"
    Tuple          ValueType = "tuple"
    NativeFunction ValueType = "native_function"
    Function       ValueType = "function"
    Break          ValueType = "break"
    Class          ValueType = "class"
    File           ValueType = "file"
)

type RuntimeValue interface {
    Get() any
    Type() ValueType
    ToString() string
    Clone() RuntimeValue
    Equals(RuntimeValue) bool
}

type InitialValue struct{}

func (v InitialValue) Equals(other RuntimeValue) bool {
    if other.Type() == Null {
        return true
    }
    return false
}

func (v InitialValue) Type() ValueType {
    return Null
}

func (v InitialValue) Get() any {
    return ""
}

func (v InitialValue) ToString() string {
    return ""
}

func (v InitialValue) Clone() RuntimeValue {
    return v
}

type NullValue struct {
    Value string
}

func (v NullValue) Equals(other RuntimeValue) bool {
    if other.Type() == Null {
        return true
    }
    return false
}

func (v NullValue) Type() ValueType {
    return Null
}

func (v NullValue) Get() any {
    return "null"
}

func (v NullValue) ToString() string {
    return "null"
}

func (v NullValue) Clone() RuntimeValue {
    return v
}

func MakeNullValue() NullValue {
    return NullValue{}
}


type NumberValue struct {
    Value float64
}

func (v NumberValue) Equals(other RuntimeValue) bool {
    if other.Type() == Number {
        return v.Value == other.Get().(float64)
    }
    return false
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

func (v NumberValue) Clone() RuntimeValue {
    return v
}

func MakeNumberValue(value float64) NumberValue {
    return NumberValue{Value: value}
}


type StringValue struct {
    Value string
}

func (v StringValue) Equals(other RuntimeValue) bool {
    if other.Type() == String {
        return v.Value == other.Get().(string)
    }
    return false
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

func (v StringValue) Clone() RuntimeValue {
    return v
}

func MakeStringValue(value string) StringValue {
    return StringValue{Value: value}
}


type BoolValue struct {
    Value bool
}

func (v BoolValue) Equals(other RuntimeValue) bool {
    if other.Type() == Bool {
        return v.Value == other.Get().(bool)
    }
    return false
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

func (v BoolValue) Clone() RuntimeValue {
    return v
}

func MakeBoolValue(value bool) BoolValue {
    return BoolValue{Value: value}
}


type ObjectValue struct {
    Properties map[string]RuntimeValue
    IsClass bool
}

func (v ObjectValue) Equals(other RuntimeValue) bool {
    if other.Type() == Object {
        otherObj := other.Get().(ObjectValue)
        if len(v.Properties) != len(otherObj.Properties) {
            return false
        }
        for key, value := range v.Properties {
            if !value.Equals(otherObj.Properties[key]) {
                return false
            }
        }
        return true
    }
    return false
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
        case Array:
            str += value.ToString()
        case Function:
            str += value.Get().(string)
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

func (v ObjectValue) Clone() RuntimeValue {
    newObject := ObjectValue{Properties: make(map[string]RuntimeValue)}
    for key, value := range v.Properties {
        newObject.Properties[key] = value.Clone()
    }
    return newObject
}


type ArrayValue struct {
    Values []RuntimeValue
}

func (v ArrayValue) Equals(other RuntimeValue) bool {
    if other.Type() == Array {
        otherArray := other.(ArrayValue)
        if len(v.Values) != len(otherArray.Values) {
            return false
        }
        for i, value := range v.Values {
            if !value.Equals(otherArray.Values[i]) {
                return false
            }
        }
        return true
    }
    return false
}

func (v ArrayValue) Type() ValueType {
    return Array
}

func (v ArrayValue) Get() any {
    str := "[ "
    for i, value := range v.Values {
        str += value.ToString()
        if i < len(v.Values)-1 {
            str += ", "
        }
    }
    str += " ]"
    return str
}

func (v ArrayValue) ToString() string {
    return v.Get().(string)
}

func (v ArrayValue) Clone() RuntimeValue {
    newArray := ArrayValue{Values: make([]RuntimeValue, len(v.Values))}
    copy(newArray.Values, v.Values)
    return newArray
}


type TupleValue struct {
    Values []RuntimeValue
}

func (v TupleValue) Equals(other RuntimeValue) bool {
    if other.Type() == Tuple {
        otherTuple := other.(TupleValue)
        if len(v.Values) != len(otherTuple.Values) {
            return false
        }
        for i, value := range v.Values {
            if !value.Equals(otherTuple.Values[i]) {
                return false
            }
        }
        return true
    }
    return false
}

func (v TupleValue) Type() ValueType {
    return Tuple
}

func (v TupleValue) Get() any {
    str := "( "
    for i, value := range v.Values {
        str += value.ToString()
        if i < len(v.Values)-1 {
            str += ", "
        }
    }
    str += " )"
    return str
}

func (v TupleValue) ToString() string {
    return v.Get().(string)
}

func (v TupleValue) Clone() RuntimeValue {
    newTuple := TupleValue{Values: make([]RuntimeValue, len(v.Values))}
    copy(newTuple.Values, v.Values)
    return newTuple
}


type FunctionCall func(args []RuntimeValue, env Environment) RuntimeValue

type NativeFunctionValue struct {
    Call FunctionCall
}

func (v NativeFunctionValue) Equals(other RuntimeValue) bool {
    return false
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

func (v NativeFunctionValue) Clone() RuntimeValue {
    return v
}

func MakeNativeFunction(call FunctionCall) NativeFunctionValue {
    return NativeFunctionValue{Call: call}
}


type FunctionValue struct {
    Name                   string
    Parameters             []string
    DeclarationEnvironment Environment
    Body                   []ast.Statement
    IsAnonymous            bool
}

func (v FunctionValue) Equals(other RuntimeValue) bool {
    return false
}

func (v FunctionValue) Type() ValueType {
    return Function
}

func (v FunctionValue) Get() any {
    str := "fn " + v.Name + "("
    for i, param := range v.Parameters {
        str += param
        if i < len(v.Parameters)-1 {
            str += ", "
        }
    }
    str += ") { ... }"
    return str
}

func (v FunctionValue) ToString() string {
    return "function"
}

func (v FunctionValue) Clone() RuntimeValue {
    name := v.Name
    parameters := make([]string, len(v.Parameters))
    copy(parameters, v.Parameters)
    body := make([]ast.Statement, len(v.Body))
    copy(body, v.Body)
    return &FunctionValue{
        Name:                   name,
        Parameters:             parameters,
        DeclarationEnvironment: v.DeclarationEnvironment,
        Body:                   body,
    }
}


type ReturnValue struct {
    Value RuntimeValue
}

func (v ReturnValue) Equals(other RuntimeValue) bool {
    return false
}

type BreakType struct{}

func (v BreakType) Equals(other RuntimeValue) bool {
    return false
}

func (v BreakType) Type() ValueType {
    return Break
}

func (v BreakType) Get() any {
    return nil
}

func (v BreakType) ToString() string {
    return "break"
}

func (v BreakType) Clone() RuntimeValue {
    return v
}

type ClassValue struct {
    Name        string
    Constructor *FunctionValue
    Methods     map[string]*FunctionValue
    Fields      map[string]RuntimeValue
}

func (v ClassValue) Equals(other RuntimeValue) bool {
    return false
}

func (v ClassValue) Type() ValueType {
    return Class
}

func (v ClassValue) Get() any {
    return v.Methods
}

func (v ClassValue) ToString() string {
    return "class"
}

func (v ClassValue) Clone() RuntimeValue {
    methods := make(map[string]*FunctionValue)
    for k, v := range v.Methods {
        methods[k] = v.Clone().(*FunctionValue)
    }
    fields := make(map[string]RuntimeValue)
    for k, v := range v.Fields {
        fields[k] = v.Clone()
    }

    constructor := v.Constructor.Clone().(*FunctionValue)
    return ClassValue{Name: v.Name, Constructor: constructor, Methods: methods, Fields: fields}
}

func MakeClassValue(name string, methods map[string]*FunctionValue) ClassValue {
    return ClassValue{Name: name, Methods: methods}
}

func MakeArrayValue(values []RuntimeValue) ArrayValue {
    return ArrayValue{Values: values}
}

func MakeTupleValue(values []RuntimeValue) TupleValue {
    return TupleValue{Values: values}
}

func MakeObjectValue(properties map[string]RuntimeValue) ObjectValue {
    return ObjectValue{Properties: properties}
}

func ToGoArrayValue(v ArrayValue) []RuntimeValue {
    return v.Values
}

func ToGoNumberValue(v NumberValue) float64 {
    return v.Value
}

func ToGoStringValue(v StringValue) string {
    return v.Value
}

func ToGoTupleValue(v TupleValue) []RuntimeValue {
    return v.Values
}

func ToGoFileValue(v FileValue) *os.File {
    return v.File
}

type FileValue struct {
    Name     string
    File     *os.File
    IsClosed bool
}

func (v FileValue) Equals(other RuntimeValue) bool {
    return false
}

func (v FileValue) Type() ValueType {
    return File
}

func (v FileValue) Get() any {
    return v
}

func (v FileValue) ToString() string {
    return "file"
}

func (v FileValue) Clone() RuntimeValue {
    return v
}

func MakeFileValue(name string, file *os.File) FileValue {
    return FileValue{Name: name, File: file, IsClosed: false}
}
