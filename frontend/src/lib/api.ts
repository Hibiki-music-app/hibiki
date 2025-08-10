const API_BASE = 'https://dab.yeet.su/api';

interface AudioQuality {
	maximumBitDepth: number;
	maximumSamplingRate: number;
	isHiRes: boolean;
}

interface Images {
	small: string;
	thumbnail: string;
	large: string;
	back?: string | null;
}

export interface Track {
	id: number;
	title: string;
	artist: string;
	artistId: number;
	albumTitle: string;
	albumCover: string;
	albumId: string;
	releaseDate: string;
	genre: string;
	duration: number;
	audioQuality: AudioQuality;
	version: string | null;
	label: string;
	labelId: number;
	upc: string;
	mediaCount: number;
	parental_warning: boolean;
	streamable: boolean;
	purchasable: boolean;
	previewable: boolean;
	genreId: number;
	genreSlug: string;
	genreColor: string;
	releaseDateStream: string;
	releaseDateDownload: string;
	maximumChannelCount: number;
	images: Images;
	isrc: string;
}

interface SearchResponse {
	tracks: Track[];
}

export async function searchTracks(query: string): Promise<Track[]> {
	if (!query) return [];

	try {
		const response = await fetch(
			`${API_BASE}/search?q=${encodeURIComponent(query)}&type=track&limit=10`
		);
		const data: SearchResponse = await response.json();

		return data.tracks.map((item) => ({
			id: item.id,
			title: item.title,
			artist: item.artist,
			artistId: item.artistId,
			albumTitle: item.albumTitle,
			albumCover: item.albumCover,
			albumId: item.albumId,
			releaseDate: item.releaseDate,
			genre: item.genre,
			duration: item.duration,
			audioQuality: item.audioQuality,
			version: item.version,
			label: item.label,
			labelId: item.labelId,
			upc: item.upc,
			mediaCount: item.mediaCount,
			parental_warning: item.parental_warning,
			streamable: item.streamable,
			purchasable: item.purchasable,
			previewable: item.previewable,
			genreId: item.genreId,
			genreSlug: item.genreSlug,
			genreColor: item.genreColor,
			releaseDateStream: item.releaseDateStream,
			releaseDateDownload: item.releaseDateDownload,
			maximumChannelCount: item.maximumChannelCount,
			images: item.images,
			isrc: item.isrc
		}));
	} catch (error) {
		console.error('Erreur de recherche:', error);
		return [];
	}
}

export async function getTrackUrl(trackId: number): Promise<string> {
	try {
		const response = await fetch(`${API_BASE}/stream?trackId=${trackId}&quality=27`);
		if (!response.ok) {
			throw new Error('Erreur lors de la lecture de la piste');
		}
		const data = await response.json();
		return data.url;
	} catch (error) {
		console.error('Erreur de lecture:', error);
		return '';
	}
}
