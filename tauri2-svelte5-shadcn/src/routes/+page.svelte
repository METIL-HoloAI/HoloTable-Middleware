<script lang="ts">
	import { GlobalState, preventDefault } from '$lib';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import * as Card from '$lib/components/ui/card/index.js';

	const gs = new GlobalState();

	$inspect(gs.greet, gs.name);

	const onsubmit = preventDefault(() => gs.nlen && gs.submit());
	const onclick = () => gs.reset();
</script>

<div
	class="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-indigo-500 via-purple-500 to-pink-500 transition-all duration-500"
>
	<Card.Root
		class="w-[380px] backdrop-blur-sm bg-white/90 dark:bg-gray-800/90 shadow-xl hover:shadow-2xl transition-all duration-300 rounded-xl [--ring:267_100%_60%]"
	>
		<Card.Header class="space-y-2">
			<Card.Title
				class="text-3xl font-bold text-center bg-gradient-to-r from-indigo-500 to-pink-500 bg-clip-text text-transparent"
			>
				{#if gs.greet}
					<p>{gs.greet}</p>
					<small class="text-sm">(from Rust side)</small>
				{:else}
					<p>Hello World!</p>
				{/if}
			</Card.Title>
		</Card.Header>
		<Card.Content class="p-6">
			<form {onsubmit} class="space-y-6">
				{#if !gs.greet}
					<Input
						type="text"
						placeholder="Enter your username"
						bind:value={gs.name}
						class="w-full px-4 py-2 rounded-lg outline-none border border-gray-200 dark:border-gray-700 focus:ring-2 focus:ring-purple-500 focus-visible:ring-purple-500 focus-visible:ring-offset-2 transition-all duration-200"
					/>
				{/if}
				{#if gs.name && !gs.greet}
					<Button
						type="submit"
						class="w-full bg-gradient-to-r from-indigo-500 to-pink-500 hover:opacity-90 transition-opacity duration-200"
					>
						Say Hello
					</Button>
				{:else if gs.greet}
					<Button
						{onclick}
						class="w-full bg-gradient-to-r from-indigo-500 to-pink-500 hover:opacity-90 transition-opacity duration-200"
						>Reset</Button
					>
				{/if}
			</form>
		</Card.Content>
	</Card.Root>
</div>
