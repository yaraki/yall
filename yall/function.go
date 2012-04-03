// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yall

import (
    "fmt"
)

var builtinFunctions = map[string]func(*Cell) Expr{

    "car": func(args *Cell) Expr {
        if cell, ok := args.Car().(*Cell); ok && cell != Empty {
            return cell.Car()
        }
        panic(NewRuntimeError("pair required, but got " + args.Car().String()))
    },

    "cdr": func(args *Cell) Expr {
        if cell, ok := args.Car().(*Cell); ok && cell != Empty {
            return cell.Cdr()
        }
        panic(NewRuntimeError("pair required, but got " + args.Car().String()))
    },

    "cons": func(arg *Cell) Expr {
        if cadr, ok := arg.Cadr().(*Cell); ok {
            return NewCell(arg.Car(), cadr)
        }
        panic(NewRuntimeError("Cons requires a cell for the second argument"))
    },

    "+": func(args *Cell) Expr {
        result := 0
        for Empty != args {
            result += args.Car().(*Integer).Value()
            args = args.Cdr()
        }
        return NewInteger(result)
    },

    "-": func(args *Cell) Expr {
        if Empty == args {
            panic(NewRuntimeError("Too few arguments to minus, at least 1 required"))
        }
        i, iok := args.Car().(*Integer)
        if !iok {
            panic(NewRuntimeError("Minus requires integers"))
        }
        result := i.Value()
        if Empty == args.Cdr() {
            return NewInteger(result * -1)
        }
        for cell := args.Cdr(); cell != Empty; cell = cell.Cdr() {
            i, iok := cell.Car().(*Integer)
            if !iok {
                panic(NewRuntimeError("Minus requires integers"))
            }
            result -= i.Value()
        }
        return NewInteger(result)
    },

    "*": func(args *Cell) Expr {
        result := 1
        for Empty != args {
            result *= args.Car().(*Integer).Value()
            args = args.Cdr()
        }
        return NewInteger(result)
    },

    "type-of": func(args *Cell) Expr {
        if args.Cdr() != Empty {
            panic(NewRuntimeError("Too many arguments to type-of"))
        }
        return typeOf(args.Car())
    },

    "println": func(args *Cell) Expr {
        args.Each(func(expr Expr) {
            if str, ok := expr.(*String); ok {
                fmt.Println(str.value)
            } else {
                fmt.Println(expr.String())
            }
        })
        return True
    },

    "empty?": func(args *Cell) Expr {
        if Empty == args.car {
            return True
        }
        return False
    },

    "list": func(args *Cell) Expr {
        return args
    },
}
