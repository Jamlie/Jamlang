package runtimelang

import (
	"fmt"
	"os"
	"strings"

	"github.com/Jamlie/Jamlang/ast"
	"github.com/Jamlie/Jamlang/internal"
)

func hasDecimalZero(num float64) bool {
	strNum := fmt.Sprintf("%f", num)
	fmt.Println(strNum)
	return strings.HasSuffix(strNum, ".0")
}

func Evaluate(astNode ast.Statement, env Environment) (RuntimeValue, error) {
	switch astNode.Kind() {
	case ast.CommentType:
		return MakeNullValue(), nil
	case ast.NumericFloatLiteralType:
		if isFloat32(astNode.(*ast.NumericFloatLiteral).Value) {
			return Float32Value{float32(astNode.(*ast.NumericFloatLiteral).Value)}, nil
		}
		return Float64Value{astNode.(*ast.NumericFloatLiteral).Value}, nil
	case ast.NumericIntegerLiteralType:
		if isInt8(float64(astNode.(*ast.NumericIntegerLiteral).Value)) || isInt16(float64(astNode.(*ast.NumericIntegerLiteral).Value)) || isInt32(float64(astNode.(*ast.NumericIntegerLiteral).Value)) {
			return Int32Value{int32(astNode.(*ast.NumericIntegerLiteral).Value)}, nil
		}
		if isInt64(float64(astNode.(*ast.NumericIntegerLiteral).Value)) {
			return Int64Value{int64(astNode.(*ast.NumericIntegerLiteral).Value)}, nil
		}
		i := astNode.(*ast.NumericIntegerLiteral).Value
		return Float64Value{float64(i)}, nil
	case ast.StringLiteralType:
		return StringValue{astNode.(*ast.StringLiteral).Value}, nil
	case ast.NullLiteralType:
		return MakeNullValue(), nil
	case ast.BinaryExpressionType:
		binaryExpression, ok := astNode.(*ast.BinaryExpression)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected BinaryExpression, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateBinaryExpression(*binaryExpression, env), nil
	case ast.UnaryExpressionType:
		unaryExpression, ok := astNode.(*ast.UnaryExpression)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected UnaryExpression, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateUnaryExpression(*unaryExpression, env), nil
	case ast.LogicalExpressionType:
		logicalExpression, ok := astNode.(*ast.LogicalExpression)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected LogicalExpression, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateLogicalExpression(*logicalExpression, env), nil
	case ast.IdentifierType:
		identifier, ok := astNode.(*ast.Identifier)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected Identifier, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateIdentifier(identifier, &env), nil
	case ast.ObjectLiteralType:
		objectLiteral, ok := astNode.(*ast.ObjectLiteral)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ObjectLiteral, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateObjectExpression(*objectLiteral, env), nil
	case ast.ArrayLiteralType:
		arrayLiteral, ok := astNode.(*ast.ArrayLiteral)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ArrayLiteral, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateArrayExpression(*arrayLiteral, env), nil
	case ast.TupleLiteralType:
		tupleLiteral, ok := astNode.(*ast.TupleLiteral)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected TupleLiteral, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateTupleExpression(*tupleLiteral, env), nil
	case ast.MemberExpressionType:
		memberExpression, ok := astNode.(*ast.MemberExpression)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected MemberExpression, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateMemberExpression(*memberExpression, env), nil
	case ast.CallExpressionType:
		callExpression, ok := astNode.(*ast.CallExpression)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected CallExpression, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateCallExpression(*callExpression, env), nil
	case ast.ReturnStatementType:
		returnStatement, ok := astNode.(*ast.ReturnStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ReturnStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateReturnStatement(*returnStatement, env), nil
	case ast.BreakStatementType:
		breakStatement, ok := astNode.(*ast.BreakStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected BreakStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateBreakStatement(*breakStatement, env), nil
	case ast.ContinueStatementType:
		continueStatement, ok := astNode.(*ast.ContinueStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ContinueStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateContinueStatement(*continueStatement, env), nil
	case ast.AssignmentExpressionType:
		assignmentExpression, ok := astNode.(*ast.AssignmentExpression)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected AssignmentExpression, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateAssignment(*assignmentExpression, env), nil
	case ast.ProgramType:
		program, ok := astNode.(*ast.Program)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected Program, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateProgram(*program, env), nil
	case ast.VariableDeclarationType:
		variableDeclaration, ok := astNode.(*ast.VariableDeclaration)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected VariableDeclaration, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}
		return EvaluateVariableDeclaration(*variableDeclaration, &env, variableDeclaration.Type), nil
	case ast.FunctionDeclarationType:
		functionDeclaration, ok := astNode.(*ast.FunctionDeclaration)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected FunctionDeclaration, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}

		fn, _ := EvaluateFunctionDeclaration(*functionDeclaration, &env, functionDeclaration.ReturnType)
		return fn, nil
	case ast.ConditionalStatementType:
		conditionalStatement, ok := astNode.(*ast.ConditionalStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ConditionalStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}

		result, err := EvaluateConditionalStatement(*conditionalStatement, &env)

		if err == IsReturnError {
			return result, IsReturnError
		}

		if err == IsBreakError {
			return result, IsBreakError
		}

		if err == IsContinueError {
			return result, IsContinueError
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
			return nil, nil
		}

		return result, nil
	case ast.WhileStatementType:
		whileStatement, ok := astNode.(*ast.WhileStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected WhileStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}

		result, err := EvaluateWhileStatement(*whileStatement, &env)

		if err == IsReturnError {
			return result, IsReturnError
		}

		if err == IsBreakError {
			return result, nil
		}

		if err == IsContinueError {
			return result, nil
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
			return nil, nil
		}

		return result, nil
	case ast.LoopStatementType:
		loopStatement, ok := astNode.(*ast.LoopStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected LoopStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}

		result, err := EvaluateLoopStatement(*loopStatement, &env)

		if err == IsReturnError {
			return result, IsReturnError
		}

		if err == IsBreakError {
			return result, nil
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
			return nil, nil
		}

		return result, nil
	case ast.ForEachStatementType:
		forEachStatement, ok := astNode.(*ast.ForEachStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ForEachStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}

		result, err := EvaluateForEachStatement(*forEachStatement, &env)

		if err == IsReturnError {
			return result, IsReturnError
		}

		if err == IsBreakError {
			return result, nil
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
			return nil, nil
		}

		return result, nil
	case ast.ForStatementType:
		forStatement, ok := astNode.(*ast.ForStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ForStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}

		result, err := EvaluateForStatement(*forStatement, &env)

		if err == IsReturnError {
			return result, IsReturnError
		}

		if err == IsBreakError {
			return result, nil
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
			return nil, nil
		}

		return result, nil
	case ast.ImportStatementType:
		importStatement, ok := astNode.(*ast.ImportStatement)
		if !ok {
			fmt.Fprintf(os.Stderr, "Error on line %d: Expected ImportStatement, got %T\n", internal.Line(), astNode)
			os.Exit(0)
			return nil, nil
		}

		result, err := EvaluateImportStatement(*importStatement, &env)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error on line %d: %s\n", internal.Line(), err.Error())
			os.Exit(0)
			return nil, nil
		}

		return result, nil
	default:
		fmt.Fprintf(os.Stderr, "Error on line %d: Unknown AST node type %T\n", internal.Line(), astNode)
		os.Exit(0)
		return nil, nil
	}
}
