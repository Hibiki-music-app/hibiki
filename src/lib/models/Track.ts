/** Réponse brute d'un item renvoyé par /search/ de l'API Monochrome */
export interface TidalTrackRaw {
	id: number;
	title: string;
	duration: number;
	version: string | null;
	isrc: string;
	audioQuality: string;
	explicit: boolean;
	popularity: number;
	artist: { id: number; name: string; picture: string | null };
	artists: Array<{ id: number; name: string }>;
	album: { id: number; title: string; cover: string; vibrantColor?: string };
	trackNumber?: number;
	volumeNumber?: number;
}

/** Réponse brute du endpoint /track/ (manifest audio) */
export interface TidalManifest {
	mimeType: string;
	codecs: string;
	encryptionType: string;
	urls: string[];
}

export interface TidalStreamRaw {
	trackId: number;
	audioQuality: string;
	audioMode: string;
	bitDepth?: number;
	sampleRate?: number;
	manifestMimeType: string;
	manifest: string; // base64 JSON contenant { urls: [...] }
}

/** Modèle Track utilisé dans l'app */
export interface Track {
	id: number;
	title: string;
	artist: string;
	artistId: number;
	albumTitle: string;
	/** URL complète de la pochette (déjà convertie depuis le cover UUID Tidal) */
	albumCover: string;
	albumId: number;
	duration: number;
	version: string | null;
	isrc?: string;
	audioQuality?: string;
	explicit?: boolean;
}

/** Convertit un UUID de cover Tidal en URL d'image */
export function tidalCoverUrl(coverUuid: string, size: 320 | 640 | 1280 = 320): string {
	return `https://resources.tidal.com/images/${coverUuid.replace(/-/g, '/')}/${size}x${size}.jpg`;
}

/** Mappe un item brut de l'API Monochrome vers notre modèle Track */
export function mapTidalTrack(raw: TidalTrackRaw): Track {
	return {
		id: raw.id,
		title: raw.title,
		artist: raw.artist.name,
		artistId: raw.artist.id,
		albumTitle: raw.album.title,
		albumCover: raw.album.cover ? tidalCoverUrl(raw.album.cover) : '',
		albumId: raw.album.id,
		duration: raw.duration,
		version: raw.version,
		isrc: raw.isrc,
		audioQuality: raw.audioQuality,
		explicit: raw.explicit,
	};
}

/** Décode le manifest base64 Tidal et retourne la première URL de stream */
export function decodeManifestUrl(base64Manifest: string): string {
	try {
		const json = JSON.parse(atob(base64Manifest)) as TidalManifest;
		return json.urls[0] ?? '';
	} catch {
		return '';
	}
}
