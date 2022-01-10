package sprites

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	animations int = 50
)

type Arrow struct {
	m, b, radians    float64
	X, Y, EndX, EndY int
	xInc, yInc       int
	image            *ebiten.Image
}

func (a *Arrow) IsOffScreen(width, height int) bool {
	if a.X < 0 || a.X > width {
		return true
	}

	if a.Y < 0 || a.Y > height {
		return true
	}

	return false
}

func (a *Arrow) Update() {
	if a.xInc != 0 {
		a.X = a.X + a.xInc
		a.Y = int((a.m * float64(a.X)) + a.b)
	} else {
		a.Y = a.Y + a.yInc
		a.X = int((float64(a.Y) - a.b) / a.m)
	}
}

func (a *Arrow) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Rotate(a.radians)
	op.GeoM.Translate(float64(a.X), float64(a.Y))
	screen.DrawImage(a.image, op)
}

func NewArrow(startX, startY, endX, endY int, image *ebiten.Image) *Arrow {
	deltaY := endY - startY
	deltaX := endX - startX
	radians := math.Atan2(float64(deltaY), float64(deltaX)) - (45 * math.Pi/180)
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
	return &Arrow{
		m:       m,
		b:       float64(startY) - float64(m * float64(startX)),
		radians: radians,
		X:       startX,
		Y:       startY,
		EndX:    endX,
		EndY:    endY,
		xInc:    xIncrement,
		yInc:    yIncrement,
		image:   image,
	}
}
