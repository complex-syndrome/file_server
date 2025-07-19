<script lang="ts">
	import { onMount, onDestroy } from 'svelte';

	import { SettingDescriptions } from '$lib/utils/consts';
	import { ConnectSocket } from '$lib/api/ws';
	import BoolSwitch from './BoolSwitch.svelte';

	let currentSettings: Record<string, any> = {};

	let socket: WebSocket;
	let selfClose: boolean = false;
	let refreshTimeout: ReturnType<typeof setTimeout> | null = null;

	function setRefreshTimeout(t: ReturnType<typeof setTimeout> | null) {
		if (refreshTimeout) clearTimeout(refreshTimeout);
		refreshTimeout = t;
	}

	// Connect to websocket (refreshfunc is already runned once on ws open)
	onMount(() => {
		setTimeout(() => {
			socket = ConnectSocket(() => selfClose, refreshSettings, setRefreshTimeout);
		}, 0);
	});

	// Close ws
	onDestroy(() => {
		selfClose = true;
		if (refreshTimeout) clearTimeout(refreshTimeout);
		if (socket && socket.readyState === WebSocket.OPEN) {
			socket.close();
		}
	});

	async function refreshSettings(): Promise<any> {
		try {
			const res = await fetch(`${import.meta.env.VITE_API_URL}/settings/list`);
			currentSettings = await res.json();
		} catch (error) {
			console.error(error);
			currentSettings = {};
		}
	}
</script>

<svelte:head>
	<title>Settings</title>
	<meta name="description" content="Access files with urls!" />
</svelte:head>

<div>
	<h1 class="m-10 text-2xl font-bold">Settings</h1>

	<div class="mx-auto w-full max-w-2xl space-y-6">
		{#each SettingDescriptions as setting}
			<div class="flex items-center justify-between gap-4 text-left">
				<div class="flex-grow">
					<h2 class="text-lg font-bold text-gray-800">{setting.title}</h2>
					<p class="text-md text-gray-500">{setting.description}</p>
				</div>
				{#if typeof currentSettings[setting.jsonKey] === 'boolean'}
					<BoolSwitch {setting} json_record={currentSettings} />
				{:else}
					<p>Unfortunately, currently there is no config provided for this setting.</p>
				{/if}
			</div>
		{/each}
	</div>
</div>
