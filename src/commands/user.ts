import { DrizzleQueryError } from "drizzle-orm";
import { createUser, getUser, getUsers, resetUsers } from "../db/queries/users.js";
import { State } from "../state.js";

export async function commandLogin(state: State, ...args: string[]) {
  const userName = args[0];

  const user = await getUser(state, userName);
  if (user === undefined) {
    throw new Error(`User ${userName} not registered`);
  }

  state.cfg.currentUserName = user.name;
  state.cfg.write();
  console.log(`${userName} set as current user`);
}

export async function commandRegister(state: State, ...args: string[]) {
  const userName = args[0];

  try {
    await createUser(state, userName);
  } catch(err) {
    if (err instanceof DrizzleQueryError) {
      throw new Error(`Register error: ${userName} already exists`);
    } else {
      throw new Error(`DB error: ${(err as Error).message}`);
    }
  }

  await commandLogin(state, userName);
}

export async function commandUsers(state: State) {
  const results = await getUsers(state);
  for (const user of results) {
    console.log(`* ${user.name}${user.name === state.cfg.currentUserName ? " (current)" : ""}`);
  }
}

export async function commandReset(state: State) {
  await resetUsers(state);
}

