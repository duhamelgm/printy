<script lang="ts">
import ShaderBackground from '$lib/components/shared/ShaderBackground.svelte';
import { onMount } from 'svelte';
import fragmentShaderManualEntry from '$lib/shaders/manual_entry_shader.glsl?raw';
import Topbar from '$lib/components/shared/Topbar.svelte';
	import Button from '$lib/components/ui/button/button.svelte';

onMount(() => {
	document.body.classList.remove('overflow-y-hidden');
});
</script>

<div class="relative h-full min-h-screen w-full bg-synthwave-gradient">
	<Topbar backHref="/" class="fixed" />

	<div class="h-full overflow-y-auto py-8">
		<!-- Main Content -->
		<div
			class="flex-1 flex flex-col items-stretch pb-8 px-4 gap-6 max-w-md mx-auto w-full overflow-y-auto pt-24"
		>
			<!-- Headline Section -->
			<div class="flex flex-col items-center justify-center text-center">
				<h1
					class="text-white tracking-tight text-[32px] font-bold leading-tight drop-shadow-[0_0_8px_rgba(13,204,242,0.5)]"
				>
					Choose Method
				</h1>
				<p class="text-gray-300 text-base font-normal leading-normal pt-2 px-4">
					Select how you would like to initiate your new ticket.
				</p>
			</div>
			<!-- Option 1: Manual Entry -->
			<div class="w-full @container group cursor-pointer">
				<div
					class="flex flex-col items-stretch justify-start rounded-xl shadow-[0_0_15px_rgba(13,204,242,0.15)] bg-card-dark border border-primary/40 hover:border-primary transition-all duration-300 transform hover:-translate-y-1"
				>
					<div
						class="w-full h-40 bg-center bg-no-repeat bg-cover rounded-t-xl relative overflow-hidden"
					>
						<ShaderBackground
							fragmentShader={fragmentShaderManualEntry}
						/>
					</div>
					<div class="flex w-full grow flex-col items-stretch justify-center gap-1 p-5">
						<div class="flex justify-between items-center mb-1">
							<p
								class="text-white text-xl font-bold leading-tight tracking-wide uppercase drop-shadow-md"
							>
								Manual Entry
							</p>
						</div>
						<p class="text-[#90c1cb] text-sm font-normal leading-relaxed mb-4">
							Fill out the ticket details manually using our standard input form.
						</p>
						<Button href="/print/form" size="lg" endIcon="arrow_forward">Create Manually</Button>
					</div>
				</div>
			</div>
			<!-- Option 2: Scan from Photo -->
			<div class="w-full @container group cursor-pointer">
				<div
					class="flex flex-col items-stretch justify-start rounded-xl shadow-[0_0_15px_rgba(13,204,242,0.15)] bg-card-dark border border-primary/40 hover:border-primary transition-all duration-300 transform hover:-translate-y-1"
				>
					<div
						class="w-full h-40 bg-center bg-no-repeat bg-cover rounded-t-xl relative overflow-hidden"
						data-alt="Futuristic neon camera lens aperture geometric pattern"
						style="background-image: url('https://lh3.googleusercontent.com/aida-public/AB6AXuASVpA1Mwq1CXclGfl-_N9TKAo6u5kn7crpEJ_MPDKGvhgjsECGBTTuxAwlkNEYNVMFOzqH6I-Dp48IbKL-HTYcmTTOOxcZK12SJrRV2Yl1xyTAREZ6y7Ye73VUJc_Fy9dgQpxs9qQRZGwi6eCplQqBNzFoXh8crWDdKZn74UxhgZ5UgjoaLXcLbCAPEBFo8L46rtKl0qe2xAmR2SR-_Kz51fRQ60QG197RnTBzE6r1bjpQ2h72Vs4x4Mz38LhBiIXEQcBt1a7PiwU');"
					>
					</div>
					<div class="flex w-full grow flex-col items-stretch justify-center gap-1 p-5">
						<div class="flex justify-between items-center mb-1">
							<p
								class="text-white text-xl font-bold leading-tight tracking-wide uppercase drop-shadow-md"
							>
								Scan from Photo
							</p>
						</div>
						<p class="text-[#90c1cb] text-sm font-normal leading-relaxed mb-4">
							Upload or take a picture to auto-generate details via OCR.
						</p>
						<button
							class="flex w-full cursor-pointer items-center justify-center overflow-hidden rounded-lg h-10 px-4 bg-transparent border border-primary text-primary text-sm font-bold leading-normal hover:bg-primary/10 transition-colors"
						>
							<span class="truncate">Open Camera</span>
							<span class="material-symbols-outlined ml-2 text-sm">center_focus_strong</span>
						</button>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
