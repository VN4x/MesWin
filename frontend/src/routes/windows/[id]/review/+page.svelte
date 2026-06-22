<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { processWindow } from '$lib/orders';
	import { windowTypeLabel } from '$lib/window-types';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let processing = $state(false);
	let result = $state<Record<string, unknown> | null>(null);
	let error = $state('');

	$effect(() => {
		result = (data.window.cv_result as Record<string, unknown>) ?? null;
	});

	async function runProcessing() {
		processing = true;
		error = '';
		try {
			const response = await processWindow(data.window.id);
			result = (response as { result: Record<string, unknown> }).result;
			await invalidateAll();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Processing failed';
		} finally {
			processing = false;
		}
	}

	const confidence = $derived(
		typeof result?.confidence === 'number' ? result.confidence : data.window.confidence
	);
</script>

<section class="mx-auto max-w-2xl space-y-6">
	<div>
		<p class="text-sm text-slate-500">Review measurements</p>
		<h1 class="text-2xl font-bold">{windowTypeLabel(data.window.window_type)}</h1>
	</div>

	<div class="card space-y-4">
		<div class="flex items-center justify-between">
			<h2 class="font-semibold">CV results</h2>
			{#if confidence != null}
				<span class="badge {confidence >= 0.7 ? 'badge-ok' : 'badge-warn'}">
					Confidence {Math.round(confidence * 100)}%
				</span>
			{/if}
		</div>

		{#if !result}
			<p class="text-sm text-slate-600">No measurements yet. Upload photos, then run processing.</p>
		{:else}
			<dl class="grid grid-cols-2 gap-3 text-sm">
				<div>
					<dt class="text-slate-500">Width</dt>
					<dd class="text-lg font-semibold">{result.width_mm} mm</dd>
				</div>
				<div>
					<dt class="text-slate-500">Height</dt>
					<dd class="text-lg font-semibold">{result.height_mm} mm</dd>
				</div>
			</dl>
			{#if Array.isArray(result.warnings) && result.warnings.length}
				<ul class="text-sm text-amber-700">
					{#each result.warnings as warning}
						<li>• {warning}</li>
					{/each}
				</ul>
			{/if}
			{#if typeof result.sketch_svg === 'string'}
				<div class="rounded border border-slate-200 bg-white p-2">
					{@html result.sketch_svg}
				</div>
			{/if}
		{/if}

		<button class="btn btn-primary" disabled={processing} onclick={runProcessing}>
			{processing ? 'Processing…' : 'Run CV processing'}
		</button>
		{#if error}
			<p class="text-sm text-red-600">{error}</p>
		{/if}
	</div>

	<a class="btn btn-secondary" href="/orders/{data.window.order}">← Back to order</a>
</section>
