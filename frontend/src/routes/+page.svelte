<script lang="ts">
    import { startRegistration, startAuthentication } from '@simplewebauthn/browser';
    import {loginFinish, loginStart, registerFinish, registerStart} from "../services/api-endpoints";
    let username: string = '';
    let message: string = '';
    let isError: boolean = false;

    // to refactor by creating a global context like react and call this "popup"
    function showMessage(msg: string, error: boolean = false) {
        message = msg;
        isError = error;
    }

    async function register() {
        try {
            // inscription server side
            const response = await fetch(registerStart, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json'},
                body: JSON.stringify({username})
            });
            if (!response.ok) {
                const msg = await response.json();
                throw new Error('User already exists or failed to get registration: ' + msg);
            }
            const options = await response.json();

            const attestationResponse = await startRegistration(options.publicKey); //passkey modal

            // validation server side
            const verificationResponse = await fetch(registerFinish, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json'},
                body: JSON.stringify(attestationResponse)
            });

            const msg = await verificationResponse.json();
            showMessage(msg, !verificationResponse.ok);
        } catch (error) {
            showMessage(`Error: ${(error as Error).message}`, true);
        }
    }

    async function login() {
        try {
            const response = await fetch(loginStart, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json'},
                body: JSON.stringify({username})
            });
            if (!response.ok) {
                const msg = await response.json()
                throw new Error('Failed to get login options: ' + msg);
            }
            const options = await response.json();

            const assertionResponse = await startAuthentication(options.publicKey);

            const verificationResponse = await fetch(loginFinish, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json'},
                body: JSON.stringify(assertionResponse)
            });
            const msg = await verificationResponse.json();
            showMessage(msg, !verificationResponse.ok);

        } catch (error) {
            showMessage(`Error: ${(error as Error).message}`, true);
        }
    }
</script>

<div class="hero bg-base-200 min-h-screen">
    <div class="hero-content flex-col lg:flex-row-reverse">
        <div class="text-center lg:text-left">
            <h1 class="text-5xl font-bold">Hibiki connect!</h1>
            <p class="py-6">
                Connect you without email and password, insane non?
            </p>
            <a href="/private">TEST</a>
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

