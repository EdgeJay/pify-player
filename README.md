# Pify Player

Web-based Spotify player powered by Svelte and Go, designed to run on Raspberry Pi devices.

## Pre-requisites

1. Docker installed
2. Node.js 22.14.0 and above installed for frontend
3. Go 1.24.0 and above installed for backend

## Getting Started

To get started in development mode, follow these steps:

1. `cd player`
2. `npm install`
3. `cd ../ && make start-dev`

Subsequent builds only execution of `make start-dev` command is needed, as `npm install` is needed to installed dependencies needed by frontend.

Open frontend application in browser using this url: `http://0.0.0.0:8080/`

## Running in production mode

`make start`
