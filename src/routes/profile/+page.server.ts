import { ClientRouter } from '$lib/services/ApiEndpoints.js';
import { redirect } from '@sveltejs/kit';

export const load = async ({ locals }) => {
	if (!locals.user) {
		throw redirect(303, ClientRouter.auth);
	}

	return {
		user: locals.user,
	};
};
