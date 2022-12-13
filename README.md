# poker

[![Go](https://github.com/fc92/poker/actions/workflows/go.yml/badge.svg)](https://github.com/fc92/poker/actions/workflows/go.yml)
[![release](https://github.com/fc92/poker/actions/workflows/release.yaml/badge.svg)](https://github.com/fc92/poker/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fc92/poker)](https://goreportcard.com/report/github.com/fc92/poker)
[![codecov](https://codecov.io/github/fc92/poker/branch/main/graph/badge.svg?token=R4OZKBC13P)](https://codecov.io/github/fc92/poker)
[![Go Reference](https://pkg.go.dev/badge/github.com/fc92/poker.svg)](https://pkg.go.dev/github.com/fc92/poker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

- [poker](#poker)
  - [Game description](#game-description)
  - [Install](#install)
    - [Deploy the game for Web based access (recommended)](#deploy-the-game-for-web-based-access-recommended)
    - [Use CLI in console mode (advanced)](#use-cli-in-console-mode-advanced)
      - [Binary](#binary)
      - [Docker image](#docker-image)
  - [Information for developers](#information-for-developers)
    - [Client console](#client-console)
    - [Server](#server)
    - [High level design](#high-level-design)
    - [Low level design](#low-level-design)
    - [Debug log](#debug-log)
    - [Ideas for the future](#ideas-for-the-future)

## Game description

*poker* is a planning poker game (or scrum poker) limiting vote influence between the players. The goals are to collect an independent vote from each player and to improve the vote experience of distributed teams.

- Vote progress is shared without showing vote values to avoid influence between players.
- Vote values are revealed only when the vote session is closed, and the vote distribution is displayed.
- Players can join or leave at any time.

![short demo](4players.gif)

## Install

The are different ways of deploying this client/server game:

- with Web based access for better user experience (recommended),
- using CLI in console mode (advanced).

### Deploy the game for Web based access (recommended)

Prerequisite: server with docker and [tty2web](https://github.com/kost/tty2web) binary  (tested on Linux x86_64).

To simplify the user experience it is recommended to:

- start the server in a docker container:

```docker run -p 192.168.0.1:8080:8080/tcp  ghcr.io/fc92/poker:main```

to expose the server on address 192.168.0.1 port TCP 8080

- provide user access in a web browser using [tty2web](https://github.com/kost/tty2web):

```tty2web --title-format Poker --permit-arguments -a 192.168.0.1 -p 8081 -w docker run -it --rm ghcr.io/fc92/poker:main /poker client -websocket 192.168.0.1:8080```

so that users can connect to <http://192.168.0.1:8081/?arg=-name&arg=Mary> to join the game in a browser with player name *Mary*.

There are multiple benefits with this tty2web deployment method:

- users only need a web browser with proper network access to play the game,
- users can play on platforms that are not natively supported like iOS or Android,
- the server and each player process run inside a restricted container for security.

### Use CLI in console mode (advanced)

Binary or docker image version can be used to start the server or a client.

#### Binary

The same binary file is used to start the server and console clients. Various platform are supported. See [latest release page](https://github.com/fc92/poker/releases/latest) for more details.

Start a single server instance (dedicated for one team):

- using default websocket:

```bash
$ ./poker server
subcommand 'server'
  websocket: localhost:8080
```

- or specifying the websocket to open:

```bash
$ ./poker server -websocket 127.0.0.1:7878
subcommand 'server'
  websocket: 127.0.0.1:7878
```

When the server is started, start a client for each player. It must be able to communicate with the server.

- example of client specifying a player name and the websocket value:

```bash
$ ./poker client -name Player1 -websocket 127.0.0.1:7878
subcommand 'client'
  name: Player1
  server websocket: 127.0.0.1:7878
```

- example of client using a generated name and default websocket:

```bash
./poker client
subcommand 'client'
  name: snowy-cloud
  server websocket: localhost:8080
```

Each player can navigate the client console and send commands using keyboard and mouse.

NB: in mobaXterm mouse events are not always correctly handled

#### Docker image

By default the docker image runs a server listening on 0.0.0.0:8080.

Using the docker command line it is possible to expose that port outside of the container and to modify the IP and port of the server.

Server:

```
docker run -p 8080:8080  ghcr.io/fc92/poker:main
subcommand 'server'
  websocket: 0.0.0.0:8080`
```

Client with IP *192.168.0.10*:

It is also possible to start a client instance.

```
docker run -it ghcr.io/fc92/poker:main /poker client -name PlayerName -websocket 192.168.0.10:8080
```

## Information for developers

This section provides technical information for developers.

### Client console

Each player starts a console client to join the team server and play the game.

Client main features:

- [X] display available commands
- [X] allow mouse and keyboard user inputs
- [X] user defined or generated player name
- [X] commands: quit game, start new vote, send vote, modify vote, close the vote session
- [X] display vote progress during vote session
- [X] display distribution of votes when vote session is closed
- [X] content is driven by the server to maintain consistency between players

### Server

The server broadcasts a shared vision to all clients in real time.

Server main features:

- [X] start single game to host client players
- [X] share vote and available commands to all users
- [X] add newly connected user
- [X] remove disconnected user
- [X] broadcast vote status per user during vote session
- [X] broadcast votes when all votes are available or vote is manually closed
- [X] reset vote values when a new vote starts
- [X] close vote when all players have voted

### High level design

HTTP websocket is used for "real-time" communication between clients and the server. Messages are in JSON format and contain `Room` and `Participant` information.

A server instance maintains a single `Room` that contains:

- a `Participant` for each player,
- the status of the game: vote open or closed.

The `Room` content changes when a client sends a command or is disconnected. Each change is broadcasted to all clients.

Each `Participant` contains information about a player:

- player name,
- commands available to the user,
- vote value.

A `Participant` can be updated by local user actions or updates sent from the server.

### Low level design

Server and console client are written in pure Go language.

The software is packaged as a single binary file for each supported platform. The same file is used with different parameters to start a server or a client instance from a text console.

It should be possible to write other client implementations using other languages supporting websocket and JSON. The focus of this project is pure GO so far.

### Debug log

Server debug logs can be activated using the `-debug` flag

### Ideas for the future

- [ ] [WebAssembly.sh](https://webassembly.sh) support: tinygo WASI support is not yes sufficient

```bash
$ tinygo build -wasm-abi=generic -target=wasi -o poker.wasm cmd/poker.go 
# golang.org/x/sys/unix
../.go/pkg/mod/golang.org/x/sys@v0.0.0-20220722155257-8c9f86f7a55f/unix/syscall_unix.go:526:17: Exec not declared by package syscall
```

- [ ] server part on ESP32: [tinygo support for ESP32](https://tinygo.org/docs/reference/microcontrollers/esp32-coreboard-v2) is not sufficient
