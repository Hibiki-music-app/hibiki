<script lang="ts">
	import * as api from '$lib/api';
	import type { Track } from '$lib/api';
	import SearchBar from '../components/SearchBar.svelte';
	import Player from '../components/Player.svelte';

	// recherche
	let searchQuery = $state('');
	let isLoading = $state(false);
	let searchResults: Track[] = $state([]);
	let hasSearched = $state(false);

	// player
	let currentTrack: Track | null = $state(null);
	let audioElement: HTMLAudioElement | null = $state(null);
	let isPlaying = $state(false);
	let player: Player;

	function searchTracks() {
		isLoading = true;
		hasSearched = true;
		api
			.searchTracks(searchQuery)
			.then((tracks) => {
				searchResults = tracks;
				isLoading = false;
			})
			.catch(() => {
				isLoading = false;
			});
	}

	function handleSearchFromBar(query: string) {
		searchQuery = query;
		searchTracks();
	}

	async function playTrack(track: Track) {
		if (player) {
			await player.playTrack(track);
		}
	}
</script>

<div class="min-h-screen bg-base-200">
	<!-- Hero Section -->
	<div class="hero min-h-[25vh] flex flex-col justify-end pb-8">
		<h1 class="text-5xl font-bold text-primary-content">On écoute quoi ?</h1>
	</div>

	<!-- Barre de recherche -->
	<SearchBar onSearch={handleSearchFromBar} />

	{#if !hasSearched && !isLoading}
		<div class="text-center">
			<div class="flex flex-wrap gap-4 justify-center">
				<button
					class="badge badge-primary badge-lg cursor-pointer hover:badge-primary-focus"
					onclick={() => handleSearchFromBar('pop music')}
				>
					🎤 Musique Pop
				</button>
				<button
					class="badge badge-secondary badge-lg cursor-pointer hover:badge-secondary-focus"
					onclick={() => handleSearchFromBar('rock classics')}
				>
					🎸 Rock Classique
				</button>
				<button
					class="badge badge-accent badge-lg cursor-pointer hover:badge-accent-focus"
					onclick={() => handleSearchFromBar('jazz')}
				>
					🎷 Jazz
				</button>
			</div>
		</div>
	{/if}

	<!-- Contenu principal -->
	<div class="container mx-auto px-4 py-8 min-h-[50vh] pb-24">
		<!-- État de chargement -->
		{#if isLoading}
			<div class="flex flex-col items-center justify-center py-16 text-center">
				<span class="loading loading-spinner loading-lg text-primary mb-4"></span>
				<p class="text-lg text-base-content/70">Recherche en cours...</p>
			</div>
		{/if}

		<!-- Message d'accueil si aucune recherche -->

		<!-- Résultats de recherche -->
		{#if searchResults.length > 0}
			<div class="mb-8">
				<div class="text-center mb-8">
					<h3 class="text-2xl font-bold text-base-content">
						Résultats pour <span class="text-primary">"{searchQuery}"</span>
					</h3>
					<div class="badge badge-primary badge-lg mt-2">
						{searchResults.length} résultat{searchResults.length > 1 ? 's' : ''}
					</div>
				</div>
				<ul class="list bg-base-100 rounded-box shadow-md">
					<li class="p-4 pb-2 text-xs opacity-60 tracking-wide">Résultats de recherche</li>
					{#each searchResults as track}
						<li class="list-row">
							<div>
								{#if track.albumCover}
									<img class="size-10 rounded-box" src={track.albumCover} alt={track.albumTitle} />
								{:else}
									<div
										class="size-10 rounded-box bg-neutral text-neutral-content flex items-center justify-center"
									>
										<span class="text-xs">🎵</span>
									</div>
								{/if}
							</div>
							<div>
								<div>{track.title}</div>
								<div class="text-xs uppercase font-semibold opacity-60">{track.artist}</div>
							</div>
							<button
								onclick={() => playTrack(track)}
								class="btn btn-square btn-ghost"
								class:text-success={currentTrack?.id === track.id && isPlaying}
								disabled={isLoading}
							>
								{#if currentTrack?.id === track.id && isPlaying}
									<svg
										class="size-[1.2em]"
										xmlns="http://www.w3.org/2000/svg"
										viewBox="0 0 24 24"
										fill="currentColor"
									>
										<rect x="6" y="4" width="4" height="16" />
										<rect x="14" y="4" width="4" height="16" />
									</svg>
								{:else}
									<svg class="size-[1.2em]" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24"
										><g
											stroke-linejoin="round"
											stroke-linecap="round"
											stroke-width="2"
											fill="none"
											stroke="currentColor"><path d="M6 3L20 12 6 21 6 3z"></path></g
										></svg
									>
								{/if}
							</button>
						</li>
					{/each}
				</ul>
			</div>
		{:else if hasSearched && !isLoading}
			<div class="text-center py-16">
				<div class="text-6xl mb-6 opacity-50">🔍</div>
				<h3 class="text-2xl font-bold text-base-content mb-4">Aucun résultat trouvé</h3>
				<p class="text-base-content/70">
					Essayez avec d'autres mots-clés ou vérifiez l'orthographe.
				</p>
				<div class="mt-8">
					<button class="btn btn-outline btn-primary" onclick={() => handleSearchFromBar('')}>
						Nouvelle recherche
					</button>
				</div>
			</div>
		{/if}
	</div>

	<!-- Player component -->
	<Player bind:this={player} bind:currentTrack bind:isPlaying bind:audioElement />
</div>

<style>
	/* Ajustements minimaux pour l'animation de rebond */
	@keyframes bounce {
		0%,
		20%,
		50%,
		80%,
		100% {
			transform: translateY(0);
		}
		40% {
			transform: translateY(-10px);
		}
		60% {
			transform: translateY(-5px);
		}
	}
</style>
