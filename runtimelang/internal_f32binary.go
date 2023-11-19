package runtimelang

import (
	"fmt"
	"math"
	"os"
)

func internal_F32Plus(lhs Float32Value, rhs RuntimeValue) RuntimeValue {

	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot add floating point values to integer values")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		return Float32Value{lhs.Value + rhs.(Float32Value).Value}
	case F64:
		var result float64 = 0
		result = float64(lhs.Value) + rhs.(Float64Value).Value
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F32Minus(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot subtract floating point values from integer values")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		return Float32Value{lhs.Value - rhs.(Float32Value).Value}
	case F64:
		var result float64 = 0
		result = float64(lhs.Value) - rhs.(Float64Value).Value
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F32Mult(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot multiply floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		return Float32Value{lhs.Value * rhs.(Float32Value).Value}
	case F64:
		var result float64 = 0
		result = float64(lhs.Value) * rhs.(Float64Value).Value
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F32Pow(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8, I16, I32, I64:
		fmt.Fprintln(os.Stderr, "Error: Cannot raise floating point values to integer powers")
		fmt.Fprintln(os.Stderr, "Consider using float32(), float64() to cast up or int8(), int16(), int32(), int64() to down.")
		os.Exit(0)
	case F32:
		return Float32Value{float32(math.Pow(float64(lhs.Value), float64(rhs.(Float32Value).Value)))}
	case F64:
		var result float64 = 0
		result = math.Pow(float64(lhs.Value), rhs.(Float64Value).Value)
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F32Div(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float32Value{lhs.Value / float32(rhs.(Int8Value).Value)}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float32Value{lhs.Value / float32(rhs.(Int16Value).Value)}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float32Value{lhs.Value / float32(rhs.(Int32Value).Value)}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float32Value{lhs.Value / float32(rhs.(Int64Value).Value)}
	case F32:
		if rhs.(Float32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		return Float32Value{lhs.Value / rhs.(Float32Value).Value}
	case F64:
		if rhs.(Float64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = float64(lhs.Value) / rhs.(Float64Value).Value
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F32IntDiv(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Floor(float64(lhs.Value) / float64(rhs.(Int8Value).Value)))
		return Float32Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Floor(float64(lhs.Value) / float64(rhs.(Int16Value).Value)))
		return Float32Value{result}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Floor(float64(lhs.Value) / float64(rhs.(Int32Value).Value)))
		return Float32Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Floor(float64(lhs.Value) / float64(rhs.(Int64Value).Value)))
		return Float32Value{result}
	case F32:
		if rhs.(Float32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Floor(float64(lhs.Value) / float64(rhs.(Float32Value).Value)))
		return Float32Value{result}
	case F64:
		if rhs.(Float64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Floor(float64(lhs.Value) / rhs.(Float64Value).Value)
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F32Mod(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Mod(float64(lhs.Value), float64(rhs.(Int8Value).Value)))
		return Float32Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Mod(float64(lhs.Value), float64(rhs.(Int16Value).Value)))
		return Float32Value{result}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Mod(float64(lhs.Value), float64(rhs.(Int32Value).Value)))
		return Float32Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Mod(float64(lhs.Value), float64(rhs.(Int64Value).Value)))
		return Float32Value{result}
	case F32:
		if rhs.(Float32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0
		result = float32(math.Mod(float64(lhs.Value), float64(rhs.(Float32Value).Value)))
		return Float32Value{result}
	case F64:
		if rhs.(Float64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0
		result = math.Mod(float64(lhs.Value), rhs.(Float64Value).Value)
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_F32GreaterThan(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isGreaterThan := lhs.Value > float32(rhs.(Int8Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isGreaterThan := lhs.Value > float32(rhs.(Int16Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isGreaterThan := lhs.Value > float32(rhs.(Int32Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isGreaterThan := lhs.Value > float32(rhs.(Int64Value).Value)
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isGreaterThan := lhs.Value > rhs.(Float32Value).Value
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isGreaterThan := float64(lhs.Value) > rhs.(Float64Value).Value
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

func internal_F32GreaterThanEqual(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isGreaterThanEqual := lhs.Value >= float32(rhs.(Int8Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isGreaterThanEqual := lhs.Value >= float32(rhs.(Int16Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isGreaterThanEqual := lhs.Value >= float32(rhs.(Int32Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isGreaterThanEqual := lhs.Value >= float32(rhs.(Int64Value).Value)
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isGreaterThanEqual := lhs.Value >= rhs.(Float32Value).Value
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isGreaterThanEqual := float64(lhs.Value) >= rhs.(Float64Value).Value
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

func internal_F32LessThan(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isLessThan := lhs.Value < float32(rhs.(Int8Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isLessThan := lhs.Value < float32(rhs.(Int16Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isLessThan := lhs.Value < float32(rhs.(Int32Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isLessThan := lhs.Value < float32(rhs.(Int64Value).Value)
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isLessThan := lhs.Value < rhs.(Float32Value).Value
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isLessThan := float64(lhs.Value) < rhs.(Float64Value).Value
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

func internal_F32LessThanEqual(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isLessThanEqual := lhs.Value <= float32(rhs.(Int8Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isLessThanEqual := lhs.Value <= float32(rhs.(Int16Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isLessThanEqual := lhs.Value <= float32(rhs.(Int32Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isLessThanEqual := lhs.Value <= float32(rhs.(Int64Value).Value)
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isLessThanEqual := lhs.Value <= rhs.(Float32Value).Value
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isLessThanEqual := float64(lhs.Value) <= rhs.(Float64Value).Value
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

func internal_F32Equal(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isEqual := lhs.Value == float32(rhs.(Int8Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isEqual := lhs.Value == float32(rhs.(Int16Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isEqual := lhs.Value == float32(rhs.(Int32Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isEqual := lhs.Value == float32(rhs.(Int64Value).Value)
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isEqual := lhs.Value == rhs.(Float32Value).Value
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isEqual := float64(lhs.Value) == rhs.(Float64Value).Value
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case Bool:
		isEqual := false
		if rhs.(BoolValue).Value {
			isEqual = lhs.Value == float32(1)
		} else {
			isEqual = lhs.Value == float32(0)
		}
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case Null:
		isEqual := lhs.Value == float32(0)
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

func internal_F32NotEqual(lhs Float32Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isNotEqual := lhs.Value != float32(rhs.(Int8Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isNotEqual := lhs.Value != float32(rhs.(Int16Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isNotEqual := lhs.Value != float32(rhs.(Int32Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isNotEqual := lhs.Value != float32(rhs.(Int64Value).Value)
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isNotEqual := lhs.Value != rhs.(Float32Value).Value
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F64:
		isNotEqual := float64(lhs.Value) != rhs.(Float64Value).Value
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
