package main

import (
	"reflect"
	"runtime"
	"strings"
)

const CPUTableSize = 0xF + 1
const CPUTable0Size = 0xE + 1
const CPUTable8Size = 0xE + 1
const CPUTableESize = 0xE + 1
const CPUTableFSize = 0x65 + 1

type CPUTableOp func(chip *Chip)
type CPUTable struct {
	Table  [CPUTableSize]CPUTableOp
	Table0 [CPUTable0Size]CPUTableOp
	Table8 [CPUTable8Size]CPUTableOp
	TableE [CPUTableESize]CPUTableOp
	TableF [CPUTableFSize]CPUTableOp
}

func ProcessTable0(chip *Chip) {
	chip.CPU.Table.Table0[chip.CPU.Op&0x000F](chip)
}

func ProcessTable8(chip *Chip) {
	chip.CPU.Table.Table8[chip.CPU.Op&0x000F](chip)
}

func ProcessTableE(chip *Chip) {
	chip.CPU.Table.TableE[chip.CPU.Op&0x000F](chip)
}

func ProcessTableF(chip *Chip) {
	chip.CPU.Table.TableF[chip.CPU.Op&0x00FF](chip)
}

func ProcessOpNull(chip *Chip) {
}

func GetTableIndex(op ChipInstruction) uint8 {
	return uint8((op & 0xF000) >> 12)
}

func GetOpName(chip *Chip, op ChipInstruction) string {
	fn := chip.CPU.Table.Table[GetTableIndex(op)]
	name := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	switch name {
	case "main.ProcessTable0":
		fn = chip.CPU.Table.Table0[op&0x000F]
	case "main.ProcessTable8":
		fn = chip.CPU.Table.Table8[op&0x000F]
	case "main.ProcessTableE":
		fn = chip.CPU.Table.TableE[op&0x000F]
	case "main.ProcessTableF":
		fn = chip.CPU.Table.TableF[op&0x00FF]
	}
	return strings.Replace(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), "main.", "", -1)
}

func SetupCPUTable(chip *Chip) {
	chip.CPU.Table.Table[0x0] = ProcessTable0
	chip.CPU.Table.Table[0x1] = ProcessOp1NNN
	chip.CPU.Table.Table[0x2] = ProcessOp2NNN
	chip.CPU.Table.Table[0x3] = ProcessOp3xkk
	chip.CPU.Table.Table[0x4] = ProcessOp4xkk
	chip.CPU.Table.Table[0x5] = ProcessOp5xy0
	chip.CPU.Table.Table[0x6] = ProcessOp6xkk
	chip.CPU.Table.Table[0x7] = ProcessOp7xkk
	chip.CPU.Table.Table[0x8] = ProcessTable8
	chip.CPU.Table.Table[0x9] = ProcessOp9xy0
	chip.CPU.Table.Table[0xA] = ProcessOpAnnn
	chip.CPU.Table.Table[0xB] = ProcessOpBnnn
	chip.CPU.Table.Table[0xC] = ProcessOpCxkk
	chip.CPU.Table.Table[0xD] = ProcessOpDxyn
	chip.CPU.Table.Table[0xE] = ProcessTableE
	chip.CPU.Table.Table[0xF] = ProcessTableF

	for i := 0; i < CPUTable0Size; i++ {
		chip.CPU.Table.Table0[i] = ProcessOpNull
		chip.CPU.Table.Table8[i] = ProcessOpNull
		chip.CPU.Table.TableE[i] = ProcessOpNull
	}

	chip.CPU.Table.Table0[0x0] = ProcessOp00E0
	chip.CPU.Table.Table0[0xE] = ProcessOp00EE

	chip.CPU.Table.Table8[0x0] = ProcessOp8xy0
	chip.CPU.Table.Table8[0x1] = ProcessOp8xy1
	chip.CPU.Table.Table8[0x2] = ProcessOp8xy2
	chip.CPU.Table.Table8[0x3] = ProcessOp8xy3
	chip.CPU.Table.Table8[0x4] = ProcessOp8xy4
	chip.CPU.Table.Table8[0x5] = ProcessOp8xy5
	chip.CPU.Table.Table8[0x6] = ProcessOp8xy6
	chip.CPU.Table.Table8[0x7] = ProcessOp8xy7
	chip.CPU.Table.Table8[0xE] = ProcessOp8xyE

	chip.CPU.Table.TableE[0x1] = ProcessOpExA1
	chip.CPU.Table.TableE[0xE] = ProcessOpEx9E

	for i := 0; i < CPUTableFSize; i++ {
		chip.CPU.Table.TableF[i] = ProcessOpNull
	}

	chip.CPU.Table.TableF[0x07] = ProcessOpFx07
	chip.CPU.Table.TableF[0x0A] = ProcessOpFx0A
	chip.CPU.Table.TableF[0x15] = ProcessOpFx15
	chip.CPU.Table.TableF[0x18] = ProcessOpFx18
	chip.CPU.Table.TableF[0x1E] = ProcessOpFx1E
	chip.CPU.Table.TableF[0x29] = ProcessOpFx29
	chip.CPU.Table.TableF[0x33] = ProcessOpFx33
	chip.CPU.Table.TableF[0x55] = ProcessOpFx55
	chip.CPU.Table.TableF[0x65] = ProcessOpFx65
}

func ProcessOp(chip *Chip) {
	tableIndex := GetTableIndex(chip.CPU.Op)
	LogOpLn(chip, "%s (raw: %d, at: 0x%02x) | Counter: %d", GetOpName(chip, chip.CPU.Op), chip.CPU.Op, tableIndex, chip.CPU.Counter)
	chip.CPU.Table.Table[tableIndex](chip)
}
