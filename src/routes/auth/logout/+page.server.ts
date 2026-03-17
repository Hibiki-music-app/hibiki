import { auth } from '$lib/auth';
import { type Actions, redirect } from '@sveltejs/kit';
import { ClientRouter } from '$lib/services/ApiEndpoints';

export const actions: Actions = {
	default: async ({ request }) => {
		await auth.api.signOut({
			headers: request.headers,
		});
		throw redirect(303, ClientRouter.auth);
	},
};
