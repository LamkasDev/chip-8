package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

const KeypadSize = 16

type ChipKeypad struct {
	States  [KeypadSize]bool
	Mapping [KeypadSize]int
}

func SetupKeypad(chip *Chip) {
	chip.Keypad.Mapping = [KeypadSize]int{
		sdl.SCANCODE_1, // 1
		sdl.SCANCODE_2, // 2
		sdl.SCANCODE_3, // 3
		sdl.SCANCODE_4, // C

		sdl.SCANCODE_Q, // 4
		sdl.SCANCODE_W, // 5
		sdl.SCANCODE_E, // 6
		sdl.SCANCODE_R, // D

		sdl.SCANCODE_A, // 7
		sdl.SCANCODE_S, // 8
		sdl.SCANCODE_D, // 9
		sdl.SCANCODE_F, // E

		sdl.SCANCODE_Y, // A
		sdl.SCANCODE_X, // 0
		sdl.SCANCODE_C, // B
		sdl.SCANCODE_V, // F
	}
}
