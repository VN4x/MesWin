<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { page } from '$app/stores';
	import { logout, pb } from '$lib/pb';

	let { children } = $props();

	const showNav = $derived($page.url.pathname !== '/login');
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>mesWin</title>
</svelte:head>

{#if showNav}
	<header class="border-b border-slate-200 bg-white">
		<div class="mx-auto flex max-w-5xl items-center justify-between px-4 py-3">
			<a href="/" class="text-lg font-bold text-[var(--color-brand)]">mesWin</a>
			<nav class="flex items-center gap-4 text-sm">
				<a href="/orders" class="text-slate-600 hover:text-slate-900">Orders</a>
				<a href="/orders/new" class="btn btn-primary">New order</a>
				{#if pb.authStore.record}
					<span class="hidden text-slate-500 sm:inline">{pb.authStore.record.email}</span>
				{/if}
				<button class="btn btn-secondary" onclick={() => logout()}>Logout</button>
			</nav>
		</div>
	</header>
{/if}

<main class="mx-auto max-w-5xl px-4 py-6">
	{@render children()}
</main>
