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
