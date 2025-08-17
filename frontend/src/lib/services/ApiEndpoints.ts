
export const API_URL: string = 'https://dab.yeet.su/api'; // keep this in .env?

const url = (path: string) => `${API_URL}${path}`;

export const ApiEndpoints = {
    auth: {
        login: () => url('auth/login'), // example
    }
}