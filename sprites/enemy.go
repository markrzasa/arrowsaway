package sprites

import (
	"image"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/markrzasa/arrowsaway/images"
)

const (
	healthMargin int     = 2
	imageWidth   int     = 32
	scaleFactor  float64 = 0.25

	HitpointIncrement int = 50
)

type enemyState int

const (
	Alive enemyState = iota
	Dead
	Buried
)

type Enemy struct {
	startX, startY, x, y, frame int
	state                       enemyState
	image                       *ebiten.Image
	stateTime                   int64
	hitpoints, totalHitpoints   int
	boss                        bool
}

func (e *Enemy) getScale() float64 {
	scale := 1.0
	if e.state == Alive {
		if e.boss {
			scale = math.Max(1, float64(e.hitpoints)/200.0)
		} else {
			scale = 1 + ((float64(e.totalHitpoints/HitpointIncrement) - 1) * scaleFactor)
		}
	}
	return scale
}

func (e *Enemy) setState(state enemyState) {
	e.state = state
	e.stateTime = time.Now().Unix()
}

func (e *Enemy) moveTowardHero(heroX, heroY int) {
	if e.x != heroX {
		if e.x < heroX {
			e.x = e.x + 1
		} else if e.x > heroX {
			e.x = e.x - 1
		}
	}
	if e.y != heroY {
		if e.y < heroY {
			e.y = e.y + 1
		} else if e.y > heroY {
			e.y = e.y - 1
		}
	}
}

func (e *Enemy) moveAwayFromHero(heroX, heroY int) {
	if e.x != heroX {
		if e.x < heroX {
			e.x = e.x - 1
		} else if e.x > heroX {
			e.x = e.x + 1
		}
	}
	if e.y != heroY {
		if e.y < heroY {
			e.y = e.y - 1
		} else if e.y > heroY {
			e.y = e.y + 1
		}
	}
}

func (e *Enemy) move(heroX, heroY int) {
	if rand.Intn(2) > 0 {
		e.moveTowardHero(heroX, heroY)
	}
}

func (e *Enemy) Update(width, height, heroX, heroY int) {
	switch e.state {
	case Alive:
		e.move(heroX, heroY)
		offset := int((time.Now().Unix() - e.stateTime) % 2)
		e.frame = offset
	case Dead:
		if time.Now().Unix() > (e.stateTime + 2) {
			e.setState(Buried)
		} else {
			e.frame = 2
		}
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	scale := e.getScale()
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(float64(e.x), float64(e.y))
	subImageRect := image.Rect(e.frame*imageWidth, 0, (e.frame+1)*imageWidth, e.image.Bounds().Dy())
	subImage := e.image.SubImage(subImageRect).(*ebiten.Image)
	screen.DrawImage(subImage, op)

	if e.IsAlive() {
		op.GeoM.Reset()
		healthImage := images.GetImages().EnemyHealth
		y := e.y + healthMargin + int(float64(e.image.Bounds().Dy())*scale)
		if y > screen.Bounds().Dy() {
			y = e.y - healthMargin - int(float64(e.image.Bounds().Dy())*scale)
		}
		op.GeoM.Scale(scale, 1)
		subImageWidth := float64(healthImage.Bounds().Dx()) * (float64(e.hitpoints) / float64(e.totalHitpoints))
		scaledHealthWidth := subImageWidth * scale
		x := e.x + int(e.Bounds().Dx()/2)
		x = x - int(scaledHealthWidth/2)
		op.GeoM.Translate(float64(x), float64(y))
		subImageRect = image.Rect(0, 0, int(subImageWidth), healthImage.Bounds().Dy())
		subImage = healthImage.SubImage(subImageRect).(*ebiten.Image)
		screen.DrawImage(subImage, op)
	}
}

func (e *Enemy) Bounds() *image.Rectangle {
	return &image.Rectangle{
		Min: image.Point{
			X: e.x,
			Y: e.y,
		},
		Max: image.Point{
			X: e.x + int(float64(imageWidth)*e.getScale()),
			Y: e.y + int(float64(e.image.Bounds().Dy())*e.getScale()),
		},
	}
}

func (e *Enemy) Shot(heroX, heroY int) {
	e.hitpoints = e.hitpoints - 1
	if e.hitpoints == 0 {
		e.setState(Dead)
	} else {
		e.moveAwayFromHero(heroX, heroY)
	}
}

func (e *Enemy) IsAlive() bool {
	return e.state == Alive
}

func (e *Enemy) IsBuried() bool {
	return e.state == Buried
}

func (e *Enemy) ToStart() {
	e.x = e.startX
	e.y = e.startY
}

func (e *Enemy) IsHit(arrow *Arrow, heroX, heroY int) (bool, bool) {
	hitWhileAlive := false
	hit := false
	b := e.Bounds()
	if b.Min.X <= arrow.X && arrow.X <= b.Max.X && b.Min.Y <= arrow.Y && arrow.Y <= b.Max.Y {
		hit = true
		if e.IsAlive() {
			hitWhileAlive = true
			e.Shot(heroX, heroY)
		}
	}

	return hitWhileAlive, hit
}

func NewEnemy(x, y, hp int, boss bool, image *ebiten.Image) *Enemy {
	enemy := &Enemy{
		startX:         x,
		startY:         y,
		x:              x,
		y:              y,
		state:          Alive,
		stateTime:      time.Now().Unix(),
		frame:          0,
		image:          image,
		hitpoints:      hp,
		totalHitpoints: hp,
		boss:           boss,
	}
	enemy.x = enemy.x - int((enemy.getScale()*float64(imageWidth))/2.0)
	enemy.startX = enemy.x
	enemy.y = enemy.y - int((enemy.getScale()*float64(enemy.image.Bounds().Dy()))/2.0)
	enemy.startY = enemy.y
	return enemy
}
