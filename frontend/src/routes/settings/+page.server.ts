import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ request, fetch }) => {
	const ip =
		request.headers.get('x-forwarded-for') ||
		request.headers.get('x-real-ip') ||
		(import.meta.env.DEV ? '127.0.0.1' : 'x.x.x.x');
	// TODO user auth

	console.log(request.headers);
	const result = await fetch(`${import.meta.env.VITE_API_URL}/settings/allow`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ ip })
	});

	if (!result.ok) redirect(303, '/403');
	return {};
};
