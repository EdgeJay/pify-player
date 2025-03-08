# Pify Player

Web-based Spotify player powered by Svelte and Go, designed to run on Raspberry Pi devices.

## Pre-requisites

1. Docker installed
2. Node.js 22.14.0 and above installed for frontend
3. Go 1.24.0 and above installed for backend
4. [SSL certs installed](#installing-ssl-certs)

## Installing SSL certs

SSL certs need to be setup for both `api` and `player` services, especially for `player` as Spotify Web SDK requires its player to be running in servers with SSL enabled.

Steps:

1. Install mkcert

### macOS

```
brew install mkcert
brew install nss # if you use Firefox
```

### Linux

```
# first install certutil
sudo apt install libnss3-tools
    -or-
sudo yum install nss-tools
    -or-
sudo pacman -S nss
    -or-
sudo zypper install mozilla-nss-tools

# install mkcert via homebrew (refer to brew.sh)
brew install mkcert
```


2. Create `.env` file in repo root folder. Can make a copy of `.env.example` file and rename it to `.env`.
3. Run `make generate-ssl`

Follow [this guide](https://words.filippo.io/mkcert-valid-https-certificates-for-localhost/) for detailed setup steps. [Github repo](https://github.com/FiloSottile/mkcert)

## Getting Started

To get started in development mode, follow these steps first to install `player` service dependencies:

1. `cd player`
2. `npm install`
3. `cd ../ && make start-dev`

Subsequent builds only execution of `make start-dev` command is needed, as `npm install` is needed to installed dependencies needed by frontend.

Open frontend application in browser using this url: `http://127.0.0.1:5173/`

## Running in production mode

`make start`
