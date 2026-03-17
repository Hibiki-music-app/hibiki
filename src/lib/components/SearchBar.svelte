<script lang="ts">
	import { onMount } from 'svelte';
	import { Search } from 'lucide-svelte';

	interface SearchBarProps {
		onSearch?: (query: string) => void;
		placeholder?: string;
	}

	let { onSearch, placeholder = 'Rechercher de la musique…' }: SearchBarProps = $props();

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

<form onsubmit={handleSearch} class="w-full max-w-2xl mx-auto">
	<div
		class="glass-medium rounded-full flex items-center px-4 py-3 gap-3 transition-colors duration-200"
		style="focus-within:border-[rgba(59,130,246,0.5)];"
	>
		<Search size={16} class="text-[#94a3b8] shrink-0" />
		<input
			bind:this={searchInput}
			bind:value={searchQuery}
			type="search"
			class="bg-transparent border-none outline-none flex-1 text-[#f8fafc] placeholder:text-[#64748b] text-sm min-w-0"
			{placeholder}
		/>
		<div class="hidden sm:flex items-center gap-1 text-xs text-[#64748b] shrink-0">
			<kbd class="px-1.5 py-0.5 rounded text-[10px] border border-[rgba(148,163,184,0.2)] bg-[rgba(15,23,42,0.4)]">Ctrl</kbd>
			<kbd class="px-1.5 py-0.5 rounded text-[10px] border border-[rgba(148,163,184,0.2)] bg-[rgba(15,23,42,0.4)]">K</kbd>
		</div>
	</div>
</form>
