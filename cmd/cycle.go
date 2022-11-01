package main

import (
	"math"

	"github.com/loov/hrtime"
)

func CycleGlobal(chip *Chip, renderer *Renderer) {
	defer CleanRenderer(renderer)

	lastCycle := hrtime.Now()
	for chip.Running {
		currentCycle := hrtime.Now()
		diffCycle := (currentCycle - lastCycle).Milliseconds()
		neededCycles := int(math.Floor(float64(diffCycle) / float64(chip.CycleDelay)))
		if neededCycles > 0 {
			lastCycle = currentCycle
		}
		for i := 0; i < neededCycles; i++ {
			CycleCPU(chip)
			CycleTimers(chip)
			CycleRenderer(chip, renderer)
		}

		CycleSDL(chip)
	}
}
