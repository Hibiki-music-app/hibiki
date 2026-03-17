import {
	type Track,
	type TidalTrackRaw,
	type TidalStreamRaw,
	mapTidalTrack,
	decodeManifestUrl,
} from '$lib/models/Track';

interface MonochromeSearchResponse {
	data: {
		items: TidalTrackRaw[];
		totalNumberOfItems: number;
		limit: number;
		offset: number;
	};
}

interface ProxyTrackResponse {
	stream: {
		data: TidalStreamRaw;
	} | null;
}

export async function searchTracks(query: string): Promise<Track[]> {
	if (!query?.trim()) return [];

	try {
		const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
		if (!response.ok) return [];

		const json = (await response.json()) as MonochromeSearchResponse;
		const items = json?.data?.items;

		if (!Array.isArray(items)) return [];
		return items.map(mapTidalTrack);
	} catch (error) {
		console.error('Erreur de recherche:', error);
		return [];
	}
}

export async function getTrackUrl(trackId: number): Promise<string> {
	try {
		const response = await fetch(`/api/track?id=${trackId}`);
		if (!response.ok) return '';

		const json = (await response.json()) as ProxyTrackResponse;
		const manifest = json?.stream?.data?.manifest;

		if (manifest) {
			return decodeManifestUrl(manifest);
		}
		return '';
	} catch (error) {
		console.error('Erreur lecture piste:', error);
		return '';
	}
}
