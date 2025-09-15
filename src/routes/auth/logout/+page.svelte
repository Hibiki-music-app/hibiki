<script lang="ts">
    import {onMount} from "svelte";
    import type {UserType} from "$lib/models/UserType";

    const {data}: {
        data: UserType | null,
    } = $props();

    onMount(()=> {
        if (!data) {
            window.location.replace("/auth/login");
        }
    });
</script>

{#if data}
    <div class="min-h-screen flex items-center justify-center p-4">
            <p>{data?.email}</p>
            <form method="POST">
                <button type="submit" class="btn btn-outline btn-warning">Se déconnecter</button>
            </form>
    </div>
{/if}
