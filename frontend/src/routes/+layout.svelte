<script lang="ts">
import '../app.css';

import favicon from '$lib/assets/favicon.svg';
import { playAudio, stopAudio } from '$lib/musicPlayer';
import { onDestroy, onMount } from 'svelte';
import { setContext } from 'svelte';
import { writable } from 'svelte/store';

let audio = writable<{
	audio: HTMLAudioElement;
	analyser: AnalyserNode | null;
	source: MediaElementAudioSourceNode | null;
	frequencyData: Float32Array | null;
} | null>(null);
let initialized = $state(false);
let { children } = $props();

onMount(() => {
	setContext('audio', audio);
});

onDestroy(() => {
	stopAudio();
	audio.set(null);
});

function onInitialize() {
	initialized = true;
	playAudio().then((audioData) => {
		audio.set(audioData);
	});
}
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<link
		href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&amp;display=swap"
		rel="stylesheet"
	/>
	<link
		href="https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400;500;700&amp;family=Noto+Sans:wght@400;500;700&amp;display=swap"
		rel="stylesheet"
	/>
	<link
		href="https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:wght,FILL@100..700,0..1&amp;display=swap"
		rel="stylesheet"
	/>
</svelte:head>

{#if initialized}
	{@render children()}
{:else}
	<div class="flex items-center justify-center h-full">
		<button class="bg-primary text-white px-4 py-2 rounded-md" onclick={onInitialize}
			>Initialize</button
		>
	</div>
{/if}
