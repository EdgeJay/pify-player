<script lang="ts">
	import { onMount } from 'svelte';
	import { getApiConnectWS, sendConnectCommand, clearSpotifyTokenFromStorage } from '$lib/ws';
	import { controlPlayback } from '$lib/device';

	let { data } = $props();
	let deviceId = $state('');

	// song track and playback status related vars
	let spotifyTrack: Spotify.Track | undefined = $state();
	let playbackPaused = $state(true);
	let position = $state('');
	let duration = $state('');
	let albumImage = $state('');
	let songTitle = $state('');
	let songArtists = $state<string[]>([]);
	let songProgress = $state(0);

	let errorMessage = $state('');

	let player: Spotify.Player;
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
			player.addListener('ready', async ({ device_id }) => {
				console.log('Ready with Device ID', device_id);
				errorMessage = '';
				deviceId = device_id;

				// check playback state, take over as playback controls if state is null
				const state = await player.getCurrentState();
				if (!state) {
					const success = await controlPlayback(token, deviceId);
					if (success) {
						console.log('took over playback control successfully');
					}
				}
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
				console.log('player_state_changed');

				spotifyTrack = current_track;
				$state.snapshot(spotifyTrack);

				// update song details
				if (spotifyTrack.album.images.length > 0) {
					albumImage = spotifyTrack.album.images[0].url;
				}
				songTitle = spotifyTrack.name;
				songArtists = spotifyTrack.artists.reduce((acc, artist) => {
					acc.push(artist.name);
					return acc;
				}, [] as string[]);

				playbackPaused = paused;

				if (intervalId) {
					clearInterval(intervalId);
				}

				if (!paused) {
					intervalId = setInterval(async () => {
						if (!playbackPaused) {
							await updatePlaybackPosition();
						}
					}, 1000);
				} else {
					if (intervalId) {
						// console.log('cleared interval');
						clearInterval(intervalId);
					}
				}
			});
		};

		const script = document.createElement('script');
		script.src = 'https://sdk.scdn.co/spotify-player.js';
		document.body.appendChild(script);
	});

	/* Playback controls */
	const updatePlaybackPosition = async () => {
		const state = await player.getCurrentState();
		if (state) {
			position = convertToMinutes(state.position);
			duration = convertToMinutes(state.duration);
			songProgress = Math.ceil((state.position / state.duration) * 100);
		}
	};

	const onPlay = () => {
		player.togglePlay();
	};

	const onNext = () => {
		player.nextTrack();
	};

	const onPrev = () => {
		player.previousTrack();
	};
</script>

<div class="player">
	<p>{errorMessage}</p>
	<div class="panel">
		<div class="album">
			<img src={albumImage} alt={songTitle} />
		</div>
		<div class="song">
			<h1>{songTitle}</h1>
			<p>{songArtists.join(', ')}</p>
		</div>
		<div class="progress">
			<span>{position}</span>
			<input type="range" step="1" style="--value:{songProgress};" value={songProgress} />
			<span>{duration}</span>
		</div>
		<div class="controls">
			<button class="sm" aria-label="Volume">
				<i class="fa fa-volume-high"></i>
			</button>
			<button class="sm" aria-label="Shuffle">
				<i class="fa fa-shuffle"></i>
			</button>
			<button aria-label="Previous" onclick={onPrev}>
				<i class="fa fa-backward"></i>
			</button>
			<button class="play" onclick={onPlay} aria-label="Play">
				<i class="fa {playbackPaused ? 'fa-play' : 'fa-pause'}"></i>
			</button>
			<button aria-label="Next" onclick={onNext}>
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

	.album img {
		width: 100%;
		height: 100%;
	}

	.song {
		color: #585858;
		height: 60px;
		padding-left: 110px;
		margin-bottom: 10px;
	}

	.progress {
		display: flex;
		flex-flow: row;
		justify-content: space-between;
		align-items: center;
		height: 25px;
		margin-bottom: 20px;
	}

	.progress span {
		color: #585858;
		font-size: 14px;
		text-align: center;
	}

	.progress input[type='range'] {
		--min: 0;
		--max: 100;
		--range: calc(var(--max) - var(--min));
		--ratio: calc((var(--value) - var(--min)) / var(--range));
		--sx: calc(0.5 * 7px + var(--ratio) * (100% - 7px));
	}

	.progress input[type='range'] {
		appearance: none;
		-webkit-appearance: none;
		width: 100%;
		height: 7px;
		background: #9a9a9a;
		border-radius: 4px;
		margin: 0 15px;
	}

	.progress input[type='range']::-webkit-slider-thumb {
		margin-top: -4px;
		appearance: none;
		-webkit-appearance: none;
		background: #585858;
		width: 16px;
		aspect-ratio: 1/1;
		border-radius: 50%;
		outline: 2px solid #fff;
		box-shadow: 0 6px 10px rgba(5, 36, 28, 0.3);
	}

	.progress input[type='range']::-webkit-slider-runnable-track {
		height: 7px;
		border: none;
		border-radius: 4px;
		background:
			linear-gradient(#585858, #585858) 0 / var(--sx) 100% no-repeat,
			#9a9a9a;
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
