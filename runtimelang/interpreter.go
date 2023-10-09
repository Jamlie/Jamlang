package runtimelang

import (
    "fmt"
    "os"

    "github.com/Jamlee977/CustomLanguage/ast"
)

func Evaluate(astNode ast.Statement, env Environment) (RuntimeValue, error) {
    switch astNode.Kind() {
    case ast.CommentType:
        return nil, nil
    case ast.NumericLiteralType:
        return NumberValue{astNode.(*ast.NumericLiteral).Value}, nil
    case ast.StringLiteralType:
        return StringValue{astNode.(*ast.StringLiteral).Value}, nil
    case ast.NullLiteralType:
        return MakeNullValue(), nil
    case ast.BinaryExpressionType:
        binaryExpression, ok := astNode.(*ast.BinaryExpression)
        if !ok {
            fmt.Printf("Error: Expected BinaryExpression, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateBinaryExpression(*binaryExpression, env), nil
    case ast.UnaryExpressionType:
        unaryExpression, ok := astNode.(*ast.UnaryExpression)
        if !ok {
            fmt.Printf("Error: Expected UnaryExpression, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateUnaryExpression(*unaryExpression, env), nil
    case ast.LogicalExpressionType:
        logicalExpression, ok := astNode.(*ast.LogicalExpression)
        if !ok {
            fmt.Printf("Error: Expected LogicalExpression, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateLogicalExpression(*logicalExpression, env), nil
    case ast.IdentifierType:
        identifier, ok := astNode.(*ast.Identifier)
        if !ok {
            fmt.Printf("Error: Expected Identifier, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateIdentifier(identifier, &env), nil
    case ast.ObjectLiteralType:
        objectLiteral, ok := astNode.(*ast.ObjectLiteral)
        if !ok {
            fmt.Printf("Error: Expected ObjectLiteral, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateObjectExpression(*objectLiteral, env), nil
    case ast.ArrayLiteralType:
        arrayLiteral, ok := astNode.(*ast.ArrayLiteral)
        if !ok {
            fmt.Printf("Error: Expected ArrayLiteral, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateArrayExpression(*arrayLiteral, env), nil
    case ast.MemberExpressionType:
        memberExpression, ok := astNode.(*ast.MemberExpression)
        if !ok {
            fmt.Printf("Error: Expected MemberExpression, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateMemberExpression(*memberExpression, env), nil
    case ast.CallExpressionType:
        callExpression, ok := astNode.(*ast.CallExpression)
        if !ok {
            fmt.Printf("Error: Expected CallExpression, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateCallExpression(*callExpression, env), nil
    case ast.ReturnStatementType:
        returnStatement, ok := astNode.(*ast.ReturnStatement)
        if !ok {
            fmt.Printf("Error: Expected ReturnStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateReturnStatement(*returnStatement, env), nil
    case ast.ClassDeclarationType:
        classDeclaration, ok := astNode.(*ast.ClassDeclaration)
        if !ok {
            fmt.Printf("Error: Expected ClassDeclaration, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }

        result, err := EvaluateClassDeclaration(*classDeclaration, &env)
        if err != nil {
            return nil, err
        }

        return result, nil
    case ast.BreakStatementType:
        breakStatement, ok := astNode.(*ast.BreakStatement)
        if !ok {
            fmt.Printf("Error: Expected BreakStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateBreakStatement(*breakStatement, env), nil
    case ast.AssignmentExpressionType:
        assignmentExpression, ok := astNode.(*ast.AssignmentExpression)
        if !ok {
            fmt.Printf("Error: Expected AssignmentExpression, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateAssignment(*assignmentExpression, env), nil
    case ast.ProgramType:
        program, ok := astNode.(*ast.Program)
        if !ok {
            fmt.Printf("Error: Expected Program, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateProgram(*program, env), nil
    case ast.VariableDeclarationType:
        variableDeclaration, ok := astNode.(*ast.VariableDeclaration)
        if !ok {
            fmt.Printf("Error: Expected VariableDeclaration, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        return EvaluateVariableDeclaration(*variableDeclaration, &env), nil
    case ast.FunctionDeclarationType:
        functionDeclaration, ok := astNode.(*ast.FunctionDeclaration)
        if !ok {
            fmt.Printf("Error: Expected FunctionDeclaration, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }
        
        fn, _ := EvaluateFunctionDeclaration(*functionDeclaration, &env)
        return fn, nil
    case ast.ConditionalStatementType:
        conditionalStatement, ok := astNode.(*ast.ConditionalStatement)
        if !ok {
            fmt.Printf("Error: Expected ConditionalStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }

        result, err := EvaluateConditionalExpression(*conditionalStatement, &env)

        if err == IsReturnError {
            return result, IsReturnError
        }

        if err == IsBreakError {
            return result, IsBreakError
        }

        if err != nil {
            fmt.Printf("Error: %s\n", err.Error())
            os.Exit(0)
            return nil, nil
        }

        return result, nil
    case ast.WhileStatementType:
        whileStatement, ok := astNode.(*ast.WhileStatement)
        if !ok {
            fmt.Printf("Error: Expected WhileStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }

        result, err := EvaluateWhileExpression(*whileStatement, &env)

        if err == IsReturnError {
            return result, IsReturnError
        }

        if err == IsBreakError {
            return result, nil
        }

        if err != nil {
            fmt.Printf("Error: %s\n", err.Error())
            os.Exit(0)
            return nil, nil
        }

        return result, nil
    case ast.LoopStatementType:
        loopStatement, ok := astNode.(*ast.LoopStatement)
        if !ok {
            fmt.Printf("Error: Expected LoopStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }

        result, err := EvaluateLoopExpression(*loopStatement, &env)

        if err == IsReturnError {
            return result, IsReturnError
        }

        if err == IsBreakError {
            return result, nil
        }

        if err != nil {
            fmt.Printf("Error: %s\n", err.Error())
            os.Exit(0)
            return nil, nil
        }

        return result, nil
    case ast.ForEachStatementType:
        forEachStatement, ok := astNode.(*ast.ForEachStatement)
        if !ok {
            fmt.Printf("Error: Expected ForEachStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }

        result, err := EvaluateForEachExpression(*forEachStatement, &env)

        if err == IsReturnError {
            return result, IsReturnError
        }

        if err == IsBreakError {
            return result, nil
        }

        if err != nil {
            fmt.Printf("Error: %s\n", err.Error())
            os.Exit(0)
            return nil, nil
        }

        return result, nil
    case ast.ForStatementType:
        forStatement, ok := astNode.(*ast.ForStatement)
        if !ok {
            fmt.Printf("Error: Expected ForStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }

        result, err := EvaluateForExpression(*forStatement, &env)

        if err == IsReturnError {
            return result, IsReturnError
        }

        if err == IsBreakError {
            return result, nil
        }

        if err != nil {
            fmt.Printf("Error: %s\n", err.Error())
            os.Exit(0)
            return nil, nil
        }

        return result, nil
    case ast.ImportStatementType:
        importStatement, ok := astNode.(*ast.ImportStatement)
        if !ok {
            fmt.Printf("Error: Expected ImportStatement, got %T\n", astNode)
            os.Exit(0)
            return nil, nil
        }

        result, err := EvaluateImportExpression(*importStatement, &env)
        if err != nil {
            fmt.Printf("Error: %s\n", err.Error())
            os.Exit(0)
            return nil, nil
        }

        return result, nil
    default:
        fmt.Printf("Unknown AST node type %T\n", astNode)
        os.Exit(0)
        return nil, nil
    }
}
