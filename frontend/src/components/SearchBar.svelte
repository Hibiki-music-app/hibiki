<script lang="ts">
	import { onMount } from 'svelte';

	interface SearchBarProps {
		onSearch?: (query: string) => void;
		placeholder?: string;
	}

	let { onSearch, placeholder = 'Rechercher de la musique...' }: SearchBarProps = $props();

	let searchInput: HTMLInputElement | null = $state(null);
	let searchQuery = $state('');

	onMount(() => {
		const handleKeydown = (event: KeyboardEvent) => {
			if (event.ctrlKey && event.key === 'k') {
				event.preventDefault();
				searchInput?.focus();
			}
		};

		document.addEventListener('keydown', handleKeydown);
		return () => document.removeEventListener('keydown', handleKeydown);
	});

	function handleSearch(event: Event) {
		event.preventDefault();
		if (searchQuery.trim()) {
			onSearch?.(searchQuery.trim());
		}
	}
</script>

<div class="search-bar">
	<form onsubmit={handleSearch}>
		<label class="input rounded-full w-full">
			<svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
				<g
					stroke-linejoin="round"
					stroke-linecap="round"
					stroke-width="2.5"
					fill="none"
					stroke="currentColor"
				>
					<circle cx="11" cy="11" r="8"></circle>
					<path d="m21 21-4.3-4.3"></path>
				</g>
			</svg>
			<input
				bind:this={searchInput}
				bind:value={searchQuery}
				type="search"
				class="grow"
				{placeholder}
			/>
			<kbd class="kbd kbd-sm hidden sm:inline-flex">Ctrl</kbd>
			<kbd class="kbd kbd-sm hidden sm:inline-flex">K</kbd>
		</label>
	</form>
</div>

<style>
	.search-bar {
		max-width: 32rem;
		margin: 0 auto;
		padding: 1rem;
		position: sticky;
		top: 0;
		z-index: 10;
	}
</style>
