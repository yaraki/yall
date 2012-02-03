// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yall

import (
    "bufio"
    "bytes"
    "io"
    "os"
    "strconv"
    "strings"
)

type reader struct {
    input *bufio.Reader
}

func (r *reader) setInput(input io.Reader) {
    r.input = bufio.NewReader(input)
}

func (r *reader) nextString() (token string, size int, err os.Error) {
    buffer := bytes.NewBufferString("\"")
    size = 1
    escaped := false
    for {
        rune, s, err := r.input.ReadRune()
        size += s
        if err != nil {
            break
        }
        switch rune {
        case '"':
            buffer.WriteRune('"')
            if !escaped {
                return buffer.String(), size, nil
            }
        case '\\':
            escaped = true
            buffer.WriteRune(rune)
        default:
            escaped = false
            buffer.WriteRune(rune)
        }
    }
    return "", 0, NewSyntaxError("Unexpected EOS in string")
}

func (r *reader) nextToken() (token string, size int, err os.Error) {
    buffer := new(bytes.Buffer)
    size = 0
    for {
        rune, s, err := r.input.ReadRune()
        size += s
        if err != nil { // EOS
            if 0 < buffer.Len() {
                return buffer.String(), size, nil
            }
            return "", 0, err
        }
        switch rune {
        case '(', ')', '[', ']', '\'', '`':
            if 0 < buffer.Len() {
                r.input.UnreadRune()
                size -= s
                return buffer.String(), size, nil
            }
            return string(rune), size, nil
        case ' ', '\t', '\r', '\n':
            if 0 < buffer.Len() {
                return buffer.String(), size, nil
            }
        case ',':
            if 0 < buffer.Len() {
                r.input.UnreadRune()
                size -= s
                return buffer.String(), size, nil
            }
            maybeAt, as, aerr := r.input.ReadRune()
            if maybeAt != '@' || aerr != nil {
                r.input.UnreadRune()
                return ",", size, nil
            }
            size += as
            return ",@", size, nil
        case '"':
            if 0 < buffer.Len() {
                r.input.UnreadRune()
                size -= s
                return buffer.String(), size, nil
            }
            return r.nextString()
        default:
            buffer.WriteRune(rune)
        }
    }
    return "", 0, NewSyntaxError("?")
}

func (r *reader) readTokens(asList bool) (Expr, int, os.Error) {
    token, size, err := r.nextToken()
    if nil != err {
        panic(NewSyntaxError("Unexpected end of input"))
    }
    var expr Expr
    if "(" == token {
        list, listSize, listErr := r.readTokens(true)
        expr = list
        size += listSize
        err = listErr
    } else if ")" == token {
        if !asList {
            panic(NewSyntaxError("Unexpected end of list"))
        }
        expr = Empty
        asList = false
    } else if "'" == token {
        rest, restSize, restErr := r.readTokens(false)
        expr = NewQuoted(rest)
        size += restSize
        err = restErr
    } else if "`" == token {
        rest, restSize, restErr := r.readTokens(false)
        expr = NewQuasiquoted(rest)
        size += restSize
        err = restErr
    } else if "," == token {
        rest, restSize, restErr := r.readTokens(false)
        expr = NewUnquoted(rest)
        size += restSize
        err = restErr
    } else if ",@" == token {
        list, listSize, listErr := r.readTokens(false)
        expr = NewSplicingUnquoted(list)
        size += listSize
        err = listErr
    } else if i, ierr := strconv.Atoi(token); ierr == nil {
        expr = NewInteger(i)
    } else if isString(token) {
        expr = NewString(token[1 : len(token)-1])
    } else {
        expr = NewSymbol(token)
    }
    if asList {
        cdr, cdrSize, cdrErr := r.readTokens(true)
        return NewCell(expr, cdr.(*Cell)), size + cdrSize, cdrErr
    }
    return expr, size, nil
}

func (r *reader) read(input io.Reader) (Expr, int, os.Error) {
    r.setInput(input)
    return r.readTokens(false)
}

func Read(input io.Reader) (expr Expr, size int, err os.Error) {
    return new(reader).read(input)
}

func ReadFromString(s string) (expr Expr, size int, err os.Error) {
    return Read(strings.NewReader(s))
}

