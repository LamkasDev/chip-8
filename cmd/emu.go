package main

func main() {
	rom, err := SelectROM()
	check(err)

	chip := SetupChip()
	LoadROM(&chip, rom)

	renderer := Renderer{}
	SetupRenderer(&chip, &renderer)

	CycleGlobal(&chip, &renderer)
}
