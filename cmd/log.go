package main

import (
	"fmt"
	"runtime"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func LogLn(format string, v ...any) {
	Print("[info]", 2, format+"\n", v...)
}

func LogOpLn(chip *Chip, format string, v ...any) {
	// fmt.Printf("%s> %s\n", strings.Repeat("-", int(chip.CPU.Stack.Pointer)), fmt.Sprintf(format, v...))
}

func Print(level string, depth int, format string, v ...any) {
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	file = file[(strings.LastIndex(file, "/") + 1):]

	fmt.Printf("[%s] %s: %s", fmt.Sprintf("%s:%v", file, line), level, fmt.Sprintf(format, v...))
}
