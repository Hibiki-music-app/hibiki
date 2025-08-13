<script lang="ts">
	import { onMount } from 'svelte';
	import userIcon from '$lib/assets/user.svg';
	import favicon from '$lib/assets/favicon.svg';

	// État pour la popup des paramètres
	let showUserMenu = $state(false);
	let searchInput: HTMLInputElement | null = $state(null);

	// État de connexion de l'utilisateur (à adapter selon votre système d'auth)
	let isUserLoggedIn = $state(true); // Changez ceci selon votre logique d'authentification
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
			// Fermer la popup si on clique ailleurs
			if (event.key === 'Escape') {
				showUserMenu = false;
			}
		};

		const handleClickOutside = (event: MouseEvent) => {
			const target = event.target as Element;
			if (!target.closest('.user-menu-container')) {
				showUserMenu = false;
			}
		};

		document.addEventListener('keydown', handleKeydown);
		document.addEventListener('click', handleClickOutside);

		return () => {
			document.removeEventListener('keydown', handleKeydown);
			document.removeEventListener('click', handleClickOutside);
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

	function toggleUserMenu() {
		showUserMenu = !showUserMenu;
	}
</script>

<header class="navbar fixed top-0 z-5 flex items-center justify-between">
	<button class="btn rounded-full" onclick={handleHomeClick} aria-label="Hibiki - Accueil">
		<img class="w-8 h-8" src={favicon} alt="Icône Hibiki" />
		<span class="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
			Hibiki
		</span>
	</button>

	<div class="user-menu-container relative">
		<button class="btn flex items-center rounded-full" onclick={toggleUserMenu} aria-label="Menu utilisateur">
			<img class="w-8 h-8 rounded-full" src={userAvatar} alt="Avatar de l'utilisateur" />
			<svg
				class="w-5 h-5 text-base-content ml-2"
				xmlns="http://www.w3.org/2000/svg"
				fill="none"
				viewBox="0 0 24 24"
				stroke="currentColor"
			>
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 9l6 6 6-6" />
			</svg>
		</button>
		{#if showUserMenu}
			<div class="dropdown-content absolute right-0 mt-2 w-48 bg-base-100 shadow-lg rounded-lg">
				<ul class="menu p-2">
					<li>
						<a href="/profile" class="flex items-center gap-2">
							<img class="w-6 h-6 rounded-full" src={userAvatar} alt="Avatar de l'utilisateur" />
							<span>{userName}</span>
						</a>
					</li>
					<li><a href="/settings">Paramètres</a></li>
					<li><a href="/logout">Déconnexion</a></li>
				</ul>
			</div>
		{/if}
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

	/* Animation pour le dropdown */
	.dropdown-content {
		animation: slideDown 0.2s ease;
	}

	@keyframes slideDown {
		from {
			opacity: 0;
			transform: translateY(-8px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}

	/* Responsive */
	@media (max-width: 768px) {
		.navbar {
			padding: 0 0.75rem;
		}

		.dropdown-content {
			width: 90vw;
			max-width: 320px;
			right: 0.75rem;
		}
	}
</style>
