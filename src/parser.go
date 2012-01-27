package yall

import (
    "bufio"
    "bytes"
    "fmt"
    "io"
    "os"
    "strconv"
    "strings"
)

type Tokenizer struct {
    reader *bufio.Reader
    saved  *string
}

func Tokenize(reader io.Reader) *Tokenizer {
    tokenizer := new(Tokenizer)
    tokenizer.reader = bufio.NewReader(reader)
    tokenizer.saved = nil
    return tokenizer
}

func (t *Tokenizer) tokenizeString() (string, os.Error) {
    var c int
    var err os.Error
    buffer := bytes.NewBufferString("\"")
    for {
        c, _, err = t.reader.ReadRune()
        if nil != err {
            break
        }
        switch c {
        case '"':
            fmt.Fprintf(buffer, "%c", '"')
            return string(buffer.Bytes()), nil
        default:
            fmt.Fprintf(buffer, "%c", c)
        }
    }
    return "", err
}

func (t *Tokenizer) Next() (string, os.Error) {
    if nil != t.saved {
        var s string = *t.saved
        t.saved = nil
        return s, nil
    }
    var c int
    var err os.Error
    buffer := bytes.NewBufferString("")
    for {
        c, _, err = t.reader.ReadRune()
        if nil != err {
            break
        }
        switch c {
        case '(', ')', '[', ']', '\'', '`':
            if 0 < buffer.Len() {
                s := string(c)
                t.saved = &s
                return string(buffer.Bytes()), nil
            } else {
                return string(c), nil
            }
        case ',':
            maybeAt, _, aerr := t.reader.ReadRune()
            var unquote string
            if nil == aerr && maybeAt == '@' {
                unquote = ",@"
            } else {
                t.reader.UnreadRune()
                unquote = ","
            }
            if 0 < buffer.Len() {
                t.saved = &unquote
                return string(buffer.Bytes()), nil
            } else {
                return unquote, nil
            }
        case ' ', '\t':
            if 0 < buffer.Len() {
                return string(buffer.Bytes()), nil
            }
        case '"':
            if 0 < buffer.Len() {
                s, serr := t.tokenizeString()
                if serr != nil {
                    return "", serr
                }
                t.saved = &s
                return string(buffer.Bytes()), nil
            } else {
                s, serr := t.tokenizeString()
                if serr != nil {
                    return "", serr
                }
                return s, nil
            }
        default:
            fmt.Fprintf(buffer, "%c", c)
        }
    }
    if 0 < buffer.Len() {
        return string(buffer.Bytes()), nil
    }
    return "", err
}

func parseTokens(tokenizer *Tokenizer, asList bool) Expr {
    var ret Expr
    token, err := tokenizer.Next()
    if nil != err {
        return Empty
    }
    if "(" == token {
        ret = parseTokens(tokenizer, true)
    } else if ")" == token {
        ret = Empty
        asList = false
    } else if "'" == token {
        ret = NewQuoted(parseTokens(tokenizer, false))
    } else if "`" == token {
        ret = NewQuasiquoted(parseTokens(tokenizer, false))
    } else if "," == token {
        ret = NewUnquoted(parseTokens(tokenizer, false))
    } else if n, nerr := strconv.Atoi(token); nil == nerr {
        ret = NewInteger(n)
    } else if isString(token) {
        return NewString(token[1 : len(token)-1])
    } else {
        ret = NewSymbol(token)
    }
    if asList {
        return NewCell(ret, parseTokens(tokenizer, true))
    }
    return ret
}

func Parse(input io.Reader) Expr {
    return parseTokens(Tokenize(input), false)
}

func isString(s string) bool {
    return strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"")
}
