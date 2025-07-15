<script lang="ts">
	import { Toaster } from 'svelte-french-toast';
	import IconLinks from './IconLinks.svelte';

	import svelte from '$lib/images/svelte.svg';
	import golang from '$lib/images/golang.svg';

	import '../app.css';

	let { children } = $props();
</script>

<div class="app">
	<main>
		{@render children()}	
	</main>

	<Toaster 
	position = "bottom-right"
	toastOptions={{
		className: "p-8 shadow-4xl text-lg rounded-xl",
	}}
	/>
	<footer>
		<IconLinks/>
		
		<div class="flex flex-col items-end gap-2 text-xs text-gray-500 mt-4">
			<span>© 2025 Mah Han Cheng</span>
			<div class="flex">
				<span class="pr-2">Powered by </span>
				<img src={svelte} alt="svelte" class="w-4 h-4"/>
				<span class="pl-2 pr-2">+</span>
				<img src={golang} alt="golang" class="w-5 h-5"/>
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
