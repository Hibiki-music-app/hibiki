import { writable, derived } from 'svelte/store';
import type { Track } from '$lib/models/Track';

export interface QueueItem {
	track: Track;
	id: string; // UUID pour identifier chaque élément de manière unique
}

export interface QueueStore {
	items: QueueItem[];
	currentIndex: number;
	history: QueueItem[];
	isShuffling: boolean;
	isLooping: boolean; // false = no loop, true = loop queue
	volume: number; // 0.0 à 1.0
}

// Génère un UUID simple
function generateUUID(): string {
	return 'xxxx-xxxx-xxxx'.replace(/[x]/g, () => {
		return Math.floor(Math.random() * 16).toString(16);
	});
}

// Store principal de la queue
const createQueueStore = () => {
	const { subscribe, set, update } = writable<QueueStore>({
		items: [],
		currentIndex: -1,
		history: [],
		isShuffling: false,
		isLooping: false,
		volume: 0.2
	});

	// Charger depuis le localStorage au démarrage
	if (typeof localStorage !== 'undefined') {
		const saved = localStorage.getItem('hibiki-queue');
		if (saved) {
			try {
				const parsed = JSON.parse(saved);
				set(parsed);
			} catch (e) {
				console.warn('Erreur lors du chargement de la queue:', e);
			}
		}
	}

	return {
		subscribe,

		// Ajouter une piste à la fin de la queue
		addTrack: (track: Track) =>
			update((state) => {
				const queueItem: QueueItem = {
					track,
					id: generateUUID()
				};
				const newState = {
					...state,
					items: [...state.items, queueItem]
				};
				saveToStorage(newState);
				return newState;
			}),

		// Ajouter plusieurs pistes
		addTracks: (tracks: Track[]) =>
			update((state) => {
				const queueItems: QueueItem[] = tracks.map((track) => ({
					track,
					id: generateUUID()
				}));
				const newState = {
					...state,
					items: [...state.items, ...queueItems]
				};
				saveToStorage(newState);
				return newState;
			}),

		// Jouer une piste spécifique et l'ajouter à la queue si elle n'y est pas
		playTrack: (track: Track) =>
			update((state) => {
				// Chercher si la piste est déjà dans la queue
				const existingIndex = state.items.findIndex((item) => item.track.id === track.id);

				if (existingIndex !== -1) {
					// La piste existe, on la sélectionne
					const newState = {
						...state,
						currentIndex: existingIndex
					};
					saveToStorage(newState);
					return newState;
				} else {
					// La piste n'existe pas, on l'ajoute et on la sélectionne
					const queueItem: QueueItem = {
						track,
						id: generateUUID()
					};
					const newState = {
						...state,
						items: [...state.items, queueItem],
						currentIndex: state.items.length // Index de la nouvelle piste
					};
					saveToStorage(newState);
					return newState;
				}
			}),

		// Jouer immédiatement (ajouter au début de la queue)
		playNext: (track: Track) =>
			update((state) => {
				const queueItem: QueueItem = {
					track,
					id: generateUUID()
				};
				const newItems = [...state.items];
				const insertIndex = state.currentIndex + 1;
				newItems.splice(insertIndex, 0, queueItem);

				const newState = {
					...state,
					items: newItems
				};
				saveToStorage(newState);
				return newState;
			}),

		// Supprimer une piste de la queue
		removeTrack: (itemId: string) =>
			update((state) => {
				const itemIndex = state.items.findIndex((item) => item.id === itemId);
				if (itemIndex === -1) return state;

				const newItems = state.items.filter((item) => item.id !== itemId);
				let newCurrentIndex = state.currentIndex;

				if (itemIndex < state.currentIndex) {
					newCurrentIndex--;
				} else if (itemIndex === state.currentIndex) {
					// Si on supprime la piste actuelle, rester sur le même index (qui pointera vers la suivante)
					if (newCurrentIndex >= newItems.length) {
						newCurrentIndex = newItems.length - 1;
					}
				}

				const newState = {
					...state,
					items: newItems,
					currentIndex: newCurrentIndex
				};
				saveToStorage(newState);
				return newState;
			}),

		// Aller à la piste suivante
		next: () =>
			update((state) => {
				if (state.items.length === 0) return state;

				let nextIndex: number;
				let newItems = [...state.items];
				let newHistory = [...state.history];

				// Ajouter la piste actuelle à l'historique si elle existe
				if (state.currentIndex !== -1 && state.items[state.currentIndex]) {
					newHistory = [...newHistory, state.items[state.currentIndex]].slice(-10); // Garder les 10 dernières
					// Retirer la piste actuelle de la queue
					newItems.splice(state.currentIndex, 1);
				}

				// Recalculer les indices après suppression
				if (state.isShuffling) {
					// Mode aléatoire
					if (newItems.length === 0) {
						nextIndex = -1;
					} else {
						nextIndex = Math.floor(Math.random() * newItems.length);
					}
				} else {
					// Mode séquentiel
					if (newItems.length === 0) {
						nextIndex = -1;
					} else {
						// Rester au même index car l'élément suivant a pris la place de l'actuel
						nextIndex = Math.min(state.currentIndex, newItems.length - 1);
					}
				}

				// Si on est en mode loop et qu'il n'y a plus d'éléments, remettre l'historique dans la queue
				if (newItems.length === 0 && state.isLooping && newHistory.length > 0) {
					newItems = [...newHistory];
					newHistory = [];
					nextIndex = 0;
				}

				const newState = {
					...state,
					items: newItems,
					currentIndex: nextIndex,
					history: newHistory
				};
				saveToStorage(newState);
				return newState;
			}),

		// Aller à la piste précédente
		previous: () =>
			update((state) => {
				if (state.history.length === 0) return state;

				// Prendre la dernière piste de l'historique
				const lastHistoryItem = state.history[state.history.length - 1];
				const newHistory = state.history.slice(0, -1);
				let newItems = [...state.items];

				// Remettre la piste de l'historique au début de la queue
				newItems.unshift(lastHistoryItem);

				const newState = {
					...state,
					items: newItems,
					currentIndex: 0, // La piste remise de l'historique devient la piste actuelle
					history: newHistory
				};
				saveToStorage(newState);
				return newState;
			}),

		// Réorganiser la queue
		reorderQueue: (startIndex: number, endIndex: number) =>
			update((state) => {
				const newItems = [...state.items];
				const [removed] = newItems.splice(startIndex, 1);
				newItems.splice(endIndex, 0, removed);

				// Ajuster l'index actuel
				let newCurrentIndex = state.currentIndex;
				if (startIndex === state.currentIndex) {
					newCurrentIndex = endIndex;
				} else if (startIndex < state.currentIndex && endIndex >= state.currentIndex) {
					newCurrentIndex--;
				} else if (startIndex > state.currentIndex && endIndex <= state.currentIndex) {
					newCurrentIndex++;
				}

				const newState = {
					...state,
					items: newItems,
					currentIndex: newCurrentIndex
				};
				saveToStorage(newState);
				return newState;
			}),

		// Basculer le mode aléatoire
		toggleShuffle: () =>
			update((state) => {
				const newState = {
					...state,
					isShuffling: !state.isShuffling
				};
				saveToStorage(newState);
				return newState;
			}),

		// Basculer le mode répétition
		toggleLoop: () =>
			update((state) => {
				const newState = {
					...state,
					isLooping: !state.isLooping
				};
				saveToStorage(newState);
				return newState;
			}),

		// Vider la queue
		clear: () =>
			update((state) => {
				const newState = {
					...state,
					items: [],
					currentIndex: -1,
					history: []
				};
				saveToStorage(newState);
				return newState;
			}),

		// Aller à un index spécifique
		setCurrentIndex: (index: number) =>
			update((state) => {
				if (index < 0 || index >= state.items.length) return state;

				// Ne pas ajouter à l'historique si c'est le même index
				if (index === state.currentIndex) return state;

				let newHistory = [...state.history];

				// Ajouter la piste actuelle à l'historique si elle existe et si on change vraiment de piste
				if (state.currentIndex !== -1 && state.items[state.currentIndex]) {
					newHistory = [...newHistory, state.items[state.currentIndex]].slice(-10);
				}

				const newState = {
					...state,
					currentIndex: index,
					history: newHistory
				};
				saveToStorage(newState);
				return newState;
			}),

		// Changer le volume (0.0 à 1.0)
		setVolume: (volume: number) =>
			update((state) => {
				const clampedVolume = Math.max(0, Math.min(1, volume));
				const newState = {
					...state,
					volume: clampedVolume
				};
				saveToStorage(newState);
				return newState;
			}),

		// Augmenter le volume
		increaseVolume: (step: number = 0.1) =>
			update((state) => {
				const newVolume = Math.min(1, state.volume + step);
				const newState = {
					...state,
					volume: newVolume
				};
				saveToStorage(newState);
				return newState;
			}),

		// Diminuer le volume
		decreaseVolume: (step: number = 0.1) =>
			update((state) => {
				const newVolume = Math.max(0, state.volume - step);
				const newState = {
					...state,
					volume: newVolume
				};
				saveToStorage(newState);
				return newState;
			}),

		// Basculer muet/non-muet
		toggleMute: () =>
			update((state) => {
				const newVolume = state.volume > 0 ? 0 : 1;
				const newState = {
					...state,
					volume: newVolume
				};
				saveToStorage(newState);
				return newState;
			})
	};
};

// Fonction pour sauvegarder dans le localStorage
function saveToStorage(state: QueueStore) {
	if (typeof localStorage !== 'undefined') {
		try {
			localStorage.setItem('hibiki-queue', JSON.stringify(state));
		} catch (e) {
			console.warn('Erreur lors de la sauvegarde de la queue:', e);
		}
	}
}

export const queueStore = createQueueStore();

// Stores dérivés pour faciliter l'utilisation
export const currentTrack = derived(queueStore, ($queue) =>
	$queue.currentIndex >= 0 && $queue.items[$queue.currentIndex]
		? $queue.items[$queue.currentIndex].track
		: null
);

export const queueLength = derived(queueStore, ($queue) => $queue.items.length);

export const hasNext = derived(
	queueStore,
	($queue) => $queue.currentIndex < $queue.items.length - 1 || $queue.isLooping
);

export const hasPrevious = derived(
	queueStore,
	($queue) => $queue.currentIndex > 0 || $queue.isLooping || $queue.history.length > 0
);

export const volume = derived(queueStore, ($queue) => $queue.volume);

export const isMuted = derived(queueStore, ($queue) => $queue.volume === 0);
