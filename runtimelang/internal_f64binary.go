package runtimelang

import (
	"fmt"
	"math"
	"os"
)

func internal_F64Plus(lhs Float64Value, rhs RuntimeValue) RuntimeValue {

	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot add floating point values to integer values")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		var result float64 = 0
		result = lhs.Value + float64(rhs.(Float32Value).Value)
		return Float64Value{result}
	case F64:
		return Float64Value{lhs.Value + rhs.(Float64Value).Value}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64Minus(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot subtract floating point values from integer values")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		var result float64 = 0
		result = lhs.Value - float64(rhs.(Float32Value).Value)
		return Float64Value{result}
	case F64:
		return Float64Value{lhs.Value - rhs.(Float64Value).Value}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64Mult(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot multiply floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		return Float32Value{float32(lhs.Value) * rhs.(Float32Value).Value}
	case F64:
		return Float64Value{lhs.Value * rhs.(Float64Value).Value}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64Pow(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot raise floating point values to integer powers")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		var result float64 = 0
		result = math.Pow(float64(lhs.Value), float64(rhs.(Float32Value).Value))
		return Float64Value{result}
	case F64:
		return Float64Value{math.Pow(lhs.Value, rhs.(Float64Value).Value)}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64Div(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = float64(lhs.Value) / float64(rhs.(Int8Value).Value)
		return Float64Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = float64(lhs.Value) / float64(rhs.(Int16Value).Value)
		return Float64Value{result}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = float64(lhs.Value) / float64(rhs.(Int32Value).Value)
		return Float64Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = float64(lhs.Value) / float64(rhs.(Int64Value).Value)
		return Float64Value{result}
	case F32:
		if rhs.(Float32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = float64(lhs.Value) / float64(rhs.(Float32Value).Value)
		return Float64Value{result}
	case F64:
		if rhs.(Float64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float64Value{lhs.Value / rhs.(Float64Value).Value}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64IntDiv(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Floor(float64(lhs.Value) / float64(rhs.(Int8Value).Value))
		return Float64Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Floor(float64(lhs.Value) / float64(rhs.(Int16Value).Value))
		return Float64Value{result}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Floor(float64(lhs.Value) / float64(rhs.(Int32Value).Value))
		return Float64Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Floor(float64(lhs.Value) / float64(rhs.(Int64Value).Value))
		return Float64Value{result}
	case F32:
		if rhs.(Float32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Floor(float64(lhs.Value) / float64(rhs.(Float32Value).Value))
		return Float64Value{result}
	case F64:
		if rhs.(Float64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float64Value{math.Floor(lhs.Value / rhs.(Float64Value).Value)}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64Mod(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Mod(float64(lhs.Value), float64(rhs.(Int8Value).Value))
		return Float64Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Mod(float64(lhs.Value), float64(rhs.(Int16Value).Value))
		return Float64Value{result}

	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Mod(float64(lhs.Value), float64(rhs.(Int32Value).Value))
		return Float64Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Mod(float64(lhs.Value), float64(rhs.(Int64Value).Value))
		return Float64Value{result}
	case F32:
		if rhs.(Float32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Mod(float64(lhs.Value), float64(rhs.(Float32Value).Value))
		return Float64Value{result}
	case F64:
		if rhs.(Float64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float64Value{math.Mod(lhs.Value, rhs.(Float64Value).Value)}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64GreaterThan(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isGreaterThan := lhs.Value > float64(rhs.(Int8Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isGreaterThan := lhs.Value > float64(rhs.(Int16Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isGreaterThan := lhs.Value > float64(rhs.(Int32Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isGreaterThan := lhs.Value > float64(rhs.(Int64Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isGreaterThan := lhs.Value > float64(rhs.(Float32Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isGreaterThan := lhs.Value > rhs.(Float64Value).Value
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64GreaterThanEqual(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isGreaterThanEqual := lhs.Value >= float64(rhs.(Int8Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isGreaterThanEqual := lhs.Value >= float64(rhs.(Int16Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isGreaterThanEqual := lhs.Value >= float64(rhs.(Int32Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isGreaterThanEqual := lhs.Value >= float64(rhs.(Int64Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isGreaterThanEqual := lhs.Value >= float64(rhs.(Float32Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isGreaterThanEqual := lhs.Value >= rhs.(Float64Value).Value
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64LessThan(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isLessThan := lhs.Value < float64(rhs.(Int8Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isLessThan := lhs.Value < float64(rhs.(Int16Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isLessThan := lhs.Value < float64(rhs.(Int32Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isLessThan := lhs.Value < float64(rhs.(Int64Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isLessThan := lhs.Value < float64(rhs.(Float32Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isLessThan := lhs.Value < rhs.(Float64Value).Value
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64LessThanEqual(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isLessThanEqual := lhs.Value <= float64(rhs.(Int8Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isLessThanEqual := lhs.Value <= float64(rhs.(Int16Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isLessThanEqual := lhs.Value <= float64(rhs.(Int32Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isLessThanEqual := lhs.Value <= float64(rhs.(Int64Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isLessThanEqual := lhs.Value <= float64(rhs.(Float32Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isLessThanEqual := lhs.Value <= rhs.(Float64Value).Value
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64Equal(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isEqual := lhs.Value == float64(rhs.(Int8Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isEqual := lhs.Value == float64(rhs.(Int16Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isEqual := lhs.Value == float64(rhs.(Int32Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isEqual := lhs.Value == float64(rhs.(Int64Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isEqual := lhs.Value == float64(rhs.(Float32Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isEqual := lhs.Value == rhs.(Float64Value).Value
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case Bool:
		isEqual := false
		if rhs.(BoolValue).Value {
			isEqual = lhs.Value == float64(1)
		} else {
			isEqual = lhs.Value == float64(0)
		}
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case Null:
		isEqual := lhs.Value == float64(0)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F64NotEqual(lhs Float64Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isNotEqual := lhs.Value != float64(rhs.(Int8Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isNotEqual := lhs.Value != float64(rhs.(Int16Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isNotEqual := lhs.Value != float64(rhs.(Int32Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isNotEqual := lhs.Value != float64(rhs.(Int64Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isNotEqual := lhs.Value != float64(rhs.(Float32Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isNotEqual := lhs.Value != rhs.(Float64Value).Value
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}
