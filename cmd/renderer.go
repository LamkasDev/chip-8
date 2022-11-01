package main

import (
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type Renderer struct {
	Window   *sdl.Window
	Surface  *sdl.Surface
	Renderer *sdl.Renderer
	Texture  *sdl.Texture
	Pitch    int
}

func SetupRenderer(chip *Chip, renderer *Renderer) {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	check(err)

	w, h := int32(ChipVideoWidth*10), int32(ChipVideoHeight*10)
	renderer.Window, err = sdl.CreateWindow("Chip-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, w, h, sdl.WINDOW_SHOWN)
	check(err)

	renderer.Surface, err = renderer.Window.GetSurface()
	check(err)

	renderer.Renderer, err = sdl.CreateRenderer(renderer.Window, -1, sdl.RENDERER_ACCELERATED)
	check(err)

	renderer.Texture, err = renderer.Renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, ChipVideoWidth, ChipVideoHeight)
	check(err)

	renderer.Pitch = int(unsafe.Sizeof(chip.Video.Buffer[0]) * ChipVideoWidth)
}

func CycleRenderer(chip *Chip, renderer *Renderer) {
	renderer.Texture.Update(nil, unsafe.Pointer(&chip.Video.Buffer), renderer.Pitch)
	renderer.Renderer.Clear()
	renderer.Renderer.Copy(renderer.Texture, nil, nil)
	renderer.Renderer.Present()
}

func CleanRenderer(renderer *Renderer) {
	renderer.Texture.Destroy()
	renderer.Renderer.Destroy()
	renderer.Window.Destroy()
	sdl.Quit()
}
