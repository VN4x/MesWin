<script lang="ts">
	import { invalidateAll } from '$app/navigation';
	import { exportOrderJson } from '$lib/orders';
	import { WINDOW_TYPES, windowTypeLabel } from '$lib/window-types';
	import { addWindow } from '$lib/orders';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();

	let selectedType = $state('single_tilt_turn');
	let label = $state('');
	let adding = $state(false);
	let exportData = $state<unknown>(null);

	async function handleAddWindow() {
		adding = true;
		try {
			await addWindow(data.order.id, selectedType, label);
			label = '';
			await invalidateAll();
		} finally {
			adding = false;
		}
	}

	async function handleExport() {
		exportData = await exportOrderJson(data.order.id);
	}
</script>

<section class="space-y-6">
	<div class="flex flex-wrap items-start justify-between gap-4">
		<div>
			<p class="text-sm text-slate-500">Measuring order</p>
			<h1 class="text-2xl font-bold">{data.order.expand?.customer?.name ?? 'Order'}</h1>
			<p class="text-sm text-slate-600">{data.order.expand?.customer?.address}</p>
		</div>
		<div class="flex gap-2">
			<button class="btn btn-secondary" onclick={handleExport}>Export JSON</button>
			<span class="badge badge-warn">{data.order.status}</span>
		</div>
	</div>

	<div class="card space-y-4">
		<h2 class="font-semibold">Add window to cart</h2>
		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<label class="mb-1 block text-sm font-medium" for="type">Window type</label>
				<select id="type" class="input" bind:value={selectedType}>
					{#each WINDOW_TYPES as wt}
						<option value={wt.slug}>{wt.label}</option>
					{/each}
				</select>
			</div>
			<div>
				<label class="mb-1 block text-sm font-medium" for="label">Label (e.g. Kitchen)</label>
				<input id="label" class="input" bind:value={label} placeholder="Optional" />
			</div>
		</div>
		<button class="btn btn-primary" disabled={adding} onclick={handleAddWindow}>
			{adding ? 'Adding…' : 'Add window'}
		</button>
	</div>

	<div class="space-y-3">
		<h2 class="font-semibold">Windows ({data.windows.length})</h2>
		{#if data.windows.length === 0}
			<p class="text-sm text-slate-500">No windows yet — add one above.</p>
		{:else}
			<ul class="space-y-3">
				{#each data.windows as window}
					<li class="card">
						<div class="flex flex-wrap items-center justify-between gap-3">
							<div>
								<p class="font-medium">
									{window.label || windowTypeLabel(window.window_type)}
								</p>
								<p class="text-sm text-slate-500">{windowTypeLabel(window.window_type)}</p>
							</div>
							<div class="flex flex-wrap gap-2">
								<span class="badge badge-warn">{window.status}</span>
								<a class="btn btn-secondary" href="/windows/{window.id}/camera">Camera</a>
								<a class="btn btn-primary" href="/windows/{window.id}/review">Review</a>
							</div>
						</div>
					</li>
				{/each}
			</ul>
		{/if}
	</div>

	{#if exportData}
		<div class="card">
			<h3 class="mb-2 font-semibold">Export preview</h3>
			<pre class="overflow-auto text-xs text-slate-700">{JSON.stringify(exportData, null, 2)}</pre>
		</div>
	{/if}
</section>
