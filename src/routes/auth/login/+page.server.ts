import { type Actions, redirect } from '@sveltejs/kit';
import { auth } from "$lib/auth";
import { ClientRouter } from "$lib/services/ApiEndpoints";


export const actions: Actions = {
	default: async ({ request }) => {
        const formData = await request.formData();
        const email = formData.get('email') as string;
        const password = formData.get('password') as string;
				const rHeaders = request.headers;
        try {
            const data = await auth.api.signInEmail({
                body: {
                    email: email,
                    password: password,
                    rememberMe: true,
                    callbackURL: ClientRouter.profile,
                },
							headers: rHeaders
            })
            return {
							success: true,
							redirect: ClientRouter.profile
						};
        } catch (error: any) {
            return { success: false, message: error.message || 'Login failed' };
        }
    }
}