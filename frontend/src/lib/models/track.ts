import type { AudioQuality } from '$lib/models/audio';
import type { Images } from '$lib/models/image';

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