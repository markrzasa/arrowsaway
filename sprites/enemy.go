package sprites

import (
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
	healthBar                 *Sprite
	Sprite                    *Sprite
	startX, startY, frame     int
	state                     enemyState
	stateTime                 int64
	hitpoints, totalHitpoints int
	boss                      bool
}

func (e *Enemy) setScale() {
	scale := 1.0
	if e.state == Alive {
		if e.boss {
			scale = math.Max(1, float64(e.hitpoints)/200.0)
		} else {
			scale = 1 + ((float64(e.totalHitpoints/HitpointIncrement) - 1) * scaleFactor)
		}
	}
	e.healthBar.ScaleX = scale
	e.Sprite.Scale(scale)
}

func (e *Enemy) setState(state enemyState) {
	e.state = state
	e.stateTime = time.Now().Unix()
}

func (e *Enemy) moveTowardHero(hero *Sprite) {
	if e.Sprite.X != hero.X {
		if e.Sprite.X < hero.X {
			e.Sprite.X = e.Sprite.X + 1
		} else if e.Sprite.X > hero.X {
			e.Sprite.X = e.Sprite.X - 1
		}
	}
	if e.Sprite.Y != hero.Y {
		if e.Sprite.Y < hero.Y {
			e.Sprite.Y = e.Sprite.Y + 1
		} else if e.Sprite.Y > hero.Y {
			e.Sprite.Y = e.Sprite.Y - 1
		}
	}
}

func (e *Enemy) moveAwayFromHero(hero *Sprite) {
	if e.Sprite.X != hero.X {
		if e.Sprite.X < hero.X {
			e.Sprite.X = e.Sprite.X - 1
		} else if e.Sprite.X > hero.X {
			e.Sprite.X = e.Sprite.X + 1
		}
	}
	if e.Sprite.Y != hero.Y {
		if e.Sprite.Y < hero.Y {
			e.Sprite.Y = e.Sprite.Y - 1
		} else if e.Sprite.Y > hero.Y {
			e.Sprite.Y = e.Sprite.Y + 1
		}
	}
}

func (e *Enemy) move(hero *Sprite) {
 	if rand.Intn(2) > 0 {
		e.moveTowardHero(hero)
	}
}

func (e *Enemy) Update(width, height int, hero *Sprite) {
	switch e.state {
	case Alive:
		e.move(hero)
		offset := int((time.Now().Unix() - e.stateTime) % 3)
		e.frame = offset
	case Dead:
		if time.Now().Unix() > (e.stateTime + 2) {
			e.setState(Buried)
		} else {
			e.frame = 3
		}
	}
	deltaX := hero.X - e.Sprite.X
	deltaY := hero.Y - e.Sprite.Y
	if e.IsAlive() {
		e.Sprite.Radians = math.Atan2(float64(deltaY), float64(deltaX)) - (math.Pi / 180)
		e.setScale()
	}
}

func (e *Enemy) Draw(screen *ebiten.Image) {
	e.Sprite.Draw(screen, e.frame)

	if e.IsAlive() {
		scaledBounds := e.Sprite.ScaledBounds()
		e.healthBar.X = e.Sprite.X
		e.healthBar.Y = scaledBounds.Max.Y + healthMargin
		if e.healthBar.Y > screen.Bounds().Dy() {
			e.healthBar.Y = scaledBounds.Min.Y - healthMargin
		}
		subImageWidth := int(float64(e.healthBar.imageWidth) * (float64(e.hitpoints)/float64(e.totalHitpoints)))
		e.healthBar.DrawSubImage(screen, subImageWidth)
	}
}

func (e *Enemy) Shot(hero *Sprite) {
	e.hitpoints = e.hitpoints - 1
	if e.hitpoints == 0 {
		e.setState(Dead)
	} else {
		e.moveAwayFromHero(hero)
	}
}

func (e *Enemy) IsAlive() bool {
	return e.state == Alive
}

func (e *Enemy) IsBuried() bool {
	return e.state == Buried
}

func (e *Enemy) ToStart() {
	e.Sprite.X = e.startX
	e.Sprite.Y = e.startY
}

func (e *Enemy) IsHit(arrow *Arrow, hero *Sprite) (bool, bool) {
	hitWhileAlive := false
	hit := false
	b := e.Sprite.ScaledBounds()
	if b.Min.X <= arrow.Sprite.X && arrow.Sprite.X <= b.Max.X && b.Min.Y <= arrow.Sprite.Y && arrow.Sprite.Y <= b.Max.Y {
		hit = true
		if e.IsAlive() {
			hitWhileAlive = true
			e.Shot(hero)
		}
	}

	return hitWhileAlive, hit
}

func NewEnemy(x, y, hp int, boss bool, image *ebiten.Image) *Enemy {
	health := images.GetImages().EnemyHealth
	enemy := &Enemy{
		startX:         x,
		startY:         y,
		state:          Alive,
		stateTime:      time.Now().Unix(),
		frame:          0,
		hitpoints:      hp,
		totalHitpoints: hp,
		boss:           boss,
		healthBar:      NewSprite(health.Bounds().Dx(), health),
		Sprite:         NewSprite(imageWidth, image),
	}
	enemy.Sprite.X = x
	enemy.startX = enemy.Sprite.X
	enemy.Sprite.Y = y
	enemy.startY = enemy.Sprite.Y
	enemy.Sprite.Radians = 0
	enemy.setScale()
	return enemy
}
