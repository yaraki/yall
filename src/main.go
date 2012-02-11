// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "bufio"
    "fmt"
    "os"
    "yall"
)

func prompt() {
    fmt.Print("yall> ")
}

func repl() {
    env := yall.NewEnv()
    for {
        prompt()
        reader := bufio.NewReader(os.Stdin)
        line, _, err := reader.ReadLine()
        if nil != err {
            return
        }
        func() {
            defer func() {
                if r := recover(); r != nil {
                    fmt.Println(r)
                }
            }()
            if result := env.EvalString(string(line)); result != nil {
                fmt.Println(result.String())
            }
        }()
    }
    fmt.Println()
}

func main() {
    repl()
}
