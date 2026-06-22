<script lang="ts">
	import { goto } from '$app/navigation';
	import { pb } from '$lib/pb';
	import { createCustomer, createOrder } from '$lib/orders';

	let name = $state('');
	let address = $state('');
	let contact = $state('');
	let email = $state('');
	let phone = $state('');
	let notes = $state('');
	let error = $state('');
	let loading = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';

		try {
			const orgId = pb.authStore.record?.organization;
			if (!orgId) {
				throw new Error('Your user account has no organization assigned.');
			}

			const customer = await createCustomer({
				organization: orgId,
				name,
				address,
				contact,
				email,
				phone,
				notes
			});

			const order = await createOrder({
				organization: orgId,
				customer: customer.id,
				notes,
				status: 'draft'
			});

			goto(`/orders/${order.id}`);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to create order';
		} finally {
			loading = false;
		}
	}
</script>

<section class="mx-auto max-w-xl space-y-4">
	<h1 class="text-2xl font-bold">New measuring order</h1>
	<p class="text-sm text-slate-600">Customer details for this site visit.</p>

	<form class="card space-y-4" onsubmit={handleSubmit}>
		<div>
			<label class="mb-1 block text-sm font-medium" for="name">Customer name</label>
			<input id="name" class="input" bind:value={name} required />
		</div>
		<div>
			<label class="mb-1 block text-sm font-medium" for="address">Address</label>
			<input id="address" class="input" bind:value={address} />
		</div>
		<div class="grid gap-4 sm:grid-cols-2">
			<div>
				<label class="mb-1 block text-sm font-medium" for="contact">Contact person</label>
				<input id="contact" class="input" bind:value={contact} />
			</div>
			<div>
				<label class="mb-1 block text-sm font-medium" for="phone">Phone</label>
				<input id="phone" class="input" bind:value={phone} />
			</div>
		</div>
		<div>
			<label class="mb-1 block text-sm font-medium" for="email">Email</label>
			<input id="email" class="input" type="email" bind:value={email} />
		</div>
		<div>
			<label class="mb-1 block text-sm font-medium" for="notes">Visit notes</label>
			<textarea id="notes" class="input min-h-24" bind:value={notes}></textarea>
		</div>
		{#if error}
			<p class="text-sm text-red-600">{error}</p>
		{/if}
		<button class="btn btn-primary w-full" disabled={loading}>
			{loading ? 'Creating…' : 'Create order & add windows'}
		</button>
	</form>
</section>
