<script lang="ts">
	import { onMount } from 'svelte';
	import { listOrders, type MeasuringOrder } from '$lib/orders';

	let orders = $state<MeasuringOrder[]>([]);
	let loading = $state(true);
	let error = $state('');

	onMount(async () => {
		try {
			orders = await listOrders();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load orders';
		} finally {
			loading = false;
		}
	});
</script>

<section class="space-y-4">
	<div class="flex items-center justify-between">
		<h1 class="text-2xl font-bold">Measuring orders</h1>
		<a href="/orders/new" class="btn btn-primary">New order</a>
	</div>

	{#if loading}
		<p class="text-slate-500">Loading…</p>
	{:else if error}
		<p class="text-red-600">{error}</p>
	{:else if orders.length === 0}
		<div class="card text-center text-slate-600">
			<p>No orders yet.</p>
			<a href="/orders/new" class="btn btn-primary mt-4">Create first order</a>
		</div>
	{:else}
		<ul class="space-y-3">
			{#each orders as order}
				<li>
					<a href="/orders/{order.id}" class="card block hover:border-[var(--color-brand)]">
						<div class="flex items-center justify-between gap-4">
							<div>
								<p class="font-medium">
									{order.expand?.customer?.name ?? 'Customer'}
								</p>
								<p class="text-sm text-slate-500">
									{order.expand?.customer?.address ?? '—'}
								</p>
							</div>
							<span class="badge badge-warn">{order.status}</span>
						</div>
					</a>
				</li>
			{/each}
		</ul>
	{/if}
</section>
