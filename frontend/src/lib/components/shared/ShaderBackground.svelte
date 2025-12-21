<script lang="ts">
import { Canvas } from '@threlte/core';
import Scene from './Scene.svelte';

type Props = {
	vertexShader?: string;
	fragmentShader?: string;
	class?: string;
};

let {
	vertexShader,
	fragmentShader,
	class: className = ''
}: Props = $props();
let container: HTMLDivElement | null = null;
</script>

<div
	bind:this={container}
	class={`absolute w-full h-full ${className}`}
>
	<Canvas>	
		<Scene
			vertexShader={vertexShader}
			fragmentShader={fragmentShader}
		/>
	</Canvas>
</div>
