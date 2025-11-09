<script lang="ts">
	import type { RecordingResponse, Recording, Release } from '$lib/models/Musicbrainz';

	let q = '';
	let results: Recording[] = [];
	let loading = false;
	let error: string | null = null;

	// Fonction utilitaire pour formater la durée a réfactor
	function formatDuration(ms: number | null): string {
		if (!ms) return 'Inconnue';
		const seconds = Math.floor(ms / 1000);
		const minutes = Math.floor(seconds / 60);
		const remainingSeconds = seconds % 60;
		return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
	}

	// Fonction pour obtenir le type de release
	function getReleaseType(release: Release): string {
		return release['release-group']?.['primary-type'] || 'Album';
	}

	async function search() {
		error = null;
		if (!q.trim()) return;
		loading = true;
		results = [];

		try {
			const res = await fetch('/api?q=' + encodeURIComponent(q));
			if (!res.ok) {
				const err = await res.json().catch(() => ({ message: 'Erreur inconnue' }));
				error = `Erreur ${res.status}: ${err.message || 'Erreur de recherche'}`;
			} else {
				const data: RecordingResponse = await res.json();
				results = data.recordings ?? [];
			}
		} catch (e) {
			error = 'Erreur de connexion';
			console.error(e);
		} finally {
			loading = false;
		}
	}

	$: hasResults = results.length > 0;
</script>

<svelte:head>
	<title>Recherche Musicale</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 py-8 px-4 sm:px-6 lg:px-8">
	<div class="max-w-7xl mx-auto">
		<!-- Header -->
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold text-gray-900 mb-3">Recherche</h1>
			<p class="text-xl text-gray-600">Trouvez des enregistrements, artistes et albums</p>
		</div>

		<!-- Search Box -->
		<div class="max-w-2xl mx-auto mb-8">
			<div class="flex flex-col sm:flex-row gap-3">
				<input
					bind:value={q}
					placeholder="Entrez un titre, un artiste..."
					on:keydown={(e) => e.key === 'Enter' && search()}
					class="flex-1 px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all disabled:bg-gray-100 disabled:cursor-not-allowed"
					disabled={loading}
				/>
				<button
					on:click={search}
					disabled={loading || !q.trim()}
					class="px-6 py-3 bg-blue-600 text-white rounded-lg font-medium hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors flex items-center justify-center gap-2 min-w-[140px]"
				>
					{#if loading}
						<span class="loading loading-dots loading-lg"></span>
					{/if}
					Rechercher
				</button>
			</div>
		</div>

		<!-- Error Message -->
		{#if error}
			<div class="max-w-2xl mx-auto mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
				<div class="flex items-center">
					<div class="flex-shrink-0">
						<svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
							<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.28 7.22a.75.75 0 00-1.06 1.06L8.94 10l-1.72 1.72a.75.75 0 101.06 1.06L10 11.06l1.72 1.72a.75.75 0 101.06-1.06L11.06 10l1.72-1.72a.75.75 0 00-1.06-1.06L10 8.94 8.28 7.22z" clip-rule="evenodd" />
						</svg>
					</div>
					<div class="ml-3">
						<h3 class="text-sm font-medium text-red-800">Erreur</h3>
						<p class="text-sm text-red-700 mt-1">{error}</p>
					</div>
				</div>
			</div>
		{/if}

		<!-- Loading State -->
		{#if loading}
			<div class="text-center py-12">
				<div class="w-12 h-12 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mx-auto mb-4"></div>
				<p class="text-gray-600 text-lg">Recherche en cours...</p>
			</div>
		{/if}

		<!-- Results -->
		{#if hasResults}
			<div class="mt-8">
				<div class="mb-6">
					<h2 class="text-2xl font-bold text-gray-900">{results.length} résultat(s) trouvé(s)</h2>
				</div>

				<div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
					{#each results as recording (recording.id)}
						<article class="bg-white rounded-xl shadow-sm border border-gray-200 hover:shadow-md transition-shadow duration-300 p-6">
							<!-- Recording Header -->
							<div class="mb-4">
								<h3 class="text-xl font-semibold text-gray-900 mb-2 line-clamp-2">{recording.title}</h3>
							</div>

							<!-- Artists -->
							{#if recording['artist-credit']?.length}
								<div class="mb-4">
									<strong class="text-sm font-medium text-gray-700 block mb-2">Artiste(s):</strong>
									<div class="flex flex-wrap gap-2">
										{#each recording['artist-credit'] as credit (credit.artist.id)}
											<span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
												{credit.name}
											</span>
										{/each}
									</div>
								</div>
							{/if}

							<!-- Recording Details -->
							<div class="flex flex-wrap gap-4 mb-4">
								{#if recording.length}
									<div class="flex items-center gap-2">
										<span class="text-xs font-medium text-gray-500">Durée:</span>
										<span class="text-sm text-gray-900 font-medium">{formatDuration(recording.length)}</span>
									</div>
								{/if}
								{#if recording.video}
									<div class="flex items-center gap-2">
										<span class="text-xs font-medium text-gray-500">Vidéo:</span>
										<span class="text-sm text-red-600 font-medium">Disponible</span>
									</div>
								{/if}
							</div>

							<!-- Releases -->
							{#if recording.releases?.length}
								<div class="border-t border-gray-200 pt-4">
									<strong class="text-sm font-medium text-gray-700 block mb-3">
										Sorties ({recording.releases.length}):
									</strong>
									<div class="space-y-3">
										{#each recording.releases as release (release.id)}
											<div class="bg-gray-50 rounded-lg p-4 border border-gray-200">
												<!-- Release Header -->
												<div class="flex items-start justify-between mb-2 gap-2">
													<span class="text-sm font-semibold text-gray-900 flex-1 line-clamp-2">
														{release.title}
													</span>
													<span class="inline-flex items-center px-2 py-1 rounded text-xs font-bold text-white {getReleaseType(release).toLowerCase() === 'album' ? 'bg-blue-500' : getReleaseType(release).toLowerCase() === 'single' ? 'bg-red-500' : getReleaseType(release).toLowerCase() === 'ep' ? 'bg-orange-500' : 'bg-green-500'}">
														{getReleaseType(release)}
													</span>
												</div>

												<!-- Release Details -->
												<div class="flex flex-wrap gap-3 text-xs text-gray-600 mb-3">
													<span class="capitalize">{release.status}</span>
													{#if release['track-count']}
														<span>{release['track-count']} pistes</span>
													{/if}
												</div>

												<!-- Media -->
												{#if release.media?.length}
													<div class="space-y-2">
														{#each release.media as medium (medium.id)}
															<div class="bg-white rounded border border-gray-200 p-3">
																<span class="text-xs font-medium text-gray-700 block mb-2">
																	Media {medium.position} ({medium['track-count']} pistes)
																</span>
																{#if medium.track?.length}
																	<div class="space-y-1">
																		{#each medium.track as track (track.id)}
																			<div class="flex items-center gap-2 py-1 px-2 rounded {track.title === recording.title ? 'bg-blue-50 border border-blue-200' : ''}">
																				<span class="text-xs text-gray-500 w-6 flex-shrink-0">
																					{track.number}.
																				</span>
																				<span class="text-sm text-gray-800 flex-1 line-clamp-1">
																					{track.title}
																				</span>
																				{#if track.length}
																					<span class="text-xs text-gray-500 flex-shrink-0">
																						{formatDuration(track.length)}
																					</span>
																				{/if}
																			</div>
																		{/each}
																	</div>
																{/if}
															</div>
														{/each}
													</div>
												{/if}
											</div>
										{/each}
									</div>
								</div>
							{/if}
						</article>
					{/each}
				</div>
			</div>
		{:else if !loading && q && !error}
			<!-- No Results -->
			<div class="text-center py-16">
				<div class="max-w-md mx-auto">
					<svg class="mx-auto h-16 w-16 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
					</svg>
					<h3 class="text-lg font-medium text-gray-900 mb-2">
						Aucun résultat trouvé pour "<strong class="text-blue-600">"{q}"</strong>"
					</h3>
					<p class="text-gray-600">Essayez avec d'autres termes de recherche.</p>
				</div>
			</div>
		{/if}
	</div>
</div>