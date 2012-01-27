package yall

import "fmt"

func Define(env *Env, args *Cell) Expr {
    if symbol, ok := args.Car().(*Symbol); ok {
        if list, lok := args.Cdr().(*Cell); lok {
            env.Intern(symbol, env.Eval(list.Car()))
            return symbol
        }
    }
    fmt.Println(";; [ERROR] Can't define")
    return nil
}

func Lambda(env *Env, args *Cell) Expr {
    formalArgs := args.Car().(*Cell)
    body := args.Cdr().(*Cell)
    return NewFunction("#lambda", func(args *Cell) Expr {
		derived := env.Derive()
		formalArgs.Each(func (e Expr) {
			symbol := e.(*Symbol)
			expr := args.Car()
			derived.Intern(symbol, expr)
			args = args.Tail()
		})
        return derived.Begin(body)
    })
}
