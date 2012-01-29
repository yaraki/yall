package yall

import "fmt"

func Car(args *Cell) Expr {
    if cell, ok := args.Car().(*Cell); ok && cell != Empty {
        return cell.Car()
    }
    fmt.Println("*** ERROR: pair required, but got " + args.Car().String())
    return nil
}

func Cdr(args *Cell) Expr {
    if cell, ok := args.Car().(*Cell); ok && cell != Empty {
        return cell.Cdr()
    }
    fmt.Println("*** ERROR: pair required, but got " + args.Car().String())
    return nil
}

func Cons(arg *Cell) Expr {
    if cadr, ok := arg.Cadr().(*Cell); ok {
        return NewCell(arg.Car(), cadr)
    }
    panic("Cons requires a cell")
}

func Plus(args *Cell) Expr {
    result := 0
    for Empty != args {
        result += args.Car().(*Integer).Value()
        args = args.Cdr()
    }
    return NewInteger(result)
}

func Multiply(args *Cell) Expr {
    result := 1
    for Empty != args {
        result *= args.Car().(*Integer).Value()
        args = args.Cdr()
    }
    return NewInteger(result)
}
