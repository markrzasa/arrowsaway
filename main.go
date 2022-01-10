package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"github.com/markrzasa/arrowsaway/fonts"
	"github.com/markrzasa/arrowsaway/images"
	"github.com/markrzasa/arrowsaway/level"
	"github.com/markrzasa/arrowsaway/sprites"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	milliBetweenShots = 250
)

type gameState int
const (
	NoClicker gameState = iota
	Running
	NextStage
	//NextLevel
	LostLife
	GameOver
	Winner
)

type ArrowsAway struct {
	height, width int

	state gameState

	gamepadIdsBuffer []ebiten.GamepadID
	gamepadIds       map[ebiten.GamepadID]bool

	hero *sprites.Hero

	levelIndex int
	levels []*level.Level

	enemies map[string]*sprites.Enemy

	arrows map[string]*sprites.Arrow

	lastShotMilli int64

	score int64

	lives int

	font font.Face
}

func (g *ArrowsAway) updateGamepads() {
	g.gamepadIdsBuffer = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIdsBuffer[:0])
	for _, id := range g.gamepadIdsBuffer {
		g.gamepadIds[id] = true
	}
	for id := range g.gamepadIds {
		if inpututil.IsGamepadJustDisconnected(id) {
			delete(g.gamepadIds, id)
		}
	}
}

func (g *ArrowsAway) isPadButtonPressed() bool {
	pressed := false
	for id := range g.gamepadIds {
		for b := ebiten.StandardGamepadButtonRightBottom; b < ebiten.StandardGamepadButtonMax; b++ {
			if ebiten.IsStandardGamepadButtonPressed(id, b) {
				pressed = true
				break
			}
		}
	}

	return pressed
}

func (g *ArrowsAway) hitEnemy(a *sprites.Arrow) bool {
	hitWhileAlive := false
	hit := false
	for _, e := range g.enemies {
		hitWhileAlive, hit = e.IsHit(a, g.hero.X, g.hero.Y)
		if hitWhileAlive {
			g.score = g.score + 10
		}	
	}

	return hit
}

func (g *ArrowsAway) updateHero() {
	noIntersect := image.Rectangle{}
	for _, e := range g.enemies {
		if e.IsAlive() {
			if e.Bounds().Intersect(*g.hero.Bounds()) != noIntersect {
				g.lives = g.lives - 1
				if g.lives == 0 {
					g.state = GameOver
				} else {
					g.state = LostLife
				}
				break
			}
		}
	}
}

func (g *ArrowsAway) updateArrows() {
	for id := range g.gamepadIds {
		x := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickHorizontal)
		y := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisRightStickVertical)
		if math.Abs(x) > 0.5 || math.Abs(y) > 0.5 {
			if g.lastShotMilli == 0 || (time.Now().UnixMilli() >= (g.lastShotMilli + milliBetweenShots)) {
				endX := g.hero.X + int(float64(g.width / 2) * x)
				endY := g.hero.Y + int(float64(g.height / 2) * y)
				id := uuid.New().String()
				g.arrows[id] = sprites.NewArrow(g.hero.X, g.hero.Y, endX, endY, images.GetImages().Arrow)
			}
		} else {
			g.lastShotMilli = 0
		}
	}

	for id, a := range g.arrows {
		if g.hitEnemy(a) {
			delete(g.arrows, id)
		} else if a.IsOffScreen(g.width, g.height) {
			delete(g.arrows, id)
		}
	}

	for _, a := range(g.arrows) {
		a.Update()
	}
}

func (g *ArrowsAway) updateEnemies() {
	for id, e := range g.enemies {
		if e.IsBuried() {
			delete(g.enemies, id)
		}
	}

	for _, e := range g.enemies {
		e.Update(g.width, g.height, g.hero.X, g.hero.Y)
	}
}

func (g *ArrowsAway) Update() error {
	g.updateGamepads()

	if len(g.gamepadIds) <= 0 {
		g.state = NoClicker
	}

	switch g.state {
	case NoClicker:
		if len(g.gamepadIds) > 0 {
			g.state = NextStage
		}
	case NextStage:
		if g.isPadButtonPressed() {
			g.hero.Center(g.width, g.height)
			g.state = Running
		}
	case Running:
		g.hero.Update(&g.gamepadIds, g.height, g.width)
		g.updateHero()
		g.updateArrows()
		g.updateEnemies()
		if len(g.enemies) == 0 {
			level := g.levels[g.levelIndex]
			if level.Complete() {
				g.levelIndex = g.levelIndex + 1
				if g.levelIndex == len(g.levels) {
					g.state = Winner
				} else {
					g.levels[g.levelIndex].PopulateEnemies(g.width, g.height, g.enemies)
					g.state = NextStage
				}
			} else {
				level.NextStage()
				level.PopulateEnemies(g.width, g.height, g.enemies)
				g.state = NextStage
			}
		}
	case LostLife:
		if g.isPadButtonPressed() {
			g.hero.Center(g.width, g.height)
			for id, e := range g.enemies {
				if !e.IsAlive() {
					delete(g.enemies, id)
				}
				e.ToStart()
			}
			g.state = Running
		}
	case GameOver:
		fallthrough
	case Winner:
		if g.isPadButtonPressed() {
			g.enemies = make(map[string]*sprites.Enemy)
			g.hero.Center(g.width, g.height)
			g.score = 0
			g.lives = 3
			g.state = NextStage
			g.levelIndex = 0
			g.levels[g.levelIndex].Reset()
			g.levels[g.levelIndex].PopulateEnemies(g.width, g.height, g.enemies)
		}
	}
	return nil
}

func (g *ArrowsAway) tileFloor(screen *ebiten.Image) {
	background := g.levels[g.levelIndex].GetBackground()
	rows := (g.height / background.Bounds().Dy()) + 1
	cols := (g.width / background.Bounds().Dx()) + 1

	op := &ebiten.DrawImageOptions{}
	for c := 0; c < cols; c++ {
		for r := 0; r < rows; r++ {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(c*background.Bounds().Dx()), float64(r*background.Bounds().Dy()))
			screen.DrawImage(background, op)
		}
	}
}

func (g *ArrowsAway) drawCentered(screen *ebiten.Image, y int, t []string) {
	for i, s := range t {
		r := text.BoundString(g.font, s)
		text.Draw(
			screen,
			s,
			g.font,
			(g.width / 2) - (r.Dx() / 2), y + (i * (r.Dy() + 10)), color.RGBA{0x00, 0x00, 0x00, 0xff})	
	}
}

func (g *ArrowsAway) Draw(screen *ebiten.Image) {
	switch g.state {
	case NoClicker:
		screen.Fill(color.RGBA{0x87, 0xCE, 0xEB, 0xff})
		g.drawCentered(screen, 40, []string{"Plugin a clicker to get started.",})
	case Running:
		g.tileFloor(screen)
		g.hero.Draw(screen)
		for _, a := range g.arrows {
			a.Draw(screen)
		}
		for _, e := range g.enemies {
			e.Draw(screen)
		}
		text.Draw(
			screen,
			fmt.Sprintf("Score: %d", g.score),
			g.font,
			10, g.height - 20, color.RGBA{0x00, 0x00, 0x00, 0xff})
		lifeImage := images.GetImages().Life
		op := &ebiten.DrawImageOptions{}
		for i := 0 ; i < g.lives ; i++ {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(g.width - 10 - (lifeImage.Bounds().Dx() * (i + 1))), float64(g.height - 20 - lifeImage.Bounds().Dy()))
			screen.DrawImage(lifeImage, op)
		}
	case NextStage:
		level := g.levels[g.levelIndex]
		screen.Fill(color.RGBA{0x87, 0xCE, 0xEB, 0xff})
		g.drawCentered(screen, 40, []string{
			level.GetName(),
			fmt.Sprintf("%d - %d", g.levelIndex + 1, level.GetStage() + 1),
			"Press a button to start",
		})
	case LostLife:
		screen.Fill(color.RGBA{0x87, 0xCE, 0xEB, 0xff})
		g.drawCentered(screen, 40, []string{fmt.Sprintf("%d lives left. Press a button to keep trying.", g.lives),})
	case Winner:
		screen.Fill(color.RGBA{0x87, 0xCE, 0xEB, 0xff})
		g.drawCentered(screen, 40, []string{"You won! Press a button to play again",})
		g.hero.Winner(screen, g.width, g.height)
	case GameOver:
		screen.Fill(color.RGBA{0x87, 0xCE, 0xEB, 0xff})
		g.drawCentered(screen, 40, []string{"Game over. Press a button to try again",})
		g.hero.GameOver(screen, g.width, g.height)
	}
}

func (g *ArrowsAway) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if (outsideHeight != g.height) || (outsideWidth != g.width) {
		g.height = outsideHeight
		g.width = outsideWidth
	}
	return outsideWidth, outsideHeight
}

func (g *ArrowsAway) initialize() {
	tt, err := opentype.Parse(fonts.PressStart2PRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	g.font, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	if g.gamepadIds == nil {
		g.gamepadIds = map[ebiten.GamepadID]bool{}
	}
	g.height = 1000
	g.width = 1000
	g.lastShotMilli = 0
	g.levelIndex = 0
	g.levels = append(g.levels, level.NewLevel("Goblins in the grass", images.GetImages().Goblin, images.GetImages().Grass, 40))
	g.levels = append(g.levels, level.NewLevel("Skeletons on the stone", images.GetImages().Skeleton, images.GetImages().Stone, 40))
	g.lives = 3
	g.score = 0
	g.state = NextStage
	g.hero = sprites.NewHero(g.width, g.height, images.GetImages().Hero)
	g.arrows = make(map[string]*sprites.Arrow)
	g.enemies = make(map[string]*sprites.Enemy)
	g.levels[g.levelIndex].PopulateEnemies(g.width, g.height, g.enemies)
}

func main() {
	game := &ArrowsAway{}
	game.initialize()
	ebiten.SetWindowSize(game.width, game.height)
	ebiten.SetWindowTitle("Arrows Away")
	ebiten.SetWindowResizable(true)
	ebiten.SetScreenTransparent(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
