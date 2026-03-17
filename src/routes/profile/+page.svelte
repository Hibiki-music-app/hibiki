<script lang="ts">
	import { RandomAvatar, ClientRouter } from '$lib/services/ApiEndpoints';
	import { User, Mail, Calendar, KeyRound, LogOut } from 'lucide-svelte';

	export let data;
	const { user } = data;

	function formatDate(date: Date | string | null | undefined): string {
		if (!date) return '—';
		return new Date(date).toLocaleDateString('fr-FR', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
		});
	}
</script>

<svelte:head>
	<title>Hibiki — Profil</title>
</svelte:head>

<div class="flex justify-center py-8" style="animation: fade-in 0.3s ease;">
	<div
		class="glass-medium rounded-[1.25rem] w-full max-w-md overflow-hidden"
		style="box-shadow: 0 25px 60px rgba(2,6,23,0.4), 0 4px 16px rgba(15,23,42,0.3), inset 0 1px 0 rgba(255,255,255,0.06);"
	>
		<!-- Header with avatar -->
		<div
			class="px-8 py-8 text-center"
			style="background: linear-gradient(135deg, rgba(59,130,246,0.08) 0%, rgba(99,102,241,0.06) 100%);
				   border-bottom: 1px solid rgba(148,163,184,0.1);"
		>
			<div class="relative inline-block mb-4">
				<img
					src={user?.image || RandomAvatar}
					alt="Avatar"
					class="w-20 h-20 rounded-full object-cover"
					style="box-shadow: 0 0 30px rgba(59,130,246,0.3), 0 0 60px rgba(59,130,246,0.1); border: 2px solid rgba(59,130,246,0.3);"
				/>
				<!-- Online indicator -->
				<div
					class="absolute bottom-1 right-1 w-4 h-4 rounded-full border-2 border-[#0f172a]"
					style="background: #22c55e;"
				></div>
			</div>

			<h1 class="text-xl font-bold text-[#f8fafc]">{user?.name ?? 'Utilisateur'}</h1>
			<p class="text-sm text-[#94a3b8] mt-1">Membre Hibiki</p>
		</div>

		<!-- Info section -->
		<div class="px-6 py-4">
			<p class="text-xs text-[#64748b] uppercase tracking-wider font-semibold mb-3">
				Informations du compte
			</p>

			<div class="space-y-1">
				<div class="glass-action rounded-xl">
					<div class="flex items-center gap-3 text-sm text-[#94a3b8]">
						<Mail size={15} />
						<span>Email</span>
					</div>
					<span class="text-sm text-[#f8fafc] truncate max-w-[200px]">{user?.email ?? '—'}</span>
				</div>

				<div class="glass-action rounded-xl">
					<div class="flex items-center gap-3 text-sm text-[#94a3b8]">
						<KeyRound size={15} />
						<span>Authentification</span>
					</div>
					<span
						class="text-xs px-2 py-1 rounded-full"
						style="background: rgba(59,130,246,0.15); color: #93c5fd; border: 1px solid rgba(59,130,246,0.3);"
					>
						Passkey
					</span>
				</div>

				<div class="glass-action rounded-xl">
					<div class="flex items-center gap-3 text-sm text-[#94a3b8]">
						<Calendar size={15} />
						<span>Membre depuis</span>
					</div>
					<span class="text-xs text-[#f8fafc]">{formatDate(user?.createdAt)}</span>
				</div>

				<div class="glass-action rounded-xl">
					<div class="flex items-center gap-3 text-sm text-[#94a3b8]">
						<User size={15} />
						<span>ID</span>
					</div>
					<span class="text-xs font-mono text-[#64748b] truncate max-w-[120px]">{user?.id}</span>
				</div>
			</div>
		</div>

		<!-- Actions -->
		<div class="px-6 pb-6">
			<a
				href={ClientRouter.logout}
				class="btn-glass w-full justify-center gap-2 py-3 touch-target"
				style="color: #f87171; border-color: rgba(239,68,68,0.3);"
			>
				<LogOut size={16} />
				Se déconnecter
			</a>
		</div>
	</div>
</div>

<!-- Bottom padding for player -->
<div class="h-28"></div>
