<script lang="ts">
	import { page } from '$app/state';

	interface pathNavElem {
		name: string
		path: string
	}

	const paths: pathNavElem[]  = [
		{name: "Home", path: "/"},
		{name: "Settings", path: "/settings"},
	]
</script>

<header class="flex justify-center">
	<nav class="flex justify-center">
		<ul class="relative p-0 m-0 h-[3em] flex justify-center items-center list-none">
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
