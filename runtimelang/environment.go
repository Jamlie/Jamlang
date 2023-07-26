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

    env.DeclareVariable("println", MakeNativeFunction(jamlangPrintln), true)
    env.DeclareVariable("print", MakeNativeFunction(jamlangPrint), true)
    env.DeclareVariable("sleep", MakeNativeFunction(jamlangSleep), true)
    env.DeclareVariable("typeof", MakeNativeFunction(jamlangTypeof), true)
    env.DeclareVariable("exit", MakeNativeFunction(jamlangExit), true)
    env.DeclareVariable("input", MakeNativeFunction(jamlangInput), true)

    return env
}

func NewEnvironment(parent *Environment) *Environment {
    return &Environment{parent, make(map[string]RuntimeValue), make(map[string]bool)}
}

func (e *Environment) DeclareVariable(name string, value RuntimeValue, constant bool) RuntimeValue {
    if _, ok := e.variables[name]; ok {
        fmt.Printf("Variable %s already declared\n", name)
        os.Exit(0)
        return nil
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
        os.Exit(0)
        return nil
    }

    if env.constants[name] {
        fmt.Printf("Variable %s is constant. Cannot reassign a constant.\n", name)
        os.Exit(0)
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
        return nil
    }

    return e.parent.Resolve(name)
}

func (e *Environment) LookupVariable(name string) RuntimeValue {
    env := e.Resolve(name)
    if env == nil {
        fmt.Printf("Variable %s not declared\n", name)
        os.Exit(0)
        return nil
    }

    return env.variables[name]
}
