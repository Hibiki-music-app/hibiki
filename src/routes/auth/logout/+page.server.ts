import { auth } from '$lib/auth';
import { type Actions, redirect } from '@sveltejs/kit';
import { ClientRouter } from '$lib/services/ApiEndpoints';



// actions call dans le form POST
export const actions: Actions = {
	default: async ({request}) => {
		const rHeaders = request.headers;
		await auth.api.signOut({
			headers: rHeaders
		});
		//cookies.delete('session', { path: '/'});
		throw redirect(303, ClientRouter.login)
	}

}
