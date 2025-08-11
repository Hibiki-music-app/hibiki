<script lang="ts">
	import { onMount } from 'svelte';
	import userIcon from '$lib/assets/user.svg';

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

<header class="sticky-header">
	<div class="header-container">
		<!-- Logo et nom Hibiki à gauche -->
		<button class="brand-logo" onclick={handleHomeClick} aria-label="Hibiki - Accueil">
			<div class="logo-icon">
				<svg
					width="28"
					height="28"
					viewBox="0 0 24 24"
					fill="none"
					stroke="currentColor"
					stroke-width="2"
					stroke-linecap="round"
					stroke-linejoin="round"
				>
					<path d="M9 18V5l12-2v13" />
					<circle cx="6" cy="18" r="3" />
					<circle cx="18" cy="16" r="3" />
				</svg>
			</div>
			<span class="brand-name">Hibiki</span>
		</button>

		<!-- Avatar utilisateur à droite -->
		<div class="user-menu-container">
			<button
				class="user-avatar"
				class:logged-out={!isUserLoggedIn}
				onclick={toggleUserMenu}
				aria-label="Menu utilisateur"
			>
				{#if isUserLoggedIn}
					<img src={userAvatar} alt="Avatar utilisateur" />
				{:else}
					<img src={userIcon} alt="Utilisateur non connecté" />
				{/if}
			</button>

			<!-- Popup des paramètres -->
			{#if showUserMenu}
				<div class="user-menu">
					{#if isUserLoggedIn}
						<div class="user-info">
							<img src={userAvatar} alt="Avatar" />
							<div>
								<div class="username">{userName}</div>
							</div>
						</div>
						<hr />
						<nav class="menu-items">
							<button class="menu-item">
								<svg
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
								>
									<circle cx="12" cy="12" r="3" />
									<path d="M12 1v6m0 6v6m11-7h-6m-6 0H1" />
								</svg>
								Paramètres
							</button>
							<button class="menu-item">
								<svg
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
								>
									<path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
									<circle cx="9" cy="7" r="4" />
									<path d="M22 21v-2a4 4 0 0 0-3-3.87M16 3.13a4 4 0 0 1 0 7.75" />
								</svg>
								Profil
							</button>
							<button class="menu-item">
								<svg
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
								>
									<path d="M3 3h18v18H3zM9 9l6 6m0-6l-6 6" />
								</svg>
								Playlists
							</button>
							<button class="menu-item">
								<svg
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
								>
									<path
										d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"
									/>
								</svg>
								Aide
							</button>
							<hr />
							<button class="menu-item danger">
								<svg
									width="18"
									height="18"
									viewBox="0 0 24 24"
									fill="none"
									stroke="currentColor"
									stroke-width="2"
								>
									<path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
									<polyline points="16,17 21,12 16,7" />
									<line x1="21" y1="12" x2="9" y2="12" />
								</svg>
								Déconnexion
							</button>
						</nav>
					{:else}
						<div class="login-prompt">
							<div class="login-icon">
								<img src={userIcon} alt="Connexion" />
							</div>
							<h3>Non connecté</h3>
							<p>Connectez-vous pour accéder à vos playlists et paramètres personnalisés.</p>
							<div class="login-actions">
								<button class="login-button primary">Se connecter</button>
								<button class="login-button secondary">S'inscrire</button>
							</div>
						</div>
					{/if}
				</div>
			{/if}
		</div>
	</div>
</header>

<style>
	.sticky-header {
		position: fixed;
		top: 0;
		left: 0;
		right: 0;
		z-index: 1000;
	}

	.header-container {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin: 0 auto;
		padding: 0 16px;
		height: 64px;
		gap: 16px;
	}

	.brand-logo {
		display: flex;
		align-items: center;
		gap: 12px;
		padding: 8px 16px;
		border: none;
		background: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(10px);
		border-radius: 12px;
		cursor: pointer;
		transition: all 0.2s ease;
		color: #374151;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
	}

	.brand-logo:hover {
		background: rgba(255, 255, 255, 1);
		transform: translateY(-1px);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	}

	.logo-icon {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		border-radius: 10px;
		color: white;
		box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
	}

	.brand-name {
		font-size: 1.25rem;
		font-weight: 700;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
		letter-spacing: -0.025em;
	}

	.app-title {
		flex: 1;
		text-align: center;
	}

	.app-title h1 {
		margin: 0;
		font-size: 1.5rem;
		font-weight: 700;
		color: #111827;
		background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
		-webkit-background-clip: text;
		-webkit-text-fill-color: transparent;
		background-clip: text;
	}

	.user-menu-container {
		position: relative;
	}

	.user-avatar {
		width: 44px;
		height: 44px;
		border: none;
		background: rgba(255, 255, 255, 0.95);
		backdrop-filter: blur(10px);
		border-radius: 50%;
		cursor: pointer;
		transition: all 0.2s ease;
		overflow: hidden;
		display: flex;
		align-items: center;
		justify-content: center;
		box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
		padding: 2px;
	}

	.user-avatar.logged-out {
		background: rgba(255, 255, 255, 0.95);
		padding: 8px;
	}

	.user-avatar.logged-out img {
		width: 24px;
		height: 24px;
		opacity: 0.7;
	}

	.user-avatar:hover {
		transform: scale(1.05);
		background: rgba(255, 255, 255, 1);
		box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
	}

	.user-avatar img {
		width: calc(100% - 4px);
		height: calc(100% - 4px);
		object-fit: cover;
		border-radius: 50%;
	}

	.user-menu {
		position: absolute;
		top: calc(100% + 8px);
		right: 0;
		width: 280px;
		background: white;
		border: 1px solid #e5e7eb;
		border-radius: 12px;
		box-shadow:
			0 10px 25px -3px rgba(0, 0, 0, 0.1),
			0 4px 6px -2px rgba(0, 0, 0, 0.05);
		padding: 16px;
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

	.user-info {
		display: flex;
		align-items: center;
		gap: 12px;
		margin-bottom: 12px;
	}

	.user-info img {
		width: 50px;
		height: 50px;
		border-radius: 50%;
		object-fit: cover;
	}

	.username {
		font-weight: 600;
		color: #111827;
		margin-bottom: 2px;
	}

	hr {
		border: none;
		height: 1px;
		background: #e5e7eb;
		margin: 12px 0;
	}

	.menu-items {
		display: flex;
		flex-direction: column;
		gap: 2px;
	}

	.menu-item {
		display: flex;
		align-items: center;
		gap: 12px;
		width: 100%;
		padding: 12px;
		border: none;
		background: none;
		border-radius: 8px;
		cursor: pointer;
		font-size: 14px;
		color: #374151;
		transition: background-color 0.2s ease;
		text-align: left;
	}

	.menu-item:hover {
		background: #f3f4f6;
	}

	.menu-item.danger {
		color: #dc2626;
	}

	.menu-item.danger:hover {
		background: #fef2f2;
	}

	.login-prompt {
		text-align: center;
		padding: 8px;
	}

	.login-icon {
		margin-bottom: 16px;
		display: flex;
		justify-content: center;
	}

	.login-icon img {
		width: 48px;
		height: 48px;
		background: #f3f4f6;
		padding: 12px;
		border-radius: 50%;
		opacity: 0.7;
	}

	.login-prompt h3 {
		margin: 0 0 8px 0;
		font-size: 18px;
		color: #111827;
	}

	.login-prompt p {
		margin: 0 0 20px 0;
		font-size: 14px;
		color: #6b7280;
		line-height: 1.4;
	}

	.login-actions {
		display: flex;
		flex-direction: column;
		gap: 8px;
	}

	.login-button {
		padding: 10px 20px;
		border: none;
		border-radius: 8px;
		font-size: 14px;
		font-weight: 500;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.login-button.primary {
		background: #6366f1;
		color: white;
	}

	.login-button.primary:hover {
		background: #5b5bf6;
	}

	.login-button.secondary {
		background: #f3f4f6;
		color: #374151;
		border: 1px solid #d1d5db;
	}

	.login-button.secondary:hover {
		background: #e5e7eb;
	}

	/* Responsive */
	@media (max-width: 768px) {
		.header-container {
			padding: 0 12px;
			gap: 12px;
		}

		.brand-name {
			font-size: 1.125rem;
		}

		.logo-icon {
			width: 36px;
			height: 36px;
		}

		.logo-icon svg {
			width: 24px;
			height: 24px;
		}

		.app-title h1 {
			font-size: 1.25rem;
		}

		.user-menu {
			width: 260px;
			right: -12px;
		}
	}
</style>
