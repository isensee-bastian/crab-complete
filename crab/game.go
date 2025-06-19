package crab

import (
	"bytes"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/isensee-bastian/crab/resources/images/sprites"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand/v2"
)

const (
	ScreenWidth  = 1000
	ScreenHeight = 800

	walkableMinY = 180 * beachScaleFactor
	walkableMaxY = 320 * beachScaleFactor

	scoreX    = 5
	scoreY    = 0
	levelX    = 850
	levelY    = 0
	gameOverX = 200
	gameOverY = walkableMinY

	beachScaleFactor   = 2
	birdScaleFactor    = 1.5
	defaultScaleFactor = 1

	spriteWidth           = 192 / 4
	spriteHeight          = 192 / 4
	animationFrameColumns = 4
	crabAnimationRow      = 0
	birdAnimationRow      = 0

	ticksPerSecond  = 60
	ticksPerFrame   = ticksPerSecond / 4
	defaultStepTick = 2

	maxBirdCount      = 3
	maxBirdStepTick   = 5
	scoreLevelDivisor = 3
)

var (
	beachImage *ebiten.Image
	fishImage  *ebiten.Image
	crabFrames []*ebiten.Image
	birdFrames []*ebiten.Image
)

func init() {
	beachImage = readImage(sprites.Beach)
	fishImage = readImage(sprites.Fish)

	crabFrames = readAnimationImages(sprites.Crab, crabAnimationRow)
	birdFrames = readAnimationImages(sprites.Bird, birdAnimationRow)
}

func readImage(rawImage []byte) *ebiten.Image {
	stdImage, _, err := image.Decode(bytes.NewReader(rawImage))
	if err != nil {
		log.Fatalf("Error while loading image: %v", err)
	}
	return ebiten.NewImageFromImage(stdImage)
}

func readAnimationImages(rawImage []byte, row int) []*ebiten.Image {
	stdAnimationImage, _, err := image.Decode(bytes.NewReader(rawImage))
	if err != nil {
		log.Fatalf("Error while loading image: %v", err)
	}
	animationImage := ebiten.NewImageFromImage(stdAnimationImage)

	var allFrames []*ebiten.Image

	for index := 0; index < animationFrameColumns; index++ {
		xOffset := index * spriteWidth
		frameImage := animationImage.SubImage(image.Rect(
			xOffset,
			spriteHeight*row,
			xOffset+spriteWidth-1,
			spriteHeight*row+spriteHeight-1,
		))
		allFrames = append(allFrames, ebiten.NewImageFromImage(frameImage))
	}

	return allFrames
}

type Sprite struct {
	x            int             // x coordinate position
	y            int             // y coordinate position
	scale        float64         // scale is used to resize the image if it is not set to 1
	image        *ebiten.Image   // image points to the current animations frame if the sprite is animated
	animation    []*ebiten.Image // animation is only relevant for animated sprites, otherwise nil
	moveStepTick int             // moveStepTick is only relevant for animated sprites and specifies the way to move per tick, 0 for static sprites
}

func (s *Sprite) Width() int {
	return int(spriteWidth * s.scale)
}

func (s *Sprite) Height() int {
	return int(spriteHeight * s.scale)
}

func (s *Sprite) Rectangle() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{
			X: s.x,
			Y: s.y,
		},
		Max: image.Point{
			X: s.x + s.Width(),
			Y: s.y + s.Height(),
		},
	}
}

func (s *Sprite) NextImage(index int) {
	if index >= len(s.animation) {
		// Do nothing for static sprites that have no animation (or not enough animation frames).
		return
	}

	s.image = s.animation[index]
}

func (s *Sprite) CollidesWith(other *Sprite) bool {
	thisRect := s.Rectangle()
	otherRect := other.Rectangle()
	overlaps := thisRect.Overlaps(otherRect)

	return overlaps
}

func (s *Sprite) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(s.scale, s.scale)
	opts.GeoM.Translate(float64(s.x), float64(s.y))

	screen.DrawImage(s.image, opts)
}

type Game struct {
	frame int
	score int
	level int
	over  bool

	crab  *Sprite
	fish  *Sprite
	birds []*Sprite
}

func (g *Game) UpdateSprites() {
	g.frame = (g.frame + 1) % ticksPerSecond
	animationIndex := g.frame / ticksPerFrame

	g.crab.NextImage(animationIndex)

	for _, bird := range g.birds {
		bird.NextImage(animationIndex)
	}
	// No need to update the fish sprite as it is not animated.
}

func NewGame() *Game {
	game := &Game{}
	game.Restart()

	return game
}

// Restart resets all game state to its initial values.
func (g *Game) Restart() {
	g.frame = 0
	g.score = 0
	g.over = false

	fishX, fishY := randomPosition()

	g.crab = &Sprite{
		x:            (ScreenWidth - spriteWidth) / 2,
		y:            walkableMinY,
		scale:        defaultScaleFactor,
		image:        crabFrames[0],
		animation:    crabFrames,
		moveStepTick: defaultStepTick,
	}
	g.birds = []*Sprite{{
		x:            0,
		y:            walkableMinY + spriteHeight,
		scale:        birdScaleFactor,
		image:        birdFrames[0],
		animation:    birdFrames,
		moveStepTick: defaultStepTick,
	}}
	g.fish = &Sprite{
		x:     fishX,
		y:     fishY,
		scale: defaultScaleFactor,
		image: fishImage,
	}
}

func (g *Game) Close() {
	// Nothing to do for cleanup so far (implement when needed).
}

func (g *Game) moveCrabLeft() {
	g.crab.x = max(g.crab.x-g.crab.moveStepTick, 0)
}

func (g *Game) moveCrabRight() {
	g.crab.x = min(g.crab.x+g.crab.moveStepTick, ScreenWidth-spriteWidth-1)
}

func (g *Game) moveCrabUp() {
	g.crab.y = max(g.crab.y-g.crab.moveStepTick, walkableMinY)
}

func (g *Game) moveCrabDown() {
	g.crab.y = min(g.crab.y+g.crab.moveStepTick, walkableMaxY-spriteWidth-1)
}

func (g *Game) updateLevel() {
	nextLevel := g.score / scoreLevelDivisor

	if nextLevel > g.level {
		// Increase difficulty, first by adding more birds, then by accelerating their speed.
		// Do not increase level and difficulty if we have reached the max count and speed of birds.
		if g.addBird() || g.speedUpRandomBird() {
			g.level = nextLevel
		}
	}
}

func (g *Game) addBird() bool {
	if len(g.birds) >= maxBirdCount {
		return false
	}

	g.birds = append(g.birds, &Sprite{
		x:            0,
		y:            walkableMinY + spriteHeight*len(g.birds)*2,
		scale:        birdScaleFactor,
		image:        birdFrames[0],
		animation:    birdFrames,
		moveStepTick: defaultStepTick,
	})

	return true
}

func (g *Game) speedUpRandomBird() bool {
	var upgradableBirds []*Sprite

	for _, bird := range g.birds {
		if bird.moveStepTick < maxBirdStepTick {
			upgradableBirds = append(upgradableBirds, bird)
		}
	}

	// Skip any more speed increases if all birds are maxed out.
	if len(upgradableBirds) <= 0 {
		return false
	}

	upgradableBirds[rand.IntN(len(upgradableBirds))].moveStepTick += 1

	return true
}

func (g *Game) moveBirds() {
	for _, bird := range g.birds {
		if bird.x >= ScreenWidth {
			bird.x = 0
		} else {
			bird.x += bird.moveStepTick
		}
	}
}

// randomPosition returns a random, but walkable position, e.g. for spawning new sprites.
func randomPosition() (int, int) {
	maxX := ScreenWidth - spriteWidth
	maxYOffset := walkableMaxY - walkableMinY - spriteHeight

	x := rand.IntN(maxX)
	y := rand.IntN(maxYOffset) + walkableMinY

	return x, y
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Signal that the game shall terminate normally.
		return ebiten.Termination
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		// Restart the game, this is also possible if it is not over yet.
		g.Restart()
		return nil
	}
	if g.over {
		// Game over, do not update the scene until the game is restarted.
		return nil
	}

	g.UpdateSprites()

	if inpututil.KeyPressDuration(ebiten.KeyArrowLeft) > 0 {
		g.moveCrabLeft()
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowRight) > 0 {
		g.moveCrabRight()
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowUp) > 0 {
		g.moveCrabUp()
	} else if inpututil.KeyPressDuration(ebiten.KeyArrowDown) > 0 {
		g.moveCrabDown()
	}

	g.moveBirds()

	for _, bird := range g.birds {
		if g.crab.CollidesWith(bird) {
			// Game over, stop the round.
			g.over = true

			return nil
		}
	}

	if g.crab.CollidesWith(g.fish) {
		// Crab got the fish, increase score, spawn a new fish and increase bird difficulty.
		g.fish.x, g.fish.y = randomPosition()
		g.score += 1
		g.updateLevel()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw static beach.
	{
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Scale(beachScaleFactor, beachScaleFactor)
		opts.GeoM.Translate(0, 0)
		screen.DrawImage(beachImage, opts)
	}

	// Draw sprites.
	g.crab.Draw(screen)
	g.fish.Draw(screen)

	for _, bird := range g.birds {
		bird.Draw(screen)
	}

	// Draw score and level indicator.
	drawBigText(screen, scoreX, scoreY, color.Black, fmt.Sprintf("Score: %d", g.score))
	drawBigText(screen, levelX, levelY, color.Black, fmt.Sprintf("Level: %d", g.level))

	// Draw game over info if game has ended.
	if g.over {
		drawBigText(screen, gameOverX, gameOverY, color.Black, fmt.Sprintf("Game Over! (Enter: restart, Esc: exit)"))
	}
}

func (g *Game) Layout(width, height int) (screenWidth, screenHeight int) {
	return width, height
}
