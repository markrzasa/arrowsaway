package sprites

import (
	"math"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/markrzasa/arrowsaway/images"
)

const (
	animations int = 50
)

type Arrow struct {
	Id         string
	m, b       float64
	EndX, EndY int
	xInc, yInc  int
	Sprite      Sprite
}

func (a *Arrow) IsOffScreen(width, height int) bool {
	if a.Sprite.X < 0 || a.Sprite.X > width {
		return true
	}

	if a.Sprite.Y < 0 || a.Sprite.Y > height {
		return true
	}

	return false
}

func (a *Arrow) Update() {
	if a.xInc != 0 {
		a.Sprite.X = a.Sprite.X + a.xInc
		a.Sprite.Y = int((a.m * float64(a.Sprite.X)) + a.b)
	} else {
		a.Sprite.Y = a.Sprite.Y + a.yInc
		a.Sprite.X = int((float64(a.Sprite.Y) - a.b) / a.m)
	}
}

func (a *Arrow) Draw(screen *ebiten.Image) {
	a.Sprite.Draw(screen, 0)
}

func NewArrow(startX, startY, endX, endY int) *Arrow {
	deltaY := endY - startY
	deltaX := endX - startX
	m := float64(deltaY) / float64(deltaX)
	xIncrement := 0
	if startX < endX {
		xIncrement = (endX - startX) / animations
	} else {
		xIncrement = -1 * ((startX - endX) / animations)
	}
	yIncrement := 0
	if startY < endY {
		yIncrement = (endY - startY) / animations
	} else {
		yIncrement = -1 * ((startY - endY) / animations)
	}
	arrowImage := images.GetImages().Arrow
	arrow := &Arrow{
		Id:      uuid.New().String(),
		m:       m,
		b:       float64(startY) - float64(m * float64(startX)),
		EndX:    endX,
		EndY:    endY,
		xInc:    xIncrement,
		yInc:    yIncrement,
		Sprite:  *NewSprite(arrowImage.Bounds().Dx(), arrowImage),
	}
	arrow.Sprite.X = startX
	arrow.Sprite.Y = startY
	arrow.Sprite.Radians = math.Atan2(float64(deltaY), float64(deltaX)) - (45 * math.Pi/180)
	return arrow
}
