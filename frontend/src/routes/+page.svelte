<script lang="ts">
	import * as api from '$lib/components/SearchTracks';
	import type { Track } from '$lib/components/SearchTracks';
	import { queueStore } from '$lib/stores/playerStore';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import { getContext } from 'svelte';

	// Récupérer le contexte du player depuis le layout
	const playerContext = getContext('player') as {
		currentTrack: Track | null;
		isPlaying: boolean;
		playTrack?: (track: Track) => Promise<void>;
		togglePlayPause?: () => void;
	};

	// recherche
	let searchQuery = $state('');
	let isLoading = $state(false);
	let searchResults: Track[] = $state([]);
	let hasSearched = $state(false);

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
		if (playerContext.playTrack) {
			await playerContext.playTrack(track);
		}
	}

	// Nouvelles fonctions pour la queue
	function addToQueue(track: Track) {
		queueStore.addTrack(track);
		// Afficher une notification ou toast ici si souhaité
	}

	function playNext(track: Track) {
		queueStore.playNext(track);
		// Afficher une notification ou toast ici si souhaité
	}

	// État pour gérer les menus contextuels
	let activeMenuTrackId: number | null = $state(null);

	function toggleMenu(trackId: number) {
		activeMenuTrackId = activeMenuTrackId === trackId ? null : trackId;
	}

	function closeMenu() {
		activeMenuTrackId = null;
	}
</script>

<div class="min-h-screen bg-base-200" onclick={closeMenu} onkeydown={(e) => e.key === 'Escape' && closeMenu()} role="button" tabindex="0">
	<!-- Hero Section -->
	<div class="hero min-h-[25vh] flex flex-col justify-end pb-8">
		<h1 class="text-4xl font-bold">On écoute quoi ?</h1>
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
				<ul class="list bg-base-100 rounded-box shadow-md">
					<li class="p-4 pb-2 text-xs opacity-60 tracking-wide">
						{searchResults.length} résultat{searchResults.length > 1 ? 's' : ''} pour "{searchQuery}"
					</li>
					{#each searchResults as track}
						<li class="list-row relative">
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
							<div class="flex-1">
								<div>{track.title}</div>
								<div class="text-xs uppercase font-semibold opacity-60">{track.artist}</div>
							</div>

							<!-- Actions -->
							<div class="flex items-center gap-1">
								<!-- Bouton play principal -->
								<button
									onclick={() => playTrack(track)}
									class="btn btn-square btn-ghost"
									class:text-success={playerContext.currentTrack?.id === track.id &&
										playerContext.isPlaying}
									disabled={isLoading}
									aria-label="Jouer"
								>
									{#if playerContext.currentTrack?.id === track.id && playerContext.isPlaying}
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

								<!-- Bouton menu options -->
								<div class="dropdown dropdown-end">
									<button
										onclick={() => toggleMenu(track.id)}
										class="btn btn-square btn-ghost btn-sm"
										aria-label="Plus d'options"
									>
										<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
											<path d="M12,16A2,2 0 0,1 14,18A2,2 0 0,1 12,20A2,2 0 0,1 10,18A2,2 0 0,1 12,16M12,10A2,2 0 0,1 14,12A2,2 0 0,1 12,14A2,2 0 0,1 10,12A2,2 0 0,1 12,10M12,4A2,2 0 0,1 14,6A2,2 0 0,1 12,8A2,2 0 0,1 10,6A2,2 0 0,1 12,4Z"/>
										</svg>
									</button>

									{#if activeMenuTrackId === track.id}
										<ul class="dropdown-content menu bg-base-100 rounded-box z-10 w-52 p-2 shadow-lg">
											<li>
												<button onclick={() => { playTrack(track); closeMenu(); }} class="flex items-center gap-2">
													<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
														<polygon points="5,3 19,12 5,21" />
													</svg>
													Jouer maintenant
												</button>
											</li>
											<li>
												<button onclick={() => { playNext(track); closeMenu(); }} class="flex items-center gap-2">
													<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
														<path d="M8,5V19L19,12M6,5V19L17,12"/>
													</svg>
													Jouer ensuite
												</button>
											</li>
											<li>
												<button onclick={() => { addToQueue(track); closeMenu(); }} class="flex items-center gap-2">
													<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
														<path d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"/>
													</svg>
													Ajouter à la file d'attente
												</button>
											</li>
										</ul>
									{/if}
								</div>
							</div>
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
