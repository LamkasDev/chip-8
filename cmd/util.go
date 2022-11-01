package main

import "math"

func WrappingAdd(a uint8, b uint8) uint8 {
	return uint8((uint16(a) + uint16(b)) % (math.MaxUint8 + 1))
}
