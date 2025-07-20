<script lang="ts">
	import { onDestroy, onMount } from 'svelte';

	import toast from 'svelte-french-toast';
	import { Trash2, DownloadCloud, FileUp } from '@lucide/svelte';

	import type { FileInfo } from '$lib/utils/interfaces';
	import { apiListFiles, apiDownloadFile, apiDeleteFile, apiUploadFile } from '$lib/api/files';
	import { ConnectSocket } from '$lib/api/ws';
	import { filterFilesFuzzy } from '$lib/utils/tools';

	let allFiles: FileInfo[] = [];
	let filteredFiles: FileInfo[] = [];
	let searchText: string = '';

	let socket: WebSocket;
	let selfClose: boolean = false;
	let refreshTimeout: ReturnType<typeof setTimeout> | null = null;

	function setRefreshTimeout(t: ReturnType<typeof setTimeout> | null) {
		if (refreshTimeout) clearTimeout(refreshTimeout);
		refreshTimeout = t;
	}

	onMount(() => {
		// Websocket to refresh (settimeout for running after call stack / main tasks)
		setTimeout(() => {
			socket = ConnectSocket(() => selfClose, refreshFileList, setRefreshTimeout);
		}, 0);

		// For pasting files to upload
		window.addEventListener('paste', pasteUpload);
		return () => {
			// Prevent multiple listeners
			window.removeEventListener('paste', pasteUpload);
		};
	});

	onDestroy(() => {
		selfClose = true;
		if (refreshTimeout) clearTimeout(refreshTimeout);
		if (socket && socket.readyState === WebSocket.OPEN) {
			socket.close();
		}
	});

	// Upload ways (button press / CTRL + V)
	let fileInput: HTMLInputElement;
	async function displayUploadFileDialog(event: Event): Promise<void> {
		const f = (event.target as HTMLInputElement).files;
		if (!f || f.length === 0) return;
		await Promise.all([...f].map((file) => apiUploadFile(file)));
		fileInput.value = '';
	}

	async function pasteUpload(event: ClipboardEvent): Promise<void> {
		const f = event.clipboardData?.files;
		if (!f || f.length === 0) return;
		await Promise.all([...f].map((file) => apiUploadFile(file)));
	}

	async function confirmAndDelete(fileName: string): Promise<void> {
		if (confirm(`Are you sure you want to delete ${fileName}?`)) {
			await apiDeleteFile(fileName);
		}
	}

	async function refreshFileList(): Promise<void> {
		try {
			allFiles = await apiListFiles();
		} catch (error) {
			toast.error('Could not refresh file list. Try reloading the page.');
			console.log('Error: ', error);
			allFiles = [];
		}
		filteredFiles = filterFilesFuzzy(searchText, allFiles);
	}
</script>

<div>
	<!-- Normal upload (hidden) -->
	<input
		class="hidden"
		type="file"
		multiple
		bind:this={fileInput}
		onchange={displayUploadFileDialog}
	/>

	<!-- Search and upload -->
	<div class="flex flex-wrap items-center gap-2 sm:m-4">
		<input
			class="w-full min-w-[200px] flex-grow rounded-xl border border-gray-300 p-4 align-middle focus:border-orange-500
        focus:ring-orange-500 sm:w-auto"
			type="text"
			bind:value={searchText}
			oninput={() => {
				filteredFiles = filterFilesFuzzy(searchText, allFiles);
			}}
			placeholder="Search your files here..."
		/>
		<button
			onclick={() => {
				fileInput.click();
			}}
			class="mb-4 flex w-full flex-grow items-center justify-center gap-2 rounded-xl border-2 border-orange-500 p-4 text-orange-500 transition hover:cursor-pointer
            hover:bg-orange-500 hover:text-white sm:mb-0
            sm:w-auto sm:max-w-[25%]"
		>
			<FileUp />
			Upload Files
		</button>
	</div>

	<!-- Display -->
	{#if filteredFiles && filteredFiles.length > 0}
		<ul>
			{#each filteredFiles as file}
				<li
					class="m-4 mt-0 flex flex-col gap-4 rounded-xl border-2 border-neutral-200 p-4 shadow-xl hover:border-orange-500 sm:m-8 sm:flex-row sm:items-center"
				>
					<div class="flex w-full min-w-0 flex-grow flex-col gap-1 p-4 sm:w-auto">
						<p class="max-w-full truncate font-semibold text-gray-800">{file.name}</p>
						<div class="flex gap-4 text-sm">
							<p class="text-gray-700">{file.size}</p>
							<p class="text-gray-500">{file.mime}</p>
						</div>
					</div>

					<div class="ml-auto flex w-full flex-row gap-2 sm:w-auto">
						<button
							aria-label="DownloadFile"
							onclick={() => apiDownloadFile(file.name)}
							class="rounded-xl p-4 text-gray-500
                        transition
                        hover:cursor-pointer hover:text-green-500
                        active:bg-gray-100"
						>
							<div class="flex gap-2">
								<DownloadCloud />
								Download
							</div>
						</button>

						<button
							aria-label="DeleteFile"
							onclick={() => confirmAndDelete(file.name)}
							class="rounded-xl p-4 text-gray-500
                    transition
                    hover:cursor-pointer hover:bg-red-500 hover:text-white"
						>
							<div class="flex gap-2">
								<Trash2 />
								Delete
							</div>
						</button>
					</div>
				</li>
			{/each}
		</ul>
	{:else}
		<p class="mt-25 text-center text-lg">No files found.</p>
	{/if}
</div>
