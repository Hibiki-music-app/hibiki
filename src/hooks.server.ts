import { auth } from "$lib/auth";
import { svelteKitHandler } from "better-auth/svelte-kit";
import { building } from '$app/environment'
import { redirect } from '@sveltejs/kit';


export async function handle({ event, resolve }) {

	const protectedRoutes = ['/profile'];
	if (protectedRoutes.includes(event.url.pathname)) {
		if (!event.locals) {
			console.log('redirect to login!')
			throw redirect(303, "/auth/login");
		}
	}
	return svelteKitHandler({ event, resolve, auth, building });
}
