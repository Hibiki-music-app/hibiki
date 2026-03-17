import type { Track } from '$lib/models/Track';

export async function searchTracks(query: string): Promise<Track[]> {
	if (!query?.trim()) return [];

	try {
		const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
		if (!response.ok) return [];
		const data = await response.json();
		// L'API Monochrome retourne { tracks: [...] } ou similaire
		return (data.tracks ?? data) as Track[];
	} catch (error) {
		console.error('Erreur de recherche:', error);
		return [];
	}
}

export async function getTrackUrl(trackId: number): Promise<string> {
	try {
		const response = await fetch(`/api/track?id=${trackId}`);
		if (!response.ok) return '';
		const data = await response.json();
		return data.stream?.url ?? data.stream?.manifest ?? '';
	} catch (error) {
		console.error('Erreur de lecture:', error);
		return '';
	}
}
