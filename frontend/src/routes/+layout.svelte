<script lang="ts">
	import { Toaster } from 'svelte-french-toast';
	import IconLinks from './IconLinks.svelte';
	import Header from './Header.svelte';

	import svelte from '$lib/images/svelte.svg';
	import golang from '$lib/images/golang.svg';

	import '../app.css';

	let { children } = $props();
</script>

<div class="app">
	<Header />

	<main>
		{@render children()}
	</main>

	<Toaster
		position="bottom-right"
		toastOptions={{
			className: 'p-8 sm:p-4 shadow-4xl text-md rounded-xl'
		}}
	/>

	<footer>
		<IconLinks />
		<div class="mt-4 flex flex-col items-end gap-2 text-xs text-gray-500">
			<span>© 2025 Mah Han Cheng</span>
			<div class="flex">
				<span class="pr-2">Powered by </span>
				<img src={svelte} alt="svelte" class="h-4 w-4" />
				<span class="pr-2 pl-2">+</span>
				<img src={golang} alt="golang" class="h-5 w-5" />
			</div>
		</div>
	</footer>
</div>

<style>
	.app {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}

	main {
		flex: 1;
		display: flex;
		flex-direction: column;
		padding: 1rem;
		width: 100%;
		max-width: 64rem;
		margin: 0 auto;
		box-sizing: border-box;
	}

	footer {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
		padding: 12px;
	}

	@media (min-width: 480px) {
		footer {
			padding: 12px 0;
		}
	}
</style>
