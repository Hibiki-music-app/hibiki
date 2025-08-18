<script lang="ts">
    import { authClient } from "$lib/auth-client";
    let email = "", password = "", error = "";

    async function signUp() {
        try {
            await authClient.signUp.email({ email, password, callbackURL: "/" });
        } catch (e) {
            error = e?.message || "Signup failed";
        }
    }
</script>

<form on:submit|preventDefault={signUp}>
    <input type="email" bind:value={email} placeholder="Email" required />
    <input type="password" bind:value={password} placeholder="Password" required />
    <button type="submit">Sign up</button>
    {#if error}<p class="error">{error}</p>{/if}
</form>
