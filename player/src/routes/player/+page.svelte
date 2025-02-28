<script lang="ts">
	import { onMount } from 'svelte';

	let deviceId = $state('');
	let spotifyTrack: Spotify.Track | undefined = $state();
	let playbackPaused = $state(true);
	let position = $state(0);
	let duration = $state(0);
	let errorMessage = $state('');

	const convertToMinutes = (ms: number): number => {
		const minutes = Math.floor((ms / 60000) * 100) / 100;
		return minutes;
	};

	interface LoginResponse {
		id: number;
		access_token: string;
		display_name: string;
	}

	onMount(() => {
		window.onSpotifyWebPlaybackSDKReady = async () => {
			let token = '';

			const player = new Spotify.Player({
				name: 'Pify Player',
				getOAuthToken: (cb) => {
					cb(token);
				},
				volume: 0.5
			});

			try {
				const response = await fetch('https://huijie-mbp.local:8080/api/auth/player', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					}
				});

				if (!response.ok) {
					throw new Error('Login failed');
				}

				const { access_token: accessToken } = (await response.json()) as LoginResponse;
				token = accessToken;

				console.log('Login successful with token:', token);

				player.connect();
			} catch (error) {
				console.error('Login error:', error);
				errorMessage = error instanceof Error ? error.message : 'Login failed';
			}

			// Player Ready
			player.addListener('ready', ({ device_id }) => {
				// console.log('Ready with Device ID', device_id);
				deviceId = device_id;
			});

			// Player Not Ready
			player.addListener('not_ready', ({ device_id }) => {
				console.log('Device ID has gone offline', device_id);
				errorMessage = `Device ID has gone offline: ${device_id}`;
			});

			player.addListener('initialization_error', ({ message }) => {
				// console.error(message);
				errorMessage = message;
			});

			player.addListener('authentication_error', ({ message }) => {
				// console.error(message);
				errorMessage = message;
			});

			player.addListener('account_error', ({ message }) => {
				// console.error(message);
				errorMessage = message;
			});

			player.addListener('autoplay_failed', () => {
				// console.log('Autoplay is not allowed by the browser autoplay rules');
				errorMessage = 'Autoplay is not allowed by the browser autoplay rules';
			});

			let intervalId: NodeJS.Timeout | undefined;

			player.addListener('player_state_changed', ({ paused, track_window: { current_track } }) => {
				spotifyTrack = current_track;
				playbackPaused = paused;

				if (intervalId) {
					clearInterval(intervalId);
				}

				if (!paused) {
					intervalId = setInterval(async () => {
						console.log('interval running');
						if (!playbackPaused) {
							await updatePlaybackPosition();
						}
					}, 1000);
				} else {
					if (intervalId) {
						console.log('cleared interval');
						clearInterval(intervalId);
					}
				}
			});

			const updatePlaybackPosition = async () => {
				const state = await player.getCurrentState();
				if (state) {
					position = convertToMinutes(state.position);
					duration = convertToMinutes(state.duration);
				}
			};
		};

		const script = document.createElement('script');
		script.src = 'https://sdk.scdn.co/spotify-player.js';
		document.body.appendChild(script);
	});
</script>

<p>Device ID: {deviceId}</p>
<p>Song Title: {spotifyTrack?.name || ''}</p>
<p>Playback Paused: {playbackPaused}</p>
<p>Position: {position} / {duration} min</p>
<p>Error Message: {errorMessage}</p>
