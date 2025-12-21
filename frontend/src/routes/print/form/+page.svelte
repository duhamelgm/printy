<script>
	import Topbar from "$lib/components/shared/Topbar.svelte";

</script>
<div
	class="relative flex h-full min-h-screen w-full flex-col bg-background-light dark:bg-background-dark font-display text-white overflow-x-hidden selection:bg-synth-pink selection:text-white"
>
<Topbar backHref="/print" />
<div class="h-full overflow-y-auto"></div>

	<div class="flex-1 flex flex-col relative z-10 p-4 space-y-6 pb-32">
		<!-- Welcome/Hero Text -->
		<div class="pt-2 pb-4">
			<h1
				class="text-3xl font-black text-transparent bg-clip-text bg-gradient-to-r from-white via-blue-100 to-white tracking-tight drop-shadow-[0_0_10px_rgba(19,91,236,0.5)]"
			>
				New Ticket
			</h1>
			<p class="text-blue-200/60 text-sm mt-1 font-medium tracking-wide">
				ENTER MISSION PARAMETERS
			</p>
		</div>
		<!-- Title Input -->
		<div class="space-y-2 group">
			<label
				class="flex items-center gap-2 text-primary text-xs font-bold uppercase tracking-widest pl-1"
			>
				<span class="material-symbols-outlined text-sm">edit</span>
				Subject Line
			</label>
			<div class="relative">
				<input
					class="w-full bg-surface-dark/50 backdrop-blur-sm border border-primary/30 rounded-xl text-white px-4 py-4 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 focus:shadow-neon-blue transition-all duration-300 placeholder:text-slate-600"
					placeholder="System Malfunction..."
					type="text"
				/>
				<div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none">
					<div
						class="h-1.5 w-1.5 rounded-full bg-primary animate-pulse shadow-[0_0_8px_#135bec]"
					></div>
				</div>
			</div>
		</div>
		<!-- Description Input -->
		<div class="space-y-2">
			<label
				class="flex items-center gap-2 text-primary text-xs font-bold uppercase tracking-widest pl-1"
			>
				<span class="material-symbols-outlined text-sm">description</span>
				Mission Details
			</label>
			<textarea
				class="w-full min-h-32 bg-surface-dark/50 backdrop-blur-sm border border-primary/30 rounded-xl text-white px-4 py-4 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 focus:shadow-neon-blue transition-all duration-300 placeholder:text-slate-600 resize-none leading-relaxed"
				placeholder="Describe the anomaly encountered in sector 7..."
			></textarea>
		</div>
		<!-- Assignee Select -->
		<div class="space-y-2">
			<label
				class="flex items-center gap-2 text-primary text-xs font-bold uppercase tracking-widest pl-1"
			>
				<span class="material-symbols-outlined text-sm">person_add</span>
				Operative
			</label>
			<div class="relative">
				<select
					class="w-full appearance-none bg-surface-dark/50 backdrop-blur-sm border border-primary/30 rounded-xl text-white px-4 py-4 pr-12 focus:outline-none focus:border-primary focus:ring-1 focus:ring-primary/50 focus:shadow-neon-blue transition-all duration-300 text-base"
					style="background-image: none;"
				>
					<option class="bg-[#050510] text-gray-500" value="">Select Unit</option>
					<option class="bg-[#050510]" value="cyber-1">Unit Alpha (Dev)</option>
					<option class="bg-[#050510]" value="cyber-2">Unit Beta (Ops)</option>
					<option class="bg-[#050510]" value="cyber-3">Commander Zero</option>
				</select>
				<div class="absolute right-4 top-1/2 -translate-y-1/2 pointer-events-none text-primary">
					<span class="material-symbols-outlined">expand_more</span>
				</div>
			</div>
		</div>
		<!-- Priority Level -->
		<div class="space-y-3 pt-2">
			<label
				class="flex items-center gap-2 text-synth-pink text-xs font-bold uppercase tracking-widest pl-1 drop-shadow-[0_0_5px_rgba(255,0,255,0.6)]"
			>
				<span class="material-symbols-outlined text-sm">warning</span>
				Urgency Level
			</label>
			<div class="grid grid-cols-3 gap-3">
				<!-- Low -->
				<label class="cursor-pointer group relative">
					<input class="peer sr-only" name="priority" type="radio" value="low" />
					<div
						class="h-12 w-full flex items-center justify-center rounded-lg border border-teal-500/30 bg-teal-500/5 hover:bg-teal-500/10 transition-all duration-300 peer-checked:border-teal-400 peer-checked:bg-teal-500/20 peer-checked:shadow-[0_0_15px_rgba(45,212,191,0.4)]"
					>
						<span class="text-teal-400 font-bold text-sm tracking-wider group-hover:text-teal-300"
							>LOW</span
						>
					</div>
				</label>
				<!-- Medium -->
				<label class="cursor-pointer group relative">
					<input checked="" class="peer sr-only" name="priority" type="radio" value="medium" />
					<div
						class="h-12 w-full flex items-center justify-center rounded-lg border border-primary/30 bg-primary/5 hover:bg-primary/10 transition-all duration-300 peer-checked:border-primary peer-checked:bg-primary/20 peer-checked:shadow-neon-blue"
					>
						<span class="text-blue-400 font-bold text-sm tracking-wider group-hover:text-blue-300"
							>MED</span
						>
					</div>
				</label>
				<!-- High -->
				<label class="cursor-pointer group relative">
					<input class="peer sr-only" name="priority" type="radio" value="high" />
					<div
						class="h-12 w-full flex items-center justify-center rounded-lg border border-synth-pink/30 bg-synth-pink/5 hover:bg-synth-pink/10 transition-all duration-300 peer-checked:border-synth-pink peer-checked:bg-synth-pink/20 peer-checked:shadow-neon-pink"
					>
						<span class="text-pink-400 font-bold text-sm tracking-wider group-hover:text-pink-300"
							>CRIT</span
						>
					</div>
				</label>
			</div>
		</div>
	</div>
	<!-- Floating Action Button Area -->
	<div
		class="fixed bottom-0 left-0 w-full p-4 bg-gradient-to-t from-[#050510] via-[#050510]/95 to-transparent z-40 pb-8"
	>
		<button
			class="relative w-full group overflow-hidden rounded-xl bg-gradient-to-r from-primary via-purple-600 to-synth-pink p-[1px] shadow-[0_0_20px_rgba(19,91,236,0.4)] hover:shadow-[0_0_30px_rgba(188,19,254,0.6)] transition-all duration-300 transform hover:-translate-y-1"
		>
			<div class="absolute inset-0 bg-white/20 group-hover:bg-transparent transition-colors"></div>
			<div
				class="relative flex items-center justify-center gap-3 bg-[#0a0c18] hover:bg-opacity-90 w-full h-14 rounded-[11px] transition-all"
			>
				<span
					class="material-symbols-outlined text-white group-hover:text-pink-200 transition-colors"
					>rocket_launch</span
				>
				<span
					class="text-white font-black text-lg tracking-widest uppercase group-hover:text-transparent group-hover:bg-clip-text group-hover:bg-gradient-to-r group-hover:from-white group-hover:to-pink-200 transition-all"
					>Transmit Ticket</span
				>
			</div>
		</button>
	</div>
</div>
