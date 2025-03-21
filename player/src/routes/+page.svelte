<script lang="ts">
	import { onMount } from 'svelte';
	import { checkSession } from '$lib/session';
	import { getAllDevices } from '$lib/device';
	import type { Device } from '$lib/device';

	let { data } = $props();
	let loggedIn = $state(false);
	let displayName = $state('');
	let devices = $state<Device[]>([]);

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
				devices = (await getAllDevices()).data.devices;
				console.log(devices);
			}
		} catch (err) {
			console.error(err);
		}
	});

	const onRefreshButton = async () => {
		// get all devices
		devices = (await getAllDevices()).data.devices;
	};

	const onLogoutButton = async () => {
		console.log('logout!');
	};
</script>

<div class="home">
	<div class="home-container">
		<header class="home-header">
			<h1>{data.playerName}</h1>
			<p>Welcome back {displayName}</p>
		</header>
		<section class="home-controls">
			<p>Devices:</p>
			{#if devices.length > 0}
				<ul>
					{#each devices as device}
						<li>{device.type} | {device.name}</li>
					{/each}
				</ul>
			{:else}
				<p>No devices found</p>
			{/if}

			<button onclick={onRefreshButton} style="margin-top:50px;">Refresh</button>
			<button onclick={onLogoutButton} style="margin-top:15px;">Logout</button>
		</section>
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

	.home-controls {
		display: flex;
		flex-direction: column;
		justify-content: flex-start;
		align-items: center;
		margin-top: 20px;
	}
</style>
