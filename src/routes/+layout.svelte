<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	import Player from '$lib/components/Player.svelte';
	import Queue from '$lib/components/Queue.svelte';
	import type { Track } from '$lib/models/Track';
	import { setContext } from 'svelte';
	import { goto } from '$app/navigation';

	let { children } = $props();

	function handleHomeClick() {
		goto('/');
	}

	// État global du player
	let currentTrack: Track | null = $state(null);
	let audioElement: HTMLAudioElement | null = $state(null);
	let isPlaying = $state(false);
	let showQueue = $state(false);
	let playerPlayTrack: ((track: Track) => Promise<void>) | undefined = $state();
	let playerTogglePlayPause: (() => void) | undefined = $state();

	// Contexte réactif pour les pages enfants
	const playerContext = {
		get currentTrack() {
			return currentTrack;
		},
		get isPlaying() {
			return isPlaying;
		},
		get playTrack() {
			return playerPlayTrack;
		},
		get togglePlayPause() {
			return playerTogglePlayPause;
		}
	};

	setContext('player', playerContext);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- DaisyUI Layout with Header -->
<div class="min-h-screen bg-base-200">
	<Header onHomeClick={handleHomeClick} />

	<main>
		{@render children?.()}
	</main>

	<Player
		bind:currentTrack
		bind:isPlaying
		bind:audioElement
		bind:playTrack={playerPlayTrack}
		bind:togglePlayPause={playerTogglePlayPause}
		bind:showQueue
	/>

	<Queue
		bind:show={showQueue}
		playTrack={playerPlayTrack}
	/>
</div>

<style>
	/* DaisyUI global overrides - minimisés */
	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>
