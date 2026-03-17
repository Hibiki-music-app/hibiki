<script lang="ts">
	import { queueStore, currentTrack } from '$lib/stores/playerStore';
	import type { QueueItem } from '$lib/stores/playerStore';
	import { Shuffle, Repeat, Trash2, X, Play, Pause, GripVertical, ListMusic } from 'lucide-svelte';

	let {
		show = $bindable(),
		playTrack,
	}: {
		show: boolean;
		playTrack?: (track: any) => Promise<void>;
	} = $props();

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
		if (event.dataTransfer) event.dataTransfer.effectAllowed = 'move';
	}

	function handleDragOver(event: DragEvent, index: number) {
		event.preventDefault();
		dragOverIndex = index;
		if (event.dataTransfer) event.dataTransfer.dropEffect = 'move';
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
		if (playTrack) playTrack(item.track);
	}

	function close() {
		show = false;
	}
</script>

{#if show}
	<!-- Backdrop -->
	<div
		class="fixed inset-0 z-40"
		style="background: rgba(2,6,23,0.6); backdrop-filter: blur(4px);"
		onclick={close}
		role="button"
		tabindex="0"
		onkeydown={(e) => e.key === 'Escape' && close()}
	></div>

	<!-- Panel -->
	<div
		class="glass-strong fixed bottom-0 left-0 right-0 z-50 flex flex-col"
		style="max-height: 72vh; border-radius: 1.25rem 1.25rem 0 0;
			   box-shadow: 0 -20px 60px rgba(2,6,23,0.5);
			   animation: slide-up 0.2s ease;"
	>
		<!-- Header -->
		<div class="flex items-center justify-between px-5 py-4 border-b border-[rgba(148,163,184,0.1)]">
			<h2 class="text-sm font-semibold text-[#f8fafc] flex items-center gap-2">
				<ListMusic size={16} class="text-[#60a5fa]" />
				File d'attente
				<span class="text-[#64748b] font-normal">({$queueStore.items.length})</span>
			</h2>

			<div class="flex items-center gap-1">
				<button
					onclick={() => queueStore.toggleShuffle()}
					class="btn-glass touch-target"
					style="padding: 0.4rem; {$queueStore.isShuffling ? 'color: #60a5fa; border-color: rgba(59,130,246,0.4);' : ''}"
					aria-label="Lecture aléatoire"
					title="Aléatoire"
				>
					<Shuffle size={15} />
				</button>

				<button
					onclick={() => queueStore.toggleLoop()}
					class="btn-glass touch-target"
					style="padding: 0.4rem; {$queueStore.isLooping ? 'color: #60a5fa; border-color: rgba(59,130,246,0.4);' : ''}"
					aria-label="Répéter"
					title="Répéter"
				>
					<Repeat size={15} />
				</button>

				{#if $queueStore.items.length > 0}
					<button
						onclick={() => queueStore.clear()}
						class="btn-glass touch-target"
						style="padding: 0.4rem; color: #f87171; border-color: rgba(239,68,68,0.3);"
						aria-label="Vider"
						title="Vider la file"
					>
						<Trash2 size={15} />
					</button>
				{/if}

				<button
					onclick={close}
					class="btn-glass touch-target ml-1"
					style="padding: 0.4rem;"
					aria-label="Fermer"
				>
					<X size={15} />
				</button>
			</div>
		</div>

		<!-- List -->
		<div class="flex-1 overflow-y-auto p-2">
			{#if $queueStore.items.length === 0}
				<div class="flex flex-col items-center justify-center py-16 text-center">
					<ListMusic size={40} class="text-[#334155] mb-4" />
					<p class="text-sm text-[#64748b] font-medium">File d'attente vide</p>
					<p class="text-xs text-[#475569] mt-1">Ajoutez des pistes depuis la recherche</p>
				</div>
			{:else}
				{#each $queueStore.items as item, index (item.id)}
					<div
						class="flex items-center gap-3 px-3 py-2.5 rounded-xl transition-all duration-150 group select-none"
						class:opacity-50={draggedIndex === index}
						style="{index === $queueStore.currentIndex
							? 'background: rgba(59,130,246,0.1); border: 1px solid rgba(59,130,246,0.25);'
							: 'border: 1px solid transparent;'}
							{dragOverIndex === index && draggedIndex !== index
							? 'border-color: rgba(59,130,246,0.5); background: rgba(59,130,246,0.08);'
							: ''}"
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
						{#if index === $queueStore.currentIndex}
							<div class="w-1 h-8 rounded-full shrink-0" style="background: #3b82f6;"></div>
						{:else}
							<div class="w-1 shrink-0"></div>
						{/if}

						<!-- Album cover -->
						<div class="w-9 h-9 rounded-lg shrink-0 overflow-hidden">
							{#if item.track.albumCover}
								<img src={item.track.albumCover} alt="Pochette" class="w-full h-full object-cover" />
							{:else}
								<div class="w-full h-full glass-subtle flex items-center justify-center text-sm">🎵</div>
							{/if}
						</div>

						<!-- Track info -->
						<button
							class="flex-1 min-w-0 text-left"
							onclick={() => playFromQueue(item, index)}
						>
							<p class="text-sm font-medium truncate
								{index === $queueStore.currentIndex ? 'text-[#93c5fd]' : 'text-[#f8fafc]'}">
								{item.track.title}
							</p>
							<p class="text-xs text-[#94a3b8] truncate">{item.track.artist}</p>
						</button>

						<!-- Duration -->
						<span class="text-xs text-[#64748b] shrink-0 tabular-nums">
							{formatDuration(item.track.duration)}
						</span>

						<!-- Actions (visible on hover) -->
						<div class="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
							<button
								onclick={() => playFromQueue(item, index)}
								class="btn-glass touch-target"
								style="padding: 0.3rem;"
								aria-label="Jouer"
							>
								{#if index === $queueStore.currentIndex}
									<Pause size={12} />
								{:else}
									<Play size={12} />
								{/if}
							</button>
							<button
								onclick={() => queueStore.removeTrack(item.id)}
								class="btn-glass touch-target"
								style="padding: 0.3rem; color: #f87171;"
								aria-label="Retirer"
							>
								<X size={12} />
							</button>
							<div class="cursor-move text-[#475569] hover:text-[#94a3b8] p-1">
								<GripVertical size={12} />
							</div>
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
{/if}

<style>
	[draggable="true"] { user-select: none; }
</style>
