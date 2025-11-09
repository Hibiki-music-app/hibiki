// the doc : https://sixnine-dotdev.github.io/dab-api-docs/#/
export const API_URL: string = 'https://dab.yeet.su/api'; // deprecated now

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

export const RandomAvatar = 'https://avatar.iran.liara.run/public'

export const AppName = 'Hibiki:0.0.1';

export const EmailRequest = 'testdev@gmail.com';

export const MBApiRoot = 'https://musicbrainz.org/ws/2';