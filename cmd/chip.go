package main

type Chip struct {
	Running    bool
	CycleDelay uint8

	CPU    ChipCPU
	Memory ChipMemory
	Timers ChipTimers
	Keypad ChipKeypad
	Video  ChipVideo
}

func SetupChip() Chip {
	chip := Chip{
		Running:    true,
		CycleDelay: 4,
	}
	SetupCPUTable(&chip)
	SetupKeypad(&chip)

	return chip
}
