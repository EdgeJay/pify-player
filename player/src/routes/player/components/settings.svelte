<script lang="ts">
	import { onMount } from 'svelte';

	interface Props {
		onClose: () => void;
	}

	let { onClose }: Props = $props();

	onMount(async () => {});

	const onSave = async () => {
		// Save settings logic here
		console.log('Settings saved');
	};
</script>

<div class="settings-dialog">
	<div class="settings-panel">
		<div class="settings-header">
			<h2>Settings</h2>
			<button aria-label="Close" onclick={onClose}>
				<i class="fa-solid fa-xmark"></i>
			</button>
		</div>
		<div class="settings-content"></div>
		<div class="settings-footer">
			<button onclick={onSave}>Save</button>
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
	}

	.settings-content {
		padding: 20px 0px;
	}

	.settings-footer {
		display: flex;
		flex-direction: row;
		justify-content: center;
	}

	.settings-header h2 {
		flex-grow: 1;
	}
</style>
