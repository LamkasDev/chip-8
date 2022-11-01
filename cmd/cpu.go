package main

type ChipPointer uint16
type ChipInstruction uint16
type ChipCPU struct {
	Reg     [16]uint8
	Op      ChipInstruction
	Index   ChipPointer
	Counter ChipPointer
	Stack   ChipCPUStack
	Table   CPUTable
}

func CycleCPU(chip *Chip) {
	chip.CPU.Op = ChipInstruction(uint16(chip.Memory.Main[chip.CPU.Counter])<<8 | uint16(chip.Memory.Main[chip.CPU.Counter+1]))
	chip.CPU.Counter += 2

	ProcessOp(chip)
}
