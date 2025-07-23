import type { FileInfo } from '$lib/utils/interfaces';
import { customHeader, tips } from './consts';

// Some tools that can be used even outside this project

// Partial fuzzy
export function filterFilesFuzzy(text: string, fileArray: FileInfo[]): FileInfo[] {
	text = text?.trim?.() || '';
	if (!text) {
		return fileArray;
	}
	const escaped = text.trim().replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
	const regexPattern = '.*' + escaped.replace(/\s+/g, '.*') + '.*';
	const regex = new RegExp(regexPattern, 'i');
	return fileArray.filter((file) => regex.test(file.name));
}

export function randomTips(): string {
	return tips[Math.floor(Math.random() * tips.length)];
}

export async function customFetch(
	input: RequestInfo | URL,
	init: RequestInit = {}
): Promise<Response> {
	return fetch(input, {
		...init,
		headers: {
			...(init.headers || {}),
			...customHeader
		}
	});
}
