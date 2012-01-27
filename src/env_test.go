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
    evalTestCase{"(cons 1 2)", "(1 . 2)"},
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

func TestEval2(t *testing.T) {
    i := -10
    env := NewEnv()
    env.internVariable("a", NewInteger(i))
    expr := env.EvalString("a")
    if expr.String() != strconv.Itoa(i) {
        t.Errorf("Received [[%v]] when expecting [[%v]]", expr.String(), i)
    }
}

func TestEval3(t *testing.T) {
    env := NewEnv()
    env.internVariable("b", NewInteger(3))
    expr := env.EvalString("`(a ,b c)")
    answer := "(a 3 c)"
    if expr.String() != answer {
        t.Errorf("Received [[%v]] when expecting [[%v]]", expr.String(), answer)
    }
}
