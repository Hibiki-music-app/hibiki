import { createAuthClient } from 'better-auth/svelte';
import { passkeyClient } from 'better-auth/client/plugins';

export const authClient = createAuthClient({
	plugins: [passkeyClient()],
});
