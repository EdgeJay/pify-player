<script lang="ts">
	interface Props {
		ytVidId?: string;
	}

	let { ytVidId }: Props = $props();
</script>

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

<style>
	.video-bg {
		position: absolute;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		z-index: 1;
	}
</style>
