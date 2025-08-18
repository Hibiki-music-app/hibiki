<script lang="ts">
    import { authClient } from "$lib/auth-client";
    let email = "", password = "", error = "";

    async function signIn() {
        try {
            await authClient.signIn.email({ email, password, callbackURL: "/" });
        } catch (e) {
            error = e?.message || "Login failed";
        }
    }
</script>

<form on:submit|preventDefault={signIn}>
    <input type="email" bind:value={email} placeholder="Email" required />
    <input type="password" bind:value={password} placeholder="Password" required />
    <button type="submit">Sign in</button>
    {#if error}<p class="error">{error}</p>{/if}
</form>
