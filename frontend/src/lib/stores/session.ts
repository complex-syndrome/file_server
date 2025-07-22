import { customFetch } from '$lib/utils/tools';
import { writable } from 'svelte/store';

export const loggedIn = writable(false);
export const storeAfterLogin = 'redirectAfterLogin';

export async function login(password: string): Promise<boolean> {
	const result = await customFetch(`api/login`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ password })
	});
	if (result.ok) {
		loggedIn.set(true);
		return true;
	}
	return false;
}
