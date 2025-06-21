package sounds

import _ "embed"

var (
	//go:embed item-pickup.mp3
	ItemPickup []byte

	//go:embed jab.mp3
	Jab []byte
)
