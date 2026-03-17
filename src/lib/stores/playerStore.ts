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
}

// Génère un UUID simple
function generateUUID(): string {
	return 'xxxx-xxxx-xxxx'.replace(/[x]/g, () => {
		return (Math.random() * 16 | 0).toString(16);
	});
}

// Store principal de la queue
const createQueueStore = () => {
	const { subscribe, set, update } = writable<QueueStore>({
		items: [],
		currentIndex: -1,
		history: [],
		isShuffling: false,
		isLooping: false
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
		addTrack: (track: Track) => update(state => {
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
		addTracks: (tracks: Track[]) => update(state => {
			const queueItems: QueueItem[] = tracks.map(track => ({
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
		playTrack: (track: Track) => update(state => {
			// Chercher si la piste est déjà dans la queue
			const existingIndex = state.items.findIndex(item => item.track.id === track.id);
			
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
		playNext: (track: Track) => update(state => {
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
		removeTrack: (itemId: string) => update(state => {
			const itemIndex = state.items.findIndex(item => item.id === itemId);
			if (itemIndex === -1) return state;

			const newItems = state.items.filter(item => item.id !== itemId);
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
		next: () => update(state => {
			if (state.items.length === 0) return state;

			let nextIndex: number;
			
			if (state.isShuffling) {
				// Mode aléatoire
				const availableIndices = state.items
					.map((_, index) => index)
					.filter(index => index !== state.currentIndex);
				
				if (availableIndices.length === 0) {
					nextIndex = state.isLooping ? 0 : state.currentIndex;
				} else {
					nextIndex = availableIndices[Math.floor(Math.random() * availableIndices.length)];
				}
			} else {
				// Mode séquentiel
				nextIndex = state.currentIndex + 1;
				
				if (nextIndex >= state.items.length) {
					nextIndex = state.isLooping ? 0 : state.currentIndex;
				}
			}

			const newState = {
				...state,
				currentIndex: nextIndex,
				history: state.currentIndex !== -1 && state.items[state.currentIndex] 
					? [...state.history, state.items[state.currentIndex]].slice(-10) // Garder les 10 dernières
					: state.history
			};
			saveToStorage(newState);
			return newState;
		}),

		// Aller à la piste précédente
		previous: () => update(state => {
			if (state.items.length === 0) return state;

			let prevIndex: number;

			if (state.history.length > 0 && !state.isShuffling) {
				// Utiliser l'historique si disponible
				const lastItem = state.history[state.history.length - 1];
				const foundIndex = state.items.findIndex(item => item.id === lastItem.id);
				
				if (foundIndex !== -1) {
					prevIndex = foundIndex;
				} else {
					prevIndex = Math.max(0, state.currentIndex - 1);
				}
			} else {
				// Mode séquentiel simple
				prevIndex = state.currentIndex - 1;
				
				if (prevIndex < 0) {
					prevIndex = state.isLooping ? state.items.length - 1 : 0;
				}
			}

			const newState = {
				...state,
				currentIndex: prevIndex,
				history: state.history.slice(0, -1) // Retirer le dernier élément de l'historique
			};
			saveToStorage(newState);
			return newState;
		}),

		// Réorganiser la queue
		reorderQueue: (startIndex: number, endIndex: number) => update(state => {
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
		toggleShuffle: () => update(state => {
			const newState = {
				...state,
				isShuffling: !state.isShuffling
			};
			saveToStorage(newState);
			return newState;
		}),

		// Basculer le mode répétition
		toggleLoop: () => update(state => {
			const newState = {
				...state,
				isLooping: !state.isLooping
			};
			saveToStorage(newState);
			return newState;
		}),

		// Vider la queue
		clear: () => update(state => {
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
		setCurrentIndex: (index: number) => update(state => {
			if (index < 0 || index >= state.items.length) return state;
			
			const newState = {
				...state,
				currentIndex: index,
				history: state.currentIndex !== -1 && state.items[state.currentIndex]
					? [...state.history, state.items[state.currentIndex]].slice(-10)
					: state.history
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
export const currentTrack = derived(queueStore, $queue => 
	$queue.currentIndex >= 0 && $queue.items[$queue.currentIndex] 
		? $queue.items[$queue.currentIndex].track 
		: null
);

export const queueLength = derived(queueStore, $queue => $queue.items.length);

export const hasNext = derived(queueStore, $queue => 
	$queue.currentIndex < $queue.items.length - 1 || $queue.isLooping
);

export const hasPrevious = derived(queueStore, $queue => 
	$queue.currentIndex > 0 || $queue.isLooping || $queue.history.length > 0
);