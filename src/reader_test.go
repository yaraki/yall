// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yall

import (
    "strings"
    "testing"
)

type nextTokenTestCase struct {
    input   string
    outputs []string
}

var nextTokenTestCases = []nextTokenTestCase{
    nextTokenTestCase{"abc", []string{"abc"}},
    nextTokenTestCase{"(a)", []string{"(", "a", ")"}},
    nextTokenTestCase{"(a  b   \n  c)", []string{"(", "a", "b", "c", ")"}},
    nextTokenTestCase{"\"hello  world\"", []string{"\"hello  world\""}},
    nextTokenTestCase{"\"a\\\"b\"", []string{"\"a\\\"b\""}},
    nextTokenTestCase{"(abc def ghi)", []string{"(", "abc", "def", "ghi", ")"}},
    nextTokenTestCase{"(abc (def ghi))", []string{"(", "abc", "(", "def", "ghi", ")", ")"}},
    nextTokenTestCase{"(\"abc def\")", []string{"(", "\"abc def\"", ")"}},
    nextTokenTestCase{"(abc def 'ghi)", []string{"(", "abc", "def", "'", "ghi", ")"}},
    nextTokenTestCase{"`(a ,b ,@(c))", []string{"`", "(", "a", ",", "b", ",@", "(", "c", ")", ")"}},
}

func TestNextToken(t *testing.T) {
    for _, tc := range nextTokenTestCases {
        r := &reader{}
        r.setInput(strings.NewReader(tc.input))
        for _, expected := range tc.outputs {
            received, _, err := r.nextToken()
            if err != nil {
                t.Errorf("expected: [[%v]], received: !!ERROR!! (%v)", expected, err)
            }
            if received != expected {
                t.Errorf("expected: [[%v]], received: [[%v]]", expected, received)
            }
        }
    }
}

type readTestCase struct {
    input  string
    output string
    size   int
}

var readTestCases = []readTestCase{
    readTestCase{"1", "1", 1},
    readTestCase{"-99", "-99", 3},
    readTestCase{"abc", "abc", 3},
    readTestCase{"(a)", "(a)", 3},
    readTestCase{"+", "+", 1},
    readTestCase{"(a  b)", "(a b)", 6},
    readTestCase{"((a)  b (c d(e)))", "((a) b (c d (e)))", 17},
    readTestCase{"\"hello\"", "\"hello\"", 7},
    readTestCase{"\"\"", "\"\"", 2},
    readTestCase{":test", ":test", 5},
    readTestCase{"`(a ,b c)", "`(a ,b c)", 9},
    readTestCase{"'a", "'a", 2},
}

func TestRead(t *testing.T) {
    for _, tc := range readTestCases {
        expr, size, err := ReadFromString(tc.input)
        if err != nil {
            t.Errorf("input: [[%v]], ERROR: [[%v]]", tc.input, err)
        } else if expr.String() != tc.output {
            t.Errorf("expected: [[%v]], received: [[%v]]", tc.output, expr)
        } else if size != tc.size {
            t.Errorf("input: [[%v]], expected size: [[%v]], received size: [[%v]]", tc.input, tc.size, size)
        }
    }
}
