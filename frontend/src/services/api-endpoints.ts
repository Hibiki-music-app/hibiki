
const API_URL: string = 'http://localhost:8080/api/'; // keep this in .env?

const url = (path: string) => `${API_URL}${path}`;

export const ApiEndpoints = {
    passkey: {
        registerStart: () => url('passkey/registerStart'),
        registerFinish: () => url('passkey/registerFinish'),
        loginStart: () => url('passkey/loginStart'),
        loginFinish: () => url('passkey/loginFinish'),
    }
}