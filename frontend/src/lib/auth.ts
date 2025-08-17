import { startRegistration, startAuthentication } from '@simplewebauthn/browser';
import { ApiEndpoints } from '../services/api-endpoints';
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
		const response = await fetch(ApiEndpoints.passkey.registerStart(), {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username })
		});
		if (!response.ok) {
			const msg = await response.json();
			throw new Error('User already exists or failed to get registration: ' + msg);
		}
		const options = await response.json();

		const attestationResponse = await startRegistration(options.publicKey); //passkey modal

		// validation server side
		const verificationResponse = await fetch(ApiEndpoints.passkey.registerFinish(), {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
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
		const response = await fetch(ApiEndpoints.passkey.loginStart(), {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify({ username })
		});
		if (!response.ok) {
			const msg = await response.json();
			throw new Error('Failed to get login options: ' + msg);
		}
		const options = await response.json();

		const assertionResponse = await startAuthentication(options.publicKey);

		const verificationResponse = await fetch(ApiEndpoints.passkey.loginFinish(), {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(assertionResponse)
		});
		const msg = await verificationResponse.json();
		showMessage(msg, !verificationResponse.ok);
	} catch (error) {
		showMessage(`Error: ${(error as Error).message}`, true);
	}
}
