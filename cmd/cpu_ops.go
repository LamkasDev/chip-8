package main

import "math"

func GetAddress(chip *Chip) ChipPointer {
	return ChipPointer(chip.CPU.Op & 0x0FFF)
}

func GetVX(chip *Chip) uint8 {
	return uint8((chip.CPU.Op & 0x0F00) >> 8)
}

func GetVY(chip *Chip) uint8 {
	return uint8((chip.CPU.Op & 0x00F0) >> 4)
}

func GetByte(chip *Chip) uint8 {
	return uint8(chip.CPU.Op & 0x00FF)
}

// CLS (Clear the display.)
func ProcessOp00E0(chip *Chip) {
	chip.Video.Buffer = [ChipVideoWidth * ChipVideoHeight]uint32{}
}

// Ret (Return from a subroutine.)
func ProcessOp00EE(chip *Chip) {
	chip.CPU.Stack.Pointer--
	chip.CPU.Counter = chip.CPU.Stack.Stack[chip.CPU.Stack.Pointer]
}

// JP addr (Jump to location nnn.)
func ProcessOp1NNN(chip *Chip) {
	chip.CPU.Counter = GetAddress(chip)
	LogOpLn(chip, "Jump to %v", chip.CPU.Counter)
}

// CALL addr (Call subroutine at nnn.)
func ProcessOp2NNN(chip *Chip) {
	address := GetAddress(chip)
	chip.CPU.Stack.Stack[chip.CPU.Stack.Pointer] = chip.CPU.Counter
	chip.CPU.Stack.Pointer++
	chip.CPU.Counter = address
	LogOpLn(chip, "Call to %v", chip.CPU.Counter)
}

// 3xkk - SE Vx, byte (Skip next instruction if Vx = kk.)
func ProcessOp3xkk(chip *Chip) {
	vx := GetVX(chip)
	b := GetByte(chip)

	if chip.CPU.Reg[vx] == b {
		chip.CPU.Counter += 2
	}
	LogOpLn(chip, "Skip next if (r_%v) %v == %v", vx, chip.CPU.Reg[vx], b)
}

// 4xkk - SNE Vx, byte (Skip next instruction if Vx != kk.)
func ProcessOp4xkk(chip *Chip) {
	vx := GetVX(chip)
	b := GetByte(chip)

	if chip.CPU.Reg[vx] != b {
		chip.CPU.Counter += 2
	}
	LogOpLn(chip, "Skip next if (r_%v) %v != %v", vx, chip.CPU.Reg[vx], b)
}

// 5xy0 - SE Vx, Vy (Skip next instruction if Vx = Vy.)
func ProcessOp5xy0(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)

	if chip.CPU.Reg[vx] == chip.CPU.Reg[vy] {
		chip.CPU.Counter += 2
	}
	LogOpLn(chip, "Skip next if (r_%v) %v == %v (r_%v)", vx, chip.CPU.Reg[vx], chip.CPU.Reg[vy], vy)
}

// 6xkk - LD Vx, byte (Set Vx = kk.)
func ProcessOp6xkk(chip *Chip) {
	vx := GetVX(chip)
	b := GetByte(chip)

	chip.CPU.Reg[vx] = b
	LogOpLn(chip, "Set r_%v to %v", vx, chip.CPU.Reg[vx])
}

// 7xkk - ADD Vx, byte (Set Vx = Vx + kk.)
func ProcessOp7xkk(chip *Chip) {
	vx := GetVX(chip)
	b := GetByte(chip)
	res := WrappingAdd(chip.CPU.Reg[vx], b)

	chip.CPU.Reg[vx] = res
	LogOpLn(chip, "Add %v to r_%v (now: %v)", b, vx, chip.CPU.Reg[vx])
}

// 8xy0 - LD Vx, Vy (Set Vx = Vy.)
func ProcessOp8xy0(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)

	chip.CPU.Reg[vx] = chip.CPU.Reg[vy]
	LogOpLn(chip, "Set r_%v to %v", vx, chip.CPU.Reg[vx])
}

// 8xy1 - OR Vx, Vy (Set Vx = Vx OR Vy.)
func ProcessOp8xy1(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)

	chip.CPU.Reg[vx] |= chip.CPU.Reg[vy]
}

// 8xy2 - AND Vx, Vy (Set Vx = Vx AND Vy.)
func ProcessOp8xy2(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)

	chip.CPU.Reg[vx] &= chip.CPU.Reg[vy]
}

// 8xy3 - XOR Vx, Vy (Set Vx = Vx XOR Vy.)
func ProcessOp8xy3(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)

	chip.CPU.Reg[vx] ^= chip.CPU.Reg[vy]
}

// 8xy4 - ADD Vx, Vy (Set Vx = Vx + Vy, set VF = carry.)
func ProcessOp8xy4(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)
	sum := uint16(chip.CPU.Reg[vx]) + uint16(chip.CPU.Reg[vy])

	if sum > math.MaxUint8 {
		chip.CPU.Reg[0xF] = 1
		sum -= 256
	} else {
		chip.CPU.Reg[0xF] = 0
	}
	chip.CPU.Reg[vx] = uint8(sum)
	LogOpLn(chip, "Add %v to r_%v (now: %v)", chip.CPU.Reg[vy], vx, chip.CPU.Reg[vx])
}

// 8xy5 - SUB Vx, Vy (Set Vx = Vx - Vy, set VF = NOT borrow.)
func ProcessOp8xy5(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)
	sub := int16(chip.CPU.Reg[vx]) - int16(chip.CPU.Reg[vy])

	if sub < 0 {
		chip.CPU.Reg[0xF] = 1
		sub += 256
	} else {
		chip.CPU.Reg[0xF] = 0
	}
	chip.CPU.Reg[vx] = uint8(sub)
	LogOpLn(chip, "Substract %v from r_%v (now: %v)", chip.CPU.Reg[vy], vx, chip.CPU.Reg[vx])
}

// 8xy6 - SHR Vx (Set Vx = Vx SHR 1.)
func ProcessOp8xy6(chip *Chip) {
	vx := GetVX(chip)

	chip.CPU.Reg[0xF] = (chip.CPU.Reg[vx] & 0x1)
	chip.CPU.Reg[vx] >>= 1
}

// 8xy7 - SUBN Vx, Vy (Set Vx = Vy - Vx, set VF = NOT borrow.)
func ProcessOp8xy7(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)
	sub := int16(chip.CPU.Reg[vy]) - int16(chip.CPU.Reg[vx])

	if sub < 0 {
		chip.CPU.Reg[0xF] = 1
		sub += 256
	} else {
		chip.CPU.Reg[0xF] = 0
	}
	chip.CPU.Reg[vx] = uint8(sub)
	LogOpLn(chip, "Substract (r_%v) %v from r_%v (now: %v)", vy, chip.CPU.Reg[vy], vx, chip.CPU.Reg[vx])
}

// 8xyE - SHL Vx {, Vy} (Set Vx = Vx SHL 1.)
func ProcessOp8xyE(chip *Chip) {
	vx := GetVX(chip)
	res := int16(chip.CPU.Reg[vx]) << 1

	if res > math.MaxUint8 {
		chip.CPU.Reg[0xF] = 1
		res -= 256
	} else {
		chip.CPU.Reg[0xF] = 0
	}
	chip.CPU.Reg[vx] = uint8(res)
}

// 9xy0 - SNE Vx, Vy (Skip next instruction if Vx != Vy.)
func ProcessOp9xy0(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)

	if chip.CPU.Reg[vx] != chip.CPU.Reg[vy] {
		chip.CPU.Counter += 2
	}
	LogOpLn(chip, "Skip next if (r_%v) %v != %v (r_%v)", vx, chip.CPU.Reg[vx], chip.CPU.Reg[vy], vy)
}

// Annn - LD I, addr (Set I = nnn.)
func ProcessOpAnnn(chip *Chip) {
	chip.CPU.Index = GetAddress(chip)
	LogOpLn(chip, "Set index to %v", chip.CPU.Index)
}

// Bnnn - JP V0, addr (Jump to location nnn + V0.)
func ProcessOpBnnn(chip *Chip) {
	address := GetAddress(chip)

	chip.CPU.Counter = ChipPointer(chip.CPU.Reg[0]) + address
	LogOpLn(chip, "Jump to %v", chip.CPU.Counter)
}

// Cxkk - RND Vx, byte (Set Vx = random byte AND kk.)
func ProcessOpCxkk(chip *Chip) {
	vx := GetVX(chip)
	b := GetByte(chip)

	chip.CPU.Reg[vx] = GetRandUint8() & b
	LogOpLn(chip, "Set r_%v to %v", vx, chip.CPU.Reg[vx])
}

// Dxyn - DRW Vx, Vy, nibble (Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.)
func ProcessOpDxyn(chip *Chip) {
	vx := GetVX(chip)
	vy := GetVY(chip)
	height := uint8(chip.CPU.Op & 0x000F)

	// Wrap if going beyond screen boundaries
	x := chip.CPU.Reg[vx] % ChipVideoWidth
	y := chip.CPU.Reg[vy] % ChipVideoHeight
	LogOpLn(chip, "vx=%v, vy=%v, h=%v, x=%v, y=%v", vx, vy, height, x, y)

	chip.CPU.Reg[0xF] = 0
	for row := uint8(0); row < height; row++ {
		spB := uint8(chip.Memory.Main[chip.CPU.Index+ChipPointer(row)])
		for col := uint8(0); col < 8; col++ {
			pixelSp := spB & (0x80 >> col)
			videoI := (uint16(y)+uint16(row))*uint16(ChipVideoWidth) + (uint16(x) + uint16(col))
			if videoI >= ChipVideoSize {
				break
			}

			// Sprite pixel is on
			pixelSc := &chip.Video.Buffer[videoI]
			if pixelSp > 0 {
				// Screen pixel also on - collision
				if *pixelSc == 0xFFFFFFFF {
					chip.CPU.Reg[0xF] = 1
				}

				// Effectively XOR with the sprite pixel
				*pixelSc ^= 0xFFFFFFFF
			}
		}
	}
}

// Ex9E - SKP Vx (Skip next instruction if key with the value of Vx is pressed.)
func ProcessOpEx9E(chip *Chip) {
	vx := GetVX(chip)
	key := chip.CPU.Reg[vx]

	if chip.Keypad.States[key] {
		chip.CPU.Counter += 2
		LogOpLn(chip, "Skipped because %v was pressed.", key)
		return
	}
	LogOpLn(chip, "Continuing because %v wasn't pressed.", key)
}

// ExA1 - SKNP Vx (Skip next instruction if key with the value of Vx is not pressed.)
func ProcessOpExA1(chip *Chip) {
	vx := GetVX(chip)
	key := chip.CPU.Reg[vx]

	if !chip.Keypad.States[key] {
		chip.CPU.Counter += 2
		LogOpLn(chip, "Skipped because %v wasn't pressed.", key)
	}
	LogOpLn(chip, "Continuing because %v was pressed.", key)
}

// Fx07 - LD Vx, DT (Set Vx = delay timer value.)
func ProcessOpFx07(chip *Chip) {
	vx := GetVX(chip)

	chip.CPU.Reg[vx] = chip.Timers.Delay
	LogOpLn(chip, "Set r_%v to %v", vx, chip.Timers.Delay)
}

// Fx0A - LD Vx, K (Wait for a key press, store the value of the key in Vx.)
func ProcessOpFx0A(chip *Chip) {
	vx := GetVX(chip)

	for key, pressed := range chip.Keypad.States {
		if pressed {
			chip.CPU.Reg[vx] = uint8(key)
			LogOpLn(chip, "Set r_%v to %v", vx, chip.CPU.Reg[vx])
			return
		}
	}
	chip.CPU.Counter -= 2
}

// Fx15 - LD DT, Vx (Set delay timer = Vx.)
func ProcessOpFx15(chip *Chip) {
	vx := GetVX(chip)

	chip.Timers.Delay = chip.CPU.Reg[vx]
	LogOpLn(chip, "Set delay to %v (r_%v)", chip.Timers.Delay, vx)
}

// Fx18 - LD ST, Vx (Set sound timer = Vx.)
func ProcessOpFx18(chip *Chip) {
	vx := GetVX(chip)

	chip.Timers.Sound = chip.CPU.Reg[vx]
	LogOpLn(chip, "Set sound to %v (r_%v)", chip.Timers.Sound, vx)
}

// Fx1E - ADD I, Vx (Set I = I + Vx.)
func ProcessOpFx1E(chip *Chip) {
	vx := GetVX(chip)

	chip.CPU.Index += ChipPointer(chip.CPU.Reg[vx])
	LogOpLn(chip, "Add (r_%v) %v to index (now: %v)", vx, chip.CPU.Reg[vx], chip.CPU.Index)
}

// Fx29 - LD F, Vx (Set I = location of sprite for digit Vx.)
func ProcessOpFx29(chip *Chip) {
	vx := GetVX(chip)
	digit := chip.CPU.Reg[vx]

	chip.CPU.Index = ChipMemoryStartFontset + ChipPointer(FontsetCharacterSize*digit)
	LogOpLn(chip, "Set %v to index", chip.CPU.Index)
}

// Fx33 - LD B, Vx (Store BCD representation of Vx in memory locations I, I+1, and I+2.)
func ProcessOpFx33(chip *Chip) {
	vx := GetVX(chip)
	value := chip.CPU.Reg[vx]

	// Ones-place
	chip.Memory.Main[chip.CPU.Index+2] = value % 10
	value /= 10

	// Tens-place
	chip.Memory.Main[chip.CPU.Index+1] = value % 10
	value /= 10

	// Hundreds-place
	chip.Memory.Main[chip.CPU.Index] = value % 10
}

// Fx55 - LD [I], Vx (Store registers V0 through Vx in memory starting at location I.)
func ProcessOpFx55(chip *Chip) {
	vx := GetVX(chip)

	for i := uint8(0); i <= vx; i++ {
		chip.Memory.Main[chip.CPU.Index+ChipPointer(i)] = chip.CPU.Reg[i]
	}
	LogOpLn(chip, "Stored %v registers in memory", vx)
}

// Fx65 - LD Vx, [I] (Read registers V0 through Vx from memory starting at location I.)
func ProcessOpFx65(chip *Chip) {
	vx := GetVX(chip)

	for i := uint8(0); i <= vx; i++ {
		chip.CPU.Reg[i] = chip.Memory.Main[chip.CPU.Index+ChipPointer(i)]
	}
	LogOpLn(chip, "Restored %v registers from memory", vx)
}
