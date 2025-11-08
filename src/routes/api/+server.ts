import {json} from '@sveltejs/kit';

export async function GET({ url }) {
	const query = url.searchParams.get('q');
	if (!query) {
		return json({error: 'Missing query'}, {status: 400});
	}
	const response = await fetch(
		`https://musicbrainz.org/ws/2/recording?query=${encodeURIComponent(query)}&fmt=json`,
		{
			headers: {
				'User-Agent': 'Hibiki/0.0.1 ( oscarodark3@gmail.com )'
			}
		}
	);
	const data = await response.json();
	return json(data);
}