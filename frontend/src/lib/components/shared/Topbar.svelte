<script lang="ts">
	let {
		backHref = null,
    class: className = ''
	}: {
		backHref?: string | null;
		class?: string;
	} = $props();
</script>

<div class="z-10 flex items-center p-6 justify-between h-16 bg-background {className} w-full">
	{#if backHref}
		<div
			class="text-white flex size-12 shrink-0 items-center justify-center rounded-full active:bg-white/10 transition-colors cursor-pointer w-12"
		>
			<a href={backHref} class="material-symbols-outlined text-primary" style="font-size: 24px;">
				arrow_back_ios_new
			</a>
		</div>
	{:else}
		<div class="size-10"></div>
	{/if}

	<!-- Logo -->
	<div class="flex flex-col items-center">
		<h2
			class="text-transparent bg-clip-text bg-linear-to-r from-cyan-400 to-fuchsia-500 text-3xl font-bold tracking-widest drop-shadow-[0_0_10px_rgba(217,70,239,0.5)]"
		>
			PRINTY
		</h2>
	</div>

	<div class="size-10"></div>
</div>
