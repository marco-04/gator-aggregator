import { State } from "../../state.js";
import { feeds } from "../schema.js";
import { getUserID } from "./users.js";

export async function createFeed(state: State, name: string, feedURL: string) {
  const userID = await getUserID(state, state.cfg.currentUserName);
  const [result] = await state.db.insert(feeds).values({
    name: name,
    url: feedURL,
    userId: userID,
  }).returning();
  return result;
}

