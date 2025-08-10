<script lang="ts">
	import * as api from '$lib/api';
	import type { Track } from '$lib/api';

	// recherche
	let searchQuery = $state('');
	let isLoading = $state(false);
	let searchResults: Track[] = $state([]);

	// player
	let currentTrack: Track | null = $state(null);
	let audioElement: HTMLAudioElement | null = $state(null);
	let isPlaying = $state(false);

    function searchTracks() {
        isLoading = true;
        api.searchTracks(searchQuery).then((tracks) => {
            searchResults = (tracks);
            isLoading = false;
        });
    }

	async function playTrack(track: Track) {
        if (currentTrack?.id === track.id && isPlaying) {
            togglePlayPause();
            return;
        }
        currentTrack = track;
        if (audioElement) {
            isPlaying = true;
            audioElement.src = await api.getTrackUrl(track.id);
            audioElement.play();
        }
    }


	function togglePlayPause() {
		if (!audioElement) return;
		
		if (isPlaying) {
			audioElement.pause();
			isPlaying = false;
		} else {
			audioElement.play();
			isPlaying = true;
		}
	}
</script>

<div class="music-player">
	<h1>🎵 Hibiki Music Player</h1>

	<!-- Barre de recherche -->
	<div class="search-section">
		<input
			type="text"
			bind:value={searchQuery}
			placeholder="Rechercher une chanson, artiste..."
			onkeydown={(e) => e.key === 'Enter' && searchTracks()}
		/>
		<button onclick={searchTracks} disabled={isLoading}>
			{isLoading ? 'Recherche...' : 'Rechercher'}
		</button>
	</div>

	<!-- Résultats de recherche -->
	{#if searchResults.length > 0}
		<div class="results">
			<h3>Résultats ({searchResults.length})</h3>
			{#each searchResults as track}
				<div class="track-item">
					<div class="track-info">
						{#if track.albumCover}
							<img src={track.albumCover} alt={track.albumTitle} class="album-cover" />
						{/if}
						<div>
							<div class="track-title">{track.title}</div>
							<div class="track-artist">{track.artist}</div>
							<div class="track-album">{track.albumTitle}</div>
						</div>
					</div>
					<button
						onclick={() => playTrack(track)}
						class="play-button"
						disabled={currentTrack?.id === track.id && isPlaying}
					>
						{currentTrack?.id === track.id && isPlaying ? '▶️ En cours' : '▶️ Jouer'}
					</button>
				</div>
			{/each}
		</div>
	{/if}

	<!-- Lecteur audio -->
	{#if currentTrack}
		<div class="player-controls">
			<div class="now-playing">
				<h3>🎵 En cours de lecture</h3>
				<div class="current-track">
					{#if currentTrack.albumCover}
						<img
							src={currentTrack.albumCover}
							alt={currentTrack.albumTitle}
							class="current-cover"
						/>
					{/if}
					<div>
						<div class="current-title">{currentTrack.title}</div>
						<div class="current-artist">{currentTrack.artist}</div>
					</div>
				</div>
			</div>

			<div class="controls">
				<button onclick={togglePlayPause} class="main-control">
					{isPlaying ? '⏸️ Pause' : '▶️ Play'}
				</button>
			</div>
		</div>
	{/if}

	<!-- Élément audio caché -->
	<audio
		bind:this={audioElement}
        onplay={() => isPlaying = true}
        onpause={() => isPlaying = false}
        onended={() => isPlaying = false}
		controls
		style="width: 100%; margin-top: 10px;"
	></audio>
</div>

<style>
	.music-player {
		max-width: 800px;
		margin: 0 auto;
		padding: 20px;
		font-family: Arial, sans-serif;
	}

	.search-section {
		display: flex;
		gap: 10px;
		margin-bottom: 20px;
	}

	.search-section input {
		flex: 1;
		padding: 12px;
		border: 2px solid #ddd;
		border-radius: 8px;
		font-size: 16px;
	}

	.search-section button {
		padding: 12px 24px;
		background: #007bff;
		color: white;
		border: none;
		border-radius: 8px;
		cursor: pointer;
		font-size: 16px;
	}

	.search-section button:disabled {
		background: #ccc;
		cursor: not-allowed;
	}

	.results {
		margin-bottom: 20px;
	}

	.track-item {
		display: flex;
		justify-content: space-between;
		align-items: center;
		padding: 12px;
		border: 1px solid #eee;
		margin-bottom: 8px;
		border-radius: 8px;
		background: #f9f9f9;
	}

	.track-info {
		display: flex;
		align-items: center;
		gap: 12px;
	}

	.album-cover {
		width: 50px;
		height: 50px;
		object-fit: cover;
		border-radius: 4px;
	}

	.track-title {
		font-weight: bold;
		margin-bottom: 4px;
	}

	.track-artist {
		color: #666;
		font-size: 14px;
	}

	.track-album {
		color: #999;
		font-size: 12px;
	}

	.play-button {
		background: #28a745;
		color: white;
		border: none;
		padding: 8px 16px;
		border-radius: 4px;
		cursor: pointer;
	}

	.play-button:disabled {
		background: #ccc;
	}

	.player-controls {
		background: #fff;
		padding: 20px;
		border: 2px solid #ddd;
		border-radius: 12px;
		margin: 20px 0;
	}

	.now-playing h3 {
		margin: 0 0 10px 0;
		color: #007bff;
	}

	.current-track {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 15px;
	}

	.current-cover {
		width: 60px;
		height: 60px;
		object-fit: cover;
		border-radius: 6px;
	}

	.current-title {
		font-weight: bold;
		font-size: 18px;
	}

	.current-artist {
		color: #666;
		margin-top: 4px;
	}

	.controls {
		text-align: center;
	}

	.main-control {
		background: #007bff;
		color: white;
		border: none;
		padding: 12px 24px;
		border-radius: 8px;
		font-size: 16px;
		cursor: pointer;
	}

	.main-control:hover {
		background: #0056b3;
	}
</style>

