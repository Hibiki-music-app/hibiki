<script lang="ts">
    import {goto} from "$app/navigation";
    import {authClient} from "$lib/client";
    let email = "", password = "", name = "";
    let message = "", isError = false;


    async function handleRegister(event: Event) {
        event.preventDefault();
        const { data, error } = await authClient.signUp.email({
            email,
            password,
            name,
        }, {
            onSuccess: () => {
                message = "Success";
                isError = false;
            },
            onError: (ctx) => {
                alert(ctx.error.message);
                message = error?.message || "Registration failed.";
                isError = true;
            },
        });
        if (data) {
            goto("/auth/login");
        }
    }

</script>

<div class="bg-cover bg-center bg-no-repeat flex items-center justify-center w-full h-svh">
	<form on:submit={handleRegister} method="POST">
		<div class="hero bg-base-200 min-h-screen">
			<div class="hero-content flex-col lg:flex-row-reverse">
				<div class="text-center lg:text-left">
					<h1 class="text-5xl font-bold">Hibiki create!</h1>
					<p class="py-6">
						<a href="/auth/login">Already created?</a>
					</p>
				</div>
				<div class="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
					<div class="card-body">
						<div class="text-center" id="message">
							<p style="color: {isError ? 'red' : 'green'}">{message}</p>
						</div>

						<label class="label" for="email">Email</label>
						<input type="email" class="input" id="email" placeholder="email" bind:value={email}/>

                        <label class="label" for="name">Name</label>
                        <input type="text" class="input" id="name" placeholder="nickname" bind:value={name}/>

						<label class="label" for="password">Password</label>
						<input type="password" class="input" id="password" placeholder="password" bind:value={password}/>

						<button class="btn bg-gradient-to-r from-blue-500 to-red-500 mt-4" id="RegisterButton"
							>Register</button
						>
					</div>
				</div>
			</div>
		</div>
	</form>
</div>
