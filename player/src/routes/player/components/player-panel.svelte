<script lang="ts">
	interface Props {
		player?: Spotify.Player;
		errorMessage?: string;
		playbackPaused?: boolean;
		position?: string;
		duration?: string;
		albumImage?: string;
		songTitle?: string;
		songArtists?: string[];
		songProgress?: number;
		volume: number;
	}

	let {
		player,
		errorMessage = '',
		playbackPaused = true,
		position = '',
		duration = '',
		albumImage = '',
		songTitle = '',
		songArtists = [],
		songProgress = 0,
		volume
	}: Props = $props();

	// volume controls
	let isTogglingVolume = $state(false);

	const onPlay = () => {
		player?.togglePlay();
	};

	const onNext = () => {
		player?.nextTrack();
	};

	const onPrev = () => {
		player?.previousTrack();
	};

	const onSeek = async (evt: Event) => {
		const state = await player?.getCurrentState();
		if (!state) {
			return;
		}
		const target = evt.target as HTMLInputElement;
		if (!target) {
			return;
		}
		const value = parseInt(target.value);
		player!.seek(state.duration * (value / 100));
	};

	const toggleVolume = () => {
		isTogglingVolume = !isTogglingVolume;
	};

	const onVolumeChange = async (evt: Event) => {
		const target = evt.target as HTMLInputElement;
		if (!target) {
			return;
		}
		const value = parseInt(target.value);
		volume = value;
		player?.setVolume(value / 100);
	};
</script>

<div class="player">
	<p>{errorMessage}</p>
	<div class="panel">
		<div class="album">
			{#if albumImage}
				<img src={albumImage} alt={songTitle} />
			{/if}
		</div>
		<div class="song">
			<h1>{songTitle}</h1>
			<p>{songArtists.join(', ')}</p>
		</div>
		<div class="progress">
			<span>{position}</span>
			<input
				type="range"
				step="1"
				style="--value:{songProgress};"
				value={songProgress}
				onchange={onSeek}
			/>
			<span>{duration}</span>
		</div>
		<div class="controls">
			<button class="sm" aria-label="Volume" onclick={toggleVolume}>
				<i class="fa fa-volume-high"></i>
			</button>
			{#if isTogglingVolume}
				<div class="volume-panel">
					<div class="progress no-margins full-width">
						<input
							type="range"
							step="1"
							style="--value:{volume};"
							value={volume}
							onchange={onVolumeChange}
						/>
						<span>{volume}</span>
					</div>
				</div>
			{:else}
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
			{/if}
		</div>
	</div>
</div>

<style>
	.player {
		position: absolute;
		display: flex;
		flex-direction: column;
		justify-content: flex-end;
		align-items: center;
		width: 100%;
		height: 100%;
		padding: 20px;
		z-index: 20;
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

	.volume-panel {
		display: flex;
		flex-direction: row;
		flex: 1;
		align-items: center;
		height: 60px;
	}

	.progress {
		display: flex;
		flex-flow: row;
		justify-content: space-between;
		align-items: center;
		height: 25px;
		margin-bottom: 20px;
	}

	.progress.no-margins {
		margin: 0;
	}

	.progress.full-width {
		flex: 1;
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
