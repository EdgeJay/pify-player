<script lang="ts">
	import { onMount } from 'svelte';
	import { checkSession } from '$lib/session';
	import { getAllDevices } from '$lib/device';

	let loggedIn = $state(false);
	let displayName = $state('');

	interface RemoteResponse {
		success: boolean;
	}

	onMount(async () => {
		// check login status first
		try {
			const { logged_in, redirect_url, user } = await checkSession();
			loggedIn = logged_in;

			if (!logged_in) {
				window.location.href = redirect_url;
			} else {
				console.log('Already logged in');
				displayName = user.display_name;

				// get all devices
				const devices = await getAllDevices();
				console.log(devices);
			}
		} catch (err) {
			console.error(err);
		}
	});

	/*
	const onPlayButton = async () => {
		const DOMAIN = window.location.hostname;
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
    */
</script>

<div class="home">
	<div class="home-container">
		<header class="home-header">
			<h1>Pify Player</h1>
			<p>Welcome back {displayName}</p>
		</header>
		<section></section>
	</div>
</div>

<style>
	.home {
		display: flex;
		flex-direction: column;
		justify-content: flex-start;
		align-items: center;
		height: 100%;
		padding-bottom: 50px;
	}
	.home-container {
		flex: 0;
		background-color: #1a1a1a;
		border: 2px solid #00ff00;
		border-radius: 10px;
		padding: 20px;
		margin-top: 100px;
		width: 90%;
		box-shadow: 0 0 20px #00ff00;
	}
	@media screen and (min-width: 800px) {
		.home-container {
			max-width: 640px;
		}
	}
	.home-header {
		text-align: center;
	}
</style>
