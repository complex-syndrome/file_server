<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { goto } from '$app/navigation';

	import { SettingDescriptions } from '$lib/utils/consts';
	import { ConnectSocket } from '$lib/connectAPI/ws';
	import BoolSwitch from './BoolSwitch.svelte';
	import { loggedIn, storeAfterLogin } from '$lib/stores/session';
	import { customFetch } from '$lib/utils/tools';
	import toast from 'svelte-french-toast';

	let currentSettings: Record<string, any> = {};

	let socket: WebSocket;
	let selfClose: boolean = false;
	let refreshTimeout: ReturnType<typeof setTimeout> | null = null;

	const linkKeyword = 'link:';

	function setRefreshTimeout(t: ReturnType<typeof setTimeout> | null) {
		if (refreshTimeout) clearTimeout(refreshTimeout);
		refreshTimeout = t;
	}

	// Connect to websocket (refreshfunc is already runned once on ws open)
	onMount(() => {
		const unsub = loggedIn.subscribe((value) => {
			if (!value) {
				sessionStorage.setItem(storeAfterLogin, window.location.pathname);
				goto('/login');
			} else {
				setTimeout(() => {
					socket = ConnectSocket(() => selfClose, refreshSettings, setRefreshTimeout);
				}, 0);
			}
		});
		return unsub();
	});

	// Close ws
	onDestroy(() => {
		selfClose = true;
		if (refreshTimeout) clearTimeout(refreshTimeout);
		if (socket && socket.readyState === WebSocket.OPEN) {
			socket.close();
		}
	});

	async function refreshSettings() {
		try {
			const res = await customFetch(`api/settings/list`);
			currentSettings = await res.json();
		} catch {
			toast.error('Error refreshing settings');
			currentSettings = {};
		}
	}

	// Settings page
</script>

<svelte:head>
	<title>Settings</title>
	<meta name="description" content="Access files with urls!" />
</svelte:head>

{#if $loggedIn}
	<div>
		<h1 class="m-10">Settings</h1>
		<div class="mx-auto w-full max-w-2xl space-y-6">
			{#each SettingDescriptions as setting}
				<div class="flex items-center justify-between gap-4 text-left">
					<div class="flex-grow">
						<h2 class="text-lg font-bold text-gray-800">{setting.title}</h2>
						{#each setting.description as d}
							{#if d.startsWith(linkKeyword)}
								<a
									href={d.replace(linkKeyword, '').trim()}
									class="text-sm text-orange-500 hover:underline"
									target="_blank"
									rel="noopener noreferrer"
								>
									{d.replace(linkKeyword, '').trim()}
								</a>
							{:else}
								<p class="text-sm text-gray-500">{d}</p>
							{/if}
						{/each}
					</div>
					{#if typeof currentSettings[setting.jsonKey] === 'boolean'}
						<BoolSwitch {setting} json_record={currentSettings} />
					{:else}
						<p>Setting for this is currently unavailable.</p>
					{/if}
				</div>
			{/each}
		</div>
	</div>
{/if}
