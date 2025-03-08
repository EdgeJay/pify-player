<script lang="ts">
	import { onMount } from 'svelte';
	import { DOMAIN } from '$env/static/private';

	let loggedIn = $state(false);

	interface LoginResponse {
		logged_in: boolean;
		redirect_url: string;
	}

	interface RemoteResponse {
		success: boolean;
	}

	onMount(async () => {
		// check login status first
		const response = await fetch(`https://${DOMAIN}:8080/api/auth/login`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		if (!response.ok) {
			throw new Error('Login failed');
		}

		const { logged_in, redirect_url } = (await response.json()) as LoginResponse;
		loggedIn = logged_in;

		if (!logged_in) {
			window.location.href = redirect_url;
		} else {
			console.log('Already logged in');
		}
	});

	const onPlayButton = async () => {
		const response = await fetch(`https://${DOMAIN}:8080/api/remote/play`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});
		const { success } = (await response.json()) as RemoteResponse;
		console.log(`success: ${success}`);
	};
</script>

<div class="player">
	<div class="equalizer">
		<div class="bar"></div>
		<div class="bar"></div>
		<div class="bar"></div>
		<div class="bar"></div>
		<div class="bar"></div>
	</div>
	<div class="player-container">
		<header class="player-header">
			<h1>Cyberpunk Music Player</h1>
		</header>
		<section class="player-controls">
			<button disabled={!loggedIn}>&#9664;&#9664;</button>
			<button onclick={onPlayButton} disabled={!loggedIn}>&#9654;</button>
			<!-- play button -->
			<button disabled={!loggedIn}>&#10074;&#10074;</button>
			<button disabled={!loggedIn}>&#9654;&#9654;</button>
		</section>
		<section class="player-progress">
			<input type="range" min="0" max="100" value="0" />
		</section>
	</div>
</div>

<style>
	.player {
		display: flex;
		flex-direction: column;
		justify-content: flex-end;
		align-items: center;
		height: 100%;
		padding-bottom: 50px;
	}
	.player-container {
		flex: 0;
		background-color: #1a1a1a;
		border: 2px solid #00ff00;
		border-radius: 10px;
		padding: 20px;
		width: 90%;
		box-shadow: 0 0 20px #00ff00;
	}
	@media screen and (min-width: 800px) {
		.player-container {
			max-width: 640px;
		}
	}
	.player-header {
		text-align: center;
		margin-bottom: 20px;
	}
	.player-controls {
		display: flex;
		justify-content: space-around;
		align-items: center;
	}
	.player-controls button {
		background-color: #00ff00;
		border: none;
		border-radius: 50%;
		width: 50px;
		height: 50px;
		color: #0f0f0f;
		font-size: 18px;
		cursor: pointer;
		transition: background-color 0.3s;
	}
	.player-controls button:hover {
		background-color: #00cc00;
	}
	.player-progress {
		margin-top: 20px;
	}
	.player-progress input[type='range'] {
		width: 100%;
		-webkit-appearance: none;
		appearance: none;
		background: #005500; /* Slightly lighter shade of dark green */
		height: 5px;
		border-radius: 5px;
		outline: none;
	}
	.player-progress input[type='range']::-webkit-slider-thumb {
		-webkit-appearance: none;
		appearance: none;
		width: 15px;
		height: 15px;
		background: #128912; /* Dark green color */
		border-radius: 50%;
		cursor: pointer;
	}
	.player-progress input[type='range']::-moz-range-thumb {
		width: 15px;
		height: 15px;
		background: #006400; /* Dark green color */
		border-radius: 50%;
		cursor: pointer;
	}

	.equalizer {
		display: flex;
		flex: 1;
		justify-content: space-around;
		align-items: center;
		margin-bottom: 50px;
	}
	.equalizer .bar {
		width: 10px;
		height: 50px;
		background-color: #00ff00;
		animation: equalizer 1s infinite;
	}
	.equalizer .bar:nth-child(1) {
		background-color: #ff4000;
		animation-delay: 0s;
	}
	.equalizer .bar:nth-child(2) {
		background-color: #ff7300;
		animation-delay: 0.2s;
	}
	.equalizer .bar:nth-child(3) {
		background-color: #ffd000;
		animation-delay: 0.4s;
	}
	.equalizer .bar:nth-child(4) {
		animation-delay: 0.6s;
	}
	.equalizer .bar:nth-child(5) {
		background-color: #00aeff;
		animation-delay: 0.8s;
	}
	@keyframes equalizer {
		0%,
		100% {
			transform: scaleY(1);
		}
		50% {
			transform: scaleY(2);
		}
	}
</style>
