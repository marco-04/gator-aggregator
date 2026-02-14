import { eq } from "drizzle-orm";
import { State } from "../../state.js";
import { feeds, users } from "../schema.js";
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

export async function listFeeds(state: State) {
  const result = await state.db.select().from(feeds).leftJoin(users, eq(users.id, feeds.userId));
  return result;
}

