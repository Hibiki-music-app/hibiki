<script lang="ts">
	import { authClient } from '$lib/client';
	import { goto } from '$app/navigation';
	import { Fingerprint, UserPlus, Loader2 } from 'lucide-svelte';

	let mode: 'signin' | 'register' = $state('signin');
	let username = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSignIn() {
		loading = true;
		error = '';
		try {
			const result = await authClient.signIn.passkey();
			if (result?.error) {
				error = result.error.message ?? 'Échec de la connexion';
			} else {
				goto('/');
			}
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : 'Échec de la connexion';
		} finally {
			loading = false;
		}
	}

	async function handleRegister() {
		if (!username.trim()) {
			error = "Veuillez entrer un nom d'utilisateur";
			return;
		}
		loading = true;
		error = '';
		try {
			// Email fantôme interne — jamais exposé à l'utilisateur
			const ghostEmail = `${username.trim().toLowerCase().replace(/\s+/g, '-')}-${crypto.randomUUID().slice(0, 8)}@passkey.local`;

			// 1. Créer le compte avec identifiants internes
			const signUpResult = await authClient.signUp.email({
				email: ghostEmail,
				name: username.trim(),
				password: crypto.randomUUID(),
			});
			if (signUpResult?.error) {
				error = signUpResult.error.message ?? "Échec de la création du compte";
				return;
			}
			// 2. Enregistrer la passkey sur le compte tout juste créé
			const passkeyResult = await authClient.passkey.addPasskey({ name: username.trim() });
			if (passkeyResult?.error) {
				error = passkeyResult.error.message ?? "Échec de l'enregistrement de la passkey";
				return;
			}
			goto('/profile');
		} catch (e: unknown) {
			error = e instanceof Error ? e.message : "Échec de l'inscription";
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Hibiki — Connexion</title>
</svelte:head>

<div class="min-h-screen flex items-center justify-center p-4">
	<!-- Branding -->
	<div class="w-full max-w-sm">
		<div class="text-center mb-8">
			<h1 class="gradient-text text-4xl font-bold mb-2">Hibiki</h1>
			<p class="text-[#94a3b8] text-sm">Musique HiFi sans compromis</p>
		</div>

		<!-- Card -->
		<div class="glass-medium rounded-[1.25rem] p-8"
			style="box-shadow: 0 25px 60px rgba(2,6,23,0.4), 0 4px 16px rgba(15,23,42,0.3), inset 0 1px 0 rgba(255,255,255,0.08);">

			<!-- Tabs -->
			<div class="flex glass-subtle rounded-full p-1 mb-8">
				<button
					class="flex-1 py-2 rounded-full text-sm font-medium transition-all duration-200
						   {mode === 'signin' ? 'glass-strong text-white shadow-sm' : 'text-[#94a3b8] hover:text-white'}"
					onclick={() => { mode = 'signin'; error = ''; }}>
					Connexion
				</button>
				<button
					class="flex-1 py-2 rounded-full text-sm font-medium transition-all duration-200
						   {mode === 'register' ? 'glass-strong text-white shadow-sm' : 'text-[#94a3b8] hover:text-white'}"
					onclick={() => { mode = 'register'; error = ''; }}>
					Inscription
				</button>
			</div>

			{#if mode === 'signin'}
				<!-- Sign In -->
				<div class="text-center space-y-6" style="animation: fade-in 0.2s ease;">
					<div class="flex flex-col items-center gap-3">
						<div class="w-16 h-16 rounded-full glass-subtle flex items-center justify-center"
							style="box-shadow: 0 0 20px rgba(59,130,246,0.2);">
							<Fingerprint size={28} class="text-[#60a5fa]" />
						</div>
						<div>
							<p class="text-[#f8fafc] font-medium">Bienvenue</p>
							<p class="text-[#94a3b8] text-xs mt-1">Utilisez votre passkey pour vous connecter</p>
						</div>
					</div>

					<button
						class="btn-glass w-full justify-center gap-2 py-3 touch-target"
						onclick={handleSignIn}
						disabled={loading}>
						{#if loading}
							<Loader2 size={16} class="animate-spin" />
							<span>Connexion en cours…</span>
						{:else}
							<Fingerprint size={16} />
							<span>Se connecter avec passkey</span>
						{/if}
					</button>
				</div>
			{:else}
				<!-- Register -->
				<div class="space-y-5" style="animation: fade-in 0.2s ease;">
					<div class="flex flex-col items-center gap-3 text-center mb-6">
						<div class="w-16 h-16 rounded-full glass-subtle flex items-center justify-center"
							style="box-shadow: 0 0 20px rgba(59,130,246,0.2);">
							<UserPlus size={28} class="text-[#60a5fa]" />
						</div>
						<div>
							<p class="text-[#f8fafc] font-medium">Créer un compte</p>
							<p class="text-[#94a3b8] text-xs mt-1">Sans email, sans mot de passe — juste une passkey</p>
						</div>
					</div>

					<div>
						<label for="username" class="block text-xs font-medium text-[#94a3b8] mb-2 uppercase tracking-wider">
							Nom d'utilisateur
						</label>
						<input
							id="username"
							type="text"
							bind:value={username}
							placeholder="Votre prénom ou pseudo"
							class="w-full glass-subtle rounded-lg px-4 py-3 text-sm text-[#f8fafc]
								   placeholder:text-[#64748b] border border-[rgba(148,163,184,0.15)]
								   focus:border-[rgba(59,130,246,0.5)] focus:outline-none transition-colors"
							onkeydown={(e) => e.key === 'Enter' && handleRegister()}
						/>
					</div>

					<button
						class="btn-glass w-full justify-center gap-2 py-3 touch-target"
						onclick={handleRegister}
						disabled={loading}>
						{#if loading}
							<Loader2 size={16} class="animate-spin" />
							<span>Création en cours…</span>
						{:else}
							<Fingerprint size={16} />
							<span>Créer avec passkey</span>
						{/if}
					</button>
				</div>
			{/if}

			{#if error}
				<div class="mt-4 p-3 rounded-lg bg-[rgba(239,68,68,0.1)] border border-[rgba(239,68,68,0.25)]">
					<p class="text-[#f87171] text-sm text-center">{error}</p>
				</div>
			{/if}
		</div>

		<p class="text-center text-[#64748b] text-xs mt-6">
			Sécurisé par WebAuthn — aucun mot de passe stocké
		</p>
	</div>
</div>
