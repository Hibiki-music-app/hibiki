<script lang="ts">
	import {RandomAvatar} from "$lib/services/ApiEndpoints";

    export let data;
	const { user } = data;
</script>

<div class="min-h-screen flex items-center justify-center p-4" style="background-image: url(/wallpaper_login_page.jpg); background-size: cover; background-position: center;">
	<div class="card w-96 shadow-xl glass">
		<figure class="px-10 pt-10">
			<div class="avatar">
				<div class="w-24 rounded-full ring ring-primary ring-offset-base-100 ring-offset-2">
					<img src={user?.image || RandomAvatar} alt="User Avatar" />
				</div>
			</div>
		</figure>
		<div class="card-body items-center text-center m-4">
			<h2 class="card-title">{user?.name || 'User'}</h2>
			
			<div class="w-full space-y-3 mt-4">
				<div class="flex justify-between items-center">
					<span class="font-semibold">Email:</span>
					<span class="badge badge-outline {user?.emailVerified ? 'badge-success' : 'badge-warning'}">
						{user?.email}
						{#if user?.emailVerified}
							<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" viewBox="0 0 20 20" fill="currentColor">
								<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
							</svg>
						{/if}
					</span>
				</div>
				
				<div class="flex justify-between items-center">
					<span class="font-semibold">User ID:</span>
					<span class="text-xs font-mono">{user?.id}</span>
				</div>
				
				<div class="flex justify-between items-center">
					<span class="font-semibold">Member since:</span>
					<span class="text-xs">{user?.createdAt.toLocaleDateString()}</span>
				</div>
				
				<div class="flex justify-between items-center">
					<span class="font-semibold">Last updated:</span>
					<span class="text-xs">{user?.updatedAt.toLocaleDateString()}</span>
				</div>
			</div>
			
			<div class="card-actions mt-6">
				<button class="btn btn-primary">Get Started</button>
				<button class="btn btn-ghost">Edit Profile</button>
			</div>
		</div>
	</div>
</div>
