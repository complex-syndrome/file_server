export async function updateSettings(json_rcd: Record<string, any>): Promise<void> {
	const response = await fetch(`${import.meta.env.VITE_API_URL}/settings/update`, {
		method: 'POST',
		body: JSON.stringify(json_rcd),
		headers: { 'Content-Type': 'application/json' }
	});
	if (!response.ok) throw Error(await response.text());
}
