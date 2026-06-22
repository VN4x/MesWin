<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { pb } from '$lib/pb';
	import { windowTypeLabel } from '$lib/window-types';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let mode = $state<'inner' | 'outer'>('inner');
	let uploading = $state(false);
	let message = $state('');

	async function capturePhoto() {
		uploading = true;
		message = '';

		try {
			const input = document.createElement('input');
			input.type = 'file';
			input.accept = 'image/*';
			input.capture = 'environment';

			const file = await new Promise<File | null>((resolve) => {
				input.onchange = () => resolve(input.files?.[0] ?? null);
				input.click();
			});

			if (!file) return;

			const formData = new FormData();
			const field = mode === 'inner' ? 'inner_photos' : 'outer_photos';
			formData.append(field, file);

			await pb.collection('window_items').update(data.window.id, formData);
			message = `${mode} photo uploaded`;
			await invalidateAll();
		} catch (err) {
			message = err instanceof Error ? err.message : 'Upload failed';
		} finally {
			uploading = false;
		}
	}
</script>

<section class="mx-auto max-w-lg space-y-6">
	<div>
		<p class="text-sm text-slate-500">Guided capture</p>
		<h1 class="text-2xl font-bold">{windowTypeLabel(data.window.window_type)}</h1>
	</div>

	<div class="card space-y-4">
		<p class="text-sm text-slate-600">
			Place the <strong>calibration target</strong> (90° triangle, 10 mm red lines) in frame. Capture
			both inner and outer reveals.
		</p>

		<div class="flex rounded-lg border border-slate-200 p-1">
			<button
				class="flex-1 rounded-md px-3 py-2 text-sm font-medium {mode === 'inner'
					? 'bg-[var(--color-brand)] text-white'
					: 'text-slate-600'}"
				onclick={() => (mode = 'inner')}
			>
				Inner
			</button>
			<button
				class="flex-1 rounded-md px-3 py-2 text-sm font-medium {mode === 'outer'
					? 'bg-[var(--color-brand)] text-white'
					: 'text-slate-600'}"
				onclick={() => (mode = 'outer')}
			>
				Outer
			</button>
		</div>

		<div
			class="flex aspect-[3/4] items-center justify-center rounded-lg border-2 border-dashed border-slate-300 bg-slate-100"
		>
			<div class="text-center text-sm text-slate-500">
				<p>Calibration overlay placeholder</p>
				<p class="mt-1">Target triangle + frame guides</p>
			</div>
		</div>

		<button class="btn btn-primary w-full" disabled={uploading} onclick={capturePhoto}>
			{uploading ? 'Uploading…' : `Capture ${mode} photo`}
		</button>

		{#if message}
			<p class="text-sm text-slate-600">{message}</p>
		{/if}

		<div class="grid grid-cols-2 gap-2 text-xs text-slate-500">
			<p>Inner: {data.window.inner_photos?.length ?? 0} photos</p>
			<p>Outer: {data.window.outer_photos?.length ?? 0} photos</p>
		</div>
	</div>

	<a class="btn btn-secondary" href="/windows/{data.window.id}/review">Continue to review →</a>
</section>
