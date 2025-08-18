import type { Actions } from "@sveltejs/kit";
import { auth } from "$lib/auth";
import { getRequestEvent } from "$app/server";

export const actions: Actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const email = data.get("email") as string;
		const password = data.get("password") as string;
		const event = getRequestEvent();

		const res = await auth.api.signUpEmail({
			body: { email, password },
			asResponse: true,
			request: event.request,
		});

		return res;
	}
};
