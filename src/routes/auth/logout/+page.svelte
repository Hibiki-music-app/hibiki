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
        <div class="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
            <p>{data?.email}</p>
            <form method="POST">
                <button type="submit" class="btn">Se déconnecter</button>
            </form>
        </div>
    </div>
{/if}