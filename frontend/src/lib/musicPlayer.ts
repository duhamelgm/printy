// Music player that randomly selects and plays tracks from assets/soundtracks
import track1 from '$lib/assets/soundtracks/1990.mp3';
import track2 from '$lib/assets/soundtracks/Echoes Of The Past.mp3';
import track3 from '$lib/assets/soundtracks/In My Garage Pt.2.mp3';
import track4 from '$lib/assets/soundtracks/Neon Rose Garden.mp3';
import track5 from '$lib/assets/soundtracks/Northern Lights.mp3';
import map1 from '$lib/assets/audio-maps/1990.png';

type AudioEntry = {
	track: string;
	map?: { src: string; fps: number };
};

// Attach pre-baked frequency map when available (currently only 1990.mp3)
const audioFiles: AudioEntry[] = [
	{ track: track1 },
	{ track: track2 },
	{ track: track3 },
	{ track: track4 },
	{ track: track5 }
];

let currentAudio: HTMLAudioElement | null = null;

export function stopAudio(): void {
	if (currentAudio) {
		currentAudio.pause();
		currentAudio = null;
	}
}

export function playAudio() {
	// Stop current track if playing
	if (currentAudio) {
		currentAudio.pause();
		currentAudio = null;
	}

	// Randomly select a track
	const randomIndex = Math.floor(Math.random() * audioFiles.length);
	const selected = audioFiles[randomIndex];

	// Create new audio element
	const audio = new Audio(selected.track);
	audio.volume = 0.2;
	audio.loop = true;

	audio.play().catch((error) => {
		console.error('Error playing audio:', error);
	});

	currentAudio = audio;

	return new Promise<{
		audio: HTMLAudioElement;
		analyser: AnalyserNode | null;
		source: MediaElementAudioSourceNode | null;
		frequencyData: Float32Array | null;
		bakedMap?: { src: string; fps: number };
	}>((resolve) => {
		audio.addEventListener('loadedmetadata', () => {
			const duration = audio.duration;

			// Calculate random start time:
			// - After 30 seconds
			// - Before 1 minute (60 seconds) from the end
			const minStartTime = 30; // 30 seconds
			const maxStartTime = Math.max(minStartTime, duration - 60); // 1 minute before end

			// If track is shorter than 90 seconds, just start at 30 seconds
			const startTime =
				duration < 90 ? minStartTime : Math.random() * (maxStartTime - minStartTime) + minStartTime;

			audio.currentTime = startTime;

			// eslint-disable-next-line @typescript-eslint/no-explicit-any
			const audioContext = new (window.AudioContext || (window as any).webkitAudioContext)();
			const analyser = audioContext.createAnalyser();
			analyser.fftSize = 256;
			analyser.smoothingTimeConstant = 0.8;
			const source = audioContext.createMediaElementSource(audio);

			source.connect(analyser);
			analyser.connect(audioContext.destination);

			const frequencyData = new Float32Array(analyser.frequencyBinCount);

			resolve({
				audio,
				analyser,
				source,
				frequencyData,
				bakedMap: selected.map
			});
		});
	});
}
