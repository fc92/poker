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
    - [Standard deployment](#standard-deployment)
    - [Custom HTTP port deployment](#custom-http-port-deployment)

## Game description

*poker* is a planning poker game (or scrum poker) limiting vote influence between the players. The goals are to collect an independent vote from each player and to improve the vote experience of distributed teams.

- Vote progress is shared without showing vote values to avoid influence between players.
- Vote values are revealed only when the vote session is closed, and the vote distribution is displayed.
- Players can join or leave at any time.

![short demo](4players.gif)

## Install

The most common usage is to deploy:

- deploy a container for one poker room on server side,
- use a modern browser on client side to join the room for the game.

### Standard deployment

- Server using default HTTP port 8081:

```bash
docker run -p 8081:8081 -td ghcr.io/fc92/poker:main
```

- Browser URL to connect as player *Mary*:
`http://server_ip:8081/?arg=-name&arg=Mary`

### Custom HTTP port deployment

The port can be modified, to add a second poker room for example:

- Server using non default HTTP port 8083:

```bash
docker run -p 8083:8083 -td ghcr.io/fc92/poker:main ./clients.sh 8083
```

- Browser URL to connect as player *Mary*:
`http://server_ip:8083/?arg=-name&arg=Mary`
