<script lang="ts">
	import { onMount } from 'svelte';
	import userIcon from '$lib/assets/user.svg';
	import favicon from '$lib/assets/favicon.svg';

	// État pour la popup des paramètres (plus nécessaire avec DaisyUI dropdown)
	let searchInput: HTMLInputElement | null = $state(null);

	// État de connexion de l'utilisateur (à adapter selon votre système d'auth)
	let isUserLoggedIn = $state(true);
	let userAvatar = $state(
		'https://s4.anilist.co/file/anilistcdn/user/avatar/large/b647930-aI0nneV0XlFa.png'
	);
	let userName = $state('Utilisateur');

	// Gestion du raccourci Ctrl+K pour focus sur la barre de recherche
	onMount(() => {
		const handleKeydown = (event: KeyboardEvent) => {
			if (event.ctrlKey && event.key === 'k') {
				event.preventDefault();
				searchInput?.focus();
			}
		};

		document.addEventListener('keydown', handleKeydown);

		return () => {
			document.removeEventListener('keydown', handleKeydown);
		};
	});

	// Props pour les événements
	interface HeaderProps {
		onHomeClick?: () => void;
	}

	let { onHomeClick }: HeaderProps = $props();

	function handleHomeClick() {
		onHomeClick?.();
	}
</script>

<header class="navbar fixed top-0 z-5 flex items-center justify-between p-4">
	<button
		class="btn border rounded-full bg-base-100 border-base-100"
		onclick={handleHomeClick}
		aria-label="Hibiki - Accueil"
	>
		<img class="w-8 h-8" src={favicon} alt="Icône Hibiki" />
		<span class="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
			Hibiki
		</span>
	</button>

	<div class="dropdown dropdown-end">
		<div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
			<div class="w-10 rounded-full">
				{#if isUserLoggedIn}
					<img src={userAvatar} alt="Avatar de l'utilisateur" />
				{:else}
					<img src={userIcon} alt="Icône utilisateur" />
				{/if}
			</div>
		</div>
		<ul class="menu menu-sm dropdown-content bg-base-100 rounded-box z-[1] mt-3 w-52 p-2 shadow">
			<li>
				<a href="/profile" class="justify-between">
					<div class="flex items-center gap-2">
						<img class="w-6 h-6 rounded-full" src={userAvatar} alt="Avatar de l'utilisateur" />
						<span>{userName}</span>
					</div>
				</a>
			</li>
			<li><a href="/settings">Paramètres</a></li>
			<li><a href="/logout">Déconnexion</a></li>
		</ul>
	</div>
</header>

<style>
	/* Ajustements pour le header DaisyUI */
	.navbar {
		height: 4rem;
	}

	/* Améliorations pour le gradient text */
	.bg-clip-text {
		background-clip: text;
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.navbar {
			padding: 0 0.75rem;
		}
	}
</style>
