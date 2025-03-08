# Pify Player - Web Player API

API backend of Pify Player, written in Go

## Pre-requisites

1. Go 1.24.0 and above must be installed

## Getting Started

1. Create `.env` file by copying `.env.example` file and edit the env vars.
2. Run `go mod download` in this folder.
3. Run `make start` command from project root to start server.

## Development

1. Run `make start-dev` command from project root to start server in development mode with live reload (via [air](https://github.com/air-verse/air)).
2. Live reload settings controlled via `.air.toml` file.
