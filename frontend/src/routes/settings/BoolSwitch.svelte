<script lang="ts">
	import type { SettingType } from '$lib/utils/interfaces';
	import { Check, X } from '@lucide/svelte';
	import { updateSettings } from '$lib/api/settings';
	import toast from 'svelte-french-toast';

	export let json_record: Record<string, any>;
	export let setting: SettingType;

	let value: boolean = json_record[setting.jsonKey];

	function toggleValue() {
		json_record[setting.jsonKey] = value;

		try {
			toast.promise(updateSettings(json_record), {
				loading: 'Updating settings...',
				success: 'Settings updated.',
				error: `An error occured while updating settings:`
			});
		} catch (error) {
			toast.error('' + error);
		}
	}
</script>

<label class="relative inline-flex cursor-pointer items-center">
	<!-- hidden -->
	<input type="checkbox" class="peer sr-only" bind:checked={value} on:change={toggleValue} />

	<!-- bg -->
	<div
		class="peer h-8 w-16 rounded-full bg-gray-300 transition-colors duration-300 peer-checked:bg-orange-500"
	></div>

	<!-- ball -->
	<div
		class="absolute top-1 left-1 flex h-6 w-6 items-center justify-center rounded-full bg-white text-xs transition-transform duration-300 peer-checked:translate-x-8"
	>
		{#if value}
			<Check class="h-4 w-4 text-orange-500" />
		{:else}
			<X class="h-4 w-4 text-gray-400" />
		{/if}
	</div>
</label>

<style></style>
