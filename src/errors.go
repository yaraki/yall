package yall

import "strings"

type SyntaxError struct {
    message string
}

func NewSyntaxError(message string) *SyntaxError {
    return &SyntaxError{message}
}

func (serr *SyntaxError) String() string {
    return serr.message
}

func isString(s string) bool {
    return strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}

type RuntimeError struct {
    message string
}

func NewRuntimeError(message string) *RuntimeError {
    return &RuntimeError{message}
}

func (err *RuntimeError) String() string {
    return "*** ERROR: " + err.message
}


