<script lang="ts">
	import { onMount } from 'svelte';
	import userIcon from '$lib/assets/user.svg';
	import favicon from '$lib/assets/favicon.svg';

	//theme
	let currentTheme = $state('auto');

	// État de connexion de l'utilisateur (à adapter selon votre système d'auth)
	let isUserLoggedIn = $state(true);
	let userAvatar = $state(
		'https://s4.anilist.co/file/anilistcdn/user/avatar/large/b647930-aI0nneV0XlFa.png'
	);
	let userName = $state('Utilisateur');

	function setTheme(theme: string) {
		currentTheme = theme;
		document.documentElement.setAttribute('data-theme', theme);
	}
</script>

<header class="navbar fixed top-0 z-5 flex items-center justify-between p-4">
	<a
		href="/"
		class="btn border rounded-full bg-base-100 border-base-300 hover:border-base-content/20 hover:shadow-lg transition-all duration-200 ease-in-out"
		aria-label="Hibiki - Accueil"
	>
		<img class="w-8 h-8" src={favicon} alt="Icône Hibiki" />
		<span class="bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
			Hibiki
		</span>
	</a>

	<div class="dropdown dropdown-end">
		<button
			tabindex="0"
			class="btn btn-ghost btn-circle avatar border border-base-300 hover:border-base-content/20 hover:shadow-lg transition-all duration-200 ease-in-out"
		>
			<div class="w-10 rounded-full">
				{#if isUserLoggedIn}
					<img src={userAvatar} alt="Avatar de l'utilisateur" />
				{:else}
					<img src={userIcon} alt="Icône utilisateur" />
				{/if}
			</div>
		</button>
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
			<li><a href="/auth/logout">Déconnexion</a></li>
			<li>
				<div class="flex items-center justify-between p-2">
					<span class="text-sm font-medium">Thème</span>
					<div class="flex items-center gap-1 border border-base-200 rounded">
						<button
							class="btn btn-ghost btn-sm p-2 {currentTheme === 'light' ? 'btn-active' : ''}"
							title="Mode clair"
							aria-label="Mode clair"
							onclick={() => setTheme('light')}
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<circle cx="12" cy="12" r="5"></circle>
								<path
									d="M12 1v2m0 18v2M4.22 4.22l1.42 1.42m12.72 12.72l1.42 1.42M1 12h2m18 0h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"
								></path>
							</svg>
						</button>
						<button
							class="btn btn-ghost btn-sm p-2 {currentTheme === 'dark' ? 'btn-active' : ''}"
							title="Mode sombre"
							aria-label="Mode sombre"
							onclick={() => setTheme('dark')}
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
							</svg>
						</button>
						<button
							class="btn btn-ghost btn-sm p-2 {currentTheme === 'auto' ? 'btn-active' : ''}"
							title="Mode automatique"
							aria-label="Mode automatique"
							onclick={() => setTheme('auto')}
						>
							<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<rect x="2" y="3" width="20" height="14" rx="2" ry="2"></rect>
								<line x1="8" y1="21" x2="16" y2="21"></line>
								<line x1="12" y1="17" x2="12" y2="21"></line>
							</svg>
						</button>
					</div>
				</div>
			</li>
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
