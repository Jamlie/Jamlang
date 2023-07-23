package runtimelang

type ValueType string

const (
    Number ValueType = "number"
    Null ValueType = "null"
)

type RuntimeValue interface {
    Get() any
    Type() ValueType
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

type NumberValue struct {
    Value float64
}

func (v NumberValue) Type() ValueType {
    return Number
}

func (v NumberValue) Get() any {
    return v.Value
}
