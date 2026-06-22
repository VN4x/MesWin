import type { LayoutLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { isLoggedIn } from '$lib/pb';

export const prerender = true;
export const ssr = false;

export const load: LayoutLoad = ({ url }) => {
	const publicPaths = ['/login'];
	if (!isLoggedIn() && !publicPaths.includes(url.pathname)) {
		redirect(303, '/login');
	}
};
