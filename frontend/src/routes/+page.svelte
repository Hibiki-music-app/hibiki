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
        api.searchTracks(searchQuery).then((tracks) => {
            searchResults = tracks;
            isLoading = false;
        }).catch(() => {
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

<div class="page-container">
    <!-- Message de bienvenue -->
    <div class="welcome-message">
        <div class="welcome-content">
            <p>On écoute quoi ?</p>
        </div>
    </div>

    <!-- Barre de recherche sticky -->
    <SearchBar onSearch={handleSearchFromBar} />

	<!-- Contenu principal -->
	<div class="main-content">
		<!-- État de chargement -->
		{#if isLoading}
			<div class="loading-state">
				<div class="spinner"></div>
				<p>Recherche en cours...</p>
			</div>
		{/if}

		<!-- Message d'accueil si aucune recherche -->
		{#if !hasSearched && !isLoading}
			<div class="welcome-state">
				<div class="welcome-icon">🎵</div>
				<h2>Bienvenue sur Hibiki</h2>
				<p>Utilisez la barre de recherche ci-dessus pour trouver vos artistes et chansons préférés.</p>
				<div class="quick-actions">
					<button class="quick-action" onclick={() => handleSearchFromBar('pop music')}>
						🎤 Musique Pop
					</button>
					<button class="quick-action" onclick={() => handleSearchFromBar('rock classics')}>
						🎸 Rock Classique
					</button>
					<button class="quick-action" onclick={() => handleSearchFromBar('jazz')}>
						🎷 Jazz
					</button>
				</div>
			</div>
		{/if}

		<!-- Résultats de recherche -->
		{#if searchResults.length > 0}
			<div class="results-section">
				<h3 class="results-title">Résultats pour "{searchQuery}" ({searchResults.length})</h3>
				<div class="results-grid">
					{#each searchResults as track}
						<div class="track-card">
							<div class="track-info">
								{#if track.albumCover}
									<img src={track.albumCover} alt={track.albumTitle} class="album-cover" />
								{:else}
									<div class="album-cover-placeholder">🎵</div>
								{/if}
								<div class="track-details">
									<div class="track-title">{track.title}</div>
									<div class="track-artist">{track.artist}</div>
									{#if track.albumTitle}
										<div class="track-album">{track.albumTitle}</div>
									{/if}
								</div>
							</div>
							<button
								onclick={() => playTrack(track)}
								class="play-button"
								class:playing={currentTrack?.id === track.id && isPlaying}
								disabled={isLoading}
							>
								{#if currentTrack?.id === track.id && isPlaying}
									<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
										<rect x="6" y="4" width="4" height="16"/>
										<rect x="14" y="4" width="4" height="16"/>
									</svg>
									Pause
								{:else}
									<svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
										<polygon points="5,3 19,12 5,21"/>
									</svg>
									{currentTrack?.id === track.id ? 'Reprendre' : 'Jouer'}
								{/if}
							</button>
						</div>
					{/each}
				</div>
			</div>
		{:else if hasSearched && !isLoading}
			<div class="no-results">
				<div class="no-results-icon">🔍</div>
				<h3>Aucun résultat trouvé</h3>
				<p>Essayez avec d'autres mots-clés ou vérifiez l'orthographe.</p>
			</div>
		{/if}
	</div>

	<!-- Player component -->
	<Player 
		bind:this={player}
		bind:currentTrack 
		bind:isPlaying 
		bind:audioElement 
	/>
</div>

<style>
	.page-container {
		min-height: 100vh;
	}

	.welcome-message {
		display: flex;
		align-items: center;
		justify-content: center;
		min-height: 50vh;
		text-align: center;
		margin-bottom: -12px;
	}

	.welcome-content {
		width: 100%;
		max-width: 600px;
	}

	.welcome-message p {
		font-size: 1.25rem;
		color: #374151;
		margin: 0;
		font-weight: 500;
	}

	.main-content {
		max-width: 1200px;
		margin: 0 auto;
		padding: 2rem 1rem;
		min-height: 50vh;
	}

	.loading-state {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		padding: 4rem 0;
		color: #6b7280;
	}

	.spinner {
		width: 40px;
		height: 40px;
		border: 4px solid #e5e7eb;
		border-top: 4px solid #6366f1;
		border-radius: 50%;
		animation: spin 1s linear infinite;
		margin-bottom: 1rem;
	}

	@keyframes spin {
		0% { transform: rotate(0deg); }
		100% { transform: rotate(360deg); }
	}

	.welcome-state {
		text-align: center;
		padding: 4rem 0;
		color: #374151;
	}

	.welcome-icon {
		font-size: 4rem;
		margin-bottom: 1.5rem;
		animation: bounce 2s infinite;
	}

	@keyframes bounce {
		0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
		40% { transform: translateY(-10px); }
		60% { transform: translateY(-5px); }
	}

	.welcome-state h2 {
		font-size: 2rem;
		margin: 0 0 1rem 0;
		color: #111827;
	}

	.welcome-state p {
		font-size: 1.1rem;
		color: #6b7280;
		margin-bottom: 2rem;
		max-width: 500px;
		margin-left: auto;
		margin-right: auto;
	}

	.quick-actions {
		display: flex;
		gap: 1rem;
		justify-content: center;
		flex-wrap: wrap;
	}

	.quick-action {
		padding: 0.75rem 1.5rem;
		background: #f3f4f6;
		border: none;
		border-radius: 50px;
		cursor: pointer;
		font-size: 1rem;
		transition: all 0.2s ease;
		color: #374151;
		border: 2px solid transparent;
	}

	.quick-action:hover {
		background: #e5e7eb;
		transform: translateY(-2px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
		border-color: #6366f1;
	}

	.results-section {
		padding: 1rem 0;
	}

	.results-title {
		font-size: 1.5rem;
		margin: 0 0 2rem 0;
		color: #111827;
		text-align: center;
	}

	.results-grid {
		display: grid;
		grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
		gap: 1.5rem;
	}

	.track-card {
		background: white;
		border-radius: 12px;
		padding: 1.5rem;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
		transition: all 0.2s ease;
		border: 1px solid #e5e7eb;
	}

	.track-card:hover {
		transform: translateY(-4px);
		box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
	}

	.track-info {
		display: flex;
		align-items: flex-start;
		gap: 1rem;
		margin-bottom: 1rem;
	}

	.album-cover {
		width: 60px;
		height: 60px;
		object-fit: cover;
		border-radius: 8px;
		flex-shrink: 0;
	}

	.album-cover-placeholder {
		width: 60px;
		height: 60px;
		background: #f3f4f6;
		border-radius: 8px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
		color: #9ca3af;
		flex-shrink: 0;
	}

	.track-details {
		flex: 1;
		min-width: 0;
	}

	.track-title {
		font-weight: 600;
		font-size: 1.1rem;
		margin-bottom: 0.25rem;
		color: #111827;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.track-artist {
		color: #6b7280;
		font-size: 0.9rem;
		margin-bottom: 0.25rem;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.track-album {
		color: #9ca3af;
		font-size: 0.8rem;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.play-button {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		width: 100%;
		padding: 0.75rem 1rem;
		background: #6366f1;
		color: white;
		border: none;
		border-radius: 8px;
		cursor: pointer;
		font-size: 0.9rem;
		font-weight: 500;
		transition: all 0.2s ease;
		justify-content: center;
	}

	.play-button:hover:not(:disabled) {
		background: #5b5bf6;
		transform: translateY(-1px);
	}

	.play-button:disabled {
		background: #d1d5db;
		cursor: not-allowed;
		transform: none;
	}

	.play-button.playing {
		background: #10b981;
	}

	.play-button.playing:hover {
		background: #059669;
	}

	.no-results {
		text-align: center;
		padding: 4rem 0;
		color: #6b7280;
	}

	.no-results-icon {
		font-size: 4rem;
		margin-bottom: 1rem;
		opacity: 0.5;
	}

	.no-results h3 {
		font-size: 1.5rem;
		margin: 0 0 1rem 0;
		color: #374151;
	}

	.no-results p {
		font-size: 1rem;
		margin: 0;
	}

	/* Ajustement du contenu principal pour éviter que le lecteur le cache */
	.page-container {
		min-height: 100vh;
		padding-bottom: 90px; /* Espace pour le lecteur fixe */
	}

	/* Responsive */
	@media (max-width: 768px) {
		.welcome-message {
			min-height: 40vh;
			padding: 0 1rem;
			margin-bottom: -8px;
		}

		.welcome-message p {
			font-size: 1.125rem;
		}
		
		.main-content {
			padding: 1rem;
		}

		.results-grid {
			grid-template-columns: 1fr;
		}

		.quick-actions {
			flex-direction: column;
			align-items: center;
		}

		.quick-action {
			width: 200px;
		}

		.page-container {
			padding-bottom: 85px;
		}
	}
</style>

