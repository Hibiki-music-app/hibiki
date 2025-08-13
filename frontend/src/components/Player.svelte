<script lang="ts">
	import type { Track } from '$lib/api';
	import * as api from '$lib/api';

	let {
		currentTrack = $bindable(),
		isPlaying = $bindable(),
		audioElement = $bindable(),
		playTrack = $bindable(),
		togglePlayPause = $bindable()
	}: {
		currentTrack: Track | null;
		isPlaying: boolean;
		audioElement: HTMLAudioElement | null;
		playTrack?: (track: Track) => Promise<void>;
		togglePlayPause?: () => void;
	} = $props();

	// États pour la progression temporelle
	let currentTime = $state(0);
	let duration = $state(0);
	let isDragging = $state(false);

	async function handlePlayTrack(track: Track) {
		if (currentTrack?.id === track.id && isPlaying) {
			handleTogglePlayPause();
			return;
		}
		currentTrack = track;
		if (audioElement) {
			isPlaying = true;
			audioElement.src = await api.getTrackUrl(track.id);
			audioElement.play();
		}
	}

	function handleTogglePlayPause() {
		if (!audioElement) return;

		if (isPlaying) {
			audioElement.pause();
			isPlaying = false;
		} else {
			audioElement.play();
			isPlaying = true;
		}
	}

	// Assigner les fonctions aux props
	playTrack = handlePlayTrack;
	togglePlayPause = handleTogglePlayPause;

	// Fonctions pour gérer la progression temporelle
	function handleTimeUpdate() {
		if (!isDragging && audioElement) {
			currentTime = audioElement.currentTime;
		}
	}

	function handleLoadedMetadata() {
		if (audioElement) {
			duration = audioElement.duration;
		}
	}

	function handleSeek(event: Event) {
		if (audioElement) {
			const target = event.target as HTMLInputElement;
			const time = (parseFloat(target.value) / 100) * duration;
			audioElement.currentTime = time;
			currentTime = time;
		}
	}

	function formatTime(seconds: number): string {
		if (!seconds || !isFinite(seconds)) return '0:00';
		const mins = Math.floor(seconds / 60);
		const secs = Math.floor(seconds % 60);
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}
</script>

{#if currentTrack}
	<div class="btm-nav bg-base-100 shadow-lg">
		<div class="flex flex-col w-full px-4 py-2">
			<!-- Barre de progression -->
			<div class="flex items-center gap-2 mb-2">
				<span class="text-xs opacity-60 min-w-[2.5rem]">
					{formatTime(currentTime)}
				</span>
				<input
					type="range"
					min="0"
					max="100"
					value={duration > 0 ? (currentTime / duration) * 100 : 0}
					class="range range-xs flex-1"
					oninput={handleSeek}
					onmousedown={() => (isDragging = true)}
					onmouseup={() => (isDragging = false)}
					ontouchstart={() => (isDragging = true)}
					ontouchend={() => (isDragging = false)}
				/>
				<span class="text-xs opacity-60 min-w-[2.5rem]">
					{formatTime(duration)}
				</span>
			</div>

			<!-- Contrôles du player -->
			<div class="flex items-center gap-3 w-full">
				<!-- Album cover -->
				<div class="avatar">
					<div class="w-10 h-10 rounded">
						{#if currentTrack.albumCover}
							<img src={currentTrack.albumCover} alt="Album cover" />
						{:else}
							<div class="bg-base-300 flex items-center justify-center">🎵</div>
						{/if}
					</div>
				</div>

				<!-- Track info -->
				<div class="flex-1 min-w-0">
					<div class="text-sm font-medium truncate">{currentTrack.title}</div>
					<div class="text-xs opacity-70 truncate">{currentTrack.artist}</div>
				</div>

				<!-- Play/Pause button -->
				<button onclick={handleTogglePlayPause} class="btn btn-ghost btn-square">
					{#if isPlaying}
						<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
							<rect x="6" y="4" width="4" height="16" />
							<rect x="14" y="4" width="4" height="16" />
						</svg>
					{:else}
						<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
							<polygon points="5,3 19,12 5,21" />
						</svg>
					{/if}
				</button>
			</div>
		</div>
	</div>
{/if}

<audio
	bind:this={audioElement}
	onplay={() => (isPlaying = true)}
	onpause={() => (isPlaying = false)}
	onended={() => (isPlaying = false)}
	ontimeupdate={handleTimeUpdate}
	onloadedmetadata={handleLoadedMetadata}
	class="hidden"
></audio>

<style>
	/* Améliorations pour le player */
	.btm-nav {
		position: fixed;
		bottom: 1rem;
		left: 1rem;
		right: 1rem;
		width: calc(100% - 2rem);
		height: auto;
		min-height: 5rem;
		padding: 0.75rem 0;
		z-index: 1000;
		border-radius: 0.7rem;
	}
</style>
