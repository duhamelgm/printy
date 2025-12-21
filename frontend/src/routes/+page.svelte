<script lang="ts">
import Topbar from '$lib/components/shared/Topbar.svelte';
import Button from '$lib/components/ui/button/button.svelte';
import ShaderBackground from '$lib/components/shared/ShaderBackground.svelte';
import fragmentShader from '$lib/shaders/hero-shader.glsl?raw';
</script>

<div
	class="relative flex h-full min-h-screen w-full flex-col bg-synthwave-gradient overflow-hidden"
>
	<Topbar />
	<div
		class="relative z-10 flex flex-col items-center justify-center pt-8 pb-12 px-6 text-center grow"
	>
		<ShaderBackground
			class="absolute inset-0 opacity-50"
			fragmentShader={fragmentShader}
		/>
    
		<div class="mb-8 relative z-10">
			<!-- Decorative Glow behind main text -->
			<div class="absolute -inset-4 bg-primary/20 rounded-full blur-3xl"></div>
			<h1
				class="relative text-white tracking-tight text-[36px] font-bold leading-tight drop-shadow-md"
			>
				Ticket generation <br /> <span class="text-cyan-400">initialized...</span>
			</h1>
			<p class="relative mt-4 text-slate-300 text-lg font-medium tracking-wide">
				Ready to process new requests.
			</p>
		</div>

		<Button href="/print" size="xl" startIcon="add_circle">Create Ticket</Button>
	</div>
</div>
