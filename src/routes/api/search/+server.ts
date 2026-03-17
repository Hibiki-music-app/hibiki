import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { MONOCHROME_API } from '$lib/services/ApiEndpoints';

export const GET: RequestHandler = async ({ url }) => {
	const query = url.searchParams.get('q');
	if (!query?.trim()) {
		return json({ error: 'Missing query' }, { status: 400 });
	}

	const response = await fetch(
		`${MONOCHROME_API}/search/?s=${encodeURIComponent(query)}`,
		{ headers: { 'User-Agent': 'Hibiki/0.1.0' } }
	);

	if (!response.ok) {
		return json({ error: 'Upstream error' }, { status: response.status });
	}

	const data: unknown = await response.json();
	return json(data);
};
