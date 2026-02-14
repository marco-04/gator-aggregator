import { eq } from "drizzle-orm";
import { State } from "../../state.js";
import { users } from "../schema.js";

export async function createUser(state: State, name: string) {
  const [result] = await state.db.insert(users).values({ name: name }).returning();
  return result;
}

export async function getUser(state: State, name: string) {
  const [result] = await state.db.select().from(users).where(eq(users.name, name));
  return result;
}

export async function getUsers(state: State) {
  const result = await state.db.select().from(users);
  return result;
}

export async function resetUsers(state: State) {
  const [result] = await state.db.delete(users);
  return result;
}

