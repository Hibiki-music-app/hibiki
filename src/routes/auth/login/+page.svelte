<script lang="ts">
	import { authClient } from '$lib/client';
	import { onMount } from 'svelte';

	let email = '',
		password = '';
	let message = '',
		isError = false;

	onMount(() => {
		authClient.oneTap(); // google one tap
	});

	// email password classique 
	async function handleLogin(event: Event) {
		event.preventDefault();
		const { data, error } = await authClient.signIn.email(
			{
				email,
				password,
				callbackURL: '/profile',
				rememberMe: false
			},
			{
				onSuccess: () => {
					message = 'Success';
					isError = false;
				},
				onError: (ctx) => {
					alert(ctx.error.message);
					message = error?.message || 'Login failed.';
					isError = true;
				}
			}
		);
	};

	// google oauth
	async function handleGoogleLogin() {
		await authClient.signIn.social(
			{
				provider: 'google',
				callbackURL: '/profile',
			},
			{
				onSuccess: () => {
					message: 'Signed in with Google';
					isError: false;
				},
				onError: (ctx) => {
					alert(ctx.error.message);
					message = ctx.error.message || 'Google login failed';
					isError = true;
				}
			},
		)
	};

</script>

<div class="bg-cover bg-center bg-no-repeat flex items-center justify-center w-full h-svh">
	<form on:submit={handleLogin} method="POST">
		<div class="hero bg-base-200 min-h-screen">
			<div class="hero-content flex-col lg:flex-row-reverse">
				<div class="text-center lg:text-left">
					<h1 class="text-5xl font-bold">Hibiki connect!</h1>
					<p class="py-6">
						<a href="/auth/signup">First time here?</a>
					</p>
				</div>
				<div class="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
					<div class="card-body">
						<div class="text-center" id="message">
							<p style="color: {isError ? 'red' : 'green'}">{message}</p>
						</div>

						<label class="label" for="email">Email</label>
						<input type="email" class="input" id="email" placeholder="email" bind:value={email} />

						<label class="label" for="password">Password</label>
						<input
							type="password"
							class="input"
							id="password"
							placeholder="password"
							bind:value={password}
						/>

						<button class="btn bg-gradient-to-r from-violet-300 to-blue-600 mt-4" id="loginButton"
							>Login</button
						>

						<button 
							class="btn bg-white text-black border-[#e5e5e5]"
							type="button"
							on:click={handleGoogleLogin}
						>
							<svg
								aria-label="Google logo"
								width="16"
								height="16"
								xmlns="http://www.w3.org/2000/svg"
								viewBox="0 0 512 512"
								><g
									><path d="m0 0H512V512H0" fill="#fff"></path><path
										fill="#34a853"
										d="M153 292c30 82 118 95 171 60h62v48A192 192 0 0190 341"
									></path><path
										fill="#4285f4"
										d="m386 400a140 175 0 0053-179H260v74h102q-7 37-38 57"
									></path><path fill="#fbbc02" d="m90 341a208 200 0 010-171l63 49q-12 37 0 73"
									></path><path
										fill="#ea4335"
										d="m153 219c22-69 116-109 179-50l55-54c-78-75-230-72-297 55"
									></path></g
								></svg
							>
							Login with Google
						</button>

					</div>
				</div>
			</div>
		</div>
	</form>
</div>
