<script lang="ts">
	import * as api from '$lib/components/SearchTracks';
	import type { Track } from '$lib/models/Track';
	import { queueStore } from '$lib/stores/playerStore';
	import SearchBar from '$lib/components/SearchBar.svelte';
	import { getContext } from 'svelte';
	import { Play, Pause, MoreVertical, ListPlus, ChevronRight, Loader2 } from 'lucide-svelte';

	const playerContext = getContext('player') as {
		currentTrack: Track | null;
		isPlaying: boolean;
		playTrack?: (track: Track) => Promise<void>;
		togglePlayPause?: () => void;
	};

	let searchQuery = $state('');
	let isLoading = $state(false);
	let searchResults: Track[] = $state([]);
	let hasSearched = $state(false);
	let activeMenuTrackId: number | null = $state(null);

	function doSearch(query: string) {
		searchQuery = query;
		isLoading = true;
		hasSearched = true;
		api.searchTracks(query)
			.then((tracks) => {
				searchResults = tracks;
				isLoading = false;
			})
			.catch(() => {
				isLoading = false;
			});
	}

	async function playTrack(track: Track) {
		if (playerContext.playTrack) {
			await playerContext.playTrack(track);
		}
	}

	function toggleMenu(trackId: number, e: Event) {
		e.stopPropagation();
		activeMenuTrackId = activeMenuTrackId === trackId ? null : trackId;
	}

	function closeMenu() {
		activeMenuTrackId = null;
	}

	const quickSearches = [
		{ label: 'Pop', query: 'pop' },
		{ label: 'Rock', query: 'rock' },
		{ label: 'Jazz', query: 'jazz' },
		{ label: 'Electronic', query: 'electronic' },
		{ label: 'Classical', query: 'classical' },
	];
</script>

<div
	onclick={closeMenu}
	onkeydown={(e) => e.key === 'Escape' && closeMenu()}
	role="presentation"
>
	<!-- Hero -->
	<div class="mb-8" style="animation: fade-in 0.3s ease;">
		<h1 class="font-bold mb-2"
			style="font-size: clamp(1.8rem, 4vw, 3rem); line-height: 1.1;">
			On écoute quoi ?
		</h1>
		<p class="text-[#94a3b8] text-sm">Recherchez parmi des millions de titres HiFi</p>
	</div>

	<!-- Search bar -->
	<div class="mb-6">
		<SearchBar onSearch={doSearch} />
	</div>

	<!-- Quick searches -->
	{#if !hasSearched}
		<div class="flex flex-wrap gap-2 mb-10" style="animation: fade-in 0.4s ease 0.1s both;">
			{#each quickSearches as { label, query }}
				<button
					class="btn-glass text-xs"
					onclick={() => doSearch(query)}
				>
					{label}
				</button>
			{/each}
		</div>
	{/if}

	<!-- Loading state -->
	{#if isLoading}
		<div class="flex flex-col items-center justify-center py-20 gap-4" style="animation: fade-in 0.2s ease;">
			<Loader2 size={32} class="text-[#3b82f6] animate-spin" />
			<p class="text-[#94a3b8] text-sm">Recherche en cours…</p>
		</div>
	{/if}

	<!-- Results -->
	{#if !isLoading && searchResults.length > 0}
		<div style="animation: slide-up 0.25s ease;">
			<p class="text-xs text-[#64748b] uppercase tracking-wider font-semibold mb-4">
				{searchResults.length} résultat{searchResults.length > 1 ? 's' : ''} pour « {searchQuery} »
			</p>

			<div class="space-y-1">
				{#each searchResults as track, i}
					<div
						class="glass-action rounded-xl group relative"
						style="animation: fade-in 0.2s ease {i * 0.03}s both;"
					>
						<!-- Cover -->
						<div class="flex items-center gap-3 flex-1 min-w-0">
							<div class="w-11 h-11 rounded-lg shrink-0 overflow-hidden">
								{#if track.albumCover}
									<img src={track.albumCover} alt={track.albumTitle} class="w-full h-full object-cover" />
								{:else}
									<div class="w-full h-full glass-subtle flex items-center justify-center text-lg">🎵</div>
								{/if}
							</div>

							<!-- Info -->
							<div class="flex-1 min-w-0">
								<p class="text-sm font-medium text-[#f8fafc] truncate
									{playerContext.currentTrack?.id === track.id ? 'text-[#93c5fd]' : ''}">
									{track.title}
								</p>
								<p class="text-xs text-[#94a3b8] truncate">{track.artist}</p>
							</div>
						</div>

						<!-- Actions -->
						<div class="flex items-center gap-1 shrink-0">
							<!-- Duration -->
							{#if track.duration}
								<span class="text-xs text-[#64748b] tabular-nums hidden sm:inline mr-2">
									{Math.floor(track.duration / 60)}:{String(Math.floor(track.duration % 60)).padStart(2, '0')}
								</span>
							{/if}

							<!-- Play -->
							<button
								onclick={() => playTrack(track)}
								class="btn-glass touch-target"
								style="padding: 0.5rem; {playerContext.currentTrack?.id === track.id && playerContext.isPlaying ? 'color: #60a5fa; border-color: rgba(59,130,246,0.4);' : ''}"
								aria-label="Jouer"
							>
								{#if playerContext.currentTrack?.id === track.id && playerContext.isPlaying}
									<Pause size={16} />
								{:else}
									<Play size={16} />
								{/if}
							</button>

							<!-- Options menu -->
							<div class="relative">
								<button
									onclick={(e) => toggleMenu(track.id, e)}
									class="btn-glass touch-target"
									style="padding: 0.5rem;"
									aria-label="Plus d'options"
								>
									<MoreVertical size={16} />
								</button>

								{#if activeMenuTrackId === track.id}
									<div
										class="glass-strong absolute right-0 top-full mt-1 w-52 rounded-xl z-10 overflow-hidden"
										style="box-shadow: 0 20px 40px rgba(2,6,23,0.4); animation: slide-up 0.15s ease;"
										role="menu"
									>
										<div class="p-1">
											<button
												class="glass-action rounded-lg w-full text-sm text-[#f8fafc] gap-2 justify-start"
												onclick={() => { playTrack(track); closeMenu(); }}
												role="menuitem"
											>
												<Play size={14} class="text-[#60a5fa]" />
												Jouer maintenant
											</button>
											<button
												class="glass-action rounded-lg w-full text-sm text-[#f8fafc] gap-2 justify-start"
												onclick={() => { queueStore.playNext(track); closeMenu(); }}
												role="menuitem"
											>
												<ChevronRight size={14} class="text-[#94a3b8]" />
												Jouer ensuite
											</button>
											<button
												class="glass-action rounded-lg w-full text-sm text-[#f8fafc] gap-2 justify-start"
												onclick={() => { queueStore.addTrack(track); closeMenu(); }}
												role="menuitem"
											>
												<ListPlus size={14} class="text-[#94a3b8]" />
												Ajouter à la file
											</button>
										</div>
									</div>
								{/if}
							</div>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/if}

	<!-- No results -->
	{#if hasSearched && !isLoading && searchResults.length === 0}
		<div class="text-center py-20" style="animation: fade-in 0.25s ease;">
			<p class="text-4xl mb-4">🔍</p>
			<h3 class="text-lg font-semibold text-[#f8fafc] mb-2">Aucun résultat</h3>
			<p class="text-sm text-[#94a3b8]">Essayez d'autres mots-clés</p>
			<button
				class="btn-glass mt-6"
				onclick={() => { hasSearched = false; searchResults = []; }}
			>
				Nouvelle recherche
			</button>
		</div>
	{/if}
</div>

<!-- Bottom padding for player bar -->
<div class="h-28"></div>
