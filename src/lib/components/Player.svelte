<script lang="ts">
	import type { Track } from '$lib/models/Track';
	import * as api from '$lib/components/SearchTracks';
	import {
		queueStore,
		currentTrack as queueCurrentTrack,
		hasNext,
		hasPrevious,
	} from '$lib/stores/playerStore';
	import { SkipBack, SkipForward, Play, Pause, Volume2, VolumeX, Volume1, ListMusic } from 'lucide-svelte';

	let {
		currentTrack = $bindable(),
		isPlaying = $bindable(),
		audioElement = $bindable(),
		playTrack = $bindable(),
		togglePlayPause = $bindable(),
		showQueue = $bindable(),
	}: {
		currentTrack: Track | null;
		isPlaying: boolean;
		audioElement: HTMLAudioElement | null;
		playTrack?: (track: Track) => Promise<void>;
		togglePlayPause?: () => void;
		showQueue?: boolean;
	} = $props();

	let currentTime = $state(0);
	let duration = $state(0);
	let isDragging = $state(false);
	let volume = $state(70);
	let showVolumeSlider = $state(false);

	$effect(() => {
		currentTrack = $queueCurrentTrack;
	});

	async function handlePlayTrack(track: Track) {
		if (currentTrack?.id === track.id && isPlaying) {
			handleTogglePlayPause();
			return;
		}
		queueStore.playTrack(track);
		if (audioElement && track) {
			isPlaying = true;
			audioElement.src = await api.getTrackUrl(track.id);
			audioElement.volume = volume / 100;
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

	async function handleNext() {
		queueStore.next();
		if ($queueCurrentTrack && audioElement) {
			isPlaying = true;
			audioElement.src = await api.getTrackUrl($queueCurrentTrack.id);
			audioElement.volume = volume / 100;
			audioElement.play();
		}
	}

	async function handlePrevious() {
		queueStore.previous();
		if ($queueCurrentTrack && audioElement) {
			isPlaying = true;
			audioElement.src = await api.getTrackUrl($queueCurrentTrack.id);
			audioElement.volume = volume / 100;
			audioElement.play();
		}
	}

	async function handleTrackEnded() {
		if ($hasNext) {
			await handleNext();
		} else {
			isPlaying = false;
		}
	}

	function toggleQueue() {
		showQueue = !showQueue;
	}

	playTrack = handlePlayTrack;
	togglePlayPause = handleTogglePlayPause;

	function handleTimeUpdate() {
		if (!isDragging && audioElement) {
			currentTime = audioElement.currentTime;
		}
	}

	function handleVolumeChange(event: Event) {
		volume = parseFloat((event.target as HTMLInputElement).value);
		if (audioElement) audioElement.volume = volume / 100;
	}

	function handleLoadedMetadata() {
		if (audioElement) {
			duration = audioElement.duration;
			audioElement.volume = volume / 100;
		}
	}

	function handleSeek(event: Event) {
		if (audioElement) {
			const val = parseFloat((event.target as HTMLInputElement).value);
			const time = (val / 100) * duration;
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

	const progressPercent = $derived(duration > 0 ? (currentTime / duration) * 100 : 0);
</script>

{#if currentTrack}
	<div
		class="glass-strong fixed bottom-4 left-4 right-4 z-50 rounded-[1.25rem] px-4 py-3"
		style="max-width: min(1100px, calc(100% - 2rem)); margin: 0 auto; left: 50%; transform: translateX(-50%); width: calc(100% - 2rem);
			   box-shadow: 0 25px 60px rgba(2,6,23,0.5), 0 4px 16px rgba(15,23,42,0.4), inset 0 1px 0 rgba(255,255,255,0.05);"
	>
		<!-- Progress bar -->
		<div class="flex items-center gap-3 mb-2">
			<span class="text-[11px] text-[#64748b] min-w-[2.5rem] text-right tabular-nums">
				{formatTime(currentTime)}
			</span>
			<div class="flex-1 relative">
				<input
					type="range"
					min="0"
					max="100"
					value={progressPercent}
					class="w-full cursor-pointer"
					oninput={handleSeek}
					onmousedown={() => (isDragging = true)}
					onmouseup={() => (isDragging = false)}
					ontouchstart={() => (isDragging = true)}
					ontouchend={() => (isDragging = false)}
					style="background: linear-gradient(to right, #3b82f6 {progressPercent}%, rgba(148,163,184,0.2) {progressPercent}%);"
					aria-label="Progression"
				/>
			</div>
			<span class="text-[11px] text-[#64748b] min-w-[2.5rem] tabular-nums">
				{formatTime(duration)}
			</span>
		</div>

		<!-- Controls row -->
		<div class="flex items-center gap-3">
			<!-- Album cover + track info -->
			<div class="flex items-center gap-3 flex-1 min-w-0">
				<div class="w-10 h-10 rounded-lg shrink-0 overflow-hidden"
					style="box-shadow: 0 0 12px rgba(59,130,246,0.15);">
					{#if currentTrack.albumCover}
						<img src={currentTrack.albumCover} alt="Pochette" class="w-full h-full object-cover" />
					{:else}
						<div class="w-full h-full glass-subtle flex items-center justify-center text-lg">🎵</div>
					{/if}
				</div>
				<div class="min-w-0">
					<p class="text-sm font-medium text-[#f8fafc] truncate">{currentTrack.title}</p>
					<p class="text-xs text-[#94a3b8] truncate">{currentTrack.artist}</p>
				</div>
			</div>

			<!-- Playback controls -->
			<div class="flex items-center gap-2 shrink-0">
				<button
					onclick={handlePrevious}
					disabled={!$hasPrevious}
					class="btn-glass btn-icon touch-target"
					aria-label="Piste précédente"
				>
					<SkipBack size={18} class={$hasPrevious ? 'text-[#f8fafc]' : 'text-[#475569]'} />
				</button>

				<button
					onclick={handleTogglePlayPause}
					class="btn-glass btn-glass-accent btn-icon touch-target"
					aria-label={isPlaying ? 'Pause' : 'Lecture'}
				>
					{#if isPlaying}
						<Pause size={20} />
					{:else}
						<Play size={20} />
					{/if}
				</button>

				<button
					onclick={handleNext}
					disabled={!$hasNext}
					class="btn-glass btn-icon touch-target"
					aria-label="Piste suivante"
				>
					<SkipForward size={18} class={$hasNext ? 'text-[#f8fafc]' : 'text-[#475569]'} />
				</button>
			</div>

			<!-- Volume + Queue -->
			<div class="flex items-center gap-2 shrink-0">
				<!-- Volume -->
				<div class="relative">
					<button
						onclick={() => (showVolumeSlider = !showVolumeSlider)}
						class="btn-glass btn-icon touch-target"
						aria-label="Volume"
					>
						{#if volume === 0}
							<VolumeX size={16} class="text-[#94a3b8]" />
						{:else if volume < 50}
							<Volume1 size={16} class="text-[#94a3b8]" />
						{:else}
							<Volume2 size={16} class="text-[#94a3b8]" />
						{/if}
					</button>
					{#if showVolumeSlider}
						<div
							class="glass-strong absolute bottom-full mb-2 right-0 p-3 rounded-xl"
							style="min-width: 8rem; animation: slide-up 0.15s ease;"
						>
							<input
								type="range"
								min="0"
								max="100"
								value={volume}
								class="w-full cursor-pointer"
								oninput={handleVolumeChange}
								style="background: linear-gradient(to right, #3b82f6 {volume}%, rgba(148,163,184,0.2) {volume}%);"
								aria-label="Volume"
							/>
						</div>
					{/if}
				</div>

				<!-- Queue toggle -->
				<button
					onclick={toggleQueue}
					class="btn-glass btn-icon touch-target {showQueue ? 'btn-glass-accent' : ''}"
					aria-label="File d'attente"
				>
					<ListMusic size={16} />
				</button>
			</div>
		</div>
	</div>
{/if}

<audio
	bind:this={audioElement}
	onplay={() => (isPlaying = true)}
	onpause={() => (isPlaying = false)}
	onended={handleTrackEnded}
	ontimeupdate={handleTimeUpdate}
	onloadedmetadata={handleLoadedMetadata}
	class="hidden"
></audio>
