package runtimelang

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Jamlie/Jamlang/parser"
)

func jamlangPrintln(args []RuntimeValue, environment Environment) RuntimeValue {
	for _, arg := range args {
		fmt.Print(arg.Get())
	}
	fmt.Println()
	return MakeNullValue()
}

func jamlangPrint(args []RuntimeValue, environment Environment) RuntimeValue {
	for _, arg := range args {
		fmt.Print(arg.Get())
	}
	return MakeNullValue()
}

func jamlangSleep(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: sleep takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 {
		fmt.Fprintln(os.Stderr, "Error: sleep takes a number - time in milliseconds")
		os.Exit(0)
	}

	time.Sleep(time.Duration(args[0].(IntValue).GetInt()) * time.Millisecond)
	return MakeNullValue()
}

func jamlangTypeof(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: typeof takes 1 argument")
		os.Exit(0)
	}

	return MakeStringValue(string(args[0].Type()))
}

func jamlangExit(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: exit takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 {
		fmt.Fprintln(os.Stderr, "Error: exit takes a number - exit code")
		os.Exit(0)
	}

	os.Exit(args[0].(IntValue).GetInt())
	return MakeNullValue()
}

func jamlangInput(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: input takes 1 argument")
		os.Exit(0)
	}

	fmt.Print(args[0].Get())
	scanner := bufio.NewReader(os.Stdin)
	input, err := scanner.ReadString('\n')
	input = strings.Trim(input, "\n")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: reading input")
		os.Exit(0)
	}
	return MakeStringValue(input)
}

// func jamlangLen(args []RuntimeValue, environment Environment) RuntimeValue {
//     if len(args) != 1 {
//         fmt.Fprintln(os.Stderr, "len takes 1 argument")
//         os.Exit(0)
//     }
//
//     if args[0].Type() == Array {
//         goArray := ToGoArrayValue(args[0].(ArrayValue))
//         return MakeInt64Value(int64(len(goArray)))
//     } else if args[0].Type() == Tuple {
//         goTuple := ToGoTupleValue(args[0].(TupleValue))
//         return MakeInt(float64(len(goTuple)))
//     } else if args[0].Type() == String {
//         goString := ToGoStringValue(args[0].(StringValue))
//         return MakeNumberValue(float64(len(goString)))
//     } else {
//         fmt.Println("len takes an array, tuple or string")
//         os.Exit(0)
//         return MakeNullValue()
//     }
// }

// func jamlangAppend(args []RuntimeValue, environment Environment) RuntimeValue {
//     if len(args) != 2 {
//         fmt.Fprintln(os.Stderr, "append takes 2 arguments")
//         os.Exit(0)
//     }
//
//     if args[0].Type() != Array {
//         fmt.Fprintln(os.Stderr, "append takes an array")
//         os.Exit(0)
//     }
//
//     goArray := ToGoArrayValue(args[0].(ArrayValue))
//     goArray = append(goArray, args[1])
//     return MakeArrayValue(goArray)
// }

// func jamlangPop(args []RuntimeValue, environment Environment) RuntimeValue {
//     if len(args) != 1 {
//         fmt.Fprintln(os.Stderr, "pop takes 1 argument")
//         os.Exit(0)
//     }
//
//     if args[0].Type() != Array {
//         fmt.Fprintln(os.Stderr, "pop takes an array")
//         os.Exit(0)
//     }
//
//     goArray := ToGoArrayValue(args[0].(ArrayValue))
//     if len(goArray) == 0 {
//         fmt.Fprintln(os.Stderr, "pop takes a non-empty array")
//         os.Exit(0)
//     }
//
//     return MakeArrayValue(goArray[:len(goArray)-1])
// }

func jamlangCopy(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: copy takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == Array {
		goArray := ToGoArrayValue(args[0].(ArrayValue))
		goArrayCopy := make([]RuntimeValue, len(goArray))
		copy(goArrayCopy, goArray)
		return MakeArrayValue(goArrayCopy)
	} else if args[0].Type() == Tuple {
		goTuple := ToGoTupleValue(args[0].(TupleValue))
		goTupleCopy := make([]RuntimeValue, len(goTuple))
		copy(goTupleCopy, goTuple)
		return MakeTupleValue(goTupleCopy)
	} else {
		fmt.Fprintln(os.Stderr, "Error: copy takes an array or tuple")
		os.Exit(0)
		return MakeNullValue()
	}
}

func jamlangArray(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: array takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == Tuple {
		goTuple := ToGoTupleValue(args[0].(TupleValue))
		return MakeArrayValue(goTuple)
	}

	if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 {
		fmt.Fprintln(os.Stderr, "Error: array takes a number")
		os.Exit(0)
	}

	var size int
	switch args[0].Type() {
	case I8:
		size = int(ToGoNumberValue(args[0].(Int8Value)))
	case I16:
		size = int(ToGoNumberValue(args[0].(Int16Value)))
	case I32:
		size = int(ToGoNumberValue(args[0].(Int32Value)))
	default:
		fmt.Fprintln(os.Stderr, "Error: array takes a number")
		os.Exit(0)
	}

	if size < 0 {
		fmt.Fprintln(os.Stderr, "Error: array takes a positive number")
		os.Exit(0)
	}

	goArray := make([]RuntimeValue, size)
	return MakeArrayValue(goArray)
}

func jamlangTuple(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: tuple takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == Array {
		goArray := ToGoArrayValue(args[0].(ArrayValue))
		return MakeTupleValue(goArray)
	}

	if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 {
		fmt.Fprintln(os.Stderr, "Error: tuple takes a number")
		os.Exit(0)
	}

	var size int
	switch args[0].Type() {
	case I8:
		size = int(ToGoNumberValue(args[0].(Int8Value)))
	case I16:
		size = int(ToGoNumberValue(args[0].(Int16Value)))
	case I32:
		size = int(ToGoNumberValue(args[0].(Int32Value)))
	default:
		fmt.Fprintln(os.Stderr, "Error: tuple takes a number")
		os.Exit(0)
	}

	if size < 0 {
		fmt.Fprintln(os.Stderr, "Error: tuple takes a positive number")
		os.Exit(0)
	}

	goArray := make([]RuntimeValue, size)
	return MakeTupleValue(goArray)
}

func jamlangToString(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: string takes 1 argument")
		os.Exit(0)
	}

	return MakeStringValue(args[0].ToString())
}

func jamlangToUint32(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: uint32 takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int8Value)))))
	}

	if args[0].Type() == I16 {
		return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int16Value)))))
	}

	if args[0].Type() == I32 {
		return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int32Value)))))
	}

	if args[0].Type() == I64 {
		return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Int64Value)))))
	}

	if args[0].Type() == F32 {
		return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Float32Value)))))
	}

	if args[0].Type() == F64 {
		return MakeInt64Value(int64(uint32(ToGoNumberValue(args[0].(Float64Value)))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: uint32 takes a string or a number")
		os.Exit(0)
	}

	uintString := args[0].ToString()

	uintUint, err := strconv.ParseUint(uintString, 10, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: uint32 takes a string or a number")
		os.Exit(0)
	}

	return MakeInt64Value(int64(uintUint))
}

func jamlangToUint64(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: uint64 takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int8Value)))))
	}

	if args[0].Type() == I16 {
		return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int16Value)))))
	}

	if args[0].Type() == I32 {
		return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int32Value)))))
	}

	if args[0].Type() == I64 {
		return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Int64Value)))))
	}

	if args[0].Type() == F32 {
		return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Float32Value)))))
	}

	if args[0].Type() == F64 {
		return MakeFloat32Value(float32(uint64(ToGoNumberValue(args[0].(Float64Value)))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: uint64 takes a string or a number")
		os.Exit(0)
	}

	uintString := args[0].ToString()

	uintUint, err := strconv.ParseUint(uintString, 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: uint64 takes a string or a number")
		os.Exit(0)
	}

	return MakeFloat32Value(float32(uintUint))
}

func jamlangToInt8(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: int8 takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value))))
	}

	if args[0].Type() == I16 {
		return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int16Value))))
	}

	if args[0].Type() == I32 {
		return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int32Value))))
	}

	if args[0].Type() == I64 {
		return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int64Value))))
	}

	if args[0].Type() == F32 {
		return MakeInt8Value(int8(ToGoNumberValue(args[0].(Float32Value))))
	}

	if args[0].Type() == F64 {
		return MakeInt8Value(int8(ToGoNumberValue(args[0].(Float64Value))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: int8 takes a string or a number")
		os.Exit(0)
	}

	getInt := args[0].Get()
	intString := getInt.(string)
	intString = strings.Trim(intString, "\r")

	intInt, err := strconv.ParseInt(intString, 10, 8)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: int8 takes a string or a number")
		os.Exit(0)
	}

	return MakeInt8Value(int8(intInt))
}

func jamlangToInt16(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: int16 takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int8Value))))
	}

	if args[0].Type() == I16 {
		return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value))))
	}

	if args[0].Type() == I32 {
		return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int32Value))))
	}

	if args[0].Type() == I64 {
		return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int64Value))))
	}

	if args[0].Type() == F32 {
		return MakeInt16Value(int16(ToGoNumberValue(args[0].(Float32Value))))
	}

	if args[0].Type() == F64 {
		return MakeInt16Value(int16(ToGoNumberValue(args[0].(Float64Value))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: int16 takes a string or a number")
		os.Exit(0)
	}

	getInt := args[0].Get()
	intString := getInt.(string)
	intString = strings.Trim(intString, "\r")

	intInt, err := strconv.ParseInt(intString, 10, 16)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: int16 takes a string or a number")
		os.Exit(0)
	}

	return MakeInt16Value(int16(intInt))
}

func jamlangToInt32(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: int32 takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int8Value))))
	}

	if args[0].Type() == I16 {
		return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int16Value))))
	}

	if args[0].Type() == I32 {
		return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value))))
	}

	if args[0].Type() == I64 {
		return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int64Value))))
	}

	if args[0].Type() == F32 {
		return MakeInt32Value(int32(ToGoNumberValue(args[0].(Float32Value))))
	}

	if args[0].Type() == F64 {
		return MakeInt32Value(int32(ToGoNumberValue(args[0].(Float64Value))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: int32 takes a string or a number")
		os.Exit(0)
	}

	getInt := args[0].Get()
	intString := getInt.(string)
	intString = strings.Trim(intString, "\r")

	intInt, err := strconv.ParseInt(intString, 10, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: int32 takes a string or a number")
		os.Exit(0)
	}

	return MakeInt32Value(int32(intInt))
}

func jamlangToInt64(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: int64 takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value))))
	}

	if args[0].Type() == I16 {
		return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value))))
	}

	if args[0].Type() == I32 {
		return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value))))
	}

	if args[0].Type() == I64 {
		return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value))))
	}

	if args[0].Type() == F32 {
		return MakeInt64Value(int64(ToGoNumberValue(args[0].(Float32Value))))
	}

	if args[0].Type() == F64 {
		return MakeInt64Value(int64(ToGoNumberValue(args[0].(Float64Value))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: int64 takes a string or a number")
		os.Exit(0)
	}

	getInt := args[0].Get()
	intString := getInt.(string)
	intString = strings.Trim(intString, "\r")

	intInt, err := strconv.ParseInt(intString, 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: int64 takes a string or a number")
		os.Exit(0)
	}

	return MakeInt64Value(int64(intInt))
}

func jamlangToFloat32(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: float32 takes a string or a numebr")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int8Value))))
	}

	if args[0].Type() == I16 {
		return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int16Value))))
	}

	if args[0].Type() == I32 {
		return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int32Value))))
	}

	if args[0].Type() == I64 {
		return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Int64Value))))
	}

	if args[0].Type() == F32 {
		return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Float32Value))))
	}

	if args[0].Type() == F64 {
		return MakeFloat32Value(float32(ToGoNumberValue(args[0].(Float64Value))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: float32 takes a string or a numebr")
		os.Exit(0)
	}

	getFloat := args[0].Get()
	floatString := getFloat.(string)
	floatString = strings.Trim(floatString, "\r")

	floatFloat, err := strconv.ParseFloat(floatString, 32)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: float32 takes a string or a numebr")
		os.Exit(0)
	}

	return MakeFloat32Value(float32(floatFloat))
}

func jamlangToFloat64(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: float64 takes a string or a numebr")
		os.Exit(0)
	}

	if args[0].Type() == I8 {
		return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int8Value))))
	}

	if args[0].Type() == I16 {
		return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int16Value))))
	}

	if args[0].Type() == I32 {
		return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int32Value))))
	}

	if args[0].Type() == I64 {
		return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Int64Value))))
	}

	if args[0].Type() == F32 {
		return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Float32Value))))
	}

	if args[0].Type() == F64 {
		return MakeFloat64Value(float64(ToGoNumberValue(args[0].(Float64Value))))
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: float64 takes a string or a numebr")
		os.Exit(0)
	}

	getFloat := args[0].Get()
	floatString := getFloat.(string)
	floatString = strings.Trim(floatString, "\r")

	floatFloat, err := strconv.ParseFloat(floatString, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: float64 takes a string or a numebr")
		os.Exit(0)
	}

	return MakeFloat64Value(float64(floatFloat))
}

func jamlangHex(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: hex takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: hex takes a string")
		os.Exit(0)
	}

	hexString := args[0].ToString()

	hexInt, err := strconv.ParseInt(hexString, 16, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: hex takes a string")
		os.Exit(0)
	}

	return MakeInt64Value(int64(hexInt))
}

func jamlangCurrentTime(args []RuntimeValue, environment Environment) RuntimeValue {
	return MakeInt64Value(time.Now().UnixMicro())
}

func jamlangBitwiseNot(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.NOT takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64 {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.NOT takes a number")
		os.Exit(0)
	}

	switch args[0].Type() {
	case I8:
		return MakeInt8Value(^int8(ToGoNumberValue(args[0].(Int8Value))))
	case I16:
		return MakeInt16Value(^int16(ToGoNumberValue(args[0].(Int16Value))))
	case I32:
		return MakeInt32Value(^int32(ToGoNumberValue(args[0].(Int32Value))))
	case I64:
		return MakeInt64Value(^int64(ToGoNumberValue(args[0].(Int64Value))))
	case F32:
		return MakeInt64Value(^int64(ToGoNumberValue(args[0].(Float32Value))))
	case F64:
		return MakeInt64Value(^int64(ToGoNumberValue(args[0].(Float64Value))))
	default:
		fmt.Fprintln(os.Stderr, "Error: Bitwise.NOT takes a number")
		os.Exit(0)
		return nil
	}
}

func jamlangBitwiseAnd(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 arguments")
		os.Exit(0)
	}

	if (args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64) || (args[1].Type() != I8 && args[1].Type() != I16 && args[1].Type() != I32 && args[1].Type() != I64 && args[1].Type() != F32 && args[1].Type() != F64) {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.AND and takes 2 numbers")
		os.Exit(0)
	}

	switch args[0].Type() {
	case I8:
		switch args[1].Type() {
		case I8:
			return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value)) & ToGoNumberValue(args[1].(Int8Value))))
		case I16:
			return MakeInt16Value(int16(int16(ToGoNumberValue(args[0].(Int8Value))) & ToGoNumberValue(args[1].(Int16Value))))
		case I32:
			return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int8Value))) & ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int8Value))) & ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) & int8(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) & int8(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I16:
		switch args[1].Type() {
		case I8:
			return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) & int16(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) & ToGoNumberValue(args[1].(Int16Value))))
		case I32:
			return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int16Value))) & ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int16Value))) & ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) & int16(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) & int16(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I32:
		switch args[1].Type() {
		case I8:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) & ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int32Value))) & ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) & int32(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I64:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) & int64(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case F32:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float32Value))) & int8(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float32Value))) & int16(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float32Value))) & int32(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) & int64(ToGoNumberValue(args[1].(Int64Value)))))
		case F32:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) & int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) & int64(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case F64:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float64Value))) & int8(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float64Value))) & int16(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float64Value))) & int32(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) & int64(ToGoNumberValue(args[1].(Int64Value)))))
		case F32:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) & int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value)))) & int64(ToGoNumberValue(args[1].(Float64Value))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
			os.Exit(0)
			return nil
		}
	default:
		fmt.Fprintln(os.Stderr, "Error: Bitwise.AND takes 2 numbers")
		os.Exit(0)
		return nil
	}
}

func jamlangBitwiseOr(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 arguments")
		os.Exit(0)
	}

	if (args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64) || (args[1].Type() != I8 && args[1].Type() != I16 && args[1].Type() != I32 && args[1].Type() != I64 && args[1].Type() != F32 && args[1].Type() != F64) {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
		os.Exit(0)
	}

	switch args[0].Type() {
	case I8:
		switch args[1].Type() {
		case I8:
			return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value)) | ToGoNumberValue(args[1].(Int8Value))))
		case I16:
			return MakeInt16Value(int16(int16(ToGoNumberValue(args[0].(Int8Value))) | ToGoNumberValue(args[1].(Int16Value))))
		case I32:
			return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int8Value))) | ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int8Value))) | ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) | int8(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) | int8(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I16:
		switch args[1].Type() {
		case I8:
			return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) | int16(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) | ToGoNumberValue(args[1].(Int16Value))))
		case I32:
			return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int16Value))) | ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int16Value))) | ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) | int16(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) | int16(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I32:
		switch args[1].Type() {
		case I8:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) | ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int32Value))) | ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) | int32(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I64:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) | int64(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case F32:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float32Value))) | int8(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float32Value))) | int16(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float32Value))) | int32(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) | int64(ToGoNumberValue(args[1].(Int64Value)))))
		case F32:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) | int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) | int64(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case F64:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float64Value))) | int8(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float64Value))) | int16(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float64Value))) | int32(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) | int64(ToGoNumberValue(args[1].(Int64Value)))))
		case F32:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) | int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value)))) | int64(ToGoNumberValue(args[1].(Float64Value))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	default:
		fmt.Fprintln(os.Stderr, "Error: Bitwise.OR takes 2 numbers")
		os.Exit(0)
		return nil
	}
}

func jamlangBitwiseXor(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 arguments")
		os.Exit(0)
	}

	if (args[0].Type() != I8 && args[0].Type() != I16 && args[0].Type() != I32 && args[0].Type() != I64 && args[0].Type() != F32 && args[0].Type() != F64) || (args[1].Type() != I8 && args[1].Type() != I16 && args[1].Type() != I32 && args[1].Type() != I64 && args[1].Type() != F32 && args[1].Type() != F64) {
		fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
		os.Exit(0)
	}

	switch args[0].Type() {
	case I8:
		switch args[1].Type() {
		case I8:
			return MakeInt8Value(int8(ToGoNumberValue(args[0].(Int8Value)) ^ ToGoNumberValue(args[1].(Int8Value))))
		case I16:
			return MakeInt16Value(int16(int16(ToGoNumberValue(args[0].(Int8Value))) ^ ToGoNumberValue(args[1].(Int16Value))))
		case I32:
			return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int8Value))) ^ ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int8Value))) ^ ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) ^ int8(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int8Value)) ^ int8(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I16:
		switch args[1].Type() {
		case I8:
			return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) ^ int16(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt16Value(int16(ToGoNumberValue(args[0].(Int16Value)) ^ ToGoNumberValue(args[1].(Int16Value))))
		case I32:
			return MakeInt32Value(int32(int32(ToGoNumberValue(args[0].(Int16Value))) ^ ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int16Value))) ^ ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) ^ int16(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int16Value)) ^ int16(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I32:
		switch args[1].Type() {
		case I8:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt32Value(int32(ToGoNumberValue(args[0].(Int32Value)) ^ ToGoNumberValue(args[1].(Int32Value))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Int32Value))) ^ ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int32Value)) ^ int32(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case I64:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ ToGoNumberValue(args[1].(Int64Value))))
		case F32:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(ToGoNumberValue(args[0].(Int64Value)) ^ int64(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case F32:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float32Value))) ^ int8(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float32Value))) ^ int16(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float32Value))) ^ int32(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) ^ int64(ToGoNumberValue(args[1].(Int64Value)))))
		case F32:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) ^ int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float32Value))) ^ int64(ToGoNumberValue(args[1].(Float64Value)))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	case F64:
		switch args[1].Type() {
		case I8:
			return MakeInt64Value(int64(int8(ToGoNumberValue(args[0].(Float64Value))) ^ int8(ToGoNumberValue(args[1].(Int8Value)))))
		case I16:
			return MakeInt64Value(int64(int16(ToGoNumberValue(args[0].(Float64Value))) ^ int16(ToGoNumberValue(args[1].(Int16Value)))))
		case I32:
			return MakeInt64Value(int64(int32(ToGoNumberValue(args[0].(Float64Value))) ^ int32(ToGoNumberValue(args[1].(Int32Value)))))
		case I64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) ^ int64(ToGoNumberValue(args[1].(Int64Value)))))
		case F32:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value))) ^ int64(ToGoNumberValue(args[1].(Float32Value)))))
		case F64:
			return MakeInt64Value(int64(int64(ToGoNumberValue(args[0].(Float64Value)))) ^ int64(ToGoNumberValue(args[1].(Float64Value))))
		default:
			fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
			os.Exit(0)
			return nil
		}
	default:
		fmt.Fprintln(os.Stderr, "Error: Bitwise.XOR takes 2 numbers")
		os.Exit(0)
		return nil
	}
}

func jamlangEval(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: eval takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "eval takes a string")
		os.Exit(0)
	}

	code := args[0].ToString()
	program := parser.NewParser().ProduceAST(code)
	newEnvironment := CreateGlobalEnvironment()
	Evaluate(&program, *newEnvironment)
	return MakeNullValue()
}

func jamlangOpen(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: open takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: open takes a string")
		os.Exit(0)
	}

	filename := args[0].(StringValue).Value

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't open file")
		os.Exit(0)
	}

	var properties = make(map[string]RuntimeValue)
	properties["close"] = jamlangClose(file)
	properties["name"] = MakeStringValue(file.Name())
	properties["file"] = MakeFileValue(file.Name(), file)
	properties["append"] = jamlangAppend(file)
	properties["read"] = jamlangRead(&file)

	return MakeObjectValue(properties)
}

func jamlangClose(file *os.File) RuntimeValue {
	return MakeNativeFunction(func(args []RuntimeValue, environment Environment) RuntimeValue {
		if len(args) != 0 {
			fmt.Fprintln(os.Stderr, "Error: close doesn't take any argument")
			os.Exit(0)
		}

		err := file.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: couldn't close file")
			os.Exit(0)
		}
		return MakeNullValue()
	}, "close")
}

func jamlangAppend(file *os.File) RuntimeValue {
	return MakeNativeFunction(func(args []RuntimeValue, environment Environment) RuntimeValue {
		if len(args) != 1 {
			fmt.Fprintln(os.Stderr, "Error: write takes 1 argument")
			os.Exit(0)
		}

		if args[0].Type() != String {
			fmt.Fprintln(os.Stderr, "Error: write takes a string")
			os.Exit(0)
		}

		_, err := file.WriteString(args[0].(StringValue).Get().(string))
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: couldn't write to file")
			os.Exit(0)
		}
		return MakeNullValue()
	}, "append")
}

func jamlangRead(file **os.File) RuntimeValue {
	return MakeNativeFunction(func(args []RuntimeValue, environment Environment) RuntimeValue {
		if len(args) != 0 {
			fmt.Fprintln(os.Stderr, "Error: read doesn't take any argument")
			os.Exit(0)
		}

		reader := bufio.NewReader(*file)
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error: couldn't read from file")
			os.Exit(0)
		}
		return MakeStringValue(line)
	}, "read")
}

func jamlangObjectKeys(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Object.keys takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != "object" {
		fmt.Fprintln(os.Stderr, "Object.keys takes an object")
		os.Exit(0)
	}

	keys := make([]RuntimeValue, 0)
	for key := range args[0].(ObjectValue).Properties {
		keys = append(keys, MakeStringValue(key))
	}

	return MakeArrayValue(keys)
}

func jamlangObjectValues(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Object.values takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != "object" {
		fmt.Fprintln(os.Stderr, "Object.values takes an object")
		os.Exit(0)
	}

	values := make([]RuntimeValue, 0)
	for _, value := range args[0].(ObjectValue).Properties {
		values = append(values, value)
	}

	return MakeArrayValue(values)
}

func jamlangObjectHas(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Object.has takes 2 arguments")
		os.Exit(0)
	}

	if args[0].Type() != "object" {
		fmt.Fprintln(os.Stderr, "Object.has takes an object")
		os.Exit(0)
	}

	if args[1].Type() != String {
		fmt.Fprintln(os.Stderr, "Object.has takes a string")
		os.Exit(0)
	}

	_, ok := args[0].(ObjectValue).Properties[args[1].(StringValue).Value]
	return MakeBoolValue(ok)
}

func jamlangHttpGet(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: http.get takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: http.get takes a string")
		os.Exit(0)
	}

	resp, err := http.Get(args[0].(StringValue).Value)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't get url")
		os.Exit(0)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't read response")
		os.Exit(0)
	}

	return MakeStringValue(string(body))
}

func jamlangHttpPost(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 2 {
		fmt.Fprintln(os.Stderr, "Error: http.post takes 2 arguments")
		os.Exit(0)
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: http.post takes a string")
		os.Exit(0)
	}

	if args[1].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: http.post takes a string")
		os.Exit(0)
	}

	resp, err := http.Post(args[0].(StringValue).Value, "application/json", strings.NewReader(args[1].(StringValue).Get().(string)))
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't post url")
		os.Exit(0)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't read response")
		os.Exit(0)
	}

	return MakeStringValue(string(body))
}

func jamlangHttpListen(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: http.listen takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: http.listen takes a string")
		os.Exit(0)
	}

	err := http.ListenAndServe(args[0].(StringValue).Value, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't listen on port")
		os.Exit(0)
	}

	return MakeNullValue()
}

func jamlangHttpNew(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 0 {
		fmt.Fprintln(os.Stderr, "Error: http.new takes 0 arguments")
		os.Exit(0)
	}

	httpObject := make(map[string]RuntimeValue)
	httpObject["listen"] = MakeNativeFunction(jamlangHttpListen, "listen")
	httpObject["get"] = MakeNativeFunction(jamlangHttpGet, "get")
	httpObject["post"] = MakeNativeFunction(jamlangHttpPost, "post")

	return MakeObjectValue(httpObject)
}

func jamlangJsonParse(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: json.parse takes 1 argument")
		os.Exit(0)
	}

	if args[0].Type() != String {
		fmt.Fprintln(os.Stderr, "Error: json.parse takes a string")
		os.Exit(0)
	}

	var data any
	err := json.Unmarshal([]byte(args[0].(StringValue).Value), &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't parse json")
		os.Exit(0)
	}

	return MakeJSONValue(data)
}

func jamlangJsonStringify(args []RuntimeValue, environment Environment) RuntimeValue {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Error: json.stringify takes 1 argument")
		os.Exit(0)
	}

	var data []byte
	var err error
	switch args[0].Type() {
	case String:
		data, err = json.Marshal(args[0].(StringValue).Value)
	case I8:
		data, err = json.Marshal(args[0].(Int8Value).Value)
	case I16:
		data, err = json.Marshal(args[0].(Int16Value).Value)
	case I32:
		data, err = json.Marshal(args[0].(Int32Value).Value)
	case I64:
		data, err = json.Marshal(args[0].(Int64Value).Value)
	case F32:
 		data, err = json.Marshal(args[0].(Float32Value).Value)
	case F64:
		data, err = json.Marshal(args[0].(Float64Value).Value)
	case Bool:
		data, err = json.Marshal(args[0].(BoolValue).Value)
	case Null:
		data, err = json.Marshal(args[0].(NullValue).Value)
	case Array:
		data, err = json.Marshal(args[0].(ArrayValue).Values)
	case Object:
		data, err = json.Marshal(args[0].(ObjectValue).Properties)
	default:
		fmt.Fprintln(os.Stderr, "Error: json.stringify takes a string, number, boolean, null, array or object")
		os.Exit(0)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: couldn't stringify json")
		os.Exit(0)
	}

	return MakeStringValue(string(data))
}
