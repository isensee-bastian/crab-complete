package crab

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	_ "image/png"
)

const (
	ScreenWidth  = 800
	ScreenHeight = 600
)

type Game struct {
}

func NewGame() *Game {
	return &Game{}
}

func (g *Game) Close() {
	// Nothing to do for cleanup so far (implement when needed).
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Signal that the game shall terminate normally.
		return ebiten.Termination
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawText(screen, 10, 10, "Hello Crab Game!")
}

func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return width, height
}
