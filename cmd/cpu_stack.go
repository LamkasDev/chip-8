package main

type ChipCPUStackPointer uint8
type ChipCPUStack struct {
	Stack   [16]ChipPointer
	Pointer ChipCPUStackPointer
}
