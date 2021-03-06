// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yall

import (
	"os"
)

func bindLambdaList(env *Env, lambdaList *Cell, args *Cell) {
	for c := lambdaList; c != Empty; c = c.cdr {
		e := c.car
		if symbol, ok := e.(*Symbol); ok {
			if symbol.name == "." { // &rest (&body)
				env.Intern(c.Cadr().(*Symbol), args)
				break
			} else { // normal arg
				expr := args.Car()
				env.Intern(symbol, expr)
				args = args.Cdr()
			}
		} else if cell, ok := e.(*Cell); ok {
			symbol := cell.car.(*Symbol)
			if Empty == args { // default arg
				defaultValue := cell.Cadr()
				env.Intern(symbol, defaultValue)
			} else { // default arg with a specified value
				expr := args.Car()
				env.Intern(symbol, expr)
				args = args.Cdr()
			}
		}
	}
}

func lambda(env *Env, args *Cell) Expr {
	lambdaList := args.Car().(*Cell)
	body := args.Cdr()
	return NewFunction("#lambda", func(args *Cell) Expr {
		derived := env.Derive()
		bindLambdaList(derived, lambdaList, args)
		return derived.Begin(body)
	})
}

func macro(env *Env, args *Cell) Expr {
	lambdaList := args.Car().(*Cell)
	body := args.Cdr()
	return NewMacro("#macro", func(args *Cell) Expr {
		derived := env.Derive()
		bindLambdaList(derived, lambdaList, args)
		return derived.Begin(body)
	})
}

var specialForms = map[string]func(*Env, *Cell) Expr{

	"def": func(env *Env, args *Cell) Expr {
		if symbol, ok := args.Car().(*Symbol); ok {
			value := env.Eval(args.Cadr())
			if function, ok := value.(*Function); ok {
				function.SetName(symbol.Name())
			}
			env.Intern(symbol, value)
			return symbol
		}
		panic(NewRuntimeError("Can't define"))
	},

	"lambda": lambda,
	"fn":     lambda,

	"macro": macro,

	"defn": func(env *Env, args *Cell) Expr {
		cell, ok := args.Car().(*Cell)
		if !ok {
			panic(NewRuntimeError("Can't define function."))
		}
		symbol := cell.Car().(*Symbol)
		lambdaArgs := cell.Cdr()
		lambdaBody := args.Cdr()
		f := lambda(env, NewCell(lambdaArgs, lambdaBody)).(*Function)
		f.SetName(symbol.Name())
		env.Intern(symbol, f)
		return symbol
	},

	"defmacro": func(env *Env, args *Cell) Expr {
		cell, ok := args.car.(*Cell)
		if !ok {
			panic(NewRuntimeError("Can't define macro."))
		}
		symbol := cell.car.(*Symbol)
		lambdaList := cell.Cdr()
		body := args.Cdr()
		m := macro(env, NewCell(lambdaList, body)).(*Macro)
		m.SetName(symbol.name)
		env.Intern(symbol, m)
		return symbol
	},

	"if": func(env *Env, args *Cell) Expr {
		condition := env.Eval(args.Car())
		if condition != False {
			return env.Eval(args.Cadr())
		}
		return env.Eval(args.Caddr())
	},

	"inc!": func(env *Env, args *Cell) Expr {
		symbol, ok := args.Car().(*Symbol)
		if !ok {
			panic(NewRuntimeError("inc! requires a symbol"))
		}
		integer, ok := env.Eval(symbol).(*Integer)
		integer.setValue(integer.Value() + 1)
		return integer
	},

	"load": func(env *Env, args *Cell) Expr {
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
	},
}
