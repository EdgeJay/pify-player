<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		basicAuthToken: string;
	}

	interface LoginQRResponse {
		data: {
			qr: string;
		};
	}

	let { basicAuthToken }: Props = $props();
	let isLoading = $state(true);
	let imageDataUrl = $state('');

	onMount(async () => {
		// get login QR code
		const domain = window.location.hostname;
		const response = await fetch(`https://${domain}:8080/api/player/login-qr`, {
			method: 'GET',
			headers: {
				Authorization: `Basic ${basicAuthToken}`
			}
		});

		if (!response.ok) {
			throw new Error('get login QR failed');
		}

		const data = (await response.json()) as LoginQRResponse;
		imageDataUrl = data.data.qr;
		isLoading = false;
	});

	const onRefreshButton = () => {
		window.location.reload();
	};
</script>

<div class="login-dialog">
	<div class="login-panel">
		{#if isLoading}
			<i class="fa fa-circle-notch"></i>
		{:else}
			<img src={imageDataUrl} alt="Login QR Code" />
			<button onclick={onRefreshButton}>Refresh</button>
		{/if}
	</div>
</div>

<style>
	.login-dialog {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.5);
		z-index: 100;
	}

	.login-panel {
		position: fixed;
		display: flex;
		flex-direction: column;
		top: 50%;
		left: 50%;
		transform: translate(-50%, -50%);
		background-color: white;
		color: #585858;
		padding: 20px;
		border-radius: 8px;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
		z-index: 101;
		transition:
			width 0.3s ease-in-out,
			height 0.3s ease-in-out;
	}

	i.fa-circle-notch {
		animation: spin 1s linear infinite;
	}

	@keyframes spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
</style>
