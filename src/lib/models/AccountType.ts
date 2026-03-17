export interface AccountType {
	id: string;
	accountId: string;
	providerId: string;
	userId: string;
	accessToken?: string;
	refreshToken?: string;
	idToken?: string;
	accessTokenExpiresAt?: Date;
	refreshTokenExpiresAt?: Date;
	scope?: string;
	password?: string;
	createdAt: Date;
	updatedAt: Date;
}
