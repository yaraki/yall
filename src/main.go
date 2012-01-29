// Copyright 2011 Yuichi Araki. All rights reserved.
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
    fmt.Print("YALL> ")
}

func repl() {
    env := yall.NewEnv()
    for {
        prompt()
        reader := bufio.NewReader(os.Stdin)
        line, _, err := reader.ReadLine()
        if nil != err {
            break
        }
        if result := env.EvalString(string(line)); result != nil {
            fmt.Println(result.String())
        }
    }
    fmt.Println()
}

func main() {
    repl()
}
