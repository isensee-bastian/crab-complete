package crab

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"log"
)

const (
	sampleRate = 48000
)

var audioContext *audio.Context

func init() {
	audioContext = audio.NewContext(sampleRate)
}

// AudioPlayer is a thin wrapper around ebitengines original audio.Player type to simplify some commonly used actions
// like rewinding and playing a sound or closing the sound stream during application shutdown.
type AudioPlayer struct {
	*audio.Player
}

func newMp3AudioPlayer(rawSound []byte) *AudioPlayer {
	sellSound, err := mp3.DecodeF32(bytes.NewReader(rawSound))
	if err != nil {
		log.Fatalf("Failed to decode raw sound as mp3: %v", err)
	}

	audioPlayer, err := audioContext.NewPlayerF32(sellSound)
	if err != nil {
		log.Fatalf("Failed to create mp3 audio player: %v", err)
	}

	return &AudioPlayer{audioPlayer}
}

func (a *AudioPlayer) Replay() {
	err := a.Rewind()

	if err != nil {
		// Logging is sufficient here as playing sounds is not critical for the overall gameplay.
		log.Printf("Error on rewinding audio: %v", err)
		return
	}

	a.Play()
}

func (a *AudioPlayer) Close() {
	if a == nil {
		// Nothing to do.
		return
	}

	err := a.Player.Close()

	// Logging is sufficient here. Keep it simple and avoid unnecessary error handling on the caller side.
	if err != nil {
		log.Printf("Error on closing audio: %v", err)
	}
}
