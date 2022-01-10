package images

import (
	"bytes"
	_ "embed"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed arrow.png
var arrow []byte

//go:embed enemyHealth.png
var enemyHealth []byte

//go:embed stone.png
var stone []byte

//go:embed grass.png
var grass []byte

//go:embed goblin.png
var goblin []byte

//go:embed hero.png
var hero []byte

//go:embed life.png
var life []byte

//go:embed skeleton.png
var skeleton []byte

type Images struct {
	Arrow       *ebiten.Image
	EnemyHealth *ebiten.Image
	Goblin      *ebiten.Image
	Grass       *ebiten.Image
	Hero        *ebiten.Image
	Life        *ebiten.Image
	Skeleton    *ebiten.Image
	Stone       *ebiten.Image
}

var images *Images = nil

func newImage(imageBytes []byte) *ebiten.Image {
	image, err := png.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(image)
}

func GetImages() *Images {
	if images == nil {
		images = &Images{
			Arrow:       newImage(arrow),
			EnemyHealth: newImage(enemyHealth),
			Goblin:      newImage(goblin),
			Grass:       newImage(grass),
			Hero:        newImage(hero),
			Life:        newImage(life),
			Skeleton:    newImage(skeleton),
			Stone:       newImage(stone),
		}	
	}

	return images
}
