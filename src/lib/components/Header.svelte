<script lang="ts">
	import favicon from '$lib/assets/favicon.svg';
	import { ClientRouter, RandomAvatar } from '$lib/services/ApiEndpoints';
	import type { UserType } from '$lib/models/UserType';
	import { User, LogOut, Settings, Home } from 'lucide-svelte';

	let {
		user = $bindable(),
		onHomeClick = $bindable(),
	}: {
		user: UserType | null;
		onHomeClick: () => void;
	} = $props();

	let userMenuOpen = $state(false);

	function handleHomeClick() {
		onHomeClick?.();
		userMenuOpen = false;
	}

	function toggleUserMenu(e: Event) {
		e.stopPropagation();
		userMenuOpen = !userMenuOpen;
	}

	function closeMenu() {
		userMenuOpen = false;
	}
</script>

<svelte:window onclick={closeMenu} />

<header
	class="glass-strong sticky top-0 z-40 flex items-center justify-between gap-4 border-b border-[rgba(148,163,184,0.1)]"
	style="padding: 0.75rem clamp(1.2rem, 2.5vw, 2rem);"
>
	<!-- Logo / Home -->
	<button class="btn-glass gap-2" onclick={handleHomeClick} aria-label="Hibiki — Accueil">
		<img class="w-5 h-5" src={favicon} alt="Hibiki" />
		<span class="gradient-text font-bold text-base">Hibiki</span>
	</button>

	<!-- User menu -->
	<div class="relative">
		<button
			class="btn-glass gap-2 touch-target"
			onclick={toggleUserMenu}
			aria-label="Menu utilisateur"
			aria-expanded={userMenuOpen}
		>
			{#if user?.image}
				<img
					src={user.image}
					alt="Avatar"
					class="w-6 h-6 rounded-full object-cover"
					style="box-shadow: 0 0 8px rgba(59,130,246,0.3);"
				/>
			{:else}
				<div
					class="w-6 h-6 rounded-full flex items-center justify-center"
					style="background: rgba(59,130,246,0.2); border: 1px solid rgba(59,130,246,0.4);"
				>
					<User size={14} class="text-[#60a5fa]" />
				</div>
			{/if}
			<span class="text-sm hidden sm:inline">{user?.name ?? 'Compte'}</span>
		</button>

		{#if userMenuOpen}
			<div
				class="glass-strong absolute right-0 top-full mt-2 w-52 rounded-[var(--radius-xl)] overflow-hidden"
				style="box-shadow: 0 20px 50px rgba(2,6,23,0.5), 0 4px 16px rgba(15,23,42,0.3);
					   animation: slide-up 0.15s ease;"
				role="menu"
			>
				{#if user}
					<!-- User info -->
					<div class="px-4 py-3 border-b border-[rgba(148,163,184,0.1)]">
						<p class="text-sm font-medium text-[#f8fafc] truncate">{user.name}</p>
						<p class="text-xs text-[#64748b] truncate">{user.email}</p>
					</div>

					<div class="p-1">
						<a
							href={ClientRouter.profile}
							class="glass-action rounded-xl w-full text-sm text-[#f8fafc] gap-3"
							role="menuitem"
							onclick={closeMenu}
						>
							<div class="flex items-center gap-3">
								<User size={15} class="text-[#94a3b8]" />
								<span>Profil</span>
							</div>
						</a>
						<a
							href={ClientRouter.logout}
							class="glass-action rounded-xl w-full text-sm text-[#f87171] gap-3"
							role="menuitem"
							onclick={closeMenu}
						>
							<div class="flex items-center gap-3">
								<LogOut size={15} class="text-[#f87171]" />
								<span>Déconnexion</span>
							</div>
						</a>
					</div>
				{:else}
					<div class="p-1">
						<a
							href={ClientRouter.auth}
							class="glass-action rounded-xl w-full text-sm text-[#60a5fa] gap-3"
							role="menuitem"
							onclick={closeMenu}
						>
							<div class="flex items-center gap-3">
								<User size={15} />
								<span>Se connecter</span>
							</div>
						</a>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</header>
