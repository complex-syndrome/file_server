<script lang="ts">
    import type { SettingType } from "$lib/utils/interfaces";
    import { Check, X } from "@lucide/svelte";
    import { updateSettings } from "$lib/api/settings";
    import toast from "svelte-french-toast";

    export let json_record: Record<string, any>;
    export let setting: SettingType;

    let value: boolean = json_record[setting.jsonKey]

    function toggleValue() {
        json_record[setting.jsonKey] = value;

        try {
          toast.promise(
              updateSettings(json_record),
              {
                  loading: 'Updating settings...',
                  success: 'Settings updated.',
                  error: `An error occured while updating settings:`,
              }
          );
        }
        catch (error) {
          toast.error('' + error)
        }
    }
</script>

<label class="relative inline-flex items-center cursor-pointer">
  <!-- hidden -->
  <input type="checkbox" class="sr-only peer" bind:checked={value} on:change={toggleValue}>

  <!-- bg -->
  <div class="w-16 h-8 bg-gray-300 peer-checked:bg-orange-500 rounded-full peer transition-colors duration-300"></div>

  <!-- ball -->
  <div class="absolute left-1 top-1 w-6 h-6 bg-white rounded-full transition-transform duration-300 peer-checked:translate-x-8 flex items-center justify-center text-xs">
    {#if value}
      <Check class="w-4 h-4 text-orange-500" />
    {:else}
      <X class="w-4 h-4 text-gray-400" />
    {/if}
  </div>
</label>

<style></style>