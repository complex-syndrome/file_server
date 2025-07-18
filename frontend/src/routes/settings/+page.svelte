<script lang="ts">
	import { SettingDescriptions } from "$lib/utils/settingType";
	import { onMount, onDestroy } from "svelte";
	import { ConnectSocket } from "$lib/api/ws";
	import BoolSwitch from "./BoolSwitch.svelte"
	
	let currentSettings: Record<string, any> = {};

	let socket: WebSocket;
    let selfClose: boolean = false;
    let refreshTimeout: ReturnType<typeof setTimeout> | null = null;
	
	async function refreshSettings() {
		const res = await fetch("api/settings");
		currentSettings = await res.json();
	}

    function setRefreshTimeout(t: ReturnType<typeof setTimeout> | null) {
        if (refreshTimeout) clearTimeout(refreshTimeout);
        refreshTimeout = t;
    }

	onMount(async () => {
		await refreshSettings()
        socket = ConnectSocket(() => { return selfClose; }, refreshSettings, setRefreshTimeout)
	})
	onDestroy(() => {
		selfClose = true;	
		if (refreshTimeout) clearTimeout(refreshTimeout);
        if (socket && socket.readyState === WebSocket.OPEN) {
            socket.close();
        }
    });
	// currentSettings = {
	//		"AllowAllIPs": true
	//		"ResourcePath": "/path/somepath"
	// }

</script>


<svelte:head>
	<title>Settings</title>
	<meta name="description" content="Access files with urls!" />
</svelte:head>

<div>

	<h1 class="text-2xl font-bold mb-10">Settings</h1>
	
	<div class="max-w-2xl w-full mx-auto space-y-6">
		{#each SettingDescriptions as setting}
		<div class="flex gap-4 text-left items-center justify-between">
			<div class="flex-grow">
				<h2 class="text-lg font-bold text-gray-800">{setting.title}</h2>
				<p class="text-md text-gray-500">{setting.description}</p>
			</div>
			{#if typeof currentSettings[setting.jsonKey] === 'boolean'}
				<BoolSwitch setting={setting} json_record={currentSettings}/>
			{:else}
				<p>Unfortunately, currently there is no config provided for this setting.</p>
			{/if}
		</div>

		{/each}
	</div>
</div>
