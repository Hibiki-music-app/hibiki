import { betterAuth } from 'better-auth';
import { prismaAdapter } from 'better-auth/adapters/prisma';
import { passkey } from 'better-auth/plugins/passkey';
import { prisma } from '$lib/db/prisma';

const BASE_URL = process.env.BETTER_AUTH_URL ?? 'http://localhost:5173';

export const auth = betterAuth({
	baseURL: BASE_URL,
	secret: process.env.BETTER_AUTH_SECRET,
	database: prismaAdapter(prisma, {
		provider: 'postgresql',
	}),
	emailAndPassword: {
		enabled: true,
	},
	plugins: [
		passkey({
			rpID: new URL(BASE_URL).hostname,
			rpName: 'Hibiki',
		}),
	],
});
