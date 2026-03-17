import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { MONOCHROME_API } from '$lib/services/ApiEndpoints';

export const GET: RequestHandler = async ({ url }) => {
	const id = url.searchParams.get('id');
	if (!id) {
		return json({ error: 'Missing id' }, { status: 400 });
	}

	const response = await fetch(`${MONOCHROME_API}/lyrics/?id=${id}`, {
		headers: { 'User-Agent': 'Hibiki/0.1.0' },
	});

	if (!response.ok) {
		return json({ error: 'Not found' }, { status: 404 });
	}

	const data: unknown = await response.json();
	return json(data);
};
