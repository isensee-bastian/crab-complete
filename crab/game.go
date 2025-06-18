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

	beachScaleFactor = 2
	birdScaleFactor  = 2
	walkableMinY     = 180 * beachScaleFactor
	walkableMaxY     = 320 * beachScaleFactor

	animationFrameWidth   = 192 / 4
	animationFrameHeight  = 192 / 4
	animationFrameColumns = 4
	crabAnimationRow      = 0
	birdAnimationRow      = 3

	ticksPerSecond = 60
	ticksPerFrame  = ticksPerSecond / 4
	moveStepTick   = 2
)

var (
	beachImage *ebiten.Image
	crabFrames []*ebiten.Image
	birdFrames []*ebiten.Image
)

func init() {
	beachStdImage, _, err := image.Decode(bytes.NewReader(sprites.Beach))
	if err != nil {
		log.Fatalf("Error while loading image: %v", err)
	}
	beachImage = ebiten.NewImageFromImage(beachStdImage)

	crabFrames = readAnimationImages(sprites.Crab, crabAnimationRow)
	birdFrames = readAnimationImages(sprites.Bird, birdAnimationRow)
}

func readAnimationImages(rawImage []byte, row int) []*ebiten.Image {
	stdAnimationImage, _, err := image.Decode(bytes.NewReader(rawImage))
	if err != nil {
		log.Fatalf("Error while loading image: %v", err)
	}
	animationImage := ebiten.NewImageFromImage(stdAnimationImage)

	var allFrames []*ebiten.Image

	for index := 0; index < animationFrameColumns; index++ {
		xOffset := index * animationFrameWidth
		frameImage := animationImage.SubImage(image.Rect(
			xOffset,
			animationFrameHeight*row,
			xOffset+animationFrameWidth-1,
			animationFrameHeight*row+animationFrameHeight-1,
		))
		allFrames = append(allFrames, ebiten.NewImageFromImage(frameImage))
	}

	return allFrames
}

type Game struct {
	frame          int
	animationIndex int

	crabX int
	crabY int

	birdX int
	birdY int
}

func NewGame() *Game {
	return &Game{
		frame:          0,
		animationIndex: 0,
		crabX:          (ScreenWidth - animationFrameWidth) / 2,
		crabY:          (ScreenHeight - animationFrameHeight) / 2,
		birdX:          0,
		birdY:          (ScreenHeight-animationFrameHeight)/2 + animationFrameHeight*2,
	}
}

func (g *Game) Close() {
	// Nothing to do for cleanup so far (implement when needed).
}

func (g *Game) moveCrabLeft() {
	g.crabX = max(g.crabX-moveStepTick, 0)
}

func (g *Game) moveCrabRight() {
	g.crabX = min(g.crabX+moveStepTick, ScreenWidth-animationFrameWidth-1)
}

func (g *Game) moveCrabUp() {
	g.crabY = max(g.crabY-moveStepTick, walkableMinY)
}

func (g *Game) moveCrabDown() {
	g.crabY = min(g.crabY+moveStepTick, walkableMaxY-animationFrameWidth-1)
}

func (g *Game) moveBird() {
	if g.birdX >= ScreenWidth {
		g.birdX = 0
	} else {
		g.birdX += moveStepTick
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Signal that the game shall terminate normally.
		return ebiten.Termination
	}

	g.frame = (g.frame + 1) % ticksPerSecond
	g.animationIndex = g.frame / ticksPerFrame

	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 {
		g.moveCrabLeft()
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 {
		g.moveCrabRight()
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowUp) > 0 {
		g.moveCrabUp()
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowDown) > 0 {
		g.moveCrabDown()
	}

	g.moveBird()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	drawText(screen, 10, 10, "Hello Crab Game!")
	// Draw beach.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(beachScaleFactor, beachScaleFactor)
		opts.GeoM.Translate(0, 0)
		screen.DrawImage(beachImage, opts)
	}
	// Draw crab.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(g.crabX), float64(g.crabY))
		screen.DrawImage(crabFrames[g.animationIndex], opts)
	}
	// Draw bird.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(birdScaleFactor, birdScaleFactor)
		opts.GeoM.Translate(float64(g.birdX), float64(g.birdY))
		screen.DrawImage(birdFrames[g.animationIndex], opts)
	}
}

func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return width, height
}
