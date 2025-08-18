<script lang="ts">

    let username: string = '';
    let message: string = '';
    let isError: boolean = false;


    async function register() {

    }

    async function login() {

    }
</script>

<div class="hero bg-base-200 min-h-screen">
    <div class="hero-content flex-col lg:flex-row-reverse">
        <div class="text-center lg:text-left">
            <h1 class="text-5xl font-bold">Hibiki connect!</h1>
            <p class="py-6">
                Connect you without email and password, insane non?
            </p>
            <a href="http://localhost:8080/private">TEST</a> <!--func test-->
        </div>
        <div class="card bg-base-100 w-full max-w-sm shrink-0 shadow-2xl">
            <div class="card-body">
                <div class="text-center" id="message">
                    <p style="color: {isError ? 'red' : 'green'}">{message}</p>
                </div>

                <label class="label" for="username">Username</label>
                <input type="text" class="input" id="username" placeholder="username" bind:value={username}>

                <button class="btn bg-gradient-to-r from-red-400 to-violet-500 mt-4" id="registerButton" on:click={register}>Register</button>
                <button class="btn bg-gradient-to-r from-violet-300 to-blue-600 mt-4" id="loginButton" on:click={login}>Login</button>

            </div>
        </div>
    </div>


</div>
