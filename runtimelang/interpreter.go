package runtimelang

import (
    "fmt"
    "os"

    "github.com/Jamlee977/CustomLanguage/ast"
)

func Evaluate(astNode ast.Statement, env Environment) (RuntimeValue, error) {
    switch astNode.Kind() {
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
            os.Exit(1)
            return nil, nil
        }
        return EvaluateBinaryExpression(*binaryExpression, env), nil
    case ast.IdentifierType:
        identifier, ok := astNode.(*ast.Identifier)
        if !ok {
            fmt.Printf("Error: Expected Identifier, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateIdentifier(identifier, env), nil
    case ast.ObjectLiteralType:
        objectLiteral, ok := astNode.(*ast.ObjectLiteral)
        if !ok {
            fmt.Printf("Error: Expected ObjectLiteral, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateObjectExpression(*objectLiteral, env), nil
    case ast.MemberExpressionType:
        memberExpression, ok := astNode.(*ast.MemberExpression)
        if !ok {
            fmt.Printf("Error: Expected MemberExpression, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateMemberExpression(*memberExpression, env), nil
    case ast.CallExpressionType:
        callExpression, ok := astNode.(*ast.CallExpression)
        if !ok {
            fmt.Printf("Error: Expected CallExpression, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateCallExpression(*callExpression, env), nil
    case ast.AssignmentExpressionType:
        assignmentExpression, ok := astNode.(*ast.AssignmentExpression)
        if !ok {
            fmt.Printf("Error: Expected AssignmentExpression, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateAssignment(*assignmentExpression, env), nil
    case ast.ProgramType:
        program, ok := astNode.(*ast.Program)
        if !ok {
            fmt.Printf("Error: Expected Program, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateProgram(*program, env), nil
    case ast.VariableDeclarationType:
        variableDeclaration, ok := astNode.(*ast.VariableDeclaration)
        if !ok {
            fmt.Printf("Error: Expected VariableDeclaration, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateVariableDeclaration(*variableDeclaration, env), nil
    case ast.FunctionDeclarationType:
        functionDeclaration, ok := astNode.(*ast.FunctionDeclaration)
        if !ok {
            fmt.Printf("Error: Expected FunctionDeclaration, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateFunctionDeclaration(*functionDeclaration, env), nil
    case ast.ConditionalStatementType:
        conditionalStatement, ok := astNode.(*ast.ConditionalStatement)
        if !ok {
            fmt.Printf("Error: Expected ConditionalStatement, got %T\n", astNode)
            os.Exit(1)
            return nil, nil
        }
        return EvaluateConditionalExpression(*conditionalStatement, env), nil
        // return EvaluateConditionalStatement(*conditionalStatement, env), nil
    default:
        fmt.Printf("Unknown AST node type %T\n", astNode)
        os.Exit(1)
        return nil, nil
    }
}
