package main

const ChipVideoWidth = 64
const ChipVideoHeight = 32
const ChipVideoSize = ChipVideoWidth * ChipVideoHeight

type ChipVideo struct {
	Buffer [ChipVideoSize]uint32
}
