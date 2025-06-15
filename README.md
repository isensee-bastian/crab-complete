# Crab Game

### About the Game

Todo

A simple 2D game built with Golang and Ebitengine.

### How to Run

* Ensure [Go is installed on your system](https://go.dev/doc/install)
* Download and extract or `git clone` this repositories content to your local machine
* Navigate into your local repository directory (e.g. in the terminal) and run `go run main.go`

### How to Play

Todo

### Troubleshooting

#### No Audio

You should hear some audio effects while playing the gameI. f you don't hear any sounds while playing, check your audio output device and volume, make sure it is not muted. If it is still not working, and you are running on Linux, you may need to apply subsequent workaround to disable a problematic audio module. This worked for me, but please use it carefully at your own risk and revert it in case of any issues:
* Open following config file for editing: `/etc/modprobe.d/alsa-base.conf`
* Append an option to disable the possibly problematic module: `options snd-hda-intel model=auto blacklist snd_soc_avs`

### Media Sources

Todo

### Tools Used

* [Go](https://go.dev/) as the general programming language.
* [Ebitengine](https://ebitengine.org/) for building a 2D game.
* [GIMP](https://www.gimp.org) for adapting images from source pictures.
* [ffmpeg](https://ffmpeg.org/) for adapting audio properties like sample size and volume.

### Improvement Ideas

Todo