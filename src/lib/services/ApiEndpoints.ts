
export const API_URL: string = 'https://dab.yeet.su/api'; // keep this in .env?

const url = (path: string) => `${API_URL}${path}`;

// endpoint of api, (only use to call the api)
export const ApiEndpoints = {
    auth: {
        login: () => url('auth/login'), // example
    }
}

// the callback routes to render to client
export const ClientRouter = {
    login: "/auth/login",
    signup: "/auth/signup",
    logout: "/auth/logout",
    profile: "/profile",
}