// Copyright 2012 Yuichi Araki. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/yaraki/yall"
	"os"
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

func loadFiles() {
	env := yall.NewEnv()
	for i := 0; i < flag.NArg(); i++ {
		file, err := os.Open(flag.Arg(i))
		if file == nil {
			fmt.Fprintf(os.Stderr, "Can't open %s: error %s\n",
				flag.Arg(i), err)
			os.Exit(1)
		}
		defer file.Close()
		env.Load(file)
	}
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		repl()
	} else {
		loadFiles()
	}
}
