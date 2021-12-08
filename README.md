# swarm-simulation

[![Test](https://github.com/H1ghBre4k3r/swarm-simulation/actions/workflows/test.yml/badge.svg)](https://github.com/H1ghBre4k3r/swarm-simulation/actions/workflows/test.yml)

Swarm simulation written in Go.

## Requirements

1. Go installed.

2. **Only for GUI version:** SDL2 (and some optional dependencies installed):

   On **Ubuntu 14.04 and above**, type:\
   `apt install libsdl2{,-image,-mixer,-ttf,-gfx}-dev`

   On **Fedora 25 and above**, type:\
   `yum install SDL2{,_image,_mixer,_ttf,_gfx}-devel`

   On **Arch Linux**, type:\
   `pacman -S sdl2{,_image,_mixer,_ttf,_gfx}`

   On **Gentoo**, type:\
   `emerge -av libsdl2 sdl2-{image,mixer,ttf,gfx}`

   On **macOS**, install SDL2 via [Homebrew](http://brew.sh) like so:\
   `brew install sdl2{,_image,_mixer,_ttf,_gfx} pkg-config`

## Usage

The simulation provides two executables: One with (`cmd/swarm-simulation`) and one without GUI (`cmd/swarm-simulation-terminal`). The GUI version requires SDL2 to be installed on your system (see above). You can compile each one via make:

```sh
$ make gui
```

and

```sh
$ make terminal
```

To run the simulation, you have to provide a configuration file via the `-c` flag.

```sh
$ ./bin/darwin_amd64/swarm-simulation -c ./examples/sample.json
```

The file has to have the following format:

```json
{
	"participants": [
		// array of participants of the simulation
		{
			"start": {
				// start coordinates for this participant
				"x": 0.1, // number between 0 and 1
				"y": 0.5 // number between 0 and 1
			},
			"radius": 0.05, // the radius of this participant (also between 0 and 1)
			"vmax": 0.001, // the maximum velocity of this participant (between 0 and 1)
			"target": {
				// coordinates of the target for this participant
				"x": 0.9, // between 0 and 1
				"y": 0.5 // betweedn 0 and 1
			},
			"script": "scripts/test.py" // path to the script for this participant, RELATIVE to the location of this configuration file
		}
		//...
	]
}
```

## Development

Simply open this folder in your favorite editor and start coding.

### Running

```sh
$ make run
```

### Running Tests

```sh
$ make test
```

### Building (for development)

```sh
$ make build
```

### Building (for production)

```sh
$ make release
```
