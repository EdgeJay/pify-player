<script lang="ts">
	import { onMount } from 'svelte';
	import { getApiConnectWS, sendConnectCommand, clearSpotifyTokenFromStorage } from '$lib/ws';

	let { data } = $props();
	let deviceId = $state('');
	let spotifyTrack: Spotify.Track | undefined = $state();
	let playbackPaused = $state(true);
	let position = $state(0);
	let duration = $state(0);
	let errorMessage = $state('');

	let player: Spotify.Player;
	let token = '';

	const convertToMinutes = (ms: number): number => {
		const minutes = Math.floor((ms / 60000) * 100) / 100;
		return minutes;
	};

	onMount(() => {
		window.onSpotifyWebPlaybackSDKReady = async () => {
			player = new Spotify.Player({
				name: data.playerName,
				getOAuthToken: (cb) => {
					cb(token);
				},
				volume: 0.5
			});
			/* Activates an HTML element in the player instance. This is typically required before any media playback
				can occur, especially in browsers that enforce user interaction before allowing audio/video playback. */
			player.activateElement();

			// establish WebSocket connection
			const ws = getApiConnectWS({
				onConnect: (accessToken: string) => {
					token = accessToken;
					player.connect();
				},
				onError: (message) => {
					errorMessage = message;
				}
			});

			const refreshAccessToken = (ws: WebSocket) => {
				// Clear expired token
				clearSpotifyTokenFromStorage();
				sendConnectCommand(ws);
			};

			// Player Ready
			player.addListener('ready', ({ device_id }) => {
				console.log('Ready with Device ID', device_id);
				deviceId = device_id;
			});

			// Player Not Ready
			player.addListener('not_ready', ({ device_id }) => {
				console.log('Device ID has gone offline', device_id);
				errorMessage = `Device ID has gone offline: ${device_id}`;
			});

			player.addListener('initialization_error', ({ message }) => {
				errorMessage = message;
			});

			player.addListener('playback_error', ({ message }) => {
				errorMessage = message;
				if (message.includes('token expired')) {
					refreshAccessToken(ws);
				}
			});

			player.addListener('authentication_error', ({ message }) => {
				console.log('authentication_error');
				errorMessage = message;
				refreshAccessToken(ws);
			});

			player.addListener('account_error', ({ message }) => {
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

<div class="player">
	<p>{errorMessage}</p>
	<div class="panel">
		<div class="album"></div>
		<div class="song"></div>
		<div class="controls">
			<button class="sm" aria-label="Volume">
				<i class="fa fa-volume-high"></i>
			</button>
			<button class="sm" aria-label="Shuffle">
				<i class="fa fa-shuffle"></i>
			</button>
			<button aria-label="Previous">
				<i class="fa fa-backward"></i>
			</button>
			<button class="play" aria-label="Play">
				<i class="fa fa-play"></i>
			</button>
			<button aria-label="Next">
				<i class="fa fa-forward"></i>
			</button>
			<button class="sm" aria-label="Repeat">
				<i class="fa fa-repeat"></i>
			</button>
			<button class="sm" aria-label="Playlist">
				<i class="fa fa-list"></i>
			</button>
		</div>
	</div>
</div>

<!-- 
<div style="margin-top:100px">
	<p>Device ID: {deviceId}</p>
	<p>Song Title: {spotifyTrack?.name || ''}</p>
	<p>Playback Paused: {playbackPaused}</p>
	<p>Position: {position} / {duration} min</p>
	<p>Error Message: {errorMessage}</p>
</div>
-->

<style>
	.player {
		display: flex;
		flex-direction: column;
		justify-content: flex-end;
		align-items: center;
		height: 100%;
		padding: 50px 20px;
	}

	.panel {
		position: relative;
		background-color: #fff;
		box-shadow: 0 30px 80px #656565;
		border-radius: 15px;
		padding: 20px 30px;
		min-width: 600px;
	}

	.album {
		position: absolute;
		top: -20px;
		left: 20px;
		width: 100px;
		height: 100px;
		background-color: #585858;
		border: 5px solid #fff;
		border-radius: 10px;
		z-index: 10;
	}

	.song {
		height: 60px;
		padding-left: 110px;
		margin-bottom: 15px;
	}

	.controls {
		display: flex;
		flex-direction: row;
		justify-content: space-between;
		align-items: center;
	}

	.controls button {
		color: #585858;
		font-size: 30px;
	}

	.controls button.play {
		background-color: #585858;
		color: #fff;
		box-shadow: 0 10px 20px #656565;
		width: 60px;
		height: 60px;
		border-radius: 30px;
	}

	.controls button.sm {
		font-size: 18px;
	}
</style>
