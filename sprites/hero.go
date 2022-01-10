package sprites

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Hero struct {
	X, Y, imageWidth int

	image *ebiten.Image
}

func (h *Hero) Update(gamepadIds *map[ebiten.GamepadID]bool, height, width int) {
	for id := range *gamepadIds {
		x := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		x = math.Round(x * 10) / 10
		y := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
		y = math.Round(y * 10) / 10
		h.X = h.X + int(x * 10)
		if h.X < (h.imageWidth / 2) {
			h.X = h.imageWidth / 2
		} else if h.X > (width - (h.imageWidth / 2)) {
			h.X = width - (h.imageWidth / 2)
		}
		h.Y = h.Y + int(y * 10)
		if h.Y < (h.image.Bounds().Dy() / 2) {
			h.Y = h.image.Bounds().Dy() / 2
		} else if h.Y > (height - (h.image.Bounds().Dy() / 2)) {
			h.Y = height - (h.image.Bounds().Dy() / 2)
		}
	}
}

func (h *Hero) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(h.X - (h.imageWidth / 2)), float64(h.Y - (h.image.Bounds().Dy() / 2)))
	subImageRect := image.Rect(0, 0, h.imageWidth, h.image.Bounds().Dy())
	subImage := h.image.SubImage(subImageRect).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}

func (h *Hero) Bounds() *image.Rectangle {
	return &image.Rectangle{
		Min: image.Point{
			X: h.X,
			Y: h.Y,
		},
		Max: image.Point{
			X: h.X + imageWidth,
			Y: h.Y + h.image.Bounds().Dy(),
		},
	}
}

func (h *Hero) Center(width, height int) {
	h.X = (width / 2) - (h.imageWidth / 2)
	h.Y = (height / 2) - (h.image.Bounds().Dy() / 2)
}

func (h *Hero) Winner(screen *ebiten.Image, width, height int) {
	scale := 10
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(scale), float64(scale))
	x := float64((width / 2) - (h.imageWidth * (scale / 2)))
	y := float64((height / 2) - (h.image.Bounds().Dy() * (scale / 2)))
	op.GeoM.Translate(x, y)
	subImageRect := image.Rect(h.imageWidth * 2, 0, h.image.Bounds().Dx(), h.image.Bounds().Dy())
	subImage := h.image.SubImage(subImageRect).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}

func (h *Hero) GameOver(screen *ebiten.Image, width, height int) {
	scale := 10
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(scale), float64(scale))
	x := float64((width / 2) - (h.imageWidth * (scale / 2)))
	y := float64((height / 2) - (h.image.Bounds().Dy() * (scale / 2)))
	op.GeoM.Translate(x, y)
	subImageRect := image.Rect(h.imageWidth, 0, h.imageWidth * 2, h.image.Bounds().Dy())
	subImage := h.image.SubImage(subImageRect).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}

func NewHero(width, height int, image *ebiten.Image) *Hero {
	h := Hero{
		imageWidth: image.Bounds().Dx() / 3,
		image: image,
	}
	h.Center(width, height)
	return &h
}
