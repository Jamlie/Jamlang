package runtimelang

type ValueType string

const (
    Number ValueType = "number"
    Null ValueType = "null"
    String ValueType = "string"
    Bool ValueType = "bool"
)

type RuntimeValue interface {
    Get() any
    Type() ValueType
}

type InitialValue struct {}

func (v InitialValue) Type() ValueType {
    return Null
}

func (v InitialValue) Get() any {
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

func MakeBoolValue(value bool) BoolValue {
    return BoolValue{Value: value}
}
