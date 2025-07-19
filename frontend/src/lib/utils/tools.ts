import type { FileInfo } from '$lib/utils/interfaces';

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
