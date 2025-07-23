<script lang="ts">
	import { page } from '$app/state';

	interface pathNavElem {
		name: string;
		path: string;
	}

	// Nav bar (for home and settings)
	const paths: pathNavElem[] = [
		{ name: 'Home', path: '/' },
		{ name: 'Settings', path: '/settings' }
	];
</script>

<header class="flex justify-center">
	<nav class="flex justify-center">
		<ul class="relative m-0 flex h-[3em] list-none items-center justify-center p-0">
			{#each paths as path}
				<li class={page.url.pathname === path.path ? 'bg-gray-50' : ''}>
					<a href={path.path}>{path.name}</a>
				</li>
			{/each}
		</ul>
	</nav>
</header>

<style>
	ul {
		background: var(--background);
		background-size: contain;
	}

	li {
		position: relative;
		height: 100%;
	}

	nav a {
		display: flex;
		height: 100%;
		align-items: center;
		padding: 0 1rem;
		color: var(--color-text);
		font-weight: 700;
		font-size: 0.8rem;
		text-transform: uppercase;
		letter-spacing: 0.1em;
		text-decoration: none;
		transition: color 0.2s linear;
	}

	a:hover {
		color: var(--color-theme-1);
	}
</style>
