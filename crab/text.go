package crab

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"log"
)

var (
	normalTextFace *text.GoTextFace
	bigTextFace    *text.GoTextFace
)

func init() {
	source, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatalf("Error on creating new text face source: %v", err)
	}

	normalTextFace = &text.GoTextFace{
		Source: source,
		Size:   24,
	}
	bigTextFace = &text.GoTextFace{
		Source: source,
		Size:   32,
	}
}

func drawText(screen *ebiten.Image, x, y int, content string) {
	drawTextWithFace(screen, x, y, content, normalTextFace)
}

func drawBigText(screen *ebiten.Image, x, y int, content string) {
	drawTextWithFace(screen, x, y, content, bigTextFace)
}

func drawTextWithFace(screen *ebiten.Image, x, y int, content string, textFace *text.GoTextFace) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	text.Draw(screen, content, textFace, op)
}
