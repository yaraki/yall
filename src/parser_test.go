package yall

import (
    "strings"
    "testing"
)

type tokenTestCase struct {
    input  string
    output []string
}

var tokenTestCases = []tokenTestCase{
    tokenTestCase{"(abc def ghi)", []string{"(", "abc", "def", "ghi", ")"}},
    tokenTestCase{"(abc (def ghi))", []string{"(", "abc", "(", "def", "ghi", ")", ")"}},
    tokenTestCase{"(\"abc def\")", []string{"(", "\"abc def\"", ")"}},
    tokenTestCase{"(abc def 'ghi)", []string{"(", "abc", "def", "'", "ghi", ")"}},
    tokenTestCase{"`(a ,b ,@(c))", []string{"`", "(", "a", ",", "b", ",@", "(", "c", ")", ")"}},
}

func TestTokenize(t *testing.T) {
    for _, tc := range tokenTestCases {
        tokenizer := Tokenize(strings.NewReader(tc.input))
        for _, answer := range tc.output {
            token, err := tokenizer.Next()
            if nil != err {
                t.Errorf("Unexpected err")
            }
            if token != answer {
                t.Errorf("Received [[%v]] when expecting [[%v]]", token, answer)
            }
        }
    }
}

type parseTestCase struct {
    input  string
    output string
}

var parseTestCases = []parseTestCase{
    parseTestCase{"12", "12"},
    parseTestCase{"-99", "-99"},
    parseTestCase{"+", "+"},
    parseTestCase{"(a  b)", "(a b)"},
    parseTestCase{"((a)  b (c d(e)))", "((a) b (c d (e)))"},
    parseTestCase{"\"hello\"", "\"hello\""},
    parseTestCase{"\"\"", "\"\""},
    parseTestCase{":test", ":test"},
}

func TestParse(t *testing.T) {
    for _, tc := range parseTestCases {
        result := Parse(strings.NewReader(tc.input))
        if result.String() != tc.output {
            t.Errorf("Received [[%v]] when expecting [[%v]]", result.String(), tc.output)
        }
    }
}

func TestIsString(t *testing.T) {
    s1 := "\"abc\""
    if !isString(s1) {
        t.Errorf("%v is expected to be a string.", s1)
    }
    s2 := "abc-def"
    if isString(s2) {
        t.Errorf("%v is not expected to be a string.", s2)
    }
}
