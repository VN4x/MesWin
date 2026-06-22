import PocketBase from 'pocketbase';
import { browser } from '$app/environment';
import { PUBLIC_PB_URL } from '$env/static/public';

const url = PUBLIC_PB_URL || 'http://127.0.0.1:8090';

export const pb = new PocketBase(url);

if (browser) {
	pb.authStore.loadFromCookie(document.cookie);
	pb.authStore.onChange(() => {
		document.cookie = pb.authStore.exportToCookie({ httpOnly: false });
	});
}

export function isLoggedIn(): boolean {
	return pb.authStore.isValid;
}

export async function login(email: string, password: string) {
	return pb.collection('users').authWithPassword(email, password);
}

export function logout() {
	pb.authStore.clear();
}
