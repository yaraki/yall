// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yall

import (
    "strconv"
    "testing"
)

type evalTestCase struct {
    input  string
    output string
}

var evalTestCases = []evalTestCase{
    evalTestCase{"()", "()"},
    evalTestCase{"123", "123"},
    evalTestCase{"(+ 1 2)", "3"},
    evalTestCase{"(+ 1 (* 2 3))", "7"},
    evalTestCase{"(cons 1 '(2))", "(1 2)"},
}

func TestEval(t *testing.T) {
    for _, tc := range evalTestCases {
        env := NewEnv()
        expr := env.EvalString(tc.input)
        if tc.output != expr.String() {
            t.Errorf("Received [[%v]] when expecting [[%v]]", expr.String(), tc.output)
        }
    }
}

// Interned symbol
func TestEval2(t *testing.T) {
    i := -10
    env := NewEnv()
    env.internVariable("a", NewInteger(i))
    expr := env.EvalString("a")
    if expr.String() != strconv.Itoa(i) {
        t.Errorf("Received [[%v]] when expecting [[%v]]", expr.String(), i)
    }
}

// Quasiquote
func TestEval3(t *testing.T) {
    env := NewEnv()
    env.internVariable("b", NewInteger(3))
    env.internVariable("c", NewCell(NewInteger(1), NewCell(NewInteger(2), Empty)))
    expr := env.EvalString("`(a ,b ,@c)")
    answer := "(a 3 1 2)"
    if expr.String() != answer {
        t.Errorf("Received [[%v]] when expecting [[%v]]", expr.String(), answer)
    }
}

// Closure (environment enclosing)
func TestEval4(t *testing.T) {
    env := NewEnv()
    env.EvalString("(def genc (lambda () ((lambda (x) (lambda () (inc! x))) 0)))")
    env.EvalString("(def a (genc))")
    if "1" != env.EvalString("(a)").String() {
        t.Errorf("Something is wrong with closure.")
    }
    if "2" != env.EvalString("(a)").String() {
        t.Errorf("Something is wrong with closure.")
    }
    if "3" != env.EvalString("(a)").String() {
        t.Errorf("Something is wrong with closure.")
    }
}
