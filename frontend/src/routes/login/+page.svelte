<script lang="ts">
	import { goto } from '$app/navigation';

	import { CircleAlert, OctagonX, Eye, EyeClosed } from '@lucide/svelte';

	import { loggedIn, storeAfterLogin, login } from '$lib/stores/session';
	import { randomTips } from '$lib/utils/tools';

	let tip = randomTips();
	let password: string;
	let loginAttempted: boolean = false;
	let shake: boolean = false;
	let showPassword: boolean = false;

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

	function toggleShowPassword() {
		showPassword = !showPassword;
	}

	// Login page before accessing settings
</script>

<svelte:head>
	<title>Login</title>
	<meta name="description" content="Access files with urls!" />
</svelte:head>

<div>
	<h1 class="m-10">Please Login Before Continuing</h1>

	<div class="mb-4 justify-center">
		<form
			class="mb-4 flex flex-wrap gap-2"
			onsubmit={(e) => {
				e.preventDefault();
				tryLogin(password);
			}}
		>
			<div class="relative w-full max-w-2xl">
				<input
					type={showPassword ? 'text' : 'password'}
					bind:value={password}
					placeholder="Password"
					autocomplete="new-password"
					class="w-full flex-grow rounded-xl border border-gray-300 p-4 pr-12 align-middle shadow-md focus:ring-gray-300"
				/>
				<button
					type="button"
					onclick={toggleShowPassword}
					class="absolute top-1/2 right-3 -translate-y-1/2 rounded-full p-2 transition hover:bg-gray-100 active:bg-gray-200"
					aria-label="Toggle password visibility"
				>
					{#if showPassword}
						<Eye class="h-5 w-5 text-gray-600" />
					{:else}
						<EyeClosed class="h-5 w-5 text-gray-600" />
					{/if}
				</button>
			</div>
			<button
				type="submit"
				class="flex w-full flex-grow items-center justify-center gap-2 rounded-xl border-2 border-orange-500 p-4 text-orange-500 transition hover:cursor-pointer
            hover:bg-orange-500 hover:text-white sm:mb-0
            sm:w-auto sm:max-w-[25%]"
			>
				Confirm
			</button>
		</form>

		{#if loginAttempted}
			<p class="mb-4 flex gap-2 text-red-500 {shake ? 'shake' : ''}">
				<OctagonX />Password incorrect. Please try again. [If the password is correct, the backend
				server may be down]
			</p>
		{/if}

		<p class="flex gap-2 text-orange-500">
			<CircleAlert />Tips: {tip}
		</p>
	</div>
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
