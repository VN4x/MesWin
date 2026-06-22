import type { PageLoad } from './$types';
import { getOrder, listWindows } from '$lib/orders';

export const load: PageLoad = async ({ params }) => {
	const [order, windows] = await Promise.all([
		getOrder(params.id),
		listWindows(params.id)
	]);
	return { order, windows };
};
