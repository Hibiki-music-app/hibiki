<script lang="ts">
	import { queueStore, currentTrack } from '$lib/stores/playerStore';
	import type { QueueItem } from '$lib/stores/playerStore';

	let { 
		show = $bindable(),
		playTrack
	}: {
		show: boolean;
		playTrack?: (track: any) => Promise<void>;
	} = $props();

	// Drag and drop
	let draggedIndex: number | null = $state(null);
	let dragOverIndex: number | null = $state(null);

	function formatDuration(seconds: number): string {
		if (!seconds || !isFinite(seconds)) return '0:00';
		const mins = Math.floor(seconds / 60);
		const secs = Math.floor(seconds % 60);
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}

	function handleDragStart(event: DragEvent, index: number) {
		draggedIndex = index;
		if (event.dataTransfer) {
			event.dataTransfer.effectAllowed = 'move';
		}
	}

	function handleDragOver(event: DragEvent, index: number) {
		event.preventDefault();
		dragOverIndex = index;
		if (event.dataTransfer) {
			event.dataTransfer.dropEffect = 'move';
		}
	}

	function handleDragLeave() {
		dragOverIndex = null;
	}

	function handleDrop(event: DragEvent, index: number) {
		event.preventDefault();
		if (draggedIndex !== null && draggedIndex !== index) {
			queueStore.reorderQueue(draggedIndex, index);
		}
		draggedIndex = null;
		dragOverIndex = null;
	}

	function handleDragEnd() {
		draggedIndex = null;
		dragOverIndex = null;
	}

	function playFromQueue(item: QueueItem, index: number) {
		queueStore.setCurrentIndex(index);
		if (playTrack) {
			playTrack(item.track);
		}
	}

	function removeFromQueue(itemId: string) {
		queueStore.removeTrack(itemId);
	}

	function clearQueue() {
		queueStore.clear();
	}

	function close() {
		show = false;
	}
</script>

{#if show}
	<!-- Overlay backdrop -->
	<div 
		class="fixed inset-0 bg-black bg-opacity-50 z-40" 
		onclick={close}
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Escape' && close()}
	></div>

	<!-- Queue Panel -->
	<div class="fixed bottom-0 left-0 right-0 bg-base-100 shadow-2xl z-50 max-h-[70vh] flex flex-col rounded-t-2xl">
		<!-- Header -->
		<div class="flex items-center justify-between p-4 border-b border-base-300">
			<h2 class="text-lg font-semibold flex items-center gap-2">
				<svg class="w-5 h-5" viewBox="0 0 24 24" fill="currentColor">
					<path d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"/>
				</svg>
				File d'attente ({$queueStore.items.length})
			</h2>
			
			<div class="flex items-center gap-2">
				<!-- Shuffle button -->
				<button 
					onclick={() => queueStore.toggleShuffle()}
					class="btn btn-ghost btn-square btn-sm"
					class:text-primary={$queueStore.isShuffling}
					aria-label="Lecture aléatoire"
				>
					<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
						<path d="M14.83,13.41L13.42,14.82L16.55,17.95L14.5,20H20V14.5L17.96,16.54L14.83,13.41M14.5,4L16.54,6.04L4,18.59L5.41,20L17.96,7.46L20,9.5V4M10.59,9.17L5.41,4L4,5.41L9.17,10.58L10.59,9.17Z"/>
					</svg>
				</button>

				<!-- Loop button -->
				<button 
					onclick={() => queueStore.toggleLoop()}
					class="btn btn-ghost btn-square btn-sm"
					class:text-primary={$queueStore.isLooping}
					aria-label="Répéter la file d'attente"
				>
					<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
						<path d="M17,17H7V14L3,18L7,22V19H19V13H17M7,7H17V10L21,6L17,2V5H5V11H7V7Z"/>
					</svg>
				</button>

				<!-- Clear button -->
				{#if $queueStore.items.length > 0}
					<button 
						onclick={clearQueue}
						class="btn btn-ghost btn-square btn-sm text-error"
						aria-label="Vider la file d'attente"
					>
						<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
							<path d="M19,4H15.5L14.5,3H9.5L8.5,4H5V6H19M6,19A2,2 0 0,0 8,21H16A2,2 0 0,0 18,19V7H6V19Z"/>
						</svg>
					</button>
				{/if}

				<!-- Close button -->
				<button onclick={close} class="btn btn-ghost btn-square btn-sm" aria-label="Fermer">
					<svg class="w-4 h-4" viewBox="0 0 24 24" fill="currentColor">
						<path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
					</svg>
				</button>
			</div>
		</div>

		<!-- Queue List -->
		<div class="flex-1 overflow-y-auto">
			{#if $queueStore.items.length === 0}
				<div class="flex flex-col items-center justify-center py-16 text-center">
					<svg class="w-16 h-16 text-base-content/30 mb-4" viewBox="0 0 24 24" fill="currentColor">
						<path d="M15,6H3V8H15V6M15,10H3V12H15V10M3,16H11V14H3V16M17,6V14.18C16.69,14.07 16.35,14 16,14A3,3 0 0,0 13,17A3,3 0 0,0 16,20A3,3 0 0,0 19,17V8H22V6H17Z"/>
					</svg>
					<h3 class="text-lg font-semibold text-base-content/60 mb-2">Aucune piste dans la file d'attente</h3>
					<p class="text-sm text-base-content/50">Ajoutez des pistes depuis la recherche pour commencer à écouter</p>
				</div>
			{:else}
				<div class="space-y-1 p-2">
					{#each $queueStore.items as item, index (item.id)}
						<div 
							class="flex items-center gap-3 p-3 rounded-lg hover:bg-base-200 transition-colors"
							class:bg-primary={index === $queueStore.currentIndex}
							class:bg-opacity-10={index === $queueStore.currentIndex}
							class:border-l-4={index === $queueStore.currentIndex}
							class:border-primary={index === $queueStore.currentIndex}
							class:opacity-60={draggedIndex === index}
							class:border-2={dragOverIndex === index}
							class:border-dashed={dragOverIndex === index}
							draggable="true"
							ondragstart={(e) => handleDragStart(e, index)}
							ondragover={(e) => handleDragOver(e, index)}
							ondragleave={handleDragLeave}
							ondrop={(e) => handleDrop(e, index)}
							ondragend={handleDragEnd}
							role="button"
							tabindex="0"
						>
							<!-- Current track indicator -->
							<div class="w-1">
								{#if index === $queueStore.currentIndex}
									<div class="w-1 h-6 bg-primary rounded-full"></div>
								{/if}
							</div>

							<!-- Album cover -->
							<div class="avatar">
								<div class="w-10 h-10 rounded">
									{#if item.track.albumCover}
										<img src={item.track.albumCover} alt="Album cover" />
									{:else}
										<div class="bg-base-300 flex items-center justify-center">🎵</div>
									{/if}
								</div>
							</div>

							<!-- Track info -->
							<button 
								class="flex-1 min-w-0 text-left" 
								onclick={() => playFromQueue(item, index)}
							>
								<div class="font-medium truncate">{item.track.title}</div>
								<div class="text-sm opacity-70 truncate">{item.track.artist}</div>
							</button>

							<!-- Duration -->
							<div class="text-xs opacity-60 mr-2">
								{formatDuration(item.track.duration)}
							</div>

							<!-- Actions -->
							<div class="flex items-center gap-1">
								<!-- Play button -->
								<button 
									onclick={() => playFromQueue(item, index)}
									class="btn btn-ghost btn-square btn-xs"
									aria-label="Jouer"
								>
									{#if index === $queueStore.currentIndex}
										<svg class="w-3 h-3" viewBox="0 0 24 24" fill="currentColor">
											<rect x="6" y="4" width="4" height="16" />
											<rect x="14" y="4" width="4" height="16" />
										</svg>
									{:else}
										<svg class="w-3 h-3" viewBox="0 0 24 24" fill="currentColor">
											<polygon points="5,3 19,12 5,21" />
										</svg>
									{/if}
								</button>

								<!-- Remove button -->
								<button 
									onclick={() => removeFromQueue(item.id)}
									class="btn btn-ghost btn-square btn-xs text-error opacity-60 hover:opacity-100"
									aria-label="Supprimer de la file d'attente"
								>
									<svg class="w-3 h-3" viewBox="0 0 24 24" fill="currentColor">
										<path d="M19,6.41L17.59,5L12,10.59L6.41,5L5,6.41L10.59,12L5,17.59L6.41,19L12,13.41L17.59,19L19,17.59L13.41,12L19,6.41Z"/>
									</svg>
								</button>

								<!-- Drag handle -->
								<div class="cursor-move opacity-30 hover:opacity-60">
									<svg class="w-3 h-3" viewBox="0 0 24 24" fill="currentColor">
										<path d="M3,15H21V13H3V15M3,19H21V17H3V19M3,11H21V9H3V11M3,5V7H21V5H3Z"/>
									</svg>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	[draggable="true"] {
		user-select: none;
	}
</style>
