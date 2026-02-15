import { asc, eq, sql } from "drizzle-orm";
import { State } from "../../state.js";
import { Feed, feeds, users } from "../schema.js";

export async function createFeed(state: State, userID: string, name: string, feedURL: string) {
  const [result] = await state.db.insert(feeds).values({
    name: name,
    url: feedURL,
    userId: userID,
  }).returning();
  return result;
}

export async function getFeed(state: State, feedURL: string) {
  const [result] = await state.db.select().from(feeds).where(eq(feeds.url, feedURL));
  return result;
}

export async function listFeeds(state: State) {
  const result = await state.db.select().from(feeds).leftJoin(users, eq(users.id, feeds.userId));
  return result;
}

export async function markFeedFetched(state: State, feedID: string) {
  await state.db.update(feeds).set({ lastFetchedAt: new Date() }).where(eq(feeds.id, feedID));
}

export async function getNextToFetch(state: State) {
  const [result] = await state.db.execute(sql`SELECT * FROM ${feeds} ORDER BY ${feeds.lastFetchedAt} ASC NULLS FIRST LIMIT 1`);
  return result as Feed;
}

