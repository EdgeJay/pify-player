<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		onClose: () => void;
	}

	let { onClose }: Props = $props();

	onMount(async () => {});

	const onShutdown = async () => {
		// Shutdown logic here
		console.log('System shutting down');
		onClose();
	};

	const onRestart = async () => {
		// Shutdown logic here
		console.log('System shutting down');
		onClose();
	};

	const onRefreshWindow = () => {
		window.location.reload();
	};

	const onSave = async () => {
		// Save settings logic here
		console.log('Settings saved');
		onClose();
	};
</script>

<div class="settings-dialog">
	<div class="settings-panel">
		<div class="settings-header">
			<h2 class="text-xl font-bold">Settings</h2>
			<button aria-label="Close" onclick={onClose}>
				<i class="fa-solid fa-xmark"></i>
			</button>
		</div>
		<div class="settings-content flex flex-col">
			<!-- System settings -->
			<h3 class="text-lg font-semibold">System settings</h3>
			<div class="flex flex-row">
				<div class="label"><span>Power Controls</span></div>
				<div class="flex grow flex-row justify-evenly">
					<button
						class="rounded-full bg-red-700 px-5 py-2 text-sm font-semibold leading-5 text-white"
						onclick={onShutdown}>Shutdown</button
					>
					<button
						class="rounded-full bg-amber-500 px-5 py-2 text-sm font-semibold leading-5 text-white"
						onclick={onRestart}>Restart</button
					>
				</div>
			</div>
			<!-- App settings -->
			<h3 class="mt-12 text-lg font-semibold">App settings</h3>
			<div class="flex flex-row">
				<div class="label"><span>Debug</span></div>
				<div class="flex grow flex-row justify-evenly">
					<button
						class="rounded-full border border-sky-500 px-5 py-2 text-sm font-semibold leading-5"
						onclick={onRefreshWindow}>Refresh Window</button
					>
				</div>
			</div>
		</div>
		<div class="settings-footer">
			<button
				class="rounded-full bg-sky-500 px-5 py-2 text-sm font-semibold leading-5 text-white"
				onclick={onSave}>Save</button
			>
		</div>
	</div>
</div>

<style>
	.settings-dialog {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.5);
		z-index: 200;
	}

	.settings-panel {
		position: fixed;
		display: flex;
		flex-direction: column;
		top: 50%;
		left: 50%;
		width: 70%;
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

	.settings-header {
		display: flex;
		flex-direction: row;
		border-bottom: var(--color-gray-200) 1px solid;
	}

	.settings-content {
		padding: 20px 0px;
	}

	.settings-content .label {
		display: flex;
		flex-direction: column;
		justify-content: center;
		width: 120px;
		font-size: var(--text-base);
	}

	.settings-content h3 {
		border-bottom: var(--color-gray-200) 1px solid;
		margin-bottom: calc(var(--spacing) * 4);
	}

	.settings-footer {
		display: flex;
		flex-direction: row;
		justify-content: center;
		padding-top: 10px;
		border-top: var(--color-gray-200) 1px solid;
	}

	.settings-header h2 {
		flex-grow: 1;
	}
</style>
