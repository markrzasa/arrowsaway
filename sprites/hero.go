package sprites

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Hero struct {
	Sprite *Sprite
}

func NewHero(image *ebiten.Image) *Hero {
	h := Hero{
		Sprite: NewSprite(image.Bounds().Dx() / 3, image),
	}
	return &h
}

func (h *Hero) Update(gamepadIds *map[ebiten.GamepadID]bool, height, width int) {
	for id := range *gamepadIds {
		prevX := h.Sprite.X
		prevY := h.Sprite.Y
		x := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		x = math.Round(x*10) / 10
		y := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
		y = math.Round(y*10) / 10
		h.Sprite.X = h.Sprite.X + int(x*10)
		if h.Sprite.X < (h.Sprite.imageWidth / 2) {
			h.Sprite.X = h.Sprite.imageWidth / 2
		} else if h.Sprite.X > (width - (h.Sprite.imageWidth / 2)) {
			h.Sprite.X = width - (h.Sprite.imageWidth / 2)
		}
		h.Sprite.Y = h.Sprite.Y + int(y*10)
		if h.Sprite.Y < (h.Sprite.image.Bounds().Dy() / 2) {
			h.Sprite.Y = h.Sprite.image.Bounds().Dy() / 2
		} else if h.Sprite.Y > (height - (h.Sprite.image.Bounds().Dy() / 2)) {
			h.Sprite.Y = height - (h.Sprite.image.Bounds().Dy() / 2)
		}

		rightX := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickHorizontal)
		rightX = math.Round(rightX*10) / 10
		rightY := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickVertical)
		rightY = math.Round(rightY*10) / 10
		if rightX == 0 && rightY == 0 {
			h.Sprite.Radians = 0
		} else {
			rightX = math.Round(rightX*10) / 10
			rightX = float64(h.Sprite.X) + (rightX * 10)
			if int(rightX) < (h.Sprite.imageWidth / 2) {
				rightX = float64(h.Sprite.imageWidth) / 2
			} else if h.Sprite.X > (width - (h.Sprite.imageWidth / 2)) {
				rightX = float64(width - (h.Sprite.imageWidth / 2))
			}
	
			rightY = math.Round(rightY*10) / 10
			rightY = float64(h.Sprite.Y + int(rightY*10))
			if int(rightY) < (h.Sprite.image.Bounds().Dy() / 2) {
				rightY = float64(h.Sprite.image.Bounds().Dy() / 2)
			} else if h.Sprite.Y > (height - (h.Sprite.image.Bounds().Dy() / 2)) {
				rightY = float64(height - (h.Sprite.image.Bounds().Dy() / 2))
			}
	
			deltaX := int(rightX) - prevX
			deltaY := int(rightY) - prevY
			h.Sprite.Radians = math.Atan2(float64(deltaY), float64(deltaX)) - (math.Pi / 180)	
		}
	}
}

func (h *Hero) Draw(screen *ebiten.Image) {
	h.Sprite.Draw(screen, 0)
}

func (h *Hero) Winner(screen *ebiten.Image, width, height int) {
	h.Sprite.Scale(10)
	h.Sprite.X = width / 2
	h.Sprite.Y = height / 2
	h.Sprite.Draw(screen, 1)
}

func (h *Hero) GameOver(screen *ebiten.Image, width, height int) {
	h.Sprite.Scale(10)
	h.Sprite.X = width / 2
	h.Sprite.Y = height / 2
	h.Sprite.Draw(screen, 2)
}
