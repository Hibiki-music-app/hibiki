import {json} from '@sveltejs/kit';
import { AppName, EmailRequest, MBApiRoot } from '$lib/services/ApiEndpoints';

export async function GET({ url }) {
	const query = url.searchParams.get('q');
	if (!query) {
		return json({error: 'Missing query'}, {status: 400});
	}
	const response = await fetch(
		`${MBApiRoot}/recording?query=${encodeURIComponent(query)}&fmt=json`, // changer l'entité recording par d'autres pour appronfondir les recherches et aussi ajouter les types dans models ;)
		{
			headers: {
				'User-Agent': `${AppName} ( ${EmailRequest} )`
			}
		}
	);
	const data = await response.json();
	return json(data);
}