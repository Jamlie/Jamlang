package runtimelang

import (
	"fmt"
	"math"
	"os"
)

func internal_I8Plus(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		var result int8 = 0
		result = lhs.Value + rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		var result int16 = 0
		result = int16(lhs.Value) + rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		var result int32 = 0
		result = int32(lhs.Value) + rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		var result int64 = 0
		result = int64(lhs.Value) + rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot add floating point values to integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8(), float32(), float64() to cast up or down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8Minus(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		var result int8 = 0
		result = lhs.Value - rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		var result int16 = 0
		result = int16(lhs.Value) - rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		var result int32 = 0
		result = int32(lhs.Value) - rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		var result int64 = 0
		result = int64(lhs.Value) - rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot subtract floating point values from integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8(), float32(), float64() to cast up or down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8Mult(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		var result int8 = 0
		result = lhs.Value * rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		var result int16 = 0
		result = int16(lhs.Value) * rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		var result int32 = 0
		result = int32(lhs.Value) * rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		var result int64 = 0
		result = int64(lhs.Value) * rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot multiply floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8(), float32(), float64() to cast up or down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8Pow(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		var result int8 = 0
		result = int8(math.Pow(float64(lhs.Value), float64(rhs.(Int8Value).Value)))
		return Int8Value{result}
	case I16:
		var result int16 = 0
		result = int16(math.Pow(float64(lhs.Value), float64(rhs.(Int16Value).Value)))
		return Int16Value{result}
	case I32:
		var result int32 = 0
		result = int32(math.Pow(float64(lhs.Value), float64(rhs.(Int32Value).Value)))
		return Int32Value{result}
	case I64:
		var result int64 = 0
		result = int64(math.Pow(float64(lhs.Value), float64(rhs.(Int64Value).Value)))
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot raise integer values to floating point powers")
		fmt.Fprintln(os.Stderr, "Consider using int8(), float32(), float64() to cast up or down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8Div(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int8 = 0
		result = lhs.Value / rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int16 = 0
		result = int16(lhs.Value) / rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int32 = 0
		result = int32(lhs.Value) / rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int64 = 0
		result = int64(lhs.Value) / rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot divide floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8(), float32(), float64() to cast up or down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8IntDiv(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int8 = 0
		result = lhs.Value / rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int16 = 0
		result = int16(lhs.Value) / rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int32 = 0
		result = int32(lhs.Value) / rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int64 = 0
		result = int64(lhs.Value) / rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot divide floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8(), float32(), float64() to cast up or down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8Mod(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		if rhs.(Int8Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int8 = 0
		result = lhs.Value % rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		if rhs.(Int16Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int16 = 0
		result = int16(lhs.Value) % rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		if rhs.(Int32Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int32 = 0
		result = int32(lhs.Value) % rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		if rhs.(Int64Value).Value == 0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result int64 = 0
		result = int64(lhs.Value) % rhs.(Int64Value).Value
		return Int64Value{result}
	case F32:
		if rhs.(Float32Value).Value == 0.0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float32 = 0.0
		result = float32(math.Mod(float64(lhs.Value), float64(rhs.(Float32Value).Value)))
		return Float32Value{result}
	case F64:
		if rhs.(Float64Value).Value == 0.0 {
			fmt.Fprintf(os.Stderr, "Error: Division by zero\n")
			os.Exit(0)
		}
		var result float64 = 0.0
		result = math.Mod(float64(lhs.Value), float64(rhs.(Float64Value).Value))
		return Float64Value{result}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8BitwiseAnd(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		var result int8 = 0
		result = lhs.Value & rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		var result int16 = 0
		result = int16(lhs.Value) & rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		var result int32 = 0
		result = int32(lhs.Value) & rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		var result int64 = 0
		result = int64(lhs.Value) & rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot bitwise and floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8() to cast down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8BitwiseOr(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		var result int8 = 0
		result = lhs.Value | rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		var result int16 = 0
		result = int16(lhs.Value) | rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		var result int32 = 0
		result = int32(lhs.Value) | rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		var result int64 = 0
		result = int64(lhs.Value) | rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot bitwise or floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8() to cast down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8BitwiseXor(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		var result int8 = 0
		result = lhs.Value ^ rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		var result int16 = 0
		result = int16(lhs.Value) ^ rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		var result int32 = 0
		result = int32(lhs.Value) ^ rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		var result int64 = 0
		result = int64(lhs.Value) ^ rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintln(os.Stderr, "Error: Cannot bitwise xor floating point values with integer values")
		fmt.Fprintln(os.Stderr, "Consider using int8() to cast down.")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8GreaterThan(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isGreaterThan := lhs.Value > rhs.(Int8Value).Value
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isGreaterThan := int16(lhs.Value) > rhs.(Int16Value).Value
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isGreaterThan := int32(lhs.Value) > rhs.(Int32Value).Value
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isGreaterThan := int64(lhs.Value) > rhs.(Int64Value).Value
		if isGreaterThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isGreaterThan := float32(lhs.Value) > rhs.(Float32Value).Value
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

func internal_I8GreaterThanEqual(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isGreaterThanEqual := lhs.Value >= rhs.(Int8Value).Value
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isGreaterThanEqual := int16(lhs.Value) >= rhs.(Int16Value).Value
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isGreaterThanEqual := int32(lhs.Value) >= rhs.(Int32Value).Value
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isGreaterThanEqual := int64(lhs.Value) >= rhs.(Int64Value).Value
		if isGreaterThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isGreaterThanEqual := float32(lhs.Value) >= rhs.(Float32Value).Value
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

func internal_I8LessThan(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isLessThan := lhs.Value < rhs.(Int8Value).Value
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isLessThan := int16(lhs.Value) < rhs.(Int16Value).Value
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isLessThan := int32(lhs.Value) < rhs.(Int32Value).Value
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isLessThan := int64(lhs.Value) < rhs.(Int64Value).Value
		if isLessThan {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isLessThan := float32(lhs.Value) < rhs.(Float32Value).Value
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

func internal_I8LessThanEqual(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isLessThanEqual := lhs.Value <= rhs.(Int8Value).Value
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isLessThanEqual := int16(lhs.Value) <= rhs.(Int16Value).Value
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isLessThanEqual := int32(lhs.Value) <= rhs.(Int32Value).Value
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isLessThanEqual := int64(lhs.Value) <= rhs.(Int64Value).Value
		if isLessThanEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isLessThanEqual := float32(lhs.Value) <= rhs.(Float32Value).Value
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

func internal_I8Equal(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isEqual := lhs.Value == rhs.(Int8Value).Value
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isEqual := int16(lhs.Value) == rhs.(Int16Value).Value
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isEqual := int32(lhs.Value) == rhs.(Int32Value).Value
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isEqual := int64(lhs.Value) == rhs.(Int64Value).Value
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isEqual := float32(lhs.Value) == rhs.(Float32Value).Value
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
			isEqual = lhs.Value == int8(1)
		} else {
			isEqual = lhs.Value == int8(0)
		}
		if isEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case Null:
		isEqual := lhs.Value == int8(0)
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

func internal_I8NotEqual(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		isNotEqual := lhs.Value != rhs.(Int8Value).Value
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I16:
		isNotEqual := int16(lhs.Value) != rhs.(Int16Value).Value
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I32:
		isNotEqual := int32(lhs.Value) != rhs.(Int32Value).Value
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case I64:
		isNotEqual := int64(lhs.Value) != rhs.(Int64Value).Value
		if isNotEqual {
			return BoolValue{true}
		}
		return BoolValue{false}
	case F32:
		isNotEqual := float32(lhs.Value) != rhs.(Float32Value).Value
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

func internal_I8LeftShift(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		result := lhs.Value << rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		result := int16(lhs.Value) << rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		result := int32(lhs.Value) << rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		result := int64(lhs.Value) << rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintf(os.Stderr, "Error: Cannot left shift by a floating point value\n")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}

func internal_I8RightShift(lhs Int8Value, rhs RuntimeValue) RuntimeValue {
	switch rhs.Type() {
	case I8:
		result := lhs.Value >> rhs.(Int8Value).Value
		return Int8Value{result}
	case I16:
		result := int16(lhs.Value) >> rhs.(Int16Value).Value
		return Int16Value{result}
	case I32:
		result := int32(lhs.Value) >> rhs.(Int32Value).Value
		return Int32Value{result}
	case I64:
		result := int64(lhs.Value) >> rhs.(Int64Value).Value
		return Int64Value{result}
	case F32, F64:
		fmt.Fprintf(os.Stderr, "Error: Cannot right shift by a floating point value\n")
		os.Exit(0)
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown type for rhs of binary expression: %s\n", rhs.Type())
		os.Exit(0)
	}
	return nil
}
