import { State } from "../../state.js";
import { feed_follows, feeds, users } from "../schema.js";
import { getFeed } from "./feeds.js";
import { and, eq } from "drizzle-orm";

export async function createFeedFollow(state: State, userID: string, feedURL: string) {
  const feed = await getFeed(state, feedURL);
  const [result] = await state.db.insert(feed_follows).values({
    userId: userID,
    feedId: feed.id
  }).returning();
  return result;
}

export async function deleteFeedFollow(state: State, userID: string, feedURL: string) {
  const feed = await getFeed(state, feedURL);
  await state.db.delete(feed_follows).where(and(eq(feed_follows.userId, userID), eq(feed_follows.feedId, feed.id)));
}

export async function getFollowsForUser(state: State, userID: string) {
  const result = await state.db.select()
    .from(feed_follows)
    .innerJoin(users, eq(users.id, feed_follows.userId))
    .innerJoin(feeds, eq(feeds.id, feed_follows.feedId))
    .where(eq(users.id, userID));
  return result;
}

