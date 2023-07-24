package runtimelang

import (
    "fmt"
    "os"
)

type Environment struct {
    parent *Environment
    variables map[string]RuntimeValue
    constants map[string]bool
}

func CreateGlobalEnvironment() *Environment {
    env := NewEnvironment(nil)
    env.DeclareVariable("true", MakeBoolValue(true), true)
    env.DeclareVariable("false", MakeBoolValue(false), true)
    env.DeclareVariable("null", MakeNullValue(), true)

    return env
}


func NewEnvironment(parent *Environment) *Environment {
    return &Environment{parent, make(map[string]RuntimeValue), make(map[string]bool)}
}

func (e *Environment) DeclareVariable(name string, value RuntimeValue, constant bool) RuntimeValue {
    if _, ok := e.variables[name]; ok {
        err := fmt.Errorf("Variable %s already declared", name)
        panic(err)
    }

    e.variables[name] = value

    if constant {
        e.constants[name] = true
    }

    return value
}

func (e *Environment) AssignVariable(name string, value RuntimeValue) RuntimeValue {
    env := e.Resolve(name)
    if env == nil {
        fmt.Printf("Variable %s not declared\n", name)
        os.Exit(1)
        return nil
    }

    if env.constants[name] {
        fmt.Printf("Variable %s is constant. Cannot reassign a constant.\n", name)
        os.Exit(1)
        return nil
    }

    env.variables[name] = value

    return value
}

func (e *Environment) Resolve(name string) *Environment {
    if _, ok := e.variables[name]; ok {
        return e
    }

    if e.parent == nil {
        fmt.Printf("Variable %s not declared\n", name)
        os.Exit(1)
        return nil
    }

    return e.parent.Resolve(name)
}

func (e *Environment) LookupVariable(name string) RuntimeValue {
    env := e.Resolve(name)
    if env == nil {
        fmt.Printf("Variable %s not declared\n", name)
        os.Exit(1)
        return nil
    }

    return env.variables[name]
}
