import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from '../index';
import { sveltekitCookies } from 'better-auth/svelte-kit';
import { getRequestEvent } from '$app/server';
import { user, account, session, verification } from '$lib/db/schema';

// better Auth instance & pluggins
export const auth = betterAuth({
	account: {
		accountLinking: {
			enabled: true,
			trustedProviders: ['google'],
		}
	},
	database: drizzleAdapter(db, {
		provider: "pg",
		schema: {user, account, session, verification}
	}),
  emailAndPassword: {
		enabled: true,
		autoSignIn: false,
	},
	socialProviders: {
		google: {
			clientId: process.env.GOOGLE_CLIENT_ID as string,
			clientSecret: process.env.GOOGLE_CLIENT_SECRET as string,
		}
	},
	plugins: [sveltekitCookies(getRequestEvent)], // put in last
	});
