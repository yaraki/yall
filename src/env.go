package yall

import (
    "fmt"
    "strings"
)

type Env struct {
    values map[string]Expr
	parent *Env
}

func NewEnv() *Env {
    env := new(Env)
    env.values = make(map[string]Expr)
	env.parent = nil
    env.internFunction("+", Plus)
    env.internFunction("*", Multiply)
    env.internSpecialForm("define", Define)
    env.internSpecialForm("lambda", Lambda)
    env.internVariable("#t", True)
    env.internVariable("#f", False)
	env.internFunction("car", Car)
	env.internFunction("cdr", Cdr)
	env.internFunction("cons", Cons)
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
		fmt.Println("*** Warning: Overwriting " + s)
	}
    env.values[s] = value
}

func (env *Env) Intern(symbol *Symbol, value Expr) {
    env.internVariable(symbol.Name(), value)
}

func (env *Env) Unintern(symbol *Symbol) {
    env.values[symbol.Name()] = nil, false
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
    fmt.Printf("*** ERROR: Unbound variable: %v\n", symbol)
    return nil
}

func (env *Env) EvalEach(cell *Cell) *Cell {
    if Empty == cell {
        return Empty
    }
    return NewCell(env.Eval(cell.Car()),
        env.EvalEach(cell.Cdr().(*Cell)))
    return Empty
}

func (env *Env) EvalCell(cell *Cell) Expr {
    head := env.Eval(cell.Car())
    if form, ok := head.(*SpecialForm); ok {
        if list, lok := cell.Cdr().(*Cell); lok {
            return form.Apply(env, list)
        }
    } else if function, ok := head.(*Function); ok {
        if list, lok := cell.Cdr().(*Cell); lok {
            return function.Apply(env.EvalEach(list))
        }
    }
    fmt.Println(";; [ERROR] Failed to eval cell")
    return nil
}

func (env *Env) EvalQuasiquoted(expr Expr) Expr {
    if unquoted, ok := expr.(*Unquoted); ok {
        return env.Eval(unquoted.expr)
    }
    if cell, ok := expr.(*Cell); ok && Empty != cell {
        return NewCell(env.EvalQuasiquoted(cell.Car()), env.EvalQuasiquoted(cell.Cdr()))
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
    fmt.Println(";; [ERROR] Failed to eval")
    return nil
}

func (env *Env) EvalString(s string) Expr {
    return env.Eval(Parse(strings.NewReader(s)))
}

func (env *Env) Begin(cell *Cell) Expr {
    var result Expr
    for Empty != cell {
        result = env.Eval(cell.Car())
        cell = cell.Cdr().(*Cell)
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
