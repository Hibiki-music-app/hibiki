import type { Actions } from "@sveltejs/kit";
import { auth } from "$lib/auth";
import { ClientRouter } from "$lib/services/ApiEndpoints";


export const actions: Actions = {
    login: async ({ request }) => {
        const formData = await request.formData();
        const email = formData.get('email') as string;
        const password = formData.get('password') as string;
        try {
            const data = await auth.api.signInEmail({
                body: {
                    email: email,
                    password: password,
                    rememberMe: true,
                    callbackURL: ClientRouter.profile,
                },
                // headers: await headers()
                headers: {}
            })
            return { success: true };
        } catch (error: any) {
            return { success: false, message: error.message || 'Login failed' };
        }
    }
}