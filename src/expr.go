package yall

import (
    "fmt"
    "strconv"
)

type Expr interface {
    String() string
}

type Cell struct {
    car Expr
    cdr Expr
}

var Empty *Cell = &Cell{nil, nil}

func NewCell(car Expr, cdr Expr) *Cell {
	return &Cell{ car, cdr }
    // cell := new(Cell)
    // cell.car = car
    // cell.cdr = cdr
    // return cell
}

func (cell *Cell) stringWithoutParens() string {
    // The cell is printed as a list when its cdr is a cell
    if tail, ok := cell.cdr.(*Cell); ok {
        if Empty == cell.cdr {
            return cell.car.String()
        }
        return fmt.Sprintf("%v %v", cell.car.String(), tail.stringWithoutParens())
    }
    // The cell is a dot cell
    return fmt.Sprintf("%v . %v", cell.car, cell.cdr)
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

func (cell *Cell) Cdr() Expr {
    return cell.cdr
}

func (cell *Cell) Tail() *Cell {
	return cell.cdr.(*Cell)
}

func (cell *Cell) Cadr() Expr {
	return cell.cdr.(*Cell).car
}

func (cell *Cell) Each(f func(Expr)) {
	for c := cell; c != Empty; c = c.Tail() {
		f(c.Car())
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

func (function *Function) Apply(args *Cell) Expr {
    return function.f(args)
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
