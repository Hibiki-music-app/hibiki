export interface Track {
	id: number;
	title: string;
	artist: string;
	artistId: number;
	albumTitle: string;
	albumCover: string;
	albumId: number;
	duration: number; // secondes
	// champs optionnels — confirmés après appel API réel
	releaseDate?: string;
	genre?: string;
	version?: string | null;
	isrc?: string;
	label?: string;
	labelId?: number;
	upc?: string;
	mediaCount?: number;
	parental_warning?: boolean;
	streamable?: boolean;
	purchasable?: boolean;
	previewable?: boolean;
	genreId?: number;
	genreSlug?: string;
	genreColor?: string;
	releaseDateStream?: string;
	releaseDateDownload?: string;
	maximumChannelCount?: number;
	audioQuality?: string;
	images?: {
		small?: string;
		thumbnail?: string;
		large?: string;
		back?: string;
	};
}

export interface TrackStream {
	url: string;
	quality: string;
	codec?: string;
	bitDepth?: number;
	sampleRate?: number;
	encryptionKey?: string;
}
