import type { PageLoad } from './$types';
import { pb } from '$lib/pb';
import type { WindowItem } from '$lib/orders';

export const load: PageLoad = async ({ params }) => {
	const window = await pb.collection('window_items').getOne<WindowItem>(params.id);
	return { window };
};
