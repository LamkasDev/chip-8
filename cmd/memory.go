package main

const ChipMemorySize = 4096
const ChipMemoryStartFontset = ChipPointer(0x50)
const ChipMemoryStartROM = ChipPointer(0x200)

type ChipMemory struct {
	Main [ChipMemorySize]byte
}
