package main

import (
	"io"
	"os"
)

func LoadROM(chip *Chip, path string) {
	f, err := os.Open(path)
	check(err)
	defer f.Close()
	n, err := f.Read(chip.Memory.Main[ChipMemoryStartROM:])
	if err != io.EOF {
		check(err)
	}
	chip.CPU.Counter = ChipMemoryStartROM

	LogLn("Loaded %d bytes into memory...", n)
}
