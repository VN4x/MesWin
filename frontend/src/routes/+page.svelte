<script lang="ts">
	import { goto } from '$app/navigation';
</script>

<section class="space-y-6">
	<div>
		<h1 class="text-3xl font-bold">mesWin</h1>
		<p class="mt-2 text-slate-600">
			Photo-based window measurement for Rehau Synego — point, shoot, export.
		</p>
	</div>

	<div class="grid gap-4 sm:grid-cols-2">
		<a href="/orders/new" class="card hover:border-[var(--color-brand)]">
			<h2 class="font-semibold">New measuring order</h2>
			<p class="mt-1 text-sm text-slate-600">Customer visit, add windows, capture photos.</p>
		</a>
		<a href="/orders" class="card hover:border-[var(--color-brand)]">
			<h2 class="font-semibold">Open orders</h2>
			<p class="mt-1 text-sm text-slate-600">Review, override dimensions, export PDF/JSON.</p>
		</a>
	</div>

	<div class="card bg-slate-900 text-slate-100">
		<h2 class="font-semibold">Accuracy target</h2>
		<p class="mt-1 text-sm text-slate-300">±5 mm with calibration target · 10 mm tolerance band</p>
	</div>
</section>
