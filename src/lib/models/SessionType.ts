export interface SessionType {
	id: string;
	userId: string;
	expiresAt: Date;
	createdAt: Date;
	updatedAt: Date;
	token: string;
	ipAddress?: string | null | undefined; // to remove no ip tracking fck this (betterAuth)
	userAgent?: string | null | undefined;
}

