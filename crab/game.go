package crab

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/isensee-bastian/crab/resources/images/sprites"
	"image"
	_ "image/png"
	"log"
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 800

	crabWidth      = 192 / 4
	crabHeight     = 192 / 4
	imagesPerRow   = 4
	ticksPerSecond = 60
	ticksPerFrame  = ticksPerSecond / 4

	moveStepTick = 4
)

var beachImage *ebiten.Image

var crabFrames []*ebiten.Image

func init() {
	beachStdImage, _, err := image.Decode(bytes.NewReader(sprites.Beach))
	if err != nil {
		log.Fatalf("Error while loading image: %v", err)
	}
	beachImage = ebiten.NewImageFromImage(beachStdImage)

	crabStdFrameImages, _, err := image.Decode(bytes.NewReader(sprites.Crab))
	if err != nil {
		log.Fatalf("Error while loading image: %v", err)
	}
	crabFrameImages := ebiten.NewImageFromImage(crabStdFrameImages)

	for index := 0; index < imagesPerRow; index++ {
		xOffset := index * crabWidth
		crabFrameImage := crabFrameImages.SubImage(image.Rect(xOffset, 0, xOffset+crabWidth-1, crabWidth-1))
		crabFrames = append(crabFrames, ebiten.NewImageFromImage(crabFrameImage))
	}
}

type Game struct {
	frame          int
	crabFrameIndex int
	crabX          int
	crabY          int
}

func NewGame() *Game {
	return &Game{
		frame:          0,
		crabFrameIndex: 0,
		crabX:          (ScreenWidth - crabWidth) / 2,
		crabY:          (ScreenHeight - crabWidth) / 2,
	}
}

func (g *Game) Close() {
	// Nothing to do for cleanup so far (implement when needed).
}

func (g *Game) moveLeft() {
	g.crabX = max(g.crabX-moveStepTick, 0)
}

func (g *Game) moveRight() {
	g.crabX = min(g.crabX+moveStepTick, ScreenWidth-crabWidth-1)
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Signal that the game shall terminate normally.
		return ebiten.Termination
	}

	g.frame = (g.frame + 1) % ticksPerSecond
	g.crabFrameIndex = g.frame / ticksPerFrame

	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 {
		g.moveLeft()
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 {
		g.moveRight()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawText(screen, 10, 10, "Hello Crab Game!")
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(2.0, 2.0)
		opts.GeoM.Translate(0, 0)
		screen.DrawImage(beachImage, opts)
	}
	{
		opts := &ebiten.DrawImageOptions{}
		//opts.GeoM.Scale(2.0, 2.0)
		opts.GeoM.Translate(float64(g.crabX), float64(g.crabY))
		screen.DrawImage(crabFrames[g.crabFrameIndex], opts)
	}
}

func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return width, height
}
