<script lang="ts">
	import FileList from "./FileList.svelte";
</script>

<svelte:head>
	<title>My File Server</title>
	<meta name="description" content="File Uploader / Downloader" />
</svelte:head>

<section>
	<h1 class="m-10">File Uploader</h1>
	<FileList/>
</section>