<script lang="ts">
	import type { Recording } from '$lib/models/Musicbrainz';

	let q = '';
	let results: Recording[] = [];
	let loading = false;
	let error: string | null = null;

	async function search() {
		error = null;
		if (!q.trim()) return;
		loading = true;

		try {
			const res = await fetch(
				'api?q=' + encodeURIComponent(q),
			);
			if (!res.ok) {
				const err = await res.json().catch(() => ({ message: 'unknown' }));
				error = `Erreur: ${res.status} ${JSON.stringify(err)}`;
				results = [];
			} else {
				const data = await res.json();
				results = data.recordings ?? [];
			}
		} catch (e) {
			error = String(e);
			results = [];
		} finally {
			loading = false;
		}
	}
</script>

<input
	bind:value={q}
	placeholder="Titre..."
	on:keydown={(e) => e.key === 'Enter' && search()}
/>

<button
	on:click={search}
	disabled={loading}>
	{loading ? '...' : 'Rechercher'}
</button>

{#if error}
	<div class="error">{error}</div>
{/if}

{#each results as r (r.id)}
	<li>
		<strong>{r.title}</strong>
		{#if r['artist-credit']?.length}
			— {r['artist-credit'][0].name}
		{/if}

		{#if r.releases?.length}
			<ul>
				{#each r.releases as release (release.id)}
					<li>{release.title} ({release.status})</li>
				{/each}
			</ul>
		{/if}
	</li>
{/each}