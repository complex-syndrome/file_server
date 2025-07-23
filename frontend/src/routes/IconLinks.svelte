<script lang="ts">
	import github from '$lib/images/github.svg';
	import linkedin from '$lib/images/linkedin.svg';

	type iconLink = {
		link: string;
		image: string;
		alt: string;
	};
	const icons: iconLink[] = [
		{ link: 'https://github.com/complex-syndrome', image: github, alt: 'GitHub' },
		{
			link: 'https://www.linkedin.com/in/mah-han-cheng-688323314/',
			image: linkedin,
			alt: 'Linkedin'
		}
	];

	// Author links to social media
</script>

<div class="flex gap-2">
	{#each icons as icon}
		<a class="h-10 w-10" href={icon.link} target="_blank" rel="noopener noreferrer">
			<img src={icon.image} alt={icon.alt} class="h-full w-full object-contain" />
		</a>
	{/each}
</div>
