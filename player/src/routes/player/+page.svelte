<script lang="ts">
	import { onMount } from 'svelte';
	import { getSpotifyTokenFromStorage } from '$lib/session';
	import { controlPlayback } from '$lib/device';
	import { refreshAccessToken } from '$lib/session';
	import { getAndSaveYoutubeVideo } from '$lib/playback';
	import PlayerPanel from './components/player-panel.svelte';
	import LoginDialog from './components/login.svelte';

	const defaultVolume = 50;

	let { data } = $props();
	let deviceId = $state('');
	let ytVidId = $state('');

	// song track and playback status related vars
	let spotifyTrack: Spotify.Track | undefined = $state();
	let playbackPaused = $state(true);
	let position = $state('');
	let duration = $state('');
	let albumImage = $state('');
	let songTitle = $state('');
	let songArtists = $state<string[]>([]);
	let songProgress = $state(0);
	// volume controls
	let volume = $state(defaultVolume);

	// error message shown to user
	let errorMessage = $state('');

	// login dialog
	let isConnected = $state(true);

	let player: Spotify.Player | undefined = $state();

	let token = '';

	const convertToMinutes = (ms: number): string => {
		const minutes = Math.floor(ms / 60000);
		const seconds = Math.floor((ms % 60000) / 1000);
		return `${minutes}:${seconds.toString().padStart(2, '0')}`;
	};

	onMount(() => {
		window.onSpotifyWebPlaybackSDKReady = async () => {
			player = new Spotify.Player({
				name: data.playerName,
				getOAuthToken: async (cb) => {
					let tokenExpired = false;
					let { accessToken, expiresAt } = getSpotifyTokenFromStorage();

					if (expiresAt) {
						tokenExpired = new Date().getTime() > expiresAt.getTime();
					}

					if (!accessToken || !expiresAt || tokenExpired) {
						try {
							const { accessToken } = await refreshAccessToken(data.basicAuthToken)();
							token = accessToken;
						} catch (err) {
							console.error('Error refreshing access token:', err);
							// show login dialog with QR code
							isConnected = false;
							return;
						}
					} else {
						isConnected = true;
						token = accessToken;
					}

					cb(token);
				},
				volume: volume / 100
			});

			/* Activates an HTML element in the player instance. This is typically required before any media playback
				can occur, especially in browsers that enforce user interaction before allowing audio/video playback. */
			player.activateElement();

			// Player Ready
			player.addListener('ready', async ({ device_id }) => {
				console.log('Ready with Device ID', device_id);
				errorMessage = '';
				deviceId = device_id;

				// check playback state, take over as playback controls if state is null
				const state = await player!.getCurrentState();
				if (!state) {
					const success = await controlPlayback(token, deviceId);
					if (success) {
						console.log('took over playback control successfully');
					}
				}
			});

			// Player Not Ready
			player.addListener('not_ready', ({ device_id }) => {
				errorMessage = `Device ID has gone offline: ${device_id}`;
			});

			player.addListener('initialization_error', ({ message }) => {
				errorMessage = message;
			});

			player.addListener('playback_error', ({ message }) => {
				errorMessage = message;
				if (message.includes('token expired')) {
					// refreshAccessToken(ws);
				}
			});

			player.addListener('authentication_error', ({ message }) => {
				console.log('authentication_error');
				errorMessage = message;
				// refreshAccessToken(ws);
			});

			player.addListener('account_error', ({ message }) => {
				errorMessage = message;
			});

			player.addListener('autoplay_failed', () => {
				errorMessage = 'Autoplay is not allowed by the browser autoplay rules';
			});

			let intervalId: NodeJS.Timeout | undefined;

			player.addListener('player_state_changed', async (playbackState) => {
				if (!playbackState) {
					return;
				}

				const {
					paused,
					track_window: { current_track }
				} = playbackState;
				const stateChanged =
					spotifyTrack?.id !== current_track?.id ||
					(spotifyTrack?.id === current_track?.id && playbackPaused !== paused);

				console.log('player_state_changed');

				spotifyTrack = current_track;
				// console.log($state.snapshot(spotifyTrack));

				// update song details
				if (spotifyTrack?.album && spotifyTrack.album.images.length > 0) {
					albumImage = spotifyTrack.album.images[0].url;
				}
				songTitle = spotifyTrack?.name || '';
				songArtists =
					spotifyTrack?.artists?.reduce((acc, artist) => {
						acc.push(artist.name);
						return acc;
					}, [] as string[]) || [];

				playbackPaused = paused;

				await updatePlaybackPosition();

				if (intervalId) {
					clearInterval(intervalId);
				}

				if (!playbackPaused) {
					// start timer to update playback position
					intervalId = setInterval(async () => {
						if (!playbackPaused) {
							await updatePlaybackPosition();
						}
					}, 1000);
				}

				if (stateChanged && data.enableYoutube) {
					console.log('fetch youtube video');
					// get additional track info and Youtube info
					try {
						// const track = await getTrack(data.basicAuthToken, spotifyTrack?.id || '');
						// console.log(track);

						const trackId = spotifyTrack?.id || '';
						console.log('can fetch youtube video', songArtists.length > 0 && songTitle);
						if (trackId) {
							const res = await getAndSaveYoutubeVideo(
								data.basicAuthToken,
								`${spotifyTrack?.artists[0].name} ${spotifyTrack?.name}`,
								trackId
							);
							ytVidId = res?.data?.video_id || '';
						}
					} catch (err) {
						console.error('Error fetching track info:', err);
					}
				}
			});

			player.connect();
		};

		// Check if the Spotify SDK script is already in the document
		const spotifyScriptExists = document.querySelector(
			'script[src="https://sdk.scdn.co/spotify-player.js"]'
		);
		if (!spotifyScriptExists) {
			const script = document.createElement('script');
			script.src = 'https://sdk.scdn.co/spotify-player.js';
			document.body.appendChild(script);
		} else if (!player) {
			window.onSpotifyWebPlaybackSDKReady();
		}
	});

	/* Playback controls */
	const updatePlaybackPosition = async () => {
		const state = await player?.getCurrentState();
		if (state) {
			position = convertToMinutes(state.position);
			duration = convertToMinutes(state.duration);
			songProgress = Math.ceil((state.position / state.duration) * 100);
		}
	};
</script>

<div class="player-page">
	<div class="album-bg" style="background-image:url({albumImage})"></div>
	<div class="video-bg">
		{#if ytVidId}
			<iframe
				width="100%"
				height="100%"
				src={`https://www.youtube.com/embed/${ytVidId}?autoplay=1&controls=0&mute=1&loop=1`}
				title="Background video"
				frameborder="0"
				allow="autoplay; clipboard-write; encrypted-media;"
				style="pointer-events: none;"
			></iframe>
		{/if}
	</div>
	<PlayerPanel
		{player}
		{errorMessage}
		{playbackPaused}
		{position}
		{duration}
		{albumImage}
		{songTitle}
		{songArtists}
		{songProgress}
		{volume}
	/>
	{#if !isConnected}
		<LoginDialog basicAuthToken={data.basicAuthToken} />
	{/if}
</div>

<style>
	.player-page {
		position: relative;
		height: 100vh;
	}

	.album-bg {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-size: cover;
		background-position: center;
		background-repeat: no-repeat;
		z-index: 0;
		opacity: 0.5;
		filter: blur(5px);
	}

	.video-bg {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		z-index: 1;
	}
</style>
