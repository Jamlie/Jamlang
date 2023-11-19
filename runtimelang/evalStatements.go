package runtimelang

import (
	"fmt"
	"io"
	"os"

	"github.com/Jamlie/Jamlang/ast"
	"github.com/Jamlie/Jamlang/internal"
	"github.com/Jamlie/Jamlang/parser"
)

func EvaluateProgram(program ast.Program, env Environment) RuntimeValue {
	var lastEvaluated RuntimeValue = &InitialValue{}
	for _, statement := range program.Body {
		lastEvaluated, _ = Evaluate(statement, env)
	}
	return lastEvaluated
}

func EvaluateImportStatement(expr ast.ImportStatement, env *Environment) (RuntimeValue, error) {
	file, err := os.Open(expr.Path)
	if err != nil {
		return MakeNullValue(), err
	}
	defer file.Close()

	parser := parser.NewParser()
	var fileString string
	if fileBytes, err := io.ReadAll(file); err != nil {
		return MakeNullValue(), err
	} else {
		fileString = string(fileBytes)
	}

	path := expr.Path
	if path[0] == '.' {
		path = path[1:]
	}
	path = path[:len(path)-4]
	if path[0] == '/' {
		path = path[1:]
	}

	program := parser.ProduceAST(fileString)
	for _, statement := range program.Body {
		if statement.Kind() == ast.FunctionDeclarationType {
			function := statement.(*ast.FunctionDeclaration)
			fn, _ := Evaluate(function, *env)
			if _, ok := fn.(FunctionValue); !ok {
				continue
			}

			if function.Name[0] >= 'A' && function.Name[0] <= 'Z' {
				env.variables[function.Name] = fn
			}
		} else if statement.Kind() == ast.VariableDeclarationType {
			variable := statement.(*ast.VariableDeclaration)
			value, _ := Evaluate(variable.Value, *env)
			if _, ok := value.(RuntimeValue); !ok {
				continue
			}

			if variable.Identifier[0] >= 'A' && variable.Identifier[0] <= 'Z' {
				env.variables[variable.Identifier] = value
			}
		}
	}

	return MakeNullValue(), nil
}

func EvaluateForStatement(expr ast.ForStatement, env *Environment) (RuntimeValue, error) {
	scope := NewEnvironment(env)

	initialVariable := expr.Init.(*ast.VariableDeclaration).Identifier
	_, err := Evaluate(expr.Init, *env)
	defer func() {
		env.RemoveVariable(initialVariable)
	}()

	if err != nil {
		return MakeNullValue(), err
	}

	for {
		if expr.Condition != nil {
			condition, err := Evaluate(expr.Condition, *scope)
			if err != nil {
				return MakeNullValue(), err
			}
			if condition.Type() == Bool {
				if !condition.(BoolValue).Value {
					break
				}
			} else {
				return MakeNullValue(), fmt.Errorf("for loop condition must be a boolean value")
			}
		}

		for _, statement := range expr.Body {
			if statement.Kind() == ast.BreakStatementType {
				return MakeNullValue(), nil
			}

			_, err := Evaluate(statement, *scope)

			if err != nil {
				return MakeNullValue(), err
			}
		}
		if expr.Update != nil {
			_, err := Evaluate(expr.Update, *scope)
			if err != nil {
				return MakeNullValue(), err
			}
		}

		for k, v := range scope.variables {
			if _, ok := env.variables[k]; ok {
				env.variables[k] = v
			}
		}

		scope = NewEnvironment(env)
	}

	return MakeNullValue(), nil
}

func EvaluateForEachStatement(expr ast.ForEachStatement, env *Environment) (RuntimeValue, error) {
	scope := NewEnvironment(env)

	collection, err := Evaluate(expr.Collection, *env)
	if err != nil {
		return MakeNullValue(), err
	}

	if collection.Type() == Array {
		scope.DeclareVariable(expr.Variable, MakeNullValue(), false, ast.AnyType)
		for _, element := range collection.(ArrayValue).Values {
			scope.AssignVariable(expr.Variable, element)

			for _, statement := range expr.Body {
				if statement.Kind() == ast.ReturnStatementType {
					return Evaluate(statement, *scope)
				}
				_, err := Evaluate(statement, *scope)
				if err != nil {
					return MakeNullValue(), err
				}
			}
		}
	} else if collection.Type() == Tuple {
		scope.DeclareVariable(expr.Variable, MakeNullValue(), false, ast.AnyType)
		for _, element := range collection.(TupleValue).Values {
			scope.AssignVariable(expr.Variable, element)

			for _, statement := range expr.Body {
				if statement.Kind() == ast.ReturnStatementType {
					return Evaluate(statement, *scope)
				}
				_, err := Evaluate(statement, *scope)
				if err != nil {
					return MakeNullValue(), err
				}
			}
		}
	} else if collection.Type() == String {
		scope.DeclareVariable(expr.Variable, MakeNullValue(), false, ast.StringType)
		for _, element := range collection.(StringValue).Value {
			scope.AssignVariable(expr.Variable, StringValue{Value: string(element)})

			for _, statement := range expr.Body {
				if statement.Kind() == ast.ReturnStatementType {
					return Evaluate(statement, *scope)
				}
				_, err := Evaluate(statement, *scope)
				if err != nil {
					return MakeNullValue(), err
				}
			}
		}
	} else if collection.Type() == Object {
		scope.DeclareVariable(expr.Key, MakeNullValue(), false, ast.AnyType)
		scope.DeclareVariable(expr.Value, MakeNullValue(), false, ast.AnyType)
		for key, value := range collection.(ObjectValue).Properties {
			scope.AssignVariable(expr.Key, StringValue{Value: key})
			scope.AssignVariable(expr.Value, value)

			for _, statement := range expr.Body {
				if statement.Kind() == ast.ReturnStatementType {
					return Evaluate(statement, *scope)
				}
				_, err := Evaluate(statement, *scope)
				if err != nil {
					return MakeNullValue(), err
				}
			}
		}
	} else {
		return MakeNullValue(), fmt.Errorf("Cannot iterate over non-iterable type")
	}

	return MakeNullValue(), nil
}

func EvaluateLoopStatement(expr ast.LoopStatement, env *Environment) (RuntimeValue, error) {
	scope := NewEnvironment(env)

	for {
		for _, statement := range expr.Body {
			if statement.Kind() == ast.ReturnStatementType {
				result, err := Evaluate(statement, *scope)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
					os.Exit(0)
				}
				return result, IsReturnError
			}
			if statement.Kind() == ast.BreakStatementType {
				return MakeNullValue(), nil
			}

			_, err := Evaluate(statement, *scope)
			if err == IsBreakError {
				return MakeNullValue(), nil
			}
			if err == IsContinueError {
				continue
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
				os.Exit(0)
			}
		}

		scope = NewEnvironment(env)
	}
}

func EvaluateWhileStatement(expr ast.WhileStatement, env *Environment) (RuntimeValue, error) {
	condition, err := Evaluate(expr.Condition, *env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
		os.Exit(0)
	}

	scope := NewEnvironment(env)

	if condition.Type() != Bool {
		return MakeNullValue(), fmt.Errorf("while statement condition must be a boolean")
	}

	for condition.Get() == true {
		for _, statement := range expr.Body {
			if statement.Kind() == ast.ReturnStatementType {
				result, err := Evaluate(statement, *scope)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
					os.Exit(0)
				}
				return result, IsReturnError
			}
			if statement.Kind() == ast.BreakStatementType {
				return MakeNullValue(), IsBreakError
			}

			_, err := Evaluate(statement, *scope)
			if err == IsBreakError {
				return MakeNullValue(), nil
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
				os.Exit(0)
			}
		}

		scope = NewEnvironment(env)
		condition, err = Evaluate(expr.Condition, *env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
		}
	}

	return MakeNullValue(), nil
}

func EvaluateConditionalStatement(expr ast.ConditionalStatement, env *Environment) (RuntimeValue, error) {
	condition, err := Evaluate(expr.Condition, *env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
		os.Exit(0)
	}

	scope := NewEnvironment(env)

	if condition.Type() != Bool {
		fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), "if statement condition must be a boolean")
		os.Exit(0)
	}

	if condition.Get() == true {
		for _, statement := range expr.Body {
			if statement.Kind() == ast.ReturnStatementType {
				result, err := Evaluate(statement, *scope)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on line %d: %s", internal.Line(), err.Error())
					os.Exit(0)
				}

				// IsReturnError is a special error that is used to indicate that a return statement has been reached
				return result, IsReturnError
			} else if statement.Kind() == ast.BreakStatementType {
				return MakeNullValue(), IsBreakError
			} else if statement.Kind() == ast.ContinueStatementType {
				return MakeNullValue(), IsContinueError
			}
			_, err := Evaluate(statement, *scope)
			if err == IsBreakError {
				return MakeNullValue(), nil
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
				os.Exit(0)
			}
		}

		scope = NewEnvironment(env)
	} else {
		trueCondIdx := -1
		for idx, elseifCond := range expr.ElseIfConditions {
			cond, err := Evaluate(elseifCond, *env)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
				os.Exit(0)
			}

			if cond.Type() != Bool {
				return MakeNullValue(), fmt.Errorf("elseif statement condition must be a boolean")
			}

			if cond.Get() == true {
				trueCondIdx = idx
				break
			}
		}

		if trueCondIdx != -1 {
			for _, statement := range expr.ElseIfBodies[trueCondIdx] {
				if statement.Kind() == ast.ReturnStatementType {
					result, err := Evaluate(statement, *scope)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
						os.Exit(0)
					}

					return result, IsReturnError
				} else if statement.Kind() == ast.BreakStatementType {
					return MakeNullValue(), IsBreakError
				}
				_, err := Evaluate(statement, *scope)
				if err == IsBreakError {
					return MakeNullValue(), nil
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
					os.Exit(0)
				}
			}

			scope = NewEnvironment(env)
		} else {
			for _, statement := range expr.Alternate {
				if statement.Kind() == ast.ReturnStatementType {
					result, err := Evaluate(statement, *scope)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
						os.Exit(0)
					}

					return result, IsReturnError
				} else if statement.Kind() == ast.BreakStatementType {
					return MakeNullValue(), IsBreakError
				} else if statement.Kind() == ast.ContinueStatementType {
					return MakeNullValue(), IsContinueError
				}
				_, err := Evaluate(statement, *scope)
				if err == IsBreakError {
					return MakeNullValue(), nil
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
					os.Exit(0)
				}
			}

			scope = NewEnvironment(env)
		}
	}

	return MakeNullValue(), nil
}

func EvaluateBreakStatement(statement ast.BreakStatement, env Environment) RuntimeValue {
	return &BreakType{}
}

func EvaluateContinueStatement(statement ast.ContinueStatement, env Environment) RuntimeValue {
	return &ContinueType{}
}

func EvaluateReturnStatement(statement ast.ReturnStatement, env Environment) RuntimeValue {
	returnValue, _ := Evaluate(statement.Value, env)
	return returnValue
}

func checkNumberTypes(value RuntimeValue, varType ast.VariableType) bool {
	return !(value.Type() == Number && varType == ast.Float64Type) && !(value.Type() == Number && varType == ast.Float32Type) && !(value.Type() == Number && varType == ast.Int64Type) && !(value.Type() == Number && varType == ast.Int32Type) && !(value.Type() == Number && varType == ast.Int16Type) && !(value.Type() == Number && varType == ast.Int8Type)
}

func makeValueWithVarType(value RuntimeValue, varType ast.VariableType) RuntimeValue {
	switch varType {
	case ast.Int8Type:
		if _, ok := value.(IntValue); !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), ast.Int8Type, value.VarType())
			os.Exit(0)
		}
		return MakeInt8Value(int8(value.(IntValue).GetInt()))
	case ast.Int16Type:
		if _, ok := value.(IntValue); !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), ast.Int16Type, value.VarType())
			os.Exit(0)
		}
		return MakeInt16Value(int16(value.(IntValue).GetInt()))
	case ast.Int32Type:
		if _, ok := value.(IntValue); !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), ast.Int32Type, value.VarType())
			os.Exit(0)
		}
		return MakeInt32Value(int32(value.(IntValue).GetInt()))
	case ast.Int64Type:
		if _, ok := value.(IntValue); !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), ast.Int64Type, value.VarType())
			os.Exit(0)
		}
		return MakeInt64Value(int64(value.(IntValue).GetInt()))
	case ast.Float32Type:
		if _, ok := value.(FloatValue); !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), ast.Float32Type, value.VarType())
			os.Exit(0)
		}
		return MakeFloat32Value(float32(value.(FloatValue).GetFloat()))
	case ast.Float64Type:
		if _, ok := value.(FloatValue); !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), ast.Float64Type, value.VarType())
			os.Exit(0)
		}
		return MakeFloat64Value(float64(value.(FloatValue).GetFloat()))
	case ast.ObjectType:
		if _, ok := value.(NullValue); ok {
			return value
		}
		if _, ok := value.(ObjectValue); !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), ast.ObjectType, value.VarType())
			os.Exit(0)
		}
		return value
	case ast.NullType:
		return value
	case ast.AnyType:
		return value
	default:
		return value
	}
}

func EvaluateVariableDeclaration(declaration ast.VariableDeclaration, env *Environment, varType ast.VariableType) RuntimeValue {
	value, _ := Evaluate(declaration.Value, *env)

	actualValue := makeValueWithVarType(value, varType)

	if value.VarType() != varType && varType != ast.AnyType && !isNumber(value) {
		fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), declaration.Type, value.VarType())
		os.Exit(0)
	}

	if isNumber(value) {
		if declaration.Type == ast.Int8Type && value.VarType() != ast.Int8Type {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), declaration.Type, value.VarType())
			os.Exit(0)
		} else if declaration.Type == ast.Int16Type && (value.VarType() != ast.Int8Type && value.VarType() != ast.Int16Type) {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), declaration.Type, value.VarType())
			os.Exit(0)
		} else if declaration.Type == ast.Int32Type && (value.VarType() != ast.Int8Type && value.VarType() != ast.Int16Type && value.VarType() != ast.Int32Type) {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), declaration.Type, value.VarType())
			os.Exit(0)
		} else if declaration.Type == ast.Int64Type && (value.VarType() != ast.Int8Type && value.VarType() != ast.Int16Type && value.VarType() != ast.Int32Type && value.VarType() != ast.Int64Type) {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), declaration.Type, value.VarType())
			os.Exit(0)
		} else if declaration.Type == ast.Float32Type && value.VarType() != ast.Float32Type {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), declaration.Type, value.VarType())
			os.Exit(0)
		} else if declaration.Type == ast.Float64Type && (value.VarType() != ast.Float32Type && value.VarType() != ast.Float64Type) {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected %s, got %s\n", internal.Line(), declaration.Type, value.VarType())
			os.Exit(0)
		}
	}

	return env.DeclareVariable(declaration.Identifier, actualValue, declaration.Constant, varType)
}

func EvaluateVariableDeclarationDeprecated(declaration ast.VariableDeclaration, env *Environment) RuntimeValue {
	value, _ := Evaluate(declaration.Value, *env)
	if value.Type() == Object && value.(ObjectValue).IsClass {
		value = value.(ObjectValue).Clone()
	}
	return env.DeclareVariable(declaration.Identifier, value, declaration.Constant, ast.AnyType)
}

func EvaluateIdentifier(identifier *ast.Identifier, env *Environment) RuntimeValue {
	if identifier == nil {
		return MakeNullValue()
	}
	value := env.LookupVariable(identifier.Symbol)
	if value == nil {
		fmt.Fprintf(os.Stderr, "Error on line %d:Undefined variable %s\n", internal.Line(), identifier.Symbol)
		os.Exit(0)
		return nil
	}

	return value
}
