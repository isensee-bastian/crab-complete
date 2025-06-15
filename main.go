package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/isensee-bastian/crab/crab"
	"log"
)

func main() {
	ebiten.SetWindowSize(crab.ScreenWidth, crab.ScreenHeight)
	ebiten.SetWindowTitle("Crab Game")

	game := crab.NewGame()
	defer game.Close()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatalf("Error while running game loop: %v", err)
	}
}
