<script lang="ts">
	import { goto } from '$app/navigation';

	import { CircleAlert, OctagonX } from '@lucide/svelte';

	import { loggedIn, storeAfterLogin, login } from '$lib/stores/session';
	import { randomTips } from '$lib/utils/tools';

	let tip = randomTips();
	let password: string;
	let loginAttempted: boolean = false;
	let shake: boolean = false;

	$: if ($loggedIn) {
		// Success: Go back
		let redirectTo: string = sessionStorage.getItem(storeAfterLogin) || '/';
		if (redirectTo.includes('..') || redirectTo.includes('//') || !redirectTo.startsWith('/'))
			redirectTo = '/';
		sessionStorage.removeItem(storeAfterLogin);
		goto(redirectTo);
	}

	async function tryLogin(password: string) {
		const success = await login(password);

		if (!success) {
			loginAttempted = true;
			shake = false; // Failed: Shake
			requestAnimationFrame(() => {
				shake = true;
			});
		} else {
			loginAttempted = false;
		}
	}
</script>

<svelte:head>
	<title>Login</title>
	<meta name="description" content="Access files with urls!" />
</svelte:head>

<div>
	<h1 class="m-10">Please Login Before Continuing</h1>

	<div class="mb-4">
		<form
			class="flex flex-wrap items-center gap-2"
			onsubmit={(e) => {
				e.preventDefault();
				tryLogin(password);
			}}
		>
			<input
				type="password"
				bind:value={password}
				placeholder="Password"
				class="w-full min-w-[200px] flex-grow rounded-xl border border-gray-300 p-4 align-middle shadow-md focus:ring-gray-300 sm:w-auto"
			/>
			<button
				type="submit"
				class="mb-4 flex w-full flex-grow items-center justify-center gap-2 rounded-xl border-2 border-orange-500 p-4 text-orange-500 transition hover:cursor-pointer
            hover:bg-orange-500 hover:text-white sm:mb-0
            sm:w-auto sm:max-w-[25%]"
			>
				Confirm
			</button>
		</form>
	</div>

	{#if loginAttempted}
		<p class="mb-4 flex gap-2 text-red-500 {shake ? 'shake' : ''}">
			<OctagonX />Password incorrect. Please try again.
		</p>
	{/if}

	<p class="flex gap-2 text-orange-500">
		<CircleAlert />Tips: {tip}
	</p>
</div>

<style>
	@keyframes shake {
		0% {
			transform: translateX(0);
		}
		25% {
			transform: translateX(-5px);
		}
		50% {
			transform: translateX(5px);
		}
		75% {
			transform: translateX(-5px);
		}
		100% {
			transform: translateX(0);
		}
	}
	.shake {
		animation: shake 0.25s linear;
	}
</style>
