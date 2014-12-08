// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yall

import (
	"bufio"
	"bytes"
	"os"
)

type Env struct {
	values map[string]Expr
	parent *Env
}

func NewEnv() *Env {
	env := new(Env)
	env.values = make(map[string]Expr)
	env.parent = nil
	env.internVariable("#t", True)
	env.internVariable("#f", False)
	for name, form := range specialForms {
		env.internSpecialForm(name, form)
	}
	for name, function := range builtinFunctions {
		env.internFunction(name, function)
	}
	f, err := os.Open(os.Getenv("GOPATH") + "/src/github.com/yaraki/yall/lisp/sys.yall")
	if err != nil {
		panic(NewRuntimeError("Failed to open sys.yall"))
	}
	defer f.Close()
	env.Load(f)
	return env
}

func (env *Env) Derive() *Env {
	derived := new(Env)
	derived.values = make(map[string]Expr)
	derived.parent = env
	return derived
}

func (env *Env) internSpecialForm(s string, f func(*Env, *Cell) Expr) {
	env.internVariable(s, NewSpecialForm(s, f))
}

func (env *Env) internFunction(s string, f func(*Cell) Expr) {
	env.internVariable(s, NewFunction(s, f))
}

func (env *Env) internVariable(s string, value Expr) {
	if nil != env.values[s] {
		panic(NewRuntimeError("Can't overwrite " + s))
	}
	env.values[s] = value
}

func (env *Env) Intern(symbol *Symbol, value Expr) {
	env.internVariable(symbol.Name(), value)
}

func (env *Env) Unintern(symbol *Symbol) {
	delete(env.values, symbol.Name())
}

func (env *Env) EvalSymbol(symbol *Symbol) Expr {
	if value, found := env.values[symbol.Name()]; found {
		return value
	}
	if nil != env.parent {
		if value := env.parent.EvalSymbol(symbol); value != nil {
			return value
		}
	}
	panic(NewRuntimeError("Unbound variable: " + symbol.String()))
}

func (env *Env) EvalEach(cell *Cell) *Cell {
	if Empty == cell {
		return Empty
	}
	return NewCell(env.Eval(cell.Car()),
		env.EvalEach(cell.Cdr()))
	return Empty
}

func (env *Env) EvalCell(cell *Cell) Expr {
	head := env.Eval(cell.Car())
	if form, ok := head.(*SpecialForm); ok {
		return form.Apply(env, cell.Cdr())
	} else if function, ok := head.(*Function); ok {
		return function.Apply(env.EvalEach(cell.Cdr()))
	} else if macro, ok := head.(*Macro); ok {
		return env.Eval(macro.Expand(cell.Cdr()))
	}
	for k, v := range env.values {
		println("# " + k + ": " + v.String())
	}
	panic(NewRuntimeError("Failed to eval cell: " + cell.String()))
}

func (env *Env) EvalQuasiquoted(expr Expr) Expr {
	if unquoted, ok := expr.(*Unquoted); ok {
		return env.Eval(unquoted.expr)
	}
	if cell, ok := expr.(*Cell); ok && cell != Empty {
		if splicing, sok := cell.Cadr().(*SplicingUnquoted); sok {
			car := env.EvalQuasiquoted(cell.car)
			cadr := env.Eval(splicing.expr)
			if c, cok := cadr.(*Cell); cok {
				return NewCell(car, c)
			}
			panic(NewRuntimeError("Invalid splicing unquote"))
		}
		return NewCell(env.EvalQuasiquoted(cell.car), env.EvalQuasiquoted(cell.cdr).(*Cell))
	}
	return expr
}

func (env *Env) Eval(expr Expr) Expr {
	if IsLiteral(expr) {
		return expr
	}
	if symbol, ok := expr.(*Symbol); ok {
		return env.EvalSymbol(symbol)
	}
	if quoted, ok := expr.(*Quoted); ok {
		return quoted.expr
	}
	if quasiquoted, ok := expr.(*Quasiquoted); ok {
		return env.EvalQuasiquoted(quasiquoted.expr)
	}
	if cell, ok := expr.(*Cell); ok {
		return env.EvalCell(cell)
	}
	panic(NewRuntimeError("Failed to eval"))
}

func (env *Env) EvalString(s string) Expr {
	if expr, _, err := ReadFromString(s); err == nil {
		return env.Eval(expr)
	}
	return nil
}

// unused
func slurp(file *os.File) string {
	buffer := bytes.NewBufferString("")
	var buf [1024]byte
	for {
		switch n, err := file.Read(buf[:]); true {
		case n < 0:
			panic(err)
		case 0 < n:
			buffer.Write(buf[0:n])
		}
	}
	return buffer.String()
}

func (env *Env) Load(file *os.File) {
	r := &reader{bufio.NewReader(file)}
	for {
		expr, _, err := r.Read()
		if err != nil {
			break
		}
		env.Eval(expr)
	}
}

func (env *Env) Begin(cell *Cell) Expr {
	var result Expr
	for Empty != cell {
		result = env.Eval(cell.Car())
		cell = cell.Cdr()
	}
	return result
}

func IsLiteral(expr Expr) bool {
	if Empty == expr {
		return true
	}
	if _, ok := expr.(*Integer); ok {
		return true
	}
	if _, ok := expr.(*String); ok {
		return true
	}
	return false
}

func typeOf(expr Expr) *Type {
	if _, ok := expr.(*Cell); ok {
		return TYPE_CELL
	}
	if _, ok := expr.(*Symbol); ok {
		return TYPE_SYMBOL
	}
	if _, ok := expr.(*Integer); ok {
		return TYPE_INTEGER
	}
	if _, ok := expr.(*String); ok {
		return TYPE_STRING
	}
	if _, ok := expr.(*Function); ok {
		return TYPE_FUNCTION
	}
	if _, ok := expr.(*Macro); ok {
		return TYPE_MACRO
	}
	if _, ok := expr.(*SpecialForm); ok {
		return TYPE_SPECIAL_FORM
	}
	if _, ok := expr.(*Bool); ok {
		return TYPE_BOOL
	}
	if _, ok := expr.(*Type); ok {
		return TYPE_TYPE
	}
	return TYPE_UNKNOWN
}
