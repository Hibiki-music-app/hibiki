import { betterAuth } from "better-auth";
import { drizzleAdapter } from "better-auth/adapters/drizzle";
import { db } from '../index';


// better Auth instance & pluggins
export const auth = betterAuth({
	account: {
		accountLinking: {
			enabled: true,
			trustedProviders: ['google'],
		}
	},
	database: drizzleAdapter(db, {
		provider: "pg"
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
	}
	});
