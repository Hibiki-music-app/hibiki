<script lang="ts">
    import {authClient} from "$lib/auth-client";
    import {onMount} from "svelte";

    const {data}: {data: {user: import("better-auth").User | null}}= $props();

    onMount(()=> {
        if (!data?.user) {
            window.location.replace("/auth/login");
        }
    });
</script>
{#if data?.user}
    <div class="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
        <p class="">{data?.user.email}</p>
        <button class="btn" onclick={async () => {
            await authClient.signOut({
                fetchOptions: {
                    onSuccess: () => {
                        window.location.reload();
                    }
                }
            });
        }}>
            Se déconnecter
        </button>
    </div>

{/if}
