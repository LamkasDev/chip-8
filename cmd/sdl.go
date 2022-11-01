package main

import "github.com/veandco/go-sdl2/sdl"

func CycleSDL(chip *Chip) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch t := event.(type) {
		case *sdl.QuitEvent:
			chip.Running = false
		case *sdl.KeyboardEvent:
			switch t.Keysym.Scancode {
			case sdl.SCANCODE_ESCAPE:
				if t.Type == sdl.KEYDOWN {
					chip.Running = false
				}
			case sdl.SCANCODE_F1:
				if t.Type == sdl.KEYDOWN {
					if chip.CycleDelay > 10 {
						chip.CycleDelay -= 10
					} else {
						chip.CycleDelay = 1
					}
				}
			case sdl.SCANCODE_F2:
				if t.Type == sdl.KEYDOWN {
					chip.CycleDelay += 10
				}
			default:
				for i, key := range chip.Keypad.Mapping {
					if key == int(t.Keysym.Scancode) {
						if t.Type == sdl.KEYDOWN {
							chip.Keypad.States[i] = true
						} else {
							chip.Keypad.States[i] = false
						}
					}
				}
			}
		}
	}
}
