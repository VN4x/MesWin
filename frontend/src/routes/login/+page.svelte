<script lang="ts">
	import { goto } from '$app/navigation';
	import { login } from '$lib/pb';

	let email = $state('');
	let password = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleLogin(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';
		try {
			await login(email, password);
			goto('/orders');
		} catch (err) {
			error = err instanceof Error ? err.message : 'Login failed';
		} finally {
			loading = false;
		}
	}
</script>

<section class="mx-auto max-w-md space-y-6 pt-12">
	<div class="text-center">
		<h1 class="text-2xl font-bold text-[var(--color-brand)]">mesWin</h1>
		<p class="mt-1 text-sm text-slate-600">Sign in to your organization</p>
	</div>

	<form class="card space-y-4" onsubmit={handleLogin}>
		<div>
			<label class="mb-1 block text-sm font-medium" for="email">Email</label>
			<input id="email" class="input" type="email" bind:value={email} required />
		</div>
		<div>
			<label class="mb-1 block text-sm font-medium" for="password">Password</label>
			<input id="password" class="input" type="password" bind:value={password} required />
		</div>
		{#if error}
			<p class="text-sm text-red-600">{error}</p>
		{/if}
		<button class="btn btn-primary w-full" disabled={loading}>
			{loading ? 'Signing in…' : 'Sign in'}
		</button>
	</form>
</section>
