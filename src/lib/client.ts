import { oneTapClient } from "better-auth/client/plugins";
import { createAuthClient } from "better-auth/svelte";


export const authClient = createAuthClient({
    plugins: [
        oneTapClient({
            clientId: import.meta.env.VITE_GOOGLE_CLIENT_ID,
            autoSelect: false,
            cancelOnTapOutside: true,
            context: "signin",
            promptOptions: {
                baseDelay: 1000,
                maxAttempts: 5,
            }
        })
    ]
});

const signIn = async () => {
    const data = await authClient.signIn.social({
        provider: "google",
    });
};