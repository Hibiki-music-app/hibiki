import type { LayoutServerLoad } from "./$types";
import { redirect } from "@sveltejs/kit";
import { auth } from "$lib/auth";

export const load: LayoutServerLoad = async ({ request }) => {
	const session = await auth.api.getSession({ headers: request.headers });
	if (!session) throw redirect(302, "/login");
	return { user: session.user };
};
