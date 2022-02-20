package sprites

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	X, Y, imageWidth int

	image *ebiten.Image

	Radians, ScaleX, ScaleY float64
}

func NewSprite(imageWidth int, image *ebiten.Image) *Sprite {
	return &Sprite{
		X:          0,
		Y:          0,
		ScaleX:     1,
		ScaleY:     1,
		imageWidth: imageWidth,
		image:      image,
		Radians:    0,
	}
}

func (s *Sprite) Scale(scale float64) {
	s.ScaleX = scale
	s.ScaleY = scale
}

func (s *Sprite) Draw(screen *ebiten.Image, frame int) {
	bounds := s.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(bounds.Dx()) / 2, -float64(bounds.Dy()) / 2)
	op.GeoM.Rotate(s.Radians)
	op.GeoM.Scale(s.ScaleX, float64(s.ScaleY))
	op.GeoM.Translate(float64(s.X), float64(s.Y))
	subImageRect := image.Rect(frame*s.imageWidth, 0, (frame+1)*s.imageWidth, s.image.Bounds().Dy())
	subImage := s.image.SubImage(subImageRect).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}

func (s *Sprite) DrawSubImage(screen *ebiten.Image, subImageWidth int) {
	bounds := s.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(subImageWidth) / 2, -float64(bounds.Dy()) / 2)
	op.GeoM.Rotate(s.Radians)
	op.GeoM.Scale(s.ScaleX, float64(s.ScaleY))
	op.GeoM.Translate(float64(s.X), float64(s.Y))
	subImageRect := image.Rect(0, 0, subImageWidth, s.image.Bounds().Dy())
	subImage := s.image.SubImage(subImageRect).(*ebiten.Image)
	screen.DrawImage(subImage, op)
}

func (s *Sprite) Bounds() *image.Rectangle {
	return &image.Rectangle{
		Min: image.Point{
			X: s.X - (s.imageWidth / 2),
			Y: s.Y - (s.image.Bounds().Dy() / 2),
		},
		Max: image.Point{
			X: s.X + (s.imageWidth / 2),
			Y: s.Y + (s.image.Bounds().Dy() / 2),
		},
	}
}

func (s *Sprite) ScaledBounds() *image.Rectangle {
	scaledWidth := float64(s.imageWidth) * s.ScaleX
	scaledHeight := float64(s.image.Bounds().Dy()) * s.ScaleY
	return &image.Rectangle{
		Min: image.Point{
			X: int(float64(s.X) - (scaledWidth / 2)),
			Y: int(float64(s.Y) - (scaledHeight / 2)),
		},
		Max: image.Point{
			X: int(float64(s.X) + (scaledWidth / 2)),
			Y: int(float64(s.Y) + (scaledHeight / 2)),
		},
	}
}

func (s *Sprite) Intersect(o *Sprite) bool {
	return s.ScaledBounds().Intersect(*o.ScaledBounds()) != image.Rectangle{}
}

func (s *Sprite) Center(width, height int) {
	s.X = (width / 2)
	s.Y = (height / 2)
}
