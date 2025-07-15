import toast from 'svelte-french-toast'

import type { FileInfo } from '$lib/utils/interfaces';


export async function apiListFiles(): Promise<FileInfo[]> {
    const response = await fetch('/api/list');
    if (!response.ok) { throw Error }
    return await response.json();
}

export async function apiDownloadFile(fileName: string): Promise<void> {
    try {
        const response = await fetch(`/api/download?file=${encodeURIComponent(fileName)}`);
        if (!response.ok) { throw new Error(`Server responded with ${response.status}`); }

        const blob = await toast.promise(
            response.blob(),
            {
                loading: 'Preparing files...',
                success: 'Download starting soon... You should be able to see your downloads.',
                error: `An error occured while downloading file: ${fileName}`,
            }
        );
        downloadViaBrowser(fileName, blob)

    } catch (error) {
        toast.error("Download failed.");
    }
}


export async function apiDeleteFile(fileName: string) {
    try {
        const response = await fetch(`/api/delete?file=${encodeURIComponent(fileName)}`, { method: 'DELETE' });
        const reply = await response.text()
        if (!response.ok) throw Error(reply)
            toast.success(reply)
    } catch (error) {
        toast.error("Delete failed.");
    }
}

export async function apiUploadFile(file: File): Promise<void> {
    const formData = new FormData();
    formData.append("file", file);
    
    try {
        const response = await fetch(`/api/upload`, { 
            method: 'POST',
            body: formData,
        });
        const reply = await response.text()

        if (!response.ok) throw Error(reply) 
        toast.success(reply);
        
    } catch (error) {
        toast.error("Upload failed.");
    }
}


function downloadViaBrowser(fileName: string, blob: Blob): void {
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');

    link.href = url;
    link.download = fileName
    
    document.body.append(link);
    link.click();

    document.body.removeChild(link);
    URL.revokeObjectURL(url)
}


