import { customFetch } from '$lib/utils/tools';

export async function updateSettings(json_rcd: Record<string, any>): Promise<void> {
	const response = await customFetch(`api/settings/update`, {
		method: 'POST',
		body: JSON.stringify(json_rcd),
		headers: {
			'Content-Type': 'application/json'
		}
	});
	if (!response.ok) throw Error(await response.text());
}

// Settings api call backend
