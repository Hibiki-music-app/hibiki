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

<Header onHomeClick={handleHomeClick} />

<main class="main-layout">
	{@render children?.()}
</main>

<style>
	:global(body) {
		margin: 0;
		padding: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
		background-color: #fafafa;
	}

	.main-layout {
		position: relative;
	}
</style>
