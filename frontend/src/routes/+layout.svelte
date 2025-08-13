<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import "../app.css";
	import Header from '../components/Header.svelte';

	let { children } = $props();

	function handleHomeClick() {
		// Navigation vers la page d'accueil
		window.location.href = '/';
	}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<!-- DaisyUI Layout with Header -->
<div class="min-h-screen bg-base-200">
	<Header onHomeClick={handleHomeClick} />
	
	<main>
		{@render children?.()}
	</main>
</div>

<style>
	/* DaisyUI global overrides - minimisés */
	:global(body) {
		margin: 0;
		padding: 0;
	}
</style>
