# poker

[![Go](https://github.com/fc92/poker/actions/workflows/go.yml/badge.svg)](https://github.com/fc92/poker/actions/workflows/go.yml)
[![release](https://github.com/fc92/poker/actions/workflows/release.yaml/badge.svg)](https://github.com/fc92/poker/actions/workflows/release.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/fc92/poker)](https://goreportcard.com/report/github.com/fc92/poker)
[![codecov](https://codecov.io/github/fc92/poker/branch/main/graph/badge.svg?token=R4OZKBC13P)](https://codecov.io/github/fc92/poker)
[![Maintainability](https://api.codeclimate.com/v1/badges/46853c43411ca445fc0d/maintainability)](https://codeclimate.com/github/fc92/poker/maintainability)
[![Go Reference](https://pkg.go.dev/badge/github.com/fc92/poker.svg)](https://pkg.go.dev/github.com/fc92/poker)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

- [poker](#poker)
  - [Game description](#game-description)
  - [Install](#install)

## Game description

*poker* is a planning poker game (or scrum poker) limiting vote influence between the players. The goals are to collect an independent vote from each player and to improve the vote experience of distributed teams.

- Vote progress is shared without showing vote values to avoid influence between players.
- Vote values are revealed only when the vote session is closed, and the vote distribution is displayed.
- Players can join or leave at any time.

## Install

The most common usage is to deploy:

- the helm chart and a tls secret `poker-tls` in a Kubernetes namespace,
- use a modern browser on client side to join the room for the game.

The list of rooms can be customized using the ROOM_LIST environment variable of the container.
