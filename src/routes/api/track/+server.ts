import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { MONOCHROME_API } from '$lib/services/ApiEndpoints';

export const GET: RequestHandler = async ({ url }) => {
	const id = url.searchParams.get('id');
	const quality = url.searchParams.get('quality') ?? 'LOSSLESS';

	if (!id) {
		return json({ error: 'Missing id' }, { status: 400 });
	}

	const headers = { 'User-Agent': 'Hibiki/0.1.0' };

	const [infoRes, streamRes] = await Promise.all([
		fetch(`${MONOCHROME_API}/info/?id=${id}`, { headers }),
		fetch(`${MONOCHROME_API}/track/?id=${id}&quality=${quality}`, { headers }),
	]);

	const info: unknown = infoRes.ok ? await infoRes.json() : null;
	const stream: unknown = streamRes.ok ? await streamRes.json() : null;

	return json({ info, stream });
};
