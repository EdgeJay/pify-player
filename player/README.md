# Pify Player - Web Player

Web player that can be accessed by mobile/desktop browsers to control Spotify playback. Application setup using Sveltekit.

## Getting Started

1. Create `.env` file by copying `.env.example` file and edit the env vars.
2. Make sure you have Node.js v22+ installed, and run `npm install` command to install dependencies.

## Developing

Start the development server:

```bash
npm run dev

# or start the server and open the app in a new browser tab
npm run dev -- --open
```

### Running HMR in Dockerized environment

Ran into minor roadblock to get HMR to work in Dockerized environment, but with some research and sample from https://github.com/woollysammoth/sveltekit-docker-nginx HMR works now.

For HMR to work properly:

- Include `--host 0.0.0.0` option when running `npm run dev`.
- Use volumes to map local files/folders into container app to make sure changes are detected.
- Make sure node_modules in container are not overridden by local copy or else it will mess with OS-specific builds of libraries (such as rollup).
- Access Svelte app via `0.0.0.0` instead of `localhost`.
- Make sure ports are exposed and mapped properly in Dockerfile and docker-compose.yml files.
- Include following snippet into `server` portion of `vite.config.ts`:

```
server: {
    hmr: {
        clientPort: 5173
    },
    host: '0.0.0.0',
    port: 5173
}
```

## Production Build

To create a production version of your app:

```bash
npm run build
```

You can preview the production build with `npm run preview`.

> To deploy your app, you may need to install an [adapter](https://svelte.dev/docs/kit/adapters) for your target environment.
