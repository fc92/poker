# Contributing

This page provide technical information to facilitate future maintenance of the project.

- [Contributing](#contributing)
  - [Building the frontend and the backed](#building-the-frontend-and-the-backed)
  - [Building backend from sources](#building-backend-from-sources)
  - [Information for developers](#information-for-developers)
    - [Web client](#web-client)
    - [Server](#server)
    - [High level design](#high-level-design)
    - [Low level design](#low-level-design)
  - [Debug log](#debug-log)
    - [Code coverage](#code-coverage)

## Building the frontend and the backed

Prerequisites:

- install node.js and ionic CLI
- install go 1.21

```bash
make ionic-build
make docker-build
```

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

## Information for developers

This section provides technical information for developers.

### Web client

Each player starts vue.js application in a browser.

Client main features:

- [X] enter player name
- [X] choose poker room
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

Server is written in pure Go language. Client is written in TypeScript using vue.js and ionic framework.

The software is packaged as a single binary file for each supported platform.

## Debug log

Server debug logs can be activated using the `-debug` flag

### Code coverage

![Codecov graph](https://codecov.io/github/fc92/poker/graphs/sunburst.svg?token=R4OZKBC13P "Codecov graph")
