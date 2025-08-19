<script>
    import { page } from '$app/stores';
</script>

<svelte:head>
    <title>Erreur {$page.status}</title>
</svelte:head>

<div class="hero min-h-screen bg-base-200">
    <div class="hero-content text-center">
        <div class="max-w-md">
            <h1 class="text-9xl font-bold text-error">{$page.status}</h1>
            <h2 class="text-2xl font-bold py-6">Oops! Une erreur s'est produite</h2>
            <p class="text-base-content/70 mb-6">
                {$page.error?.message || 'Une erreur inattendue s\'est produite'}
            </p>
            <a href="/" class="btn btn-primary">Retour à l'accueil</a>
        </div>
    </div>
</div>
