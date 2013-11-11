// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package core

import (
    "fmt"
    "strconv"
)

type Expr interface {
    String() string
}

type Cell struct {
    car Expr
    cdr *Cell
}

var Empty *Cell = &Cell{nil, nil}

func NewCell(car Expr, cdr *Cell) *Cell {
    return &Cell{car, cdr}
}

func (cell *Cell) stringWithoutParens() string {
    // The cell is printed as a list when its cdr is a cell
    if Empty == cell.cdr {
        return cell.car.String()
    }
    return fmt.Sprintf("%v %v", cell.car.String(), cell.cdr.stringWithoutParens())
    // The cell is a dot cell
    // return fmt.Sprintf("%v . %v", cell.car, cell.cdr)
}

func (cell *Cell) String() string {
    if cell == Empty {
        return "()"
    }
    return "(" + cell.stringWithoutParens() + ")"
}

func (cell *Cell) Car() Expr {
    return cell.car
}

func (cell *Cell) Cdr() *Cell {
    return cell.cdr
}

func (cell *Cell) Cadr() Expr {
    return cell.cdr.car
}

func (cell *Cell) Cddr() Expr {
    return cell.cdr.cdr
}

func (cell *Cell) Caddr() Expr {
    return cell.cdr.cdr.car
}

func (cell *Cell) Each(f func(Expr)) {
    for c := cell; c != Empty; c = c.cdr {
        f(c.car)
    }
}

type Symbol struct {
    name string
}

func NewSymbol(name string) *Symbol {
    symbol := new(Symbol)
    symbol.name = name
    return symbol
}

func (symbol *Symbol) String() string {
    return symbol.name
}

func (symbol *Symbol) Name() string {
    return symbol.name
}

type Integer struct {
    value int
}

func NewInteger(value int) *Integer {
    integer := new(Integer)
    integer.value = value
    return integer
}

func (integer *Integer) String() string {
    return strconv.Itoa(integer.value)
}

func (integer *Integer) Value() int {
    return integer.value
}

func (integer *Integer) setValue(value int) {
    integer.value = value
}

type String struct {
    value string
}

func NewString(value string) *String {
    s := new(String)
    s.value = value
    return s
}

func (s *String) String() string {
    return "\"" + s.value + "\""
}

type Quoted struct {
    expr Expr
}

func NewQuoted(expr Expr) *Quoted {
    quoted := new(Quoted)
    quoted.expr = expr
    return quoted
}

func (quoted *Quoted) String() string {
    return "'" + quoted.expr.String()
}

type Quasiquoted struct {
    expr Expr
}

func NewQuasiquoted(expr Expr) *Quasiquoted {
    quasiquoted := new(Quasiquoted)
    quasiquoted.expr = expr
    return quasiquoted
}

func (quasiquoted *Quasiquoted) String() string {
    return "`" + quasiquoted.expr.String()
}

type Unquoted struct {
    expr Expr
}

func NewUnquoted(expr Expr) *Unquoted {
    unquoted := new(Unquoted)
    unquoted.expr = expr
    return unquoted
}

func (unquoted *Unquoted) String() string {
    return "," + unquoted.expr.String()
}

type SplicingUnquoted struct {
    expr Expr
}

func NewSplicingUnquoted(expr Expr) *SplicingUnquoted {
    unquoted := new(SplicingUnquoted)
    unquoted.expr = expr
    return unquoted
}

func (unquoted *SplicingUnquoted) String() string {
    return ",@" + unquoted.expr.String()
}

type Function struct {
    name string
    f    func(*Cell) Expr
}

func NewFunction(name string, f func(*Cell) Expr) *Function {
    function := new(Function)
    function.name = name
    function.f = f
    return function
}

func (function *Function) String() string {
    return "<function " + function.name + ">"
}

func (function *Function) SetName(name string) {
    function.name = name
}

func (function *Function) Apply(args *Cell) Expr {
    return function.f(args)
}

type Macro struct {
    name string
    f    func(*Cell) Expr
}

func NewMacro(name string, f func(*Cell) Expr) *Macro {
    macro := new(Macro)
    macro.name = name
    macro.f = f
    return macro
}

func (macro *Macro) String() string {
    return "<macro " + macro.name + ">"
}

func (macro *Macro) SetName(name string) {
    macro.name = name
}

func (macro *Macro) Expand(args *Cell) Expr {
    return macro.f(args)
}

type SpecialForm struct {
    name string
    f    func(*Env, *Cell) Expr
}

func NewSpecialForm(name string, f func(*Env, *Cell) Expr) *SpecialForm {
    form := new(SpecialForm)
    form.name = name
    form.f = f
    return form
}

func (form *SpecialForm) String() string {
    return "<special-form " + form.name + ">"
}

func (form *SpecialForm) Apply(env *Env, args *Cell) Expr {
    return form.f(env, args)
}

type Bool struct {
    value bool
}

var True *Bool = &Bool{true}
var False *Bool = &Bool{false}

func (b *Bool) String() string {
    if b.value {
        return "#t"
    }
    return "#f"
}

type Type struct {
    name string
}

func NewType(name string) *Type {
    return &Type{name}
}

func (t *Type) String() string {
    return "<" + t.name + ">"
}

var TYPE_CELL *Type = NewType("cell")
var TYPE_SYMBOL *Type = NewType("symbol")
var TYPE_INTEGER *Type = NewType("integer")
var TYPE_STRING *Type = NewType("string")
var TYPE_FUNCTION *Type = NewType("function")
var TYPE_MACRO *Type = NewType("macro")
var TYPE_SPECIAL_FORM *Type = NewType("special-form")
var TYPE_BOOL *Type = NewType("bool")
var TYPE_TYPE *Type = NewType("type")
var TYPE_UNKNOWN *Type = NewType("unknown")
