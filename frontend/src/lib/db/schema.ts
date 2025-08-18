import { integer, pgTable, varchar } from "drizzle-orm/pg-core";

// testing drizzle setup (to change !!)
export const usersTable = pgTable("users", {
	id: integer().primaryKey().generatedAlwaysAsIdentity(),
	name: varchar({ length: 255 }).notNull(),
	age: integer().notNull(),
	email: varchar({ length: 255 }).notNull().unique(),
});
