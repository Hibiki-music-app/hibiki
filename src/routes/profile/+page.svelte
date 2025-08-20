<script>
	export let data;
	const { user } = data;
</script>

<div class="hero min-h-screen" style="background-image: url(/wallpaper_login_page.jpg);">
	<div class="hero-overlay"></div>
	<div class="hero-content text-neutral-content text-center">
		<div class="max-w-md">
			<h1 class="mb-5 text-5xl font-bold">Hello there</h1>
			<p class="mb-5">welcome to your profile page.</p>
			{#if user}
				<p class="mb-5">Email : {user.email}</p>
				<p class="mb-5">Username : {user.username}</p>
			{/if}
			<button class="btn btn-primary">Get Started</button>
		</div>
	</div>
</div>
