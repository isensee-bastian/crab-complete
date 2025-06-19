# Crab Game

### About the Game

You are a crab on a sunny day at the beach. Suddenly, you feel a strong appetite for dead fish. Hence, you start looking around for delicious fish to devour. Easier said than done, since those nasty seagulls are strolling around here. Avoid getting eaten yourself by dodging the birds and eat as much fish as you can to make this a perfect day.

A simple 2D game built with Golang and Ebitengine.

### How to Run

* Ensure [Go is installed on your system](https://go.dev/doc/install)
* Download and extract or `git clone` this repositories content to your local machine
* Navigate into your local repository directory (e.g. in the terminal) and run `go run main.go`

### How to Play

* Use the arrow keys to move the crab left, right, up or down.
* For exiting the game, press the escape key.

### Troubleshooting

#### No Audio

You should hear some audio effects while playing the gameI. f you don't hear any sounds while playing, check your audio output device and volume, make sure it is not muted. If it is still not working, and you are running on Linux, you may need to apply subsequent workaround to disable a problematic audio module. This worked for me, but please use it carefully at your own risk and revert it in case of any issues:
* Open following config file for editing: `/etc/modprobe.d/alsa-base.conf`
* Append an option to disable the possibly problematic module: `options snd-hda-intel model=auto blacklist snd_soc_avs`

### Media Sources

All source images are AI generated using the following sites:
* [Retro Diffusion](https://www.retrodiffusion.ai/) for generating pixel art images, especially animations.
* [ideogram](https://ideogram.ai) for general image generation, including some pixel art images.

### Tools Used

* [Go](https://go.dev/) as the general programming language.
* [Ebitengine](https://ebitengine.org/) for building a 2D game.
* [GIMP](https://www.gimp.org) for adapting images from source pictures.
* [ffmpeg](https://ffmpeg.org/) for adapting audio properties like sample size and volume.

### Improvement Ideas

* Introduce sound effects.
* Ensure fish does not spawn directly at crab location or in too near distance.
* Increase difficulty over time by adding more birds or increasing their speed.
* Consider adding crabs that follow the players crab to make dodging more difficult over time.
* Improve collision detection (avoid collision when only outside transparent sprite parts overlap).