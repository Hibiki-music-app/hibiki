export interface RecordingResponse {
	created: string; // date ISO
	count: number;
	offset: number;
	recordings: Recording[];
}

export interface Recording {
	id: string;
	score: number;
	'artist-credit-id': string;
	title: string;
	length: number | null;
	video: string | null;
	'artist-credit'?: ArtistCredit[];
	releases?: Release[];
}

export interface ArtistCredit {
	name: string;
	artist: Artist;
}

export interface Artist {
	id: string;
	name: string;
	'sort-name': string;
}

export interface Release {
	id: string;
	'status-id': string;
	count: number;
	title: string;
	status: string;
	'release-group': ReleaseGroup;
	'track-count': number;
	media?: Media[];
}

export interface ReleaseGroup {
	id: string;
	'type-id': string;
	'title': string;
	'primary-type': string;
	'primary-type-id': string;
}

export interface Media {
	id: string;
	position: number;
	track?: Track[];
	'track-count': number;
	'track-offset': number;
}

export interface Track {
	id: string;
	number: string;
	title: string;
	length: number | null;
}

export type SearchEntityType =
	| 'area'
	| 'artist'
	| 'event'
	| 'genre'
	| 'instrument'
	| 'label'
	| 'place'
	| 'recording'
	| 'release'
	| 'release-group'
	| 'series'
	| 'work'
	| 'url';