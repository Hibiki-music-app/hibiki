<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import '../app.css';
	import Header from '$lib/components/Header.svelte';
	import Player from '$lib/components/Player.svelte';
	import Queue from '$lib/components/Queue.svelte';
	import { setContext } from 'svelte';
	import { goto } from '$app/navigation';
	import type { Track } from '$lib/models/Track';
	import type { UserType } from '$lib/models/UserType';

	function handleHomeClick() {
		goto('/');
	}

	let {
		children,
		data,
	}: {
		children: any;
		data: { user: UserType | null };
	} = $props();

	const { user } = data;

	let currentTrack: Track | null = $state(null);
	let audioElement: HTMLAudioElement | null = $state(null);
	let isPlaying = $state(false);
	let showQueue = $state(false);
	let playerPlayTrack: ((track: Track) => Promise<void>) | undefined = $state();
	let playerTogglePlayPause: (() => void) | undefined = $state();

	const playerContext = {
		get currentTrack() { return currentTrack; },
		get isPlaying() { return isPlaying; },
		get playTrack() { return playerPlayTrack; },
		get togglePlayPause() { return playerTogglePlayPause; },
	};

	setContext('player', playerContext);
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- Ambient background blobs -->
<div class="fixed inset-0 pointer-events-none overflow-hidden" aria-hidden="true">
	<div
		class="absolute -top-40 -left-40 w-[700px] h-[700px] gpu-accelerated"
		style="background: radial-gradient(ellipse at 20% 30%, rgba(59,130,246,0.12) 0%, transparent 55%);
			   animation: ambient-float 20s ease-in-out infinite;"
	></div>
	<div
		class="absolute -bottom-40 -right-40 w-[600px] h-[600px] gpu-accelerated"
		style="background: radial-gradient(ellipse at 80% 70%, rgba(99,102,241,0.10) 0%, transparent 55%);
			   animation: ambient-float 20s ease-in-out infinite reverse;"
	></div>
</div>

<div class="min-h-screen relative flex flex-col">
	<Header onHomeClick={handleHomeClick} {user} />

	<main
		class="flex-1 relative z-[1]"
		style="max-width: min(1100px, 100%); width: 100%; margin: 0 auto; padding: clamp(1.5rem, 2.4vw, 2.75rem);"
	>
		{#if children}
			{@render children()}
		{/if}
	</main>

	<Player
		bind:currentTrack
		bind:isPlaying
		bind:audioElement
		bind:playTrack={playerPlayTrack}
		bind:togglePlayPause={playerTogglePlayPause}
		bind:showQueue
	/>

	<Queue bind:show={showQueue} playTrack={playerPlayTrack} />
</div>
