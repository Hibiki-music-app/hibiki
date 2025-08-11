<script lang="ts">
    import * as api from '$lib/components/trackApi';
    import type { Track } from '$lib/models/track';

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

<div class="max-w-4xl mx-auto p-5">
    <h1 class="text-3xl font-bold mb-6">🎵 Hibiki Music Player</h1>

    <!-- Search Section -->
    <div class="join w-full mb-8">
        <input
                type="text"
                bind:value={searchQuery}
                placeholder="Rechercher une chanson, artiste..."
                class="input input-bordered join-item w-full"
                onkeydown={(e) => e.key === 'Enter' && searchTracks()}
        />
        <button
                onclick={searchTracks}
                disabled={isLoading}
                class="btn btn-primary join-item"
        >
            {#if isLoading}
                <span class="loading loading-spinner"></span>
                Recherche...
            {:else}
                Rechercher
            {/if}
        </button>
    </div>

    <!-- Search Results -->
    {#if searchResults.length > 0}
        <div class="mb-8">
            <h3 class="text-xl font-semibold mb-4">Résultats ({searchResults.length})</h3>
            <div class="space-y-2">
                {#each searchResults as track}
                    <div class="card card-side bg-base-100 shadow-sm">
                        <figure class="w-16">
                            {#if track.albumCover}
                                <img src={track.albumCover} alt={track.albumTitle} class="w-full h-full object-cover" />
                            {/if}
                        </figure>
                        <div class="card-body p-4">
                            <h4 class="card-title text-base">{track.title}</h4>
                            <p class="text-sm text-base-content/70">{track.artist}</p>
                            <p class="text-xs text-base-content/50">{track.albumTitle}</p>
                        </div>
                        <div class="card-actions p-4">
                            <button
                                    onclick={() => playTrack(track)}
                                    class="btn btn-sm btn-success"
                                    disabled={currentTrack?.id === track.id && isPlaying}
                            >
                                {currentTrack?.id === track.id && isPlaying ? '▶️ En cours' : '▶️ Jouer'}
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
    {/if}

    <!-- Player Controls -->
    {#if currentTrack}
        <div class="card bg-base-100 shadow-md mb-8">
            <div class="card-body">
                <h3 class="card-title text-primary">🎵 En cours de lecture</h3>
                <div class="flex items-center gap-4 mt-4">
                    {#if currentTrack.albumCover}
                        <figure class="w-16 h-16">
                            <img
                                    src={currentTrack.albumCover}
                                    alt={currentTrack.albumTitle}
                                    class="w-full h-full object-cover rounded-lg"
                            />
                        </figure>
                    {/if}
                    <div>
                        <h4 class="font-bold text-lg">{currentTrack.title}</h4>
                        <p class="text-base-content/70">{currentTrack.artist}</p>
                    </div>
                </div>
                <div class="card-actions justify-center mt-4">
                    <button onclick={togglePlayPause} class="btn btn-primary">
                        {isPlaying ? '⏸️ Pause' : '▶️ Play'}
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <!-- Hidden Audio Element -->
    <audio
            bind:this={audioElement}
            onplay={() => isPlaying = true}
            onpause={() => isPlaying = false}
            onended={() => isPlaying = false}
            controls
            class="w-full mt-4"
    ></audio>
</div>