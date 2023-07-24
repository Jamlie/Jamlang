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
            err := fmt.Errorf("Expected BinaryExpression, got %T", astNode)
            panic(err)
        }
        return EvaluateBinaryExpression(*binaryExpression, env), nil
    case ast.IdentifierType:
        identifier, ok := astNode.(*ast.Identifier)
        if !ok {
            err := fmt.Errorf("Expected Identifier, got %T", astNode)
            panic(err)
        }

        return EvaluateIdentifier(identifier, env), nil
    case ast.AssignmentExpressionType:
        assignmentExpression, ok := astNode.(*ast.AssignmentExpression)
        if !ok {
            err := fmt.Errorf("Expected AssignmentExpression, got %T", astNode)
            panic(err)
        }
        return EvaluateAssignment(*assignmentExpression, env), nil
    case ast.ProgramType:
        program, ok := astNode.(*ast.Program)
        if !ok {
            err := fmt.Errorf("Expected Program, got %T", astNode)
            panic(err)
        }
        return EvaluateProgram(*program, env), nil
    case ast.VariableDeclarationType:
        variableDeclaration, ok := astNode.(*ast.VariableDeclaration)
        if !ok {
            err := fmt.Errorf("Expected VariableDeclaration, got %T", astNode)
            panic(err)
        }
        return EvaluateVariableDeclaration(*variableDeclaration, env), nil
    default:
        fmt.Printf("Unknown AST node type %T\n", astNode)
        os.Exit(1)
        return nil, nil
    }
}
