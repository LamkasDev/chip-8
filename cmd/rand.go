package main

import (
	"math"
	"math/rand"
)

func GetRandUint8() uint8 {
	return uint8(rand.Intn(math.MaxUint8 + 1))
}
