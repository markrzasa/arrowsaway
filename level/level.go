package level

import (
	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/markrzasa/arrowsaway/sprites"
)

const (
	numStages int = 1
)

type Level struct {
	name       string
	enemyImage *ebiten.Image
	numEnemies int
	stage      int
	bgImage    *ebiten.Image
}

func (l *Level) GetBackground() *ebiten.Image {
	return l.bgImage
}

func (l *Level) GetName() string {
	return l.name
}

func (l *Level) GetNumEnemies() int {
	return l.numEnemies
}

func (l *Level) GetStage() int {
	return l.stage
}

func (l *Level) Complete() bool {
	return l.stage == numStages
}

func (l *Level) NextStage() {
	l.stage = l.stage + 1
}

func (l *Level) Reset() {
	l.stage = 0
}

func (l *Level) getHitpoints(i int) int {
	hp := sprites.HitpointIncrement
	switch l.stage {
	case 0:
	case 1:
		hp = hp + (hp * (i % 2))
	case 2:
		hp = hp + (hp * (i % 3))
	default:
		hp = hp + (hp * (i % 4))
	}

	return hp
}

func (l *Level) PopulateEnemies(width, height int, enemies map[string]*sprites.Enemy) {
	if l.stage == numStages {
		enemies[uuid.NewString()] = sprites.NewEnemy(0, 0, 1000, true, l.enemyImage)
	} else {
		enemiesPerSide := l.GetNumEnemies() / 4
		for i := 0 ; i < enemiesPerSide ; i++ {
			x := 0
			y := (i * (height / enemiesPerSide))
			enemies[uuid.New().String()] = sprites.NewEnemy(x, y, l.getHitpoints(i), false, l.enemyImage)
		}
 		for i := 0 ; i < enemiesPerSide ; i++ {
			x := (i * (width / enemiesPerSide))
			y := 0
			enemies[uuid.New().String()] = sprites.NewEnemy(x, y, l.getHitpoints(i), false, l.enemyImage)
		}
		for i := 0 ; i < enemiesPerSide ; i++ {
			x := width - (l.enemyImage.Bounds().Dx() / 5)
			y := (i * (height / enemiesPerSide))
			enemies[uuid.New().String()] = sprites.NewEnemy(x, y, l.getHitpoints(i), false, l.enemyImage)
		}
		for i := 0 ; i < enemiesPerSide ; i++ {
			x := (i * (width / enemiesPerSide))
			y := height - l.enemyImage.Bounds().Dy()
			enemies[uuid.New().String()] = sprites.NewEnemy(x, y, l.getHitpoints(i), false, l.enemyImage)
		}	
 	}
}

func NewLevel(name string, enemyImage *ebiten.Image,bgImage *ebiten.Image, numEnemies int) *Level {
	return &Level{
		name:       name,
		bgImage:    bgImage,
		enemyImage: enemyImage,
		numEnemies: numEnemies,
		stage:      0,
	}
}
