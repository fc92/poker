# Contributing

This page provide technical information to facilitate future maintenance of the project.

- [Contributing](#contributing)
  - [Building backend from sources](#building-backend-from-sources)
  - [Deployments methods](#deployments-methods)
    - [Use CLI in console mode](#use-cli-in-console-mode)
    - [Deploy the game for Web based access](#deploy-the-game-for-web-based-access)
  - [Information for developers](#information-for-developers)
    - [Client console](#client-console)
    - [Server](#server)
    - [High level design](#high-level-design)
    - [Low level design](#low-level-design)
  - [Debug log](#debug-log)

## Building backend from sources

Prerequisites:

- internet access
- recent go version (tested with 1.19.x to 1.21)

Get the sources and build the `poker` binary:

```bash
git clone https://github.com/fc92/poker
cd poker/backend
go build -o poker cmd/poker.go
```

This binary contains can be used to start a server or a client process. It needs a terminal to run properly.

Other common go checks and tests can be run with:

```bash
go vet -v ./...
go test -v ./...
```

## Deployments methods

The are different ways of deploying this client/server game:

- using CLI in console mode (recommended for development or native client on a wide range of platforms),
- with Web based access for better user experience (recommended for end users),

### Use CLI in console mode

The same binary file is used to start the server and console clients. Various platform are supported. See [latest release page](https://github.com/fc92/poker/releases/latest) for more details. For example to use the Windows terminal natively without Docker use the corresponding binary file.

Start a single server instance (dedicated for one team):

- using default websocket:

```bash
./poker server
subcommand 'server'
  websocket: localhost:8080
```

- or specifying the websocket to open:

```bash
./poker server -websocket 127.0.0.1:7878
subcommand 'server'
  websocket: 127.0.0.1:7878
```

When the server is started, start a client for each player. It must be able to communicate with the server.

- example of client specifying a player name and the websocket value:

```bash
./poker client -name Player1 -websocket 127.0.0.1:7878
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

This binary file is located in `/poker` in the docker image and can be used with `docker run` on a limited number of supported platforms.

### Deploy the game for Web based access

In production the docker image is used to provide both server and client processes. This is described in the [README page](/README.md).

For development outside of docker the manual steps to work with the `poker` binary are described here. This binary is the result of the `go build -o poker cmd/poker.go` command and can also be found in `/poker` of the docker image.

Prerequisite: server with docker and [tty2web](https://github.com/kost/tty2web) binary  (tested on Linux x86_64).

To simplify the user experience it is recommended to:

- start the server in a docker container:

```bash
docker run -p 192.168.0.1:8080:8080/tcp  ghcr.io/fc92/poker:main /poker server -websocket 8080
```

to expose the server on address 192.168.0.1 port TCP 8080

- provide user access in a web browser using [tty2web](https://github.com/kost/tty2web):

```bash
tty2web --title-format Poker --permit-arguments -a 192.168.0.1 -p 8081 -w docker run -it --rm ghcr.io/fc92/poker:main /poker client -websocket 192.168.0.1:8080
```

so that users can connect to <http://192.168.0.1:8081/?arg=-name&arg=Mary> to join the game in a browser with player name *Mary*.

There are multiple benefits with this tty2web deployment method:

- users only need a web browser with proper network access to play the game,
- users can play on platforms that are not natively supported like iOS or Android,
- the server and each player process run inside a restricted container for security.


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

## Debug log

Server debug logs can be activated using the `-debug` flag

### Code coverage
![Codecov graph](https://codecov.io/github/fc92/poker/graphs/sunburst.svg?token=R4OZKBC13P "Codecov graph")
