<script lang="ts">
    import { onDestroy, onMount } from 'svelte';
    
    import toast from 'svelte-french-toast'
    import { Trash2, DownloadCloud, FileUp } from '@lucide/svelte'

    import type { FileInfo } from '$lib/utils/interfaces';
    import { apiListFiles, apiDownloadFile, apiDeleteFile, apiUploadFile } from '$lib/api/files';
    import { ConnectSocket } from '$lib/api/ws';
    import { filterFilesFuzzy } from '$lib/utils/tools';

    let allFiles: FileInfo[] = []
    let filteredFiles: FileInfo[] = []
    let searchText: string = ''

    let socket: WebSocket;
    let selfClose: boolean = false;
    let refreshTimeout: ReturnType<typeof setTimeout> | null = null;

    function setRefreshTimeout(t: ReturnType<typeof setTimeout> | null) {
        if (refreshTimeout) clearTimeout(refreshTimeout);
        refreshTimeout = t;
    }

    onMount(() => {
        // For pasting files to upload
        window.addEventListener('paste', pasteUpload);
        
        // Load files
        (async () => {
            try {
                allFiles = await toast.promise(
                    apiListFiles(),
                    {
                        loading: 'Loading files...',
                        success: 'Files loaded successfully!',
                        error: 'Error: Could not load files.',
                    }
                );
                
            } catch (error) {
                allFiles = []
            }
            filteredFiles = allFiles;

        })();
        
        // Websocket to refresh
        socket = ConnectSocket(() => { return selfClose; }, refreshFileList, setRefreshTimeout)

        // Prevent multiple listeners
        return () => {
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
        await Promise.all([...f].map(file => apiUploadFile(file)));
        fileInput.value = ""
    }
    
    async function pasteUpload(event: ClipboardEvent): Promise<void> {
        const f = event.clipboardData?.files
        if (!f || f.length === 0) return;
        await Promise.all([...f].map(file => apiUploadFile(file)));
    }

    async function confirmAndDelete(fileName: string): Promise<void> {
        if (confirm(`Are you sure you want to delete ${fileName}?`)) {
            await apiDeleteFile(fileName);
        }
    }

    async function refreshFileList(): Promise<void> {
        try {
            const newAllFiles = await apiListFiles()
            allFiles = newAllFiles
            filteredFiles = filterFilesFuzzy(searchText, allFiles);

        } catch (error) {
            toast.error("Could not refresh file list. Try reloading the page.")
            console.log("Error: ", error)
            return
        }
    }

</script>

<div>
    <!-- Normal upload (hidden) -->
    <input class="hidden" type="file" multiple bind:this={fileInput} onchange={displayUploadFileDialog}/>

    <!-- Search and upload -->
    <div class="flex items-center gap-2 m-8">
        <input 
            class="p-4 flex-grow align-middle rounded-xl border border-gray-300 
            focus:ring-orange-500 focus:border-orange-500"
            type="text"
            bind:value={searchText}
            oninput={() => { filteredFiles = filterFilesFuzzy(searchText, allFiles); }}
            placeholder="Search your files here...">
        <button 
            onclick={() => { fileInput.click() }}
            class="p-4 transition rounded-xl border-2 flex items-center justify-center gap-2 min-w-3xs
                text-orange-500 border-orange-500 hover:cursor-pointer 
                hover:bg-orange-500 hover:text-white">
            <FileUp/>
            Upload Files
        </button>
    </div>

    <!-- Display -->
    {#if filteredFiles && filteredFiles.length > 0}
        <ul>
            {#each filteredFiles as file}
                <li class="p-4 m-8 mt-0 border-2 border-neutral-200 rounded-xl shadow-xl flex items-center hover:border-orange-500">

                    <div class="p-4 flex flex-col gap-1">
                    <p class="font-semibold text-gray-800 truncate max-w-[500px] min-w-0">{file.name}</p>
                    <div class="flex gap-4 text-sm">
                        <p class="text-gray-700">{file.size}</p>
                        <p class="text-gray-500">{file.mime}</p>
                    </div>
                    </div>
                    
                    <div class="ml-auto flex gap-2">
                        <button aria-label="DownloadFile"
                        onclick={() => apiDownloadFile(file.name)}
                        class="p-4 transition rounded-xl 
                        text-gray-500
                        hover:cursor-pointer hover:text-green-500
                        active:bg-gray-100">

                        <div class="flex gap-2">
                            <DownloadCloud/>
                            Download
                        </div>    
                    </button>
                    
                    <button 
                    aria-label="DeleteFile"
                    onclick={() => confirmAndDelete(file.name)} 
                    class="p-4 transition rounded-xl
                    text-gray-500 
                    hover:cursor-pointer hover:bg-red-500 hover:text-white">
                    
                        <div class="flex gap-2">
                            <Trash2/>
                            Delete
                        </div>    
                        </button>
                    </div>
                </li>
            {/each}
        </ul>
    {:else}
        <p class="mt-25 text-lg text-center">No files found.</p>
    {/if}
</div>