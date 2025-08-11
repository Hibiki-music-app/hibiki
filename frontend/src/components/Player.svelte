<script lang="ts">
	import type { Track } from '$lib/api';
	import * as api from '$lib/api';

	let {
		currentTrack = $bindable(),
		isPlaying = $bindable(),
		audioElement = $bindable()
	}: {
		currentTrack: Track | null,
		isPlaying: boolean,
		audioElement: HTMLAudioElement | null
	} = $props();

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

	// Exposer les fonctions au parent
	export { playTrack, togglePlayPause };
</script>

<!-- Lecteur audio fixe en bas de l'écran -->
{#if currentTrack}
	<div class="audio-player">
		<div class="player-content">
			<div class="now-playing">
				<div class="track-info-player">
					{#if currentTrack.albumCover}
						<img
							src={currentTrack.albumCover}
							alt={currentTrack.albumTitle || 'Album cover'}
							class="current-cover"
						/>
					{:else}
						<div class="current-cover-placeholder">🎵</div>
					{/if}
					<div class="current-details">
						<div class="current-title">{currentTrack.title}</div>
						<div class="current-artist">{currentTrack.artist}</div>
					</div>
				</div>
			</div>

			<div class="player-controls">
				<button onclick={togglePlayPause} class="play-control" aria-label={isPlaying ? 'Pause' : 'Play'}>
					{#if isPlaying}
						<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
							<rect x="6" y="4" width="4" height="16"/>
							<rect x="14" y="4" width="4" height="16"/>
						</svg>
					{:else}
						<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
							<polygon points="5,3 19,12 5,21"/>
						</svg>
					{/if}
				</button>
			</div>

			<div class="volume-section">
				<button class="volume-button" aria-label="Volume">
					<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<polygon points="11,5 6,9 2,9 2,15 6,15 11,19"/>
						<path d="M15.54 8.46a5 5 0 0 1 0 7.07M18.07 6.93a10 10 0 0 1 0 10.14"/>
					</svg>
				</button>
			</div>
		</div>

		<!-- Élément audio caché -->
		<audio
			bind:this={audioElement}
			onplay={() => isPlaying = true}
			onpause={() => isPlaying = false}
			onended={() => isPlaying = false}
			style="display: none;"
		></audio>
	</div>
{/if}

<style>
	/* Lecteur audio fixe en bas */
	.audio-player {
		position: fixed;
		bottom: 0;
		left: 0;
		right: 0;
		background: rgba(255, 255, 255, 0.98);
		backdrop-filter: blur(20px);
		border-top: 1px solid rgba(229, 231, 235, 0.8);
		box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
		z-index: 1000;
		padding: 16px;
	}

	.player-content {
		max-width: 1200px;
		margin: 0 auto;
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 16px;
	}

	.track-info-player {
		display: flex;
		align-items: center;
		gap: 12px;
		flex: 1;
		min-width: 0;
	}

	.current-cover {
		width: 50px;
		height: 50px;
		object-fit: cover;
		border-radius: 8px;
		flex-shrink: 0;
	}

	.current-cover-placeholder {
		width: 50px;
		height: 50px;
		background: #f3f4f6;
		border-radius: 8px;
		display: flex;
		align-items: center;
		justify-content: center;
		font-size: 1.5rem;
		color: #9ca3af;
		flex-shrink: 0;
	}

	.current-details {
		flex: 1;
		min-width: 0;
	}

	.current-title {
		font-weight: 600;
		font-size: 1rem;
		color: #111827;
		margin-bottom: 2px;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.current-artist {
		font-size: 0.875rem;
		color: #6b7280;
		overflow: hidden;
		text-overflow: ellipsis;
		white-space: nowrap;
	}

	.player-controls {
		display: flex;
		align-items: center;
		justify-content: center;
		flex: 0 0 auto;
	}

	.play-control {
		width: 48px;
		height: 48px;
		border: none;
		background: #6366f1;
		color: white;
		border-radius: 50%;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
		box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
	}

	.play-control:hover {
		background: #5b5bf6;
		transform: scale(1.05);
		box-shadow: 0 4px 12px rgba(99, 102, 241, 0.4);
	}

	.play-control:active {
		transform: scale(0.95);
	}

	.volume-section {
		display: flex;
		align-items: center;
		flex: 0 0 auto;
	}

	.volume-button {
		width: 40px;
		height: 40px;
		border: none;
		background: transparent;
		color: #6b7280;
		border-radius: 50%;
		cursor: pointer;
		display: flex;
		align-items: center;
		justify-content: center;
		transition: all 0.2s ease;
	}

	.volume-button:hover {
		background: #f3f4f6;
		color: #374151;
	}

	/* Responsive */
	@media (max-width: 768px) {
		/* Lecteur audio mobile */
		.player-content {
			gap: 12px;
		}

		.current-details {
			max-width: 150px;
		}

		.current-title {
			font-size: 0.9rem;
		}

		.current-artist {
			font-size: 0.8rem;
		}

		.play-control {
			width: 44px;
			height: 44px;
		}

		.volume-section {
			display: none; /* Cache le bouton volume sur mobile */
		}
	}
</style>