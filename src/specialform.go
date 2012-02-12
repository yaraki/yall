// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yall

import (
    "os"
)

func Def(env *Env, args *Cell) Expr {
    if symbol, ok := args.Car().(*Symbol); ok {
        value := env.Eval(args.Cadr())
        if function, ok := value.(*Function); ok {
            function.SetName(symbol.Name())
        }
        env.Intern(symbol, value)
        return symbol
    } else if cell, ok := args.Car().(*Cell); ok {
        symbol := cell.Car().(*Symbol)
        lambdaArgs := cell.Cdr()
        lambdaBody := args.Cdr()
        lambda := Lambda(env, NewCell(lambdaArgs, lambdaBody)).(*Function)
        lambda.SetName(symbol.Name())
        env.Intern(symbol, lambda)
        return symbol
    }
    panic(NewRuntimeError("Can't define"))
}

func Lambda(env *Env, args *Cell) Expr {
    formalArgs := args.Car().(*Cell)
    body := args.Cdr()
    return NewFunction("#lambda", func(args *Cell) Expr {
        derived := env.Derive()
        formalArgs.Each(func(e Expr) {
            symbol := e.(*Symbol)
            expr := args.Car()
            derived.Intern(symbol, expr)
            args = args.Cdr()
        })
        return derived.Begin(body)
    })
}

var specialForms = map[string]func(*Env, *Cell) Expr{
    "macro": func(env *Env, args *Cell) Expr {
        formalArgs := args.Car().(*Cell)
        body := args.Cdr()
        return NewMacro("#macro", func(args *Cell) Expr {
            derived := env.Derive()
            formalArgs.Each(func(e Expr) {
                symbol := e.(*Symbol)
                expr := args.Car()
                derived.Intern(symbol, expr)
                args = args.Cdr()
            })
            return derived.Begin(body)
        })
    },
}

func If(env *Env, args *Cell) Expr {
    condition := env.Eval(args.Car())
    if condition != False {
        return env.Eval(args.Cadr())
    }
    return env.Eval(args.Caddr())
}

func Incf(env *Env, args *Cell) Expr {
    symbol, ok := args.Car().(*Symbol)
    if !ok {
        panic(NewRuntimeError("inc! requires a symbol"))
    }
    integer, ok := env.Eval(symbol).(*Integer)
    integer.setValue(integer.Value() + 1)
    return integer
}

func Load(env *Env, args *Cell) Expr {
    args.Each(func(expr Expr) {
        if filename, ok := expr.(*String); ok {
            file, err := os.Open(filename.value)
            if nil != err {
                panic(NewRuntimeError("Cannot load: " + filename.String()))
            }
            defer file.Close()
            env.Load(file)
        } else {
            panic(NewRuntimeError("Cannot load: " + expr.String()))
        }
    })
    return True
}
