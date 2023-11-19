package runtimelang

import (
	"os"
	"strconv"

	"github.com/Jamlie/Jamlang/ast"
)

type ValueType string

const (
	I8             ValueType = "i8"
	I16            ValueType = "i16"
	I32            ValueType = "i32"
	I64            ValueType = "i64"
	F32            ValueType = "f32"
	F64            ValueType = "f64"
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
	Continue       ValueType = "continue"
	Class          ValueType = "class"
	File           ValueType = "file"
	Type           ValueType = "type"
)

type RuntimeValue interface {
	Get() any
	Type() ValueType
	ToString() string
	Clone() RuntimeValue
	Equals(RuntimeValue) bool
	VarType() ast.VariableType
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

func (v InitialValue) VarType() ast.VariableType {
	return ast.NullType
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
	return Object
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

func (v NullValue) VarType() ast.VariableType {
	return ast.ObjectType
}

func MakeNullValue() NullValue {
	return NullValue{}
}

type NumberValue[T any] interface {
	GetV() T
}

type IntValue interface {
	GetInt() int
}

type Int8Value struct {
	Value int8
}

func (v Int8Value) Equals(other RuntimeValue) bool {
	if other.Type() == I8 {
		return v.Value == other.Get().(int8)
	}
	return false
}

func (v Int8Value) Type() ValueType {
	return I8
}

func (v Int8Value) Get() any {
	return v.Value
}

func (v Int8Value) ToString() string {
	return strconv.FormatInt(int64(v.Value), 10)
}

func (v Int8Value) Clone() RuntimeValue {
	return v
}

func (v Int8Value) VarType() ast.VariableType {
	return ast.Int8Type
}

func (v Int8Value) GetV() int8 {
	return v.Value
}

func (v Int8Value) GetInt() int {
	return int(v.Value)
}

func MakeInt8Value(value int8) Int8Value {
	return Int8Value{Value: value}
}

type Int16Value struct {
	Value int16
}

func (v Int16Value) Equals(other RuntimeValue) bool {
	if other.Type() == I16 || other.Type() == I8 {
		return v.Value == other.Get().(int16)
	}
	return false
}

func (v Int16Value) Type() ValueType {
	return I16
}

func (v Int16Value) Get() any {
	return v.Value
}

func (v Int16Value) ToString() string {
	return strconv.FormatInt(int64(v.Value), 10)
}

func (v Int16Value) Clone() RuntimeValue {
	return v
}

func (v Int16Value) VarType() ast.VariableType {
	return ast.Int16Type
}

func (v Int16Value) GetV() int16 {
	return v.Value
}

func (v Int16Value) GetInt() int {
	return int(v.Value)
}

func MakeInt16Value(value int16) Int16Value {
	return Int16Value{Value: value}
}

type Int32Value struct {
	Value int32
}

func (v Int32Value) Equals(other RuntimeValue) bool {
	if other.Type() == I32 || other.Type() == I16 || other.Type() == I8 {
		return v.Value == other.Get().(int32)
	}
	return false
}

func (v Int32Value) Type() ValueType {
	return I32
}

func (v Int32Value) Get() any {
	return v.Value
}

func (v Int32Value) ToString() string {
	return strconv.FormatInt(int64(v.Value), 10)
}

func (v Int32Value) Clone() RuntimeValue {
	return v
}

func (v Int32Value) VarType() ast.VariableType {
	return ast.Int32Type
}

func (v Int32Value) GetV() int32 {
	return v.Value
}

func (v Int32Value) GetInt() int {
	return int(v.Value)
}

func MakeInt32Value(value int32) Int32Value {
	return Int32Value{Value: value}
}

type Int64Value struct {
	Value int64
}

func (v Int64Value) Equals(other RuntimeValue) bool {
	if other.Type() == I64 || other.Type() == I32 || other.Type() == I16 || other.Type() == I8 {
		return v.Value == other.Get().(int64)
	}
	return false
}

func (v Int64Value) Type() ValueType {
	return I64
}

func (v Int64Value) Get() any {
	return v.Value
}

func (v Int64Value) ToString() string {
	return strconv.FormatInt(v.Value, 10)
}

func (v Int64Value) Clone() RuntimeValue {
	return v
}

func (v Int64Value) VarType() ast.VariableType {
	return ast.Int64Type
}

func (v Int64Value) GetV() int64 {
	return v.Value
}

func (v Int64Value) GetInt() int {
	return int(v.Value)
}

func MakeInt64Value(value int64) Int64Value {
	return Int64Value{Value: value}
}

type FloatValue interface {
	GetFloat() float64
}

type Float32Value struct {
	Value float32
}

func (v Float32Value) Equals(other RuntimeValue) bool {
	if other.Type() == F32 {
		return v.Value == other.Get().(float32)
	}
	return false
}

func (v Float32Value) Type() ValueType {
	return F32
}

func (v Float32Value) Get() any {
	return v.Value
}

func (v Float32Value) ToString() string {
	return strconv.FormatFloat(float64(v.Value), 'f', -1, 32)
}

func (v Float32Value) Clone() RuntimeValue {
	return v
}

func (v Float32Value) VarType() ast.VariableType {
	return ast.Float32Type
}

func (v Float32Value) GetV() float32 {
	return v.Value
}

func (v Float32Value) GetFloat() float64 {
	return float64(v.Value)
}

func MakeFloat32Value(value float32) Float32Value {
	return Float32Value{Value: value}
}

type Float64Value struct {
	Value float64
}

func (v Float64Value) Equals(other RuntimeValue) bool {
	if other.Type() == F64 || other.Type() == F32 {
		return v.Value == other.Get().(float64)
	}
	return false
}

func (v Float64Value) Type() ValueType {
	return F64
}

func (v Float64Value) Get() any {
	return v.Value
}

func (v Float64Value) ToString() string {
	return strconv.FormatFloat(v.Value, 'f', -1, 64)
}

func (v Float64Value) Clone() RuntimeValue {
	return v
}

func (v Float64Value) VarType() ast.VariableType {
	return ast.Float64Type
}

func (v Float64Value) GetV() float64 {
	return v.Value
}

func (v Float64Value) GetFloat() float64 {
	return v.Value
}

func MakeFloat64Value(value float64) Float64Value {
	return Float64Value{Value: value}
}

// func (v NumberValue) Equals(other RuntimeValue) bool {
//     if other.Type() == Number {
//         return v.Value == other.Get().(float64)
//     }
//     return false
// }
//
// func (v NumberValue) Type() ValueType {
//     return Number
// }
//
// func (v NumberValue) Get() any {
//     return v.Value
// }
//
// func (v NumberValue) ToString() string {
//     return strconv.FormatFloat(v.Value, 'f', -1, 64)
// }
//
// func (v NumberValue) Clone() RuntimeValue {
//     return v
// }
//
// func (v NumberValue) VarType() ast.VariableType {
//     return ast.Float64Type
// }
//
// func MakeNumberValue(value float64) NumberValue {
//     return NumberValue{Value: value}
// }
//
// func MakeNumberValueAsInt16(value int16) NumberValue {
//     return NumberValue{Value: float64(value)}
// }

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
	str := v.Value
	decodedStr, err := strconv.Unquote(`"` + str + `"`)
	if err != nil {
		return str
	}
	return decodedStr
}

func (v StringValue) ToString() string {
	return v.Value
}

func (v StringValue) Clone() RuntimeValue {
	return v
}

func (v StringValue) VarType() ast.VariableType {
	return ast.StringType
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

func (v BoolValue) VarType() ast.VariableType {
	return ast.BoolType
}

func MakeBoolValue(value bool) BoolValue {
	return BoolValue{Value: value}
}

type ObjectValue struct {
	Properties map[string]RuntimeValue
	IsClass    bool
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
		case I8, I16, I32, I64, F32, F64:
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
		case NativeFunction:
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

func (v ObjectValue) VarType() ast.VariableType {
	return ast.ObjectType
}

type ArrayValue struct {
	Values []RuntimeValue
}

func (v *ArrayValue) Push(value RuntimeValue) {
	v.Values = append(v.Values, value)
}

func (v *ArrayValue) Pop() RuntimeValue {
	if len(v.Values) > 0 {
		value := v.Values[len(v.Values)-1]
		v.Values = v.Values[:len(v.Values)-1]
		return value
	}
	return MakeNullValue()
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

func (v ArrayValue) VarType() ast.VariableType {
	return ast.ArrayType
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

func (v TupleValue) VarType() ast.VariableType {
	return ast.TupleType
}

type FunctionCall func(args []RuntimeValue, env Environment) RuntimeValue

type NativeFunctionValue struct {
	Call FunctionCall
	Name string
}

func (v NativeFunctionValue) Equals(other RuntimeValue) bool {
	return false
}

func (v NativeFunctionValue) Type() ValueType {
	return NativeFunction
}

func (v NativeFunctionValue) Get() any {
	return "fn " + v.Name + "(...)" + " { [native code] }"
}

func (v NativeFunctionValue) ToString() string {
	return "native function"
}

func (v NativeFunctionValue) Clone() RuntimeValue {
	return v
}

func (v NativeFunctionValue) VarType() ast.VariableType {
	return ast.AnyType
}

func MakeNativeFunction(call FunctionCall, name string) NativeFunctionValue {
	return NativeFunctionValue{Call: call, Name: name}
}

type FunctionValue struct {
	Name                   string
	Parameters             []string
	DeclarationEnvironment Environment
	Body                   []ast.Statement
	IsAnonymous            bool
	ReturnType             ast.VariableType
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

func (v FunctionValue) VarType() ast.VariableType {
	return ast.FunctionType
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

func (v BreakType) VarType() ast.VariableType {
	return ast.AnyType
}

type ContinueType struct{}

func (v ContinueType) Equals(other RuntimeValue) bool {
	return false
}

func (v ContinueType) Type() ValueType {
	return Continue
}

func (v ContinueType) Get() any {
	return nil
}

func (v ContinueType) ToString() string {
	return "continue"
}

func (v ContinueType) Clone() RuntimeValue {
	return v
}

func (v ContinueType) VarType() ast.VariableType {
	return ast.AnyType
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

func (v ClassValue) VarType() ast.VariableType {
	return ast.AnyType
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

func ToGoNumberValue[T any](v NumberValue[T]) T {
	return v.GetV()
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

func (v FileValue) VarType() ast.VariableType {
	return ast.AnyType
}

func MakeFileValue(name string, file *os.File) FileValue {
	return FileValue{Name: name, File: file, IsClosed: false}
}

type TypeValue struct {
	Name  string
	Value ast.Expression
}

func (v TypeValue) Equals(other RuntimeValue) bool {
	return false
}

func (v TypeValue) Type() ValueType {
	return Type
}

func (v TypeValue) Get() any {
	return v.Value
}

func (v TypeValue) ToString() string {
	return "type"
}

func (v TypeValue) Clone() RuntimeValue {
	return v
}

func (v TypeValue) VarType() ast.VariableType {
	return ast.AnyType
}

func MakeTypeValue(name string, value ast.Expression) TypeValue {
	return TypeValue{Name: name, Value: value}
}
