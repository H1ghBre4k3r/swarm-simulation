# Swarm Simulation

[![Test](https://github.com/H1ghBre4k3r/swarm-simulation/actions/workflows/test.yml/badge.svg)](https://github.com/H1ghBre4k3r/swarm-simulation/actions/workflows/test.yml)

Swarm simulation written in Go.

## Roadmap

- [x] sub-process per entity
- [x] rendering using SDL2
- [x] participant configuration
- [x] simulation configuration
- [x] terminal only version (without SDL2 dependency)
- [ ] GLFW + OpenGL/BGFX for rendering
- [ ] support for 3D

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
	"settings": {
		// optional settings for the simulation
		"tickLength": 10, // [OPTIONAL, default=1] minimum length of a tick in the simulation
		"tau": 120, // [OPTIONAL, default=1] 'look ahead range' for participants
		"noise": 0.01 // [OPTIONAL, default=0] noise for the simulation (between 0 and 1 - everything else is not useful)
	},
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
			"script": "scripts/test.py", // path to the script for this participant, RELATIVE to the location of this configuration file
			"ignoreFinish": true // [OPTIONAL, default=false] flag for indicating, if the simulation shall ignore the movement (and process) of this participant
		}
		//...
	]
}
```

Additionally, the simulation supports a number of command line flags:

- `-h`: Show usage of the simulation
- `-no-gui`: Don't start a GUI
- `-no-grid`: Hide the grid within the simulation
- `-c`: Path to the configuration file
- `-n`: Noise for the simulation (this overwrites the value in the configuration file)

The headless executable (only terminal - see **Headless**) supports only a part of the flags.

## Development

Simply open this folder in your favorite editor and start coding.

### Running

```sh
$ make run
```

### Running Tests

> There are currently no tests :^)

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

### Headless

If you want to build a headless version (i.e., without GUI support and SDL2), you can compile a specific, console-only version of the simulation.

```sh
$ make terminal
```

## Communication

The simulation starts a new process for each participant and communicates with it via `stdin`, `stdou`, and `stderr`:

- `stdin` is used to pass information from the simulation to the participants
- `stdout` is used to receive information from the participants
- `stderr` _can_ be used to log arbitrary information to the output of the simulation

### Procedure

At startup, the simulation sends some initial information to each participant:

```json
{
	"position": {
		// coordinates of the position of this participant
		"x": 0.2,
		"y": 0.5
	},
	"radius": 0.015, // radius of this participant
	, // radius of the stddev around this participant
	"vmax": 0.0005, // maximum velocity of this participant
	"target": {
		// coordinates of the target of this participant
		"x": 0.8,
		"y": 0.5
	}
}
```

During each tick of the simulation, it sends its current position and information about all other participants to each participant:

```json
{
	"position": {
		// coordinates of the position of this participant
		// ...
	},
	"participants": [
		{
			"position": {
				// position of the participant
			},
			"velocity": {
				// current velocity of this participant
				"x": 0.0012,
				"y": 0.000004
			},
			"distance": 0.412, // relative distance to this participant
			"radius": 0.015, // radius of this participant
			5 // stddev for this participant
		}
	]
}
```

After this, it expects an answer of each participant:

```json
{
	"action": "move",
	"payload": {
		// new velocity of this participant
		"x": 0.00123,
		"y": 0.002
	}
}
```
