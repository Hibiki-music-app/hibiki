<script lang="ts">
	import type { Track } from '$lib/models/Track';
	import * as api from '$lib/components/SearchTracks';
	import {
		queueStore,
		currentTrack as queueCurrentTrack,
		hasNext,
		hasPrevious
	} from '$lib/stores/playerStore';

	let {
		currentTrack = $bindable(),
		isPlaying = $bindable(),
		audioElement = $bindable(),
		playTrack = $bindable(),
		togglePlayPause = $bindable(),
		showQueue = $bindable()
	}: {
		currentTrack: Track | null;
		isPlaying: boolean;
		audioElement: HTMLAudioElement | null;
		playTrack?: (track: Track) => Promise<void>;
		togglePlayPause?: () => void;
		showQueue?: boolean;
	} = $props();

	// États pour la progression temporelle
	let currentTime = $state(0);
	let duration = $state(0);
	let isDragging: boolean = $state(false);
	let volume: number = $state(50);
	let showVolumeSlider: boolean = $state(false);

	// Synchroniser avec le store de queue
	$effect(() => {
		currentTrack = $queueCurrentTrack;
		if (audioElement && $queueCurrentTrack) {
			api
				.getMusicUrl($queueCurrentTrack.id)
				.then((url) => {
					if (audioElement) {
						console.log('Setting audio source:', url);
						audioElement.src = url;
						audioElement.volume = volume / 100;
						// Si on était en train de jouer, reprendre la lecture
						if (isPlaying) {
							audioElement.play().catch((error) => {
								console.error('Error playing audio:', error);
								isPlaying = false;
							});
						}
					}
				})
				.catch((error) => {
					console.error('Error loading music URL:', error);
				});
		}
	});

	async function handlePlayTrack(track: Track) {
		if (currentTrack?.id === track.id && isPlaying) {
			handleTogglePlayPause();
			return;
		}

		// Utiliser le store de queue pour jouer la piste
		queueStore.playTrack(track);

		// La piste sera automatiquement mise à jour via l'effect
		if (audioElement && track) {
			isPlaying = true;
			audioElement.volume = volume / 100;
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

	// Nouvelle fonction pour piste suivante
	async function handleNext() {
		queueStore.next();
		isPlaying = true;
	}

	// Nouvelle fonction pour piste précédente
	async function handlePrevious() {
		queueStore.previous();
		isPlaying = true;
	}

	// Nouvelle fonction pour gérer la fin de piste
	async function handleTrackEnded() {
		if ($hasNext) {
			await handleNext();
		} else {
			isPlaying = false;
		}
	}

	// Fonction pour basculer l'affichage de la queue
	function toggleQueue() {
		showQueue = !showQueue;
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

	function handleVolumeChange(event: Event) {
		const target = event.target as HTMLInputElement;
		volume = parseFloat(target.value);
		if (audioElement) {
			audioElement.volume = volume / 100; // Convertir en 0-1
		}
	}

	function toggleVolumeSlider() {
		showVolumeSlider = !showVolumeSlider;
	}

	function handleLoadedMetadata() {
		if (audioElement) {
			duration = audioElement.duration;
			audioElement.volume = volume / 100; // Convertir en 0-1
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

				<!-- Volume control -->
				<div class="relative">
					<button onclick={toggleVolumeSlider} class="btn btn-ghost btn-square btn-sm">
						{#if volume === 0}
							<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
								<path
									d="M16.5 12c0-1.77-1.02-3.29-2.5-4.03v2.21l2.45 2.45c.03-.2.05-.41.05-.63zm2.5 0c0 .94-.2 1.82-.54 2.64l1.51 1.51C20.63 14.91 21 13.5 21 12c0-4.28-2.99-7.86-7-8.77v2.06c2.89.86 5 3.54 5 6.71zM4.27 3L3 4.27 7.73 9H3v6h4l5 5v-6.73l4.25 4.25c-.67.52-1.42.93-2.25 1.18v2.06c1.38-.31 2.63-.95 3.69-1.81L19.73 21 21 19.73l-9-9L4.27 3zM12 4L9.91 6.09 12 8.18V4z"
								/>
							</svg>
						{:else if volume < 50}
							<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
								<path
									d="M18.5 12c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM5 9v6h4l5 5V4L9 9H5z"
								/>
							</svg>
						{:else}
							<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
								<path
									d="M3 9v6h4l5 5V4L7 9H3zm13.5 3c0-1.77-1.02-3.29-2.5-4.03v8.05c1.48-.73 2.5-2.25 2.5-4.02zM14 3.23v2.06c2.89.86 5 3.54 5 6.71s-2.11 5.85-5 6.71v2.06c4.01-.91 7-4.49 7-8.77s-2.99-7.86-7-8.77z"
								/>
							</svg>
						{/if}
					</button>
					{#if showVolumeSlider}
						<div
							class="absolute bottom-full mb-2 left-1/2 transform -translate-x-1/2 bg-base-100 p-2 rounded shadow-lg"
						>
							<input
								type="range"
								min="0"
								max="100"
								value={volume}
								class="range range-xs w-20"
								oninput={handleVolumeChange}
							/>
						</div>
					{/if}
				</div>

				<!-- Contrôles de navigation -->
				<div class="flex items-center gap-2">
					<!-- Bouton précédent -->
					<button
						onclick={handlePrevious}
						class="btn btn-ghost btn-square btn-sm"
						disabled={!$hasPrevious}
						class:opacity-50={!$hasPrevious}
					>
						<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
							<path d="M6 6h2v12H6zm3.5 6l8.5 6V6l-8.5 6z" />
						</svg>
					</button>

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

					<!-- Bouton suivant -->
					<button
						onclick={handleNext}
						class="btn btn-ghost btn-square btn-sm"
						disabled={!$hasNext}
						class:opacity-50={!$hasNext}
					>
						<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
							<path d="M6 18l8.5-6L6 6v12zM16 6v12h2V6h-2z" />
						</svg>
					</button>
				</div>

				<!-- Bouton Queue -->
				<button
					onclick={toggleQueue}
					class="btn btn-ghost btn-square btn-sm"
					class:text-primary={showQueue}
				>
					<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
						<path
							d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"
						/>
					</svg>
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
