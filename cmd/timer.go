package main

type ChipTimers struct {
	Delay uint8
	Sound uint8
}

func CycleTimers(chip *Chip) {
	if chip.Timers.Delay > 0 {
		chip.Timers.Delay--
	}
	if chip.Timers.Sound > 0 {
		chip.Timers.Sound--
	}
}
