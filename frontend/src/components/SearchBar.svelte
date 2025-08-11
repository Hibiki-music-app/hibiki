<script lang="ts">
	import { onMount } from 'svelte';

	// Props
	interface SearchBarProps {
		onSearch?: (query: string) => void;
		placeholder?: string;
		showShortcutHint?: boolean;
	}

	let {
		onSearch,
		placeholder = 'Rechercher de la musique...',
		showShortcutHint = true
	}: SearchBarProps = $props();

	let searchInput: HTMLInputElement | null = $state(null);
	let searchQuery = $state('');
	let isSticky = $state(false);

	// Gestion du raccourci Ctrl+K
	onMount(() => {
		const handleKeydown = (event: KeyboardEvent) => {
			if (event.ctrlKey && event.key === 'k') {
				event.preventDefault();
				searchInput?.focus();
			}
		};

		// Observer pour détecter quand la barre devient sticky
		const observer = new IntersectionObserver(
			([entry]) => {
				isSticky = !entry.isIntersecting;
			},
			{ threshold: 0, rootMargin: '-1px 0px 0px 0px' }
		);

		const sentinel = document.querySelector('.sticky-sentinel');
		if (sentinel) {
			observer.observe(sentinel);
		}

		document.addEventListener('keydown', handleKeydown);

		return () => {
			document.removeEventListener('keydown', handleKeydown);
			observer.disconnect();
		};
	});

	function handleSearch(event: Event) {
		event.preventDefault();
		if (searchQuery.trim()) {
			onSearch?.(searchQuery.trim());
		}
	}
</script>

<!-- Sentinel pour détecter le sticky -->
<div class="sticky-sentinel"></div>

<div class="search-bar-container" class:is-sticky={isSticky}>
	<div class="search-bar-content">
		<form class="search-form" onsubmit={handleSearch}>
			<div class="search-input-wrapper">
				<svg
					class="search-icon"
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<circle cx="11" cy="11" r="8" />
					<path d="m21 21-4.35-4.35" />
				</svg>
				<input
					bind:this={searchInput}
					bind:value={searchQuery}
					type="text"
					{placeholder}
					class="search-input"
					aria-label="Barre de recherche"
				/>
				{#if showShortcutHint}
					<kbd class="shortcut-hint">Ctrl+K</kbd>
				{/if}
			</div>
		</form>
	</div>
</div>

<style>
	.sticky-sentinel {
		height: 1px;
		position: absolute;
		top: 0;
		width: 100%;
	}

	.search-bar-container {
		position: sticky;
		top: 0;
		z-index: 1100;
		background: transparent;
		border-bottom: 1px solid transparent;
		padding: 12px 0;
	}

	.search-bar-container.is-sticky {
		border-bottom-color: rgba(229, 231, 235, 0.8);
		box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
	}

	.search-bar-container.is-sticky .search-input {
		background: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(10px);
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}

	.search-bar-content {
		max-width: 1200px;
		margin: 0 auto;
		padding: 0 16px;
	}

	.search-form {
		max-width: 600px;
		margin: 0 auto;
	}

	.search-input-wrapper {
		position: relative;
		display: flex;
		align-items: center;
	}

	.search-icon {
		position: absolute;
		left: 12px;
		width: 18px;
		height: 18px;
		color: #9ca3af;
		pointer-events: none;
		z-index: 1;
	}

	.search-input {
		width: 100%;
		height: 44px;
		padding: 0 80px 0 44px;
		border: 2px solid #e5e7eb;
		border-radius: 22px;
		background: #f9fafb;
		font-size: 16px;
		outline: none;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
	}

	.search-input:focus {
		border-color: #6366f1;
		background: white;
		box-shadow:
			0 0 0 3px rgba(99, 102, 241, 0.1),
			0 4px 12px rgba(0, 0, 0, 0.15);
	}

	.shortcut-hint {
		position: absolute;
		right: 12px;
		background: #e5e7eb;
		color: #6b7280;
		padding: 4px 8px;
		border-radius: 6px;
		font-size: 12px;
		font-family: monospace;
		pointer-events: none;
		font-weight: 500;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.search-bar-content {
			padding: 0 12px;
		}

		.search-input {
			font-size: 16px;
		}

		.shortcut-hint {
			display: none;
		}
	}
</style>
