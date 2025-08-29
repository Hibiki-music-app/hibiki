export interface SessionType {
	id: string;
	expiresAt: Date;
	token: string;
	createdAt: Date;
	updatedAt: Date;
	ipAddress?: string; // to remove, no logger on this app (i write this just to try better auth)
	userAgent?: string;
	userId: string;
}
