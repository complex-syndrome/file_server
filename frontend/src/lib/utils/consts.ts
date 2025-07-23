import type { SettingType } from './interfaces';

// Some constants

// Secret header between frontend and backend
export const customHeader = {
	'X-From-Frontend': import.meta.env.VITE_CUSTOM_VALUE || ''
};

export const SettingDescriptions: SettingType[] = [
	{
		title: 'Allow Direct Access of API From Other IPs',
		description: [
			`Allow other IPs to directly connect to the API without having to go through the web UI.`,
			`For api docs please refer to README.md at`,
			`link:https://github.com/complex-syndrome/file_server/blob/main/README.md`
		],
		jsonKey: 'AllowOtherIPs'
	}
];

export const tips: string[] = [
	'Login is refreshed after reloading the page.',
	'To refresh the password, run the backend / whole program again after changing .env',
	'Changes to the uploaded files and settings are refreshed immediately.',
	'Use strong passwords!',
	'You can paste files using CTRL+V to upload them!'
];
